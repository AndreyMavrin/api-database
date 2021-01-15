package server

import (
	"park_2020/api-database/models"
)

func InsertForum(forum models.Forum) (models.Forum, error) {
	row := models.DB.QueryRow(`INSERT INTO forums(slug, title, username) VALUES ($1, $2, $3) RETURNING *;`,
		forum.Slug, forum.Title, forum.User)
	var f models.Forum
	err := row.Scan(&f.User, &f.Posts, &f.Threads, &f.Slug, &f.Title)
	return f, err
}

func SelectForum(slug string) (models.Forum, error) {
	row := models.DB.QueryRow(`SELECT username, posts, threads, slug, title FROM forums WHERE slug=$1 LIMIT 1;`, slug)
	var f models.Forum
	err := row.Scan(&f.User, &f.Posts, &f.Threads, &f.Slug, &f.Title)
	return f, err
}

func StatusForum() models.Status {
	var status models.Status
	models.DB.QueryRow(`SELECT COUNT(*) FROM users;`).Scan(&status.User)
	models.DB.QueryRow(`SELECT COUNT(*) FROM forums;`).Scan(&status.Forum)
	models.DB.QueryRow(`SELECT COUNT(*) FROM threads;`).Scan(&status.Thread)
	models.DB.QueryRow(`SELECT COUNT(*) FROM posts;`).Scan(&status.Post)
	return status
}

func ClearDB() error {
	var err error
	_, err = models.DB.Exec(`TRUNCATE users, forums, threads, posts, votes, users_forum;`)
	return err
}
