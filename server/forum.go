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

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status := StatusForum()
	body, err := json.Marshal(status)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func ClearHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := ClearDB()
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("null"))
}

func CreateForum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var forum models.Forum
	err := json.NewDecoder(r.Body).Decode(&forum)
	if err != nil {
		log.Println(err)
		return
	}

	user, err := SelectUserByNickname(forum.User)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonToMessage("Can't find user"))
		return
	}

	forum.User = user.Nickname

	_, err = InsertForum(forum)
	if err != nil {
		forum, err = SelectForum(forum.Slug)
		if err != nil {
			log.Println(err)
			return
		}

		body, err := json.Marshal(forum)
		if err != nil {
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusConflict)
		w.Write(body)
		return
	}

	body, err := json.Marshal(forum)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func ForumDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	RequestUrl := r.URL.Path
	RequestUrl = strings.TrimPrefix(RequestUrl, "/api/forum/")
	slug := strings.TrimSuffix(RequestUrl, "/details")

	forum, err := SelectForum(slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonToMessage("Can't find forum"))
		return
	}

	body, err := json.Marshal(forum)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)

}

func ForumUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	RequestUrl := r.URL.Path
	RequestUrl = strings.TrimPrefix(RequestUrl, "/api/forum/")
	slug := strings.TrimSuffix(RequestUrl, "/users")

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 0
	}

	since := r.URL.Query().Get("since")

	desc, err := strconv.ParseBool(r.URL.Query().Get("desc"))
	if err != nil {
		desc = false
	}

	if !CheckForum(slug) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonToMessage("Can't find forum"))
		return
	}

	users, err := SelectUsersByForum(slug, since, limit, desc)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonToMessage("Can't find forum"))
		return
	}

	if len(users) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	body, err := json.Marshal(users)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func CreateForumSlug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	RequestUrl := r.URL.Path
	RequestUrl = strings.TrimPrefix(RequestUrl, "/api/forum/")
	slug := strings.TrimSuffix(RequestUrl, "/create")

	var thread models.Thread
	err := json.NewDecoder(r.Body).Decode(&thread)
	if err != nil {
		log.Println(err)
		return
	}

	forum, err := SelectForum(slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonToMessage("Can't find thread forum"))
		return
	}

	thread.Forum = forum.Slug
	threadInsert, err := InsertThread(thread)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			thread, err := SelectThread(thread.Slug.String)
			if err != nil {
				log.Println(err)
				return
			}

			body, err := json.Marshal(thread)
			if err != nil {
				log.Println(err)
				return
			}

			w.WriteHeader(http.StatusConflict)
			w.Write(body)

			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonToMessage("Can't find thread author"))
		return
	}

	body, err := json.Marshal(threadInsert)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}
