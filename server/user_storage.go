package server

import (
	"park_2020/api-database/models"
)

func CheckUserByEmail(email string) bool {
	var count int
	models.DB.QueryRow(`SELECT COUNT(id) FROM users WHERE email ILIKE $1;`, email).Scan(&count)
	return count > 0
}

func CheckUserByNickname(nickname string) bool {
	var count int
	models.DB.QueryRow(`SELECT COUNT(id) FROM users WHERE nickname ILIKE $1;`, nickname).Scan(&count)
	return count > 0
}

func InsertUser(user models.User) error {
	_, err := models.DB.Exec(`INSERT INTO users(about, email, fullname, nickname) VALUES ($1, $2, $3, $4);`,
		user.About, user.Email, user.Fullname, user.Nickname)
	return err
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
	if err != nil {
		return u, err
	}
	return u, nil
}

func UpdateUser(user models.User) (models.User, error) {
	if user.About != "" {
		_, err := models.DB.Exec(`UPDATE users SET about=$1 WHERE nickname=$2;`, user.About, user.Nickname)
		if err != nil {
			return user, err
		}
	}
	if user.Email != "" {
		_, err := models.DB.Exec(`UPDATE users SET email=$1 WHERE nickname=$2;`, user.Email, user.Nickname)
		if err != nil {
			return user, err
		}
	}
	if user.Fullname != "" {
		_, err := models.DB.Exec(`UPDATE users SET fullname=$1 WHERE nickname=$2;`, user.Fullname, user.Nickname)
		if err != nil {
			return user, err
		}
	}
	if user.Nickname != "" {
		_, err := models.DB.Exec(`UPDATE users SET nickname=$1 WHERE nickname=$2;`, user.Nickname, user.Nickname)
		if err != nil {
			return user, err
		}
	}
	return user, nil
}
