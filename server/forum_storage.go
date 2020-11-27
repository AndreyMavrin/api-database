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
