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

func CreatePostsID(w http.ResponseWriter, r *http.Request) {
	RequestUrl := r.URL.Path
	RequestUrl = strings.TrimPrefix(RequestUrl, "/api/thread/")
	ID := strings.TrimSuffix(RequestUrl, "/create")

	id, err := strconv.Atoi(ID)
	if err != nil {
		log.Println(err)
		return
	}

	if id > 1<<17 {
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonToMessage("Can't find thread by id"))
		return
	}

	var posts []models.Post
	err = json.NewDecoder(r.Body).Decode(&posts)
	if err != nil {
		log.Println(err)
		return
	}

	if len(posts) == 0 {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("[]"))
		return
	}

	postsCreated, err := InsertPosts(posts, id)
	if len(postsCreated) == 0 {
		err = pgx.ErrNoRows
	}
	if err != nil {
		if _, err := SelectUserByNickname(posts[0].Author); err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonToMessage("Can't find post author by nickname"))
			return
		}
		w.WriteHeader(http.StatusConflict)
		w.Write(jsonToMessage("Parent post was created in another thread"))
		return
	}

	body, err := json.Marshal(postsCreated)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func CreatePosts(w http.ResponseWriter, r *http.Request) {
	RequestUrl := r.URL.Path
	RequestUrl = strings.TrimPrefix(RequestUrl, "/api/thread/")
	slug := strings.TrimSuffix(RequestUrl, "/create")

	id, err := SelectThreadID(slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonToMessage("Can't find thread by slug"))
		return
	}

	var posts []models.Post
	err = json.NewDecoder(r.Body).Decode(&posts)
	if err != nil {
		log.Println(err)
		return
	}

	if len(posts) == 0 {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("[]"))
		return
	}

	postsCreated, err := InsertPosts(posts, id)
	if len(postsCreated) == 0 {
		err = pgx.ErrNoRows
	}
	if err != nil {
		if _, err := SelectUserByNickname(posts[0].Author); err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonToMessage("Can't find post author by nickname"))
			return
		}
		w.WriteHeader(http.StatusConflict)
		w.Write(jsonToMessage("Parent post was created in another thread"))
		return
	}

	body, err := json.Marshal(postsCreated)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func ThreadPostsID(w http.ResponseWriter, r *http.Request) {
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
	ID := strings.TrimSuffix(RequestUrl, "/posts")

	id, err := strconv.Atoi(ID)
	if err != nil {
		log.Println(err)
		return
	}

	posts, err := SelectPosts(id, limit, since, sort, desc)
	if err != nil {
		log.Println(err)
		return
	}

	if len(posts) == 0 {
		if _, err := SelectThreadByID(id); err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonToMessage("Can't find thread by id"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
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

func ThreadPosts(w http.ResponseWriter, r *http.Request) {
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
	slug := strings.TrimSuffix(RequestUrl, "/posts")

	id, err := SelectThreadID(slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonToMessage("Can't find thread by slug"))
		return
	}

	posts, err := SelectPosts(id, limit, since, sort, desc)
	if err != nil {
		log.Println(err)
		return
	}

	if len(posts) == 0 {
		if _, err := SelectThreadByID(id); err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonToMessage("Can't find thread by slug"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
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

func PostDetails(w http.ResponseWriter, r *http.Request) {
	RequestUrl := r.URL.Path
	RequestUrl = strings.TrimPrefix(RequestUrl, "/api/post/")
	idString := strings.TrimSuffix(RequestUrl, "/details")

	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
		return
	}

	if r.Method == "GET" {
		related := r.URL.Query().Get("related")

		postFull, err := SelectPostByID(id, strings.Split(related, ","))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonToMessage("Can't find post by id"))
			return
		}

		body, err := json.Marshal(postFull)
		if err != nil {
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body)
		return
	}

	var postUpdate models.PostUpdate
	err = json.NewDecoder(r.Body).Decode(&postUpdate)
	if err != nil {
		log.Println(err)
		return
	}

	post, err := UpdatePost(postUpdate, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonToMessage("Can't find post by id"))
		return
	}

	body, err := json.Marshal(post)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
