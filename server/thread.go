package server

import (
	"encoding/json"
	"log"
	"net/http"
	"park_2020/api-database/models"
	"strconv"
	"strings"

	"github.com/jackc/pgx"
)

func ForumThreads(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 0
	}

	since := r.URL.Query().Get("since")

	desc, err := strconv.ParseBool(r.URL.Query().Get("desc"))
	if err != nil {
		desc = false
	}

	RequestUrl := r.URL.Path
	RequestUrl = strings.TrimPrefix(RequestUrl, "/api/forum/")
	forum := strings.TrimSuffix(RequestUrl, "/threads")

	threads, err := SelectThreads(forum, since, limit, desc)
	if err == pgx.ErrNoRows || len(threads) == 0 {
		if !CheckForum(forum) {
			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonToMessage("Can't find forum"))
			return
		}
	}

	if len(threads) == 0 {
		w.Write([]byte("[]"))
		return
	}

	body, err := json.Marshal(threads)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)

}

func VoteThread(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var vote models.Vote
	err := json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		log.Println(err)
		return
	}

	RequestUrl := r.URL.Path
	RequestUrl = strings.TrimPrefix(RequestUrl, "/api/thread/")
	slugOrID := strings.TrimSuffix(RequestUrl, "/vote")

	var thread models.Thread
	id, errInt := strconv.Atoi(slugOrID)
	if errInt != nil {
		slug := slugOrID

		thread, err = SelectThread(slug)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonToMessage("Can't find thread by slug"))
			return
		}
	}

	if thread.ID != 0 {
		id = int(thread.ID)
	}

	vote.Thread = int64(id)
	err = InsertVote(vote)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			err = UpdateVote(vote)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write(jsonToMessage("Can't find thread by slug"))
				return
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonToMessage("Can't find thread by slug"))
			return
		}
	}

	threadUpdate, err := SelectThreadByID(int32(id))
	if err != nil {
		log.Println(err)
		return
	}

	body, err := json.Marshal(threadUpdate)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)

}

func ThreadDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	RequestUrl := r.URL.Path
	RequestUrl = strings.TrimPrefix(RequestUrl, "/api/thread/")
	slugOrID := strings.TrimSuffix(RequestUrl, "/details")

	var thread models.Thread

	id, errInt := strconv.Atoi(slugOrID)
	var err error
	if errInt != nil {
		slug := slugOrID

		thread, err = SelectThread(slug)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonToMessage("Can't find thread"))
			return
		}

	} else {
		thread, err = SelectThreadByID(int32(id))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonToMessage("Can't find thread by id"))
			return
		}
	}

	if r.Method == "GET" {
		body, err := json.Marshal(thread)
		if err != nil {
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body)
		return
	}

	var threadUpdate models.Thread
	err = json.NewDecoder(r.Body).Decode(&threadUpdate)
	if err != nil {
		log.Println(err)
		return
	}

	threadUpdate.ID = thread.ID
	thread, err = UpdateThread(threadUpdate)
	if err != nil {
		log.Println(err)
		return
	}

	body, err := json.Marshal(thread)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
