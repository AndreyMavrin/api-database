package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func ForumThreads(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 0
	}

	since := r.URL.Query().Get("since")

	desc, err := strconv.ParseBool(r.URL.Query().Get("desc"))
	if err != nil {
		desc = false
	}

	RequestUrl := r.URL.Path
	RequestUrl = strings.TrimPrefix(RequestUrl, "/api/forum/")
	forum := strings.TrimSuffix(RequestUrl, "/threads")

	if !CheckForum(forum) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonToMessage("Can't find forum"))
		return
	}

	threads, err := SelectThreads(forum, since, limit, desc)
	if err != nil {
		log.Println(err)
		return
	}

	body, err := json.Marshal(threads)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)

}
