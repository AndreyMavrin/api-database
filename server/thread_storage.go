package server

import (
	"park_2020/api-database/models"
	"time"

	"github.com/jackc/pgx"
)

func InsertThread(thread models.Thread) (models.Thread, error) {
	var row *pgx.Row
	timeCreated := time.Now()
	if thread.Created == timeCreated {
		row = models.DB.QueryRow(`INSERT INTO threads(author, forum, message, slug, title) VALUES ($1, $2, $3, $4, $5) RETURNING *;`,
			thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Title)
	} else {
		row = models.DB.QueryRow(`INSERT INTO threads(author, created, forum, message, slug, title) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;`,
			thread.Author, thread.Created, thread.Forum, thread.Message, thread.Slug, thread.Title)
	}
	var th models.Thread
	err := row.Scan(&th.Author, &th.Created, &th.Forum, &th.ID, &th.Message, &th.Slug, &th.Title, &th.Votes)
	return th, err
}

func SelectThread(slug string) (models.Thread, error) {
	row := models.DB.QueryRow(`SELECT author, created, forum, id, message, slug, title, votes FROM threads WHERE slug=$1 LIMIT 1;`, slug)
	var th models.Thread
	err := row.Scan(&th.Author, &th.Created, &th.Forum, &th.ID, &th.Message, &th.Slug, &th.Title, &th.Votes)
	return th, err
}

func SelectThreadByID(id int32) (models.Thread, error) {
	row := models.DB.QueryRow(`SELECT author, created, forum, id, message, slug, title, votes FROM threads WHERE id = $1 LIMIT 1;`, id)
	var th models.Thread
	err := row.Scan(&th.Author, &th.Created, &th.Forum, &th.ID, &th.Message, &th.Slug, &th.Title, &th.Votes)
	return th, err
}

func SelectThreads(forum, since string, limit int, desc bool) ([]models.Thread, error) {
	var threads []models.Thread
	var rows *pgx.Rows
	var err error

	if since != "" {
		if desc {
			rows, err = models.DB.Query(`SELECT author, created, forum, id, message, slug, title, votes FROM threads
		WHERE forum=$1 AND created <= $2 ORDER BY created DESC LIMIT $3;`, forum, since, limit)
		} else {
			rows, err = models.DB.Query(`SELECT author, created, forum, id, message, slug, title, votes FROM threads
		WHERE forum=$1 AND created >= $2 ORDER BY created ASC LIMIT $3;`, forum, since, limit)
		}
	} else {
		if desc {
			rows, err = models.DB.Query(`SELECT author, created, forum, id, message, slug, title, votes FROM threads
		WHERE forum=$1 ORDER BY created DESC LIMIT $2;`, forum, limit)
		} else {
			rows, err = models.DB.Query(`SELECT author, created, forum, id, message, slug, title, votes FROM threads
		WHERE forum=$1 ORDER BY created ASC LIMIT $2;`, forum, limit)
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

func UpdateThread(thread models.Thread) (models.Thread, error) {
	row := models.DB.QueryRow(`UPDATE threads SET message=COALESCE(NULLIF($1, ''), message),
		title=COALESCE(NULLIF($2, ''), title) WHERE id = $3 RETURNING *;`, thread.Message, thread.Title, thread.ID)

	var th models.Thread
	err := row.Scan(&th.Author, &th.Created, &th.Forum, &th.ID, &th.Message, &th.Slug, &th.Title, &th.Votes)
	return th, err
}

func InsertVote(vote models.Vote) error {
	_, err := models.DB.Exec(`INSERT INTO votes(nickname, voice, thread) VALUES ($1, $2, NULLIF($3, 0));`, vote.Nickname, vote.Voice, vote.Thread)
	return err
}

func UpdateVote(vote models.Vote) error {
	_, err := models.DB.Exec(`UPDATE votes SET voice=$1 WHERE nickname=$2 AND thread=$3;`, vote.Voice, vote.Nickname, vote.Thread)
	return err
}
