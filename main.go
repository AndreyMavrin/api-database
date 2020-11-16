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

var conf models.Config

func init() {
	models.LoadConfig(&conf)
}

func DBConnection(conf *models.Config) *sql.DB {
	connString := fmt.Sprintf("host=%v user=%v password=%v dbname=%v sslmode=disable",
		conf.SQLDataBase.Server,
		conf.SQLDataBase.UserID,
		conf.SQLDataBase.Password,
		conf.SQLDataBase.Database,
	)

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

func main() {
	models.DB = DBConnection(&conf)
	router := mux.NewRouter()
	router.HandleFunc("/health", server.HealthHandler)
	router.HandleFunc("/api/user/{nickname}/create", server.CreateUser).Methods(http.MethodPost)
	http.Handle("/", router)

	fmt.Println("Starting server at: 5000")
	http.ListenAndServe(":5000", nil)
}