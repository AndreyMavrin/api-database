package server

import (
	"fmt"
	"park_2020/api-database/models"

	"github.com/jackc/pgx"
)

func InsertUser(user models.User) error {
	_, err := models.DB.Exec(`INSERT INTO users(about, email, fullname, nickname) VALUES ($1, $2, $3, $4);`,
		user.About, user.Email, user.Fullname, user.Nickname)
	return err
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
	var rows *pgx.Rows
	var err error

	if desc {
		if since != "" {
			rows, err = models.DB.Query(`SELECT users.about, users.email, users.fullName, users.nickname FROM users
				join users_forum uf on users.nickname = uf.nickname
				WHERE uf.slug=$1 AND uf.nickname < $2 COLLATE "C"
				ORDER BY users.nickname COLLATE "C" DESC LIMIT NULLIF($3, 0);`, slug, since, limit)
		} else {
			rows, err = models.DB.Query(`SELECT users.about, users.email, users.fullName, users.nickname FROM users
				join users_forum uf on users.nickname = uf.nickname
				WHERE uf.slug=$1 ORDER BY users.nickname COLLATE "C" DESC LIMIT NULLIF($2, 0);`, slug, limit)
		}
	} else {
		rows, err = models.DB.Query(`SELECT users.about, users.email, users.fullName, users.nickname FROM users
				join users_forum uf on users.nickname = uf.nickname
				WHERE uf.slug=$1 AND uf.nickname > $2 COLLATE "C"
				ORDER BY users.nickname COLLATE "C" LIMIT NULLIF($3, 0);`, slug, since, limit)
	}

	if err != nil {
		fmt.Println(err)
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
