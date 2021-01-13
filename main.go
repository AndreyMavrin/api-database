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

	router.HandleFunc("/health", server.HealthHandler)
	router.HandleFunc("/api/user/{nickname}/create", server.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/api/user/{nickname}/profile", server.UserProfile).Methods(http.MethodGet, http.MethodPost)

	router.HandleFunc("/api/forum/create", server.CreateForum).Methods(http.MethodPost)
	router.HandleFunc("/api/forum/{slug}/details", server.ForumDetails).Methods(http.MethodGet)
	router.HandleFunc("/api/forum/{slug}/create", server.CreateForumSlug).Methods(http.MethodPost)
	router.HandleFunc("/api/forum/{slug}/threads", server.ForumThreads).Methods(http.MethodGet)
	router.HandleFunc("/api/forum/{slug}/users", server.ForumUsers).Methods(http.MethodGet)

	router.HandleFunc("/api/thread/{slug_or_id}/create", server.CreatePosts).Methods(http.MethodPost)
	router.HandleFunc("/api/thread/{slug_or_id}/vote", server.VoteThread).Methods(http.MethodPost)
	router.HandleFunc("/api/thread/{slug_or_id}/details", server.ThreadDetails).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/api/thread/{slug_or_id}/posts", server.ThreadPosts).Methods(http.MethodGet)

	router.HandleFunc("/api/post/{id}/details", server.PostDetails).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/api/service/status", server.StatusHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/service/clear", server.ClearHandler).Methods(http.MethodPost)

	http.Handle("/", router)

	fmt.Println("Starting server at: 5000")
	http.ListenAndServe(":5000", nil)
}
