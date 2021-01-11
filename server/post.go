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

	RequestUrl := r.URL.Path
	RequestUrl = strings.TrimPrefix(RequestUrl, "/api/thread/")
	slugOrID := strings.TrimSuffix(RequestUrl, "/create")

	var thread models.Thread
	var err error
	id, errInt := strconv.Atoi(slugOrID)
	if errInt != nil {
		slug := slugOrID

		thread, err = SelectThread(slug)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonToMessage("Can't find post thread by slug"))
			return
		}

	} else {
		thread, err = SelectThreadByID(int32(id))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonToMessage("Can't find post thread by id"))
			return
		}
	}

	var posts []models.Post
	err = json.NewDecoder(r.Body).Decode(&posts)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("[]"))
		return
	}

	var postsCreated []models.Post
	for _, post := range posts {

		post.Thread = thread.ID
		post.Forum = thread.Forum

		post, err = InsertPost(post)
		if err != nil {
			if !CheckUserByNickname(posts[0].Author) {
				w.WriteHeader(http.StatusNotFound)
				w.Write(jsonToMessage("Can't find post author by nickname"))
				return
			}
			w.WriteHeader(http.StatusConflict)
			w.Write(jsonToMessage("Parent post was created in another thread"))
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

	} else {
		thread, err = SelectThreadByID(int32(id))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonToMessage("Can't find thread by id"))
			return
		}
	}

	posts, err := SelectPosts(int(thread.ID), limit, since, sort, desc)
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
	if len(posts) != 0 {
		w.Write(body)
	} else {
		w.Write([]byte("[]"))
	}
}

func PostDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	RequestUrl := r.URL.Path
	RequestUrl = strings.TrimPrefix(RequestUrl, "/api/post/")
	idString := strings.TrimSuffix(RequestUrl, "/details")

	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
		return
	}

	related := r.URL.Query().Get("related")

	if r.Method == "GET" {
		var postFull models.PostFull
		post, err := SelectPostByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonToMessage("Can't find post by id"))
			return
		}

		if strings.Contains(related, "user") {
			user, err := SelectUserByPost(id)
			if err != nil {
				log.Println(err)
				return
			}
			postFull.Author = &user
		}

		if strings.Contains(related, "forum") {
			forum, err := SelectForumByPost(id)
			if err != nil {
				log.Println(err)
				return
			}
			postFull.Forum = &forum
		}

		if strings.Contains(related, "thread") {
			thread, err := SelectThreadByPost(id)
			if err != nil {
				log.Println(err)
				return
			}
			postFull.Thread = &thread
		}

		postFull.Post = &post

		body, err := json.Marshal(postFull)
		if err != nil {
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body)
		return
	}

	post, err := SelectPostByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonToMessage("Can't find post by id"))
		return
	}

	var postUpdate models.PostUpdate
	err = json.NewDecoder(r.Body).Decode(&postUpdate)
	if err != nil {
		log.Println(err)
		return
	}

	post, err = UpdatePost(post, postUpdate)
	if err != nil {
		log.Println(err)
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
