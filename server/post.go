package server

import (
	"encoding/json"
	"log"
	"net/http"
	"park_2020/api-database/models"
	"strconv"
	"strings"
)

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
			thread, err := SelectThreadByAuthor(post.Author)
			if err != nil {
				log.Println(err)
				return
			}
			post.Thread = 1
			post.ID = 1
			post.Forum = thread.Forum

			postsCreated = append(postsCreated, post)
		}

	} else {
		for _, post := range posts {
			thread, err := SelectThreadByAuthor(post.Author)
			if err != nil {
				log.Println(err)
				return
			}
			post.Thread = id
			post.ID = 1
			post.Forum = thread.Forum

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
