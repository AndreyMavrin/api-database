package server

import "park_2020/api-database/models"

func InsertPost(post models.Post) error {
	_, err := models.DB.Exec(`INSERT INTO posts(author, created, forum, message) VALUES ($1, $2, $3, $4);`,
		post.Author, post.Created, post.Forum, post.Message)
	return err
}

func CheckPost(slug string) bool {
	var count int
	models.DB.QueryRow(`SELECT COUNT(*) FROM posts WHERE slug ILIKE $1;`, slug).Scan(&count)
	return count > 0
}
