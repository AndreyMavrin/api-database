package server

import (
	"encoding/json"
	"log"
	"net/http"
	"park_2020/api-database/models"
)

func CreateForum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var forum models.Forum
	err := json.NewDecoder(r.Body).Decode(&forum)
	if err != nil {
		log.Println(err)
		return
	}

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
