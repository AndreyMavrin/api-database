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
		return
	}

	if len(posts) == 0 {
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
	w.Write(body)
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

/*
08:24:48.299 CRIT Unexpected data (1754-08-30 22:43:41.128654848 +0000 UTC != 0001-01-01 00:00:00 +0000 UTC): Post[0].Created
panic: Unexpected data (1754-08-30 22:43:41.128654848 +0000 UTC != 0001-01-01 00:00:00 +0000 UTC): Post[0].Created

goroutine 16 [running, locked to thread]:
github.com/op/go-logging.(*Logger).Panicf(0xc0005ce180, 0xaab301, 0x1e, 0xc0006dbd40, 0x3, 0x3)
	/var/jenkins_home/workspace/tech-db-forum/vendor/github.com/op/go-logging/logger.go:194 +0x105
github.com/bozaro/tech-db-forum/tests.(*PerfSession).CheckDate(0xc011f10909, 0xc006b83528, 0xc0003cd200, 0xc011fe2bc0, 0xf)
	/var/jenkins_home/workspace/tech-db-forum/tests/perf_session.go:62 +0x2f9
github.com/bozaro/tech-db-forum/tests.(*PPost).Validate(0xc006b83500, 0xb8fb20, 0xc011f10909, 0xc00034bb60, 0x0, 0xc011fe2960, 0x7)
	/var/jenkins_home/workspace/tech-db-forum/tests/check_post_get_one.go:137 +0x5d3
github.com/bozaro/tech-db-forum/tests.PerfThreadGetPostsSuccess.func3(0xb8fb20, 0xc011f10909)
	/var/jenkins_home/workspace/tech-db-forum/tests/check_thread_get_posts.go:489 +0x34b
github.com/bozaro/tech-db-forum/tests.(*PerfSession).Validate(...)
	/var/jenkins_home/workspace/tech-db-forum/tests/perf_session.go:15
github.com/bozaro/tech-db-forum/tests.PerfThreadGetPostsSuccess(0xc00000e220, 0xc02cb0e008)
	/var/jenkins_home/workspace/tech-db-forum/tests/check_thread_get_posts.go:454 +0x52d
github.com/bozaro/tech-db-forum/tests.(*Perf).Run.func1(0xc0241bd438, 0xc00000e220, 0xc0241bd440, 0xc0241bd450)
	/var/jenkins_home/workspace/tech-db-forum/tests/perf.go:87 +0x8b
created by github.com/bozaro/tech-db-forum/tests.(*Perf).Run
	/var/jenkins_home/workspace/tech-db-forum/tests/perf.go:79 +0xeb
*/

/*
10:49:33.901 CRIT Unexpected data (1754-08-30 22:43:41.128654848 +0000 UTC != 0001-01-01 00:00:00 +0000 UTC): Post.Created
panic: Unexpected data (1754-08-30 22:43:41.128654848 +0000 UTC != 0001-01-01 00:00:00 +0000 UTC): Post.Created

goroutine 48 [running, locked to thread]:
github.com/op/go-logging.(*Logger).Panicf(0xc0003256b0, 0xaab301, 0x1e, 0xc007b4e270, 0x3, 0x3)
	/var/jenkins_home/workspace/tech-db-forum/vendor/github.com/op/go-logging/logger.go:194 +0x105
github.com/bozaro/tech-db-forum/tests.(*PerfSession).CheckDate(0xc0072e3fdc, 0xc00bd9bde8, 0xc007b2a9a0, 0xc010e066a0, 0xc)
	/var/jenkins_home/workspace/tech-db-forum/tests/perf_session.go:62 +0x2f9
github.com/bozaro/tech-db-forum/tests.(*PPost).Validate(0xc00bd9bdc0, 0xb8fb20, 0xc0072e3fdc, 0xc007638ae0, 0x0, 0xa90c57, 0x4)
	/var/jenkins_home/workspace/tech-db-forum/tests/check_post_get_one.go:137 +0x5d3
github.com/bozaro/tech-db-forum/tests.PerfPostGetOneSuccess.func1(0xb8fb20, 0xc0072e3fdc)
	/var/jenkins_home/workspace/tech-db-forum/tests/check_post_get_one.go:160 +0xa4
github.com/bozaro/tech-db-forum/tests.(*PerfSession).Validate(...)
	/var/jenkins_home/workspace/tech-db-forum/tests/perf_session.go:15
github.com/bozaro/tech-db-forum/tests.PerfPostGetOneSuccess(0xc00050a200, 0xc0000a0010)
	/var/jenkins_home/workspace/tech-db-forum/tests/check_post_get_one.go:158 +0x2e8
github.com/bozaro/tech-db-forum/tests.(*Perf).Run.func1(0xc02a765b78, 0xc00050a200, 0xc02a765b80, 0xc02a765b90)
	/var/jenkins_home/workspace/tech-db-forum/tests/perf.go:87 +0x8b
created by github.com/bozaro/tech-db-forum/tests.(*Perf).Run
	/var/jenkins_home/workspace/tech-db-forum/tests/perf.go:79 +0xeb

*/
