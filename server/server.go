package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"park_2020/api-database/models"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	RequestUrl := r.URL.Path
	RequestUrl = strings.TrimPrefix(RequestUrl, "/api/user/")
	nickname := strings.TrimSuffix(RequestUrl, "/create")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		return
	}
	user.Nickname = nickname

	if CheckUserByEmail(user.Email) || CheckUserByNickname(nickname) {
		users, err := SelectUser(user.Email, user.Nickname)
		if err != nil {
			log.Println(err)
			return
		}

		body, err := json.Marshal(users)
		if err != nil {
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusConflict)
		w.Write(body)
		return
	}

	err = InsertUser(user)
	if err != nil {
		log.Println(err)
		return
	}

	body, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}
