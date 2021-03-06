package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"park_2020/api-database/models"

	"github.com/jackc/pgx"
)

func jsonToMessage(message string) []byte {
	jsonError, err := json.Marshal(models.Error{Message: message})
	if err != nil {
		return []byte("")
	}
	return jsonError
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
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

	err = InsertUser(user)
	if err != nil {
		users, err := SelectUsers(user.Email, user.Nickname)
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

	body, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func UserProfile(w http.ResponseWriter, r *http.Request) {
	RequestUrl := r.URL.Path
	RequestUrl = strings.TrimPrefix(RequestUrl, "/api/user/")
	nickname := strings.TrimSuffix(RequestUrl, "/profile")

	if r.Method == "GET" {
		user, err := SelectUserByNickname(nickname)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonToMessage("Can't find user"))
			return
		}

		body, err := json.Marshal(user)
		if err != nil {
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body)
		return
	}

	var userUpdate models.User
	err := json.NewDecoder(r.Body).Decode(&userUpdate)
	if err != nil {
		log.Println(err)
		return
	}

	userUpdate.Nickname = nickname
	user, err := UpdateUser(userUpdate)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			w.WriteHeader(http.StatusConflict)
			w.Write(jsonToMessage("This email is already registered"))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonToMessage("Can't find user"))
		return
	}

	body, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
