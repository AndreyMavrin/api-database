package server

import (
	"database/sql"
	"park_2020/api-database/models"
)

func InsertThread(thread models.Thread) (models.Thread, error) {
	row := models.DB.QueryRow(`INSERT INTO threads(author, created, forum, message, slug, title) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;`,
		thread.Author, thread.Created, thread.Forum, thread.Message, thread.Slug, thread.Title)
	var th models.Thread
	err := row.Scan(&th.Author, &th.Created, &th.Forum, &th.ID, &th.Message, &th.Slug, &th.Title, &th.Votes)
	return th, err
}

func CheckThread(slug string) bool {
	var count int
	models.DB.QueryRow(`SELECT COUNT(id) FROM threads WHERE slug ILIKE $1;`, slug).Scan(&count)
	return count > 0
}

func CheckThreadByID(id int) bool {
	var count int
	models.DB.QueryRow(`SELECT COUNT(id) FROM threads WHERE id = $1;`, id).Scan(&count)
	return count > 0
}

func SelectThread(slug string) (models.Thread, error) {
	row := models.DB.QueryRow(`SELECT author, created, forum, id, message, slug, title, votes FROM threads WHERE slug ILIKE $1;`, slug)
	var th models.Thread
	err := row.Scan(&th.Author, &th.Created, &th.Forum, &th.ID, &th.Message, &th.Slug, &th.Title, &th.Votes)
	return th, err
}

func SelectThreadByID(id int) (models.Thread, error) {
	row := models.DB.QueryRow(`SELECT author, created, forum, id, message, slug, title, votes FROM threads WHERE id = $1;`, id)
	var th models.Thread
	err := row.Scan(&th.Author, &th.Created, &th.Forum, &th.ID, &th.Message, &th.Slug, &th.Title, &th.Votes)
	return th, err
}

func SelectThreads(forum, since string, limit int, desc bool) ([]models.Thread, error) {
	var threads []models.Thread
	var rows *sql.Rows
	var err error

	if since != "" {
		if desc {
			rows, err = models.DB.Query(`SELECT author, created, forum, id, message, slug, title, votes FROM threads
		WHERE forum ILIKE $1 AND created <= $2 ORDER BY created DESC LIMIT $3;`, forum, since, limit)
		} else {
			rows, err = models.DB.Query(`SELECT author, created, forum, id, message, slug, title, votes FROM threads
		WHERE forum ILIKE $1 AND created >= $2 ORDER BY created ASC LIMIT $3;`, forum, since, limit)
		}
	} else {
		if desc {
			rows, err = models.DB.Query(`SELECT author, created, forum, id, message, slug, title, votes FROM threads
		WHERE forum ILIKE $1 ORDER BY created DESC LIMIT $2;`, forum, limit)
		} else {
			rows, err = models.DB.Query(`SELECT author, created, forum, id, message, slug, title, votes FROM threads
		WHERE forum ILIKE $1 ORDER BY created ASC LIMIT $2;`, forum, limit)
		}
	}

	if err != nil {
		return threads, err
	}
	defer rows.Close()

	for rows.Next() {
		var th models.Thread
		err = rows.Scan(&th.Author, &th.Created, &th.Forum, &th.ID, &th.Message, &th.Slug, &th.Title, &th.Votes)
		if err != nil {
			continue
		}
		threads = append(threads, th)
	}
	return threads, nil
}

func SelectThreadByPost(id int) (models.Thread, error) {
	var thread models.Thread
	row := models.DB.QueryRow(`SELECT thread FROM posts WHERE id = $1;`, id)
	err := row.Scan(&thread.ID)
	if err != nil {
		return thread, err
	}

	thread, err = SelectThreadByID(thread.ID)
	if err != nil {
		return thread, err
	}

	return thread, nil
}

func UpdateThread(thread models.Thread, threadUpdate models.ThreadUpdate) error {
	if threadUpdate.Message != "" {
		_, err := models.DB.Exec(`UPDATE threads SET message=$1 WHERE message=$2;`, threadUpdate.Message, thread.Message)
		if err != nil {
			return err
		}
	}
	if threadUpdate.Title != "" {
		_, err := models.DB.Exec(`UPDATE threads SET title=$1 WHERE title=$2;`, threadUpdate.Title, thread.Title)
		if err != nil {
			return err
		}
	}
	return nil
}

func InsertVote(vote models.Vote) error {
	_, err := models.DB.Exec(`INSERT INTO votes(nickname, voice, thread) VALUES ($1, $2, NULLIF($3, 0));`, vote.Nickname, vote.Voice, vote.Thread)
	return err
}

func UpdateVote(vote models.Vote) error {
	_, err := models.DB.Exec(`UPDATE votes SET voice=$1 WHERE nickname=$2 AND thread=$3;`, vote.Voice, vote.Nickname, vote.Thread)
	if err != nil {
		return err
	}
	return nil
}
