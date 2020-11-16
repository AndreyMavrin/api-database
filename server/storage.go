package server

import "park_2020/api-database/models"

func CheckUser(nickname string) bool {
	var count int
	models.DB.QueryRow(`SELECT COUNT(nickname) FROM users WHERE nickname=$1;`, nickname).Scan(&count)
	return count > 0
}

func InsertUser(user models.User) error {
	_, err := models.DB.Exec(`INSERT INTO users(about, email, fullname, nickname) VALUES ($1, $2, $3, $4);`,
		user.About, user.Email, user.Fullname, user.Nickname)
	return err
}
