package server

import (
	"database/sql"
	"park_2020/api-database/models"
)

func InsertPost(post models.Post) (models.Post, error) {
	row := models.DB.QueryRow(`INSERT INTO posts(author, created, forum, message, parent, thread) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;`,
		post.Author, post.Created, post.Forum, post.Message, post.Parent, post.Thread)
	var p models.Post
	err := row.Scan(&p.Author, &p.Created, &p.Forum, &p.ID, &p.Message, &p.Parent, &p.Thread, &p.Path)
	return p, err
}

func SelectPosts(author string, limit, since int, sort string, desc bool) ([]models.Post, error) {
	var posts []models.Post
	var rows *sql.Rows
	var err error

	if since == 0 {
		if sort == "flat" || sort == "" {
			if desc {
				rows, err = models.DB.Query(`SELECT author, created, forum, id, message, parent, thread FROM posts
		WHERE author ILIKE $1 ORDER BY created DESC, id DESC LIMIT $2;`, author, limit)
			} else {
				rows, err = models.DB.Query(`SELECT author, created, forum, id, message, parent, thread FROM posts
		WHERE author ILIKE $1 ORDER BY created ASC, id LIMIT $2;`, author, limit)
			}
		} else if sort == "tree" {
			if desc {
				rows, err = models.DB.Query(`SELECT author, created, forum, id, message, parent, thread FROM posts
		WHERE author ILIKE $1 ORDER BY path DESC, id  DESC LIMIT $2;`, author, limit)
			} else {
				rows, err = models.DB.Query(`SELECT author, created, forum, id, message, parent, thread FROM posts
		WHERE author ILIKE $1 ORDER BY path, id LIMIT $2;`, author, limit)
			}
		} else {
			if desc {
				rows, err = models.DB.Query(`SELECT author, created, forum, id, message, parent, thread FROM posts
			WHERE path[1] IN (SELECT id FROM posts WHERE author = $1 AND parent IS NULL ORDER BY id DESC LIMIT $2)
			ORDER BY path[1] DESC, path, id;`, author, limit)
			} else {
				rows, err = models.DB.Query(`SELECT author, created, forum, id, message, parent, thread FROM posts
			WHERE path[1] IN (SELECT id FROM posts WHERE author = $1 AND parent IS NULL ORDER BY id LIMIT $2)
			ORDER BY path, id;`, author, limit)
			}
		}
	} else {
		if sort == "flat" || sort == "" {
			if desc {
				rows, err = models.DB.Query(`SELECT author, created, forum, id, message, parent, thread FROM posts
			WHERE author ILIKE $1 AND id < $2 ORDER BY created DESC, id DESC LIMIT $3;`, author, since, limit)
			} else {
				rows, err = models.DB.Query(`SELECT author, created, forum, id, message, parent, thread FROM posts
			WHERE author ILIKE $1 AND id > $2 ORDER BY created ASC, id LIMIT $3;`, author, since, limit)
			}
		} else if sort == "tree" {
			if desc {
				rows, err = models.DB.Query(`SELECT author, created, forum, id, message, parent, thread FROM posts
			WHERE author ILIKE $1 AND PATH < (SELECT path FROM posts WHERE id = $2)
			ORDER BY path DESC, id  DESC LIMIT $3;`, author, since, limit)
			} else {
				rows, err = models.DB.Query(`SELECT author, created, forum, id, message, parent, thread FROM posts
			WHERE author ILIKE $1 AND PATH > (SELECT path FROM posts WHERE id = $2)
			ORDER BY path, id LIMIT $3;`, author, since, limit)
			}
		} else {
			if desc {
				rows, err = models.DB.Query(`SELECT author, created, forum, id, message, parent, thread FROM posts
				WHERE path[1] IN (SELECT id FROM posts WHERE author = $1 AND parent IS NULL AND PATH[1] <
				(SELECT path[1] FROM posts WHERE id = $2) ORDER BY id DESC LIMIT $3) ORDER BY path[1] DESC, path, id;`, author, since, limit)
			} else {
				rows, err = models.DB.Query(`SELECT author, created, forum, id, message, parent, thread FROM posts
				WHERE path[1] IN (SELECT id FROM posts WHERE author = $1 AND parent IS NULL AND PATH[1] >
				(SELECT path[1] FROM posts WHERE id = $2) ORDER BY id LIMIT $3) ORDER BY path, id;`, author, since, limit)
			}
		}
	}

	if err != nil {
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Post
		err = rows.Scan(&p.Author, &p.Created, &p.Forum, &p.ID, &p.Message, &p.Parent, &p.Thread)
		if err != nil {
			return posts, err
		}

		posts = append(posts, p)
	}
	return posts, nil
}
