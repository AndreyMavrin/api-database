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

	// mu := &sync.Mutex{}
	var postsCreated []models.Post
	for _, post := range posts {
		thread, err := SelectThreadByAuthor(post.Author)
		if err != nil {
			log.Println(err)
			return
		}
		post.Thread = 1
		post.ID = 1
		post.Forum = thread.Forum

		err = InsertPost(post)
		if err != nil {
			log.Println(err)
			return
		}

		postsCreated = append(postsCreated, post)

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

func ThreadPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 0
	}

	since, err := strconv.Atoi(r.URL.Query().Get("since"))
	if err != nil {
		since = 0
	}
	sort := r.URL.Query().Get("sort")

	desc, err := strconv.ParseBool(r.URL.Query().Get("desc"))
	if err != nil {
		desc = false
	}

	RequestUrl := r.URL.Path
	RequestUrl = strings.TrimPrefix(RequestUrl, "/api/thread/")
	slugOrID := strings.TrimSuffix(RequestUrl, "/posts")

	_, errInt := strconv.Atoi(slugOrID)
	if errInt != nil {
		slug := slugOrID
		if !CheckThread(slug) {
			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonToMessage("Can't find thread"))
			return
		}

		thread, err := SelectThread(slug)
		if err != nil {
			log.Println(err)
			return
		}

		posts, err := SelectPosts(thread.Author, limit, since, sort, desc)
		if err != nil {
			log.Println(err)
			return
		}

		body, err := json.Marshal(posts)
		if err != nil {
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body)

	}

}
