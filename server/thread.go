package server

import (
	"encoding/json"
	"log"
	"net/http"
	"park_2020/api-database/models"
	"strconv"
	"strings"
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

	if !CheckForum(forum) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonToMessage("Can't find forum"))
		return
	}

	threads, err := SelectThreads(forum, since, limit, desc)
	if err != nil {
		log.Println(err)
		return
	}

	body, err := json.Marshal(threads)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if len(threads) != 0 {
		w.Write(body)
	} else {
		w.Write([]byte("[]"))
	}
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

	if !CheckUserByNickname(vote.Nickname) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonToMessage("Can't find user"))
		return
	}

	var thread models.Thread
	id, errInt := strconv.Atoi(slugOrID)
	if errInt != nil {
		slug := slugOrID

		thread, err = SelectThread(slug)
		if err != nil {
			log.Println(err)
			return
		}

	} else {
		thread, err = SelectThreadByID(id)
		if err != nil {
			log.Println(err)
			return
		}
	}

	vote.Thread = thread.ID

	if !CheckVote(vote.Nickname) {
		err = InsertVote(vote)
		if err != nil {
			log.Println(err)
			return
		}
		thread.Votes += vote.Voice
	} else {
		lastVote, err := LastVote(vote.Nickname)
		if err != nil {
			log.Println(err)
			return
		}

		err = UpdateVote(vote)
		if err != nil {
			log.Println(err)
			return
		}

		vote, err := SelectVote(vote.Nickname)
		if err != nil {
			log.Println(err)
			return
		}

		if lastVote != vote.Voice {
			if vote.Voice == 1 {
				thread.Votes += 2
			} else {
				thread.Votes -= 2
			}
		}
	}

	err = AddVoiceToThread(thread, thread.Votes)
	if err != nil {
		log.Println(err)
		return
	}

	thread, err = SelectThread(thread.Slug)
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

func ThreadDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	RequestUrl := r.URL.Path
	RequestUrl = strings.TrimPrefix(RequestUrl, "/api/thread/")
	slugOrID := strings.TrimSuffix(RequestUrl, "/details")

	if r.Method == "GET" {
		var thread models.Thread

		id, errInt := strconv.Atoi(slugOrID)
		if errInt != nil {
			slug := slugOrID
			if !CheckThread(slug) {
				w.WriteHeader(http.StatusNotFound)
				w.Write(jsonToMessage("Can't find thread"))
				return
			}

			var err error
			thread, err = SelectThread(slug)
			if err != nil {
				log.Println(err)
				return
			}

		} else {
			var err error
			thread, err = SelectThreadByID(id)
			if err != nil {
				log.Println(err)
				return
			}
		}

		body, err := json.Marshal(thread)
		if err != nil {
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body)
		return
	}

	var thread models.Thread
	id, errInt := strconv.Atoi(slugOrID)
	if errInt != nil {
		slug := slugOrID
		if !CheckThread(slug) {
			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonToMessage("Can't find thread"))
			return
		}

		var err error
		thread, err = SelectThread(slug)
		if err != nil {
			log.Println(err)
			return
		}

	} else {
		var err error
		thread, err = SelectThreadByID(id)
		if err != nil {
			log.Println(err)
			return
		}
	}

	var threadUpdate models.ThreadUpdate
	err := json.NewDecoder(r.Body).Decode(&threadUpdate)
	if err != nil {
		log.Println(err)
		return
	}

	err = UpdateThread(thread, threadUpdate)
	if err != nil {
		log.Println(err)
		return
	}

	thread, err = SelectThread(thread.Slug)
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
