package server

import (
	"park_2020/api-database/models"
)

func InsertForum(forum models.Forum) error {
	_, err := models.DB.Exec(`INSERT INTO forums(slug, title, username) VALUES ($1, $2, $3);`,
		forum.Slug, forum.Title, forum.User)
	return err
}

func CheckForum(slug string) bool {
	var count int
	models.DB.QueryRow(`SELECT COUNT(id) FROM forums WHERE slug ILIKE $1;`, slug).Scan(&count)
	return count > 0
}

func SelectForum(slug string) (models.Forum, error) {
	row := models.DB.QueryRow(`SELECT slug, title, username FROM forums WHERE slug ILIKE $1;`, slug)
	var f models.Forum
	err := row.Scan(&f.Slug, &f.Title, &f.User)
	return f, err
}

func InsertThread(thread models.Thread) error {
	_, err := models.DB.Exec(`INSERT INTO threads(author, created, forum, message, slug, title) VALUES ($1, $2, $3, $4, $5, $6);`,
		thread.Author, thread.Created, thread.Forum, thread.Message, thread.Slug, thread.Title)
	return err
}

func CheckThreadByAuthor(author string) bool {
	var count int
	models.DB.QueryRow(`SELECT COUNT(id) FROM users WHERE author ILIKE $1;`, author).Scan(&count)
	return count > 0
}

func CheckThreadByForum(slug string) bool {
	var count int
	models.DB.QueryRow(`SELECT COUNT(id) FROM forums WHERE slug ILIKE $1;`, slug).Scan(&count)
	return count > 0
}

func CheckThread(slug string) bool {
	var count int
	models.DB.QueryRow(`SELECT COUNT(id) FROM threads WHERE slug ILIKE $1;`, slug).Scan(&count)
	return count > 0
}

func SelectThread(slug string) (models.Thread, error) {
	row := models.DB.QueryRow(`SELECT author, created, forum, message, slug, title FROM threads WHERE slug ILIKE $1;`, slug)
	var th models.Thread
	err := row.Scan(&th.Author, &th.Created, &th.Forum, &th.Message, &th.Slug, &th.Title)
	return th, err
}
