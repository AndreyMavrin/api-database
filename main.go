package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"park_2020/api-database/models"
	"park_2020/api-database/server"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

// type AccessLogger struct {
// 	LogrusLogger *logrus.Entry
// }

func DBConnection() *sql.DB {
	connString := "host=localhost user=amavrin password=root dbname=forums sslmode=disable"

	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(10)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// func (ac *AccessLogger) AccessLogMiddleware(_ *mux.Router) mux.MiddlewareFunc {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			next.ServeHTTP(w, r)

// 			ac.LogrusLogger.WithFields(logrus.Fields{
// 				"method":      r.Method,
// 				"remote_addr": r.RemoteAddr,
// 			}).Info(r.URL.Path)
// 		})
// 	}
// }

func main() {
	models.DB = DBConnection()

	router := mux.NewRouter()

	// logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true})
	// logrus.WithFields(logrus.Fields{
	// 	"logger": "LOGRUS",
	// }).Info("Starting server")

	// AccessLogOut := new(AccessLogger)

	// contextLogger := logrus.WithFields(logrus.Fields{
	// 	"mode":   "[access_log]",
	// 	"logger": "LOGRUS",
	// })
	// logrus.SetFormatter(&logrus.JSONFormatter{})
	// AccessLogOut.LogrusLogger = contextLogger

	// router.Use(AccessLogOut.AccessLogMiddleware(router))

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
