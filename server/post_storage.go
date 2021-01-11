package server

import (
	"database/sql"
	"park_2020/api-database/models"
)

func InsertPost(post models.Post) (models.Post, error) {
	row := models.DB.QueryRow(`INSERT INTO posts(author, created, forum, message, parent, thread) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;`,
		post.Author, post.Created, post.Forum, post.Message, post.Parent, post.Thread)
	var p models.Post
	err := row.Scan(&p.Author, &p.Created, &p.Forum, &p.ID, &p.IsEdited, &p.Message, &p.Parent, &p.Thread, &p.Path)
	return p, err
}

func SelectPosts(threadID int, limit, since int, sort string, desc bool) ([]models.Post, error) {
	var posts []models.Post
	var rows *sql.Rows
	var err error

	if since == 0 {
		if sort == "flat" || sort == "" {
			if desc {
				rows, err = models.DB.Query(`SELECT * FROM posts WHERE thread=$1 ORDER BY created DESC, id DESC LIMIT $2;`, threadID, limit)
			} else {
				rows, err = models.DB.Query(`SELECT * FROM posts WHERE thread=$1 ORDER BY created ASC, id LIMIT $2;`, threadID, limit)
			}
		} else if sort == "tree" {
			if desc {
				rows, err = models.DB.Query(`SELECT * FROM posts WHERE thread=$1 ORDER BY path DESC, id  DESC LIMIT $2;`, threadID, limit)
			} else {
				rows, err = models.DB.Query(`SELECT * FROM posts WHERE thread=$1 ORDER BY path, id LIMIT $2;`, threadID, limit)
			}
		} else {
			if desc {
				rows, err = models.DB.Query(`SELECT * FROM posts WHERE path[1] IN
				(SELECT id FROM posts WHERE thread=$1 AND parent IS NULL ORDER BY id DESC LIMIT $2)
				ORDER BY path[1] DESC, path, id;`, threadID, limit)
			} else {
				rows, err = models.DB.Query(`SELECT * FROM posts WHERE path[1] IN
				(SELECT id FROM posts WHERE thread=$1 AND parent IS NULL ORDER BY id LIMIT $2)
				ORDER BY path, id;`, threadID, limit)
			}
		}
	} else {
		if sort == "flat" || sort == "" {
			if desc {
				rows, err = models.DB.Query(`SELECT * FROM posts WHERE thread=$1 AND id < $2
				ORDER BY created DESC, id DESC LIMIT $3;`, threadID, since, limit)
			} else {
				rows, err = models.DB.Query(`SELECT * FROM posts WHERE thread=$1 AND id > $2
				ORDER BY created ASC, id LIMIT $3;`, threadID, since, limit)
			}
		} else if sort == "tree" {
			if desc {
				rows, err = models.DB.Query(`SELECT * FROM posts WHERE thread=$1 AND PATH < (SELECT path FROM posts WHERE id = $2)
				ORDER BY path DESC, id  DESC LIMIT $3;`, threadID, since, limit)
			} else {
				rows, err = models.DB.Query(`SELECT * FROM posts WHERE thread=$1 AND PATH > (SELECT path FROM posts WHERE id = $2)
				ORDER BY path, id LIMIT $3;`, threadID, since, limit)
			}
		} else {
			if desc {
				rows, err = models.DB.Query(`SELECT * FROM posts WHERE path[1] IN (SELECT id FROM posts WHERE thread=$1 AND parent IS NULL AND PATH[1] <
				(SELECT path[1] FROM posts WHERE id = $2) ORDER BY id DESC LIMIT $3) ORDER BY path[1] DESC, path, id;`, threadID, since, limit)
			} else {
				rows, err = models.DB.Query(`SELECT * FROM posts WHERE path[1] IN (SELECT id FROM posts WHERE thread=$1 AND parent IS NULL AND PATH[1] >
				(SELECT path[1] FROM posts WHERE id = $2) ORDER BY id LIMIT $3) ORDER BY path, id;`, threadID, since, limit)
			}
		}
	}

	if err != nil {
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Post
		err = rows.Scan(&p.Author, &p.Created, &p.Forum, &p.ID, &p.IsEdited, &p.Message, &p.Parent, &p.Thread, &p.Path)
		if err != nil {
			return posts, err
		}

		posts = append(posts, p)
	}
	return posts, nil
}

func SelectPostByID(id int) (models.Post, error) {
	var post models.Post
	row := models.DB.QueryRow(`SELECT author, created, forum, id, is_edited, message, parent, thread FROM posts WHERE id = $1;`, id)
	err := row.Scan(&post.Author, &post.Created, &post.Forum, &post.ID, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)
	if err != nil {
		return post, err
	}
	return post, nil
}

func UpdatePost(post models.Post, postUpdate models.PostUpdate) (models.Post, error) {
	if postUpdate.Message != "" && postUpdate.Message != post.Message {
		row := models.DB.QueryRow(`UPDATE posts SET message=$1, is_edited=true WHERE message=$2 RETURNING *;`, postUpdate.Message, post.Message)
		err := row.Scan(&post.Author, &post.Created, &post.Forum, &post.ID, &post.IsEdited, &post.Message, &post.Parent, &post.Thread, &post.Path)
		if err != nil {
			return post, err
		}
	}
	return post, nil
}
