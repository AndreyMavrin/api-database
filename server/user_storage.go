package server

import (
	"database/sql"
	"park_2020/api-database/models"

	"github.com/lib/pq"
)

func InsertUser(user models.User) error {
	_, err := models.DB.Exec(`INSERT INTO users(about, email, fullname, nickname) VALUES ($1, $2, $3, $4);`,
		user.About, user.Email, user.Fullname, user.Nickname)
	return err
}

func CheckUserByEmail(email string) bool {
	var count int
	models.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE email ILIKE $1;`, email).Scan(&count)
	return count > 0
}

func CheckUserByNickname(nickname string) bool {
	var count int
	models.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE nickname ILIKE $1;`, nickname).Scan(&count)
	return count > 0
}

func SelectUsers(email, nickname string) ([]models.User, error) {
	var users []models.User
	rows, err := models.DB.Query(`SELECT about, email, fullname, nickname FROM users WHERE email ILIKE $1 OR nickname ILIKE $2;`, email, nickname)
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
	row := models.DB.QueryRow(`SELECT about, email, fullname, nickname FROM users WHERE nickname ILIKE $1;`, nickname)
	var u models.User
	err := row.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname)
	return u, err
}

func UpdateUser(user models.User, userUpdate models.UserUpdate) error {
	if userUpdate.About != "" {
		_, err := models.DB.Exec(`UPDATE users SET about=$1 WHERE about=$2;`, userUpdate.About, user.About)
		if err != nil {
			return err
		}
	}
	if userUpdate.Email != "" {
		_, err := models.DB.Exec(`UPDATE users SET email=$1 WHERE email=$2;`, userUpdate.Email, user.Email)
		if err != nil {
			return err
		}
	}
	if userUpdate.Fullname != "" {
		_, err := models.DB.Exec(`UPDATE users SET fullname=$1 WHERE fullname=$2;`, userUpdate.Fullname, user.Fullname)
		if err != nil {
			return err
		}
	}
	return nil
}

func SelectUsersByForum(slug, since string, limit int, desc bool) ([]models.User, error) {
	var users []models.User
	var usernames []string
	var rows *sql.Rows
	var err error
	rows, err = models.DB.Query(`SELECT author FROM threads WHERE forum ILIKE $1 UNION SELECT author FROM posts WHERE forum ILIKE $1;`, slug)
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
