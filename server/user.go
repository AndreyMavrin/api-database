package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"park_2020/api-database/models"
)

func jsonToMessage(message string) []byte {
	jsonError, err := json.Marshal(models.Error{Message: message})
	if err != nil {
		return []byte("")
	}
	return jsonError
}

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

func UserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	RequestUrl := r.URL.Path
	RequestUrl = strings.TrimPrefix(RequestUrl, "/api/user/")
	nickname := strings.TrimSuffix(RequestUrl, "/profile")

	if r.Method == "GET" {
		if !CheckUserByNickname(nickname) {
			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonToMessage("Can't find user"))
			return
		}

		user, err := SelectUserByNickname(nickname)
		if err != nil {
			log.Println(err)
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

	var userUpdate models.UserUpdate
	err := json.NewDecoder(r.Body).Decode(&userUpdate)
	if err != nil {
		log.Println(err)
		return
	}

	if !CheckUserByNickname(nickname) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonToMessage("Can't find user"))
		return
	}

	if CheckUserByEmail(userUpdate.Email) {
		w.WriteHeader(http.StatusConflict)
		w.Write(jsonToMessage("This email is already registered"))
		return
	}

	user, err := SelectUserByNickname(nickname)
	if err != nil {
		log.Println(err)
		return
	}

	err = UpdateUser(user, userUpdate)
	if err != nil {
		log.Println(err)
		return
	}

	userUpdated, err := SelectUserByNickname(nickname)
	if err != nil {
		log.Println(err)
		return
	}

	body, err := json.Marshal(userUpdated)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
	return

}
