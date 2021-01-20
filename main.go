package main

import (
	"fmt"
	"net/http"

	"park_2020/api-database/models"
	"park_2020/api-database/server"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"

	_ "github.com/lib/pq"
)

func contentTypeMiddleware(_ *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	connString := "host=localhost user=amavrin password=root dbname=forums sslmode=disable"
	pgxConn, err := pgx.ParseConnectionString(connString)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	pgxConn.PreferSimpleProtocol = true

	config := pgx.ConnPoolConfig{
		ConnConfig:     pgxConn,
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}

	models.DB, err = pgx.NewConnPool(config)
	if err != nil {
		fmt.Println(err.Error())
	}

	router := mux.NewRouter()

	router.Use(contentTypeMiddleware(router))

	router.HandleFunc("/health", server.HealthHandler)
	router.HandleFunc("/api/user/{nickname}/create", server.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/api/user/{nickname}/profile", server.UserProfile).Methods(http.MethodGet, http.MethodPost)

	router.HandleFunc("/api/forum/create", server.CreateForum).Methods(http.MethodPost)
	router.HandleFunc("/api/forum/{slug}/details", server.ForumDetails).Methods(http.MethodGet)
	router.HandleFunc("/api/forum/{slug}/create", server.CreateForumSlug).Methods(http.MethodPost)
	router.HandleFunc("/api/forum/{slug}/threads", server.ForumThreads).Methods(http.MethodGet)
	router.HandleFunc("/api/forum/{slug}/users", server.ForumUsers).Methods(http.MethodGet)

	router.HandleFunc("/api/thread/{id:[0-9]+}/create", server.CreatePostsID).Methods(http.MethodPost)
	router.HandleFunc("/api/thread/{slug}/create", server.CreatePosts).Methods(http.MethodPost)
	router.HandleFunc("/api/thread/{id:[0-9]+}/vote", server.VoteThreadID).Methods(http.MethodPost)
	router.HandleFunc("/api/thread/{slug}/vote", server.VoteThread).Methods(http.MethodPost)
	router.HandleFunc("/api/thread/{id:[0-9]+}/details", server.ThreadDetailsID).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/api/thread/{slug}/details", server.ThreadDetails).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/api/thread/{id:[0-9]+}/posts", server.ThreadPostsID).Methods(http.MethodGet)
	router.HandleFunc("/api/thread/{slug}/posts", server.ThreadPosts).Methods(http.MethodGet)

	router.HandleFunc("/api/post/{id}/details", server.PostDetails).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/api/service/status", server.StatusHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/service/clear", server.ClearHandler).Methods(http.MethodPost)

	http.Handle("/", router)

	fmt.Println("Starting server at: 5000")
	http.ListenAndServe(":5000", nil)
}
