package server

import (
	"database/sql"
	"park_2020/api-database/models"
)

func InsertPost(post models.Post) error {
	_, err := models.DB.Exec(`INSERT INTO posts(author, created, forum, message, parent, thread) VALUES ($1, $2, $3, $4, nullif($5, 0), $6);`,
		post.Author, post.Created, post.Forum, post.Message, post.Parent, post.Thread)
	return err
}

func SelectPosts(author string, limit, since int, sort string, desc bool) ([]models.Post, error) {
	var posts []models.Post
	var rows *sql.Rows
	var err error

	if sort == "flat" {
		if desc {
			rows, err = models.DB.Query(`SELECT author, created, forum, message, COALESCE(parent, 0) FROM posts
		WHERE author ILIKE $1 ORDER BY created DESC, id LIMIT $2;`, author, limit)
		} else {
			rows, err = models.DB.Query(`SELECT author, created, forum, message, COALESCE(parent, 0) FROM posts
		WHERE author ILIKE $1 ORDER BY created ASC, id LIMIT $2;`, author, limit)
		}
	} else {
		if desc {
			rows, err = models.DB.Query(`SELECT author, created, forum, message, COALESCE(parent, 0) FROM posts
		WHERE author ILIKE $1 ORDER BY path DESC, id  DESC LIMIT $2;`, author, limit)
		} else {
			rows, err = models.DB.Query(`SELECT author, created, forum, message, COALESCE(parent, 0) FROM posts
		WHERE author ILIKE $1 ORDER BY path, id LIMIT $2;`, author, limit)
		}
	}

	if err != nil {
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Post
		err = rows.Scan(&p.Author, &p.Created, &p.Forum, &p.Message, &p.Parent)
		if err != nil {
			continue
		}
		p.ID = 1
		p.Thread = 1

		posts = append(posts, p)
	}
	return posts, nil
}
