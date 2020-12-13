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

func CheckThreadByForum(slug string) bool {
	var count int
	models.DB.QueryRow(`SELECT COUNT(*) FROM forums WHERE slug ILIKE $1;`, slug).Scan(&count)
	return count > 0
}

func CheckThread(slug string) bool {
	var count int
	models.DB.QueryRow(`SELECT COUNT(id) FROM threads WHERE slug ILIKE $1;`, slug).Scan(&count)
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
