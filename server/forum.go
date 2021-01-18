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

	forumInserted, err := InsertForum(forum)
	if pgErr, ok := err.(pgx.PgError); ok {
		switch pgErr.Code {
		case "23505":
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
	}

	if err == pgx.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonToMessage("Can't find user"))
		return
	}

	body, err := json.Marshal(forumInserted)
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

	users, err := SelectUsersByForum(slug, since, limit, desc)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonToMessage("Can't find forum"))
		return
	}

	if len(users) == 0 {
		if _, err := SelectForum(slug); err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonToMessage("Can't find forum"))
			return
		}
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
