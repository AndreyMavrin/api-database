package server

import (
	"encoding/json"
	"log"
	"net/http"
	"park_2020/api-database/models"
	"strings"
)

func CreateForum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var forum models.Forum
	err := json.NewDecoder(r.Body).Decode(&forum)
	if err != nil {
		log.Println(err)
		return
	}

	if !CheckUserByNickname(forum.User) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonToMessage("Can't find user"))
		return
	}

	user, err := SelectUserByNickname(forum.User)
	if err != nil {
		log.Println(err)
		return
	}

	forum.User = user.Nickname

	if CheckForum(forum.Slug) {
		forum, err := SelectForum(forum.Slug)
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

	err = InsertForum(forum)
	if err != nil {
		log.Println(err)
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

	if !CheckForum(slug) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonToMessage("Can't find forum"))
		return
	}

	forum, err := SelectForum(slug)
	if err != nil {
		log.Println(err)
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

	if !CheckThreadByAuthor(thread.Author) && !CheckThreadByForum(slug) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonToMessage("Can't find thread"))
		return
	}

	if CheckThread(slug) {
		thread, err := SelectThread(slug)
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

	err = InsertThread(thread)
	if err != nil {
		log.Println(err)
		return
	}

	thread.ID = 1

	body, err := json.Marshal(thread)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}
