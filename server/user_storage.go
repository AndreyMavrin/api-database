package server

import (
	"park_2020/api-database/models"

	"github.com/jackc/pgx"
	"github.com/lib/pq"
)

func InsertUser(user models.User) error {
	_, err := models.DB.Exec(`INSERT INTO users(about, email, fullname, nickname) VALUES ($1, $2, $3, $4);`,
		user.About, user.Email, user.Fullname, user.Nickname)
	return err
}

func CheckUserByNickname(nickname string) bool {
	var count int
	models.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE nickname=$1;`, nickname).Scan(&count)
	return count > 0
}

func SelectUsers(email, nickname string) ([]models.User, error) {
	var users []models.User
	rows, err := models.DB.Query(`SELECT about, email, fullname, nickname FROM users WHERE email=$1 OR nickname=$2 LIMIT 2;`, email, nickname)
	if err != nil {
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var u models.User
		err = rows.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname)
		if err != nil {
			continue
		}
		users = append(users, u)
	}
	return users, nil
}

func SelectUserByNickname(nickname string) (models.User, error) {
	row := models.DB.QueryRow(`SELECT about, email, fullname, nickname FROM users WHERE nickname=$1;`, nickname)
	var u models.User
	err := row.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname)
	return u, err
}

func UpdateUser(user models.User) (models.User, error) {
	row := models.DB.QueryRow(`UPDATE users SET about=COALESCE(NULLIF($1, ''), about),
				email=COALESCE(NULLIF($2, ''), email), 	fullname=COALESCE(NULLIF($3, ''), fullname)
				WHERE nickname ILIKE $4 RETURNING *;`, user.About, user.Email, user.Fullname, user.Nickname)

	var u models.User
	err := row.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname)
	return u, err
}

func SelectUsersByForum(slug, since string, limit int, desc bool) ([]models.User, error) {
	var users []models.User
	var usernames []string
	var rows *pgx.Rows
	var err error
	rows, err = models.DB.Query(`SELECT author FROM threads WHERE forum=$1 UNION SELECT author FROM posts WHERE forum=$1;`, slug)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var u string
		err = rows.Scan(&u)
		if err != nil {
			return users, err
		}
		usernames = append(usernames, u)
	}

	if limit == 0 {
		if since == "" {
			if desc {
				rows, err = models.DB.Query(`SELECT about, email, fullname, nickname FROM users WHERE nickname = ANY($1)
			ORDER BY LOWER(nickname) COLLATE "C" DESC;`, pq.Array(usernames))
			} else {
				rows, err = models.DB.Query(`SELECT about, email, fullname, nickname FROM users WHERE nickname = ANY($1)
			ORDER BY LOWER(nickname) COLLATE "C";`, pq.Array(usernames))
			}
		} else {
			if desc {
				rows, err = models.DB.Query(`SELECT about, email, fullname, nickname FROM users WHERE nickname = ANY($1)
			AND LOWER(nickname) < LOWER($2) COLLATE "C" ORDER BY LOWER(nickname) COLLATE "C" DESC;`, pq.Array(usernames), since)
			} else {
				rows, err = models.DB.Query(`SELECT about, email, fullname, nickname FROM users WHERE nickname = ANY($1)
			AND LOWER(nickname) > LOWER($2) COLLATE "C" ORDER BY LOWER(nickname) COLLATE "C";`, pq.Array(usernames), since)
			}
		}
	} else {
		if since == "" {
			if desc {
				rows, err = models.DB.Query(`SELECT about, email, fullname, nickname FROM users WHERE nickname = ANY($1)
			ORDER BY LOWER(nickname) COLLATE "C" DESC LIMIT $2;`, pq.Array(usernames), limit)
			} else {
				rows, err = models.DB.Query(`SELECT about, email, fullname, nickname FROM users WHERE nickname = ANY($1)
			ORDER BY LOWER(nickname) COLLATE "C" LIMIT $2;`, pq.Array(usernames), limit)
			}
		} else {
			if desc {
				rows, err = models.DB.Query(`SELECT about, email, fullname, nickname FROM users WHERE nickname = ANY($1)
			AND LOWER(nickname) < LOWER($2) COLLATE "C" ORDER BY LOWER(nickname) COLLATE "C" DESC LIMIT $3;`, pq.Array(usernames), since, limit)
			} else {
				rows, err = models.DB.Query(`SELECT about, email, fullname, nickname FROM users WHERE nickname = ANY($1)
			AND LOWER(nickname) > LOWER($2) COLLATE "C" ORDER BY LOWER(nickname) COLLATE "C" LIMIT $3;`, pq.Array(usernames), since, limit)
			}
		}
	}

	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var u models.User
		err = rows.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}

	return users, nil
}

func SelectUserByPost(id int) (models.User, error) {
	var user models.User
	row := models.DB.QueryRow(`SELECT author FROM posts WHERE id = $1;`, id)
	var author string
	err := row.Scan(&author)
	if err != nil {
		return user, err
	}

	row = models.DB.QueryRow(`SELECT about, email, fullname, nickname FROM users WHERE nickname=$1;`, author)
	err = row.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
	if err != nil {
		return user, err
	}

	return user, nil
}
