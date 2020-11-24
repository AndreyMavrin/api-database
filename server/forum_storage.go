package server

import (
	"park_2020/api-database/models"
)

func InsertForum(forum models.Forum) error {
	_, err := models.DB.Exec(`INSERT INTO forums(slug, title, users) VALUES ($1, $2, $3);`,
		forum.Slug, forum.Title, forum.User)
	return err
}
