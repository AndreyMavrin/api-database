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

func CreatePosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var posts []models.Post
	err := json.NewDecoder(r.Body).Decode(&posts)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("[]"))
		return
	}

	RequestUrl := r.URL.Path
	RequestUrl = strings.TrimPrefix(RequestUrl, "/api/thread/")
	slugOrID := strings.TrimSuffix(RequestUrl, "/create")

	var postsCreated []models.Post
	id, errInt := strconv.Atoi(slugOrID)
	if errInt != nil {
		for _, post := range posts {

			slug := slugOrID
			post.Author = slug
			post.ID = 1
			post.Forum = "asdf"
			postsCreated = append(postsCreated, post)
		}

	} else {
		for _, post := range posts {
			post.Thread = id
			post.ID = 1
			post.Forum = "asdf"
			postsCreated = append(postsCreated, post)
		}
	}

	body, err := json.Marshal(postsCreated)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if len(postsCreated) != 0 {
		w.Write(body)
	} else {
		w.Write([]byte("[]"))
	}
}
