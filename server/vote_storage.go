package server

import (
	"park_2020/api-database/models"
)

func CheckVote(user string) bool {
	var count int
	models.DB.QueryRow(`SELECT COUNT(id) FROM votes WHERE nickname ILIKE $1;`, user).Scan(&count)
	return count > 0
}

func InsertVote(vote models.Vote) error {
	_, err := models.DB.Exec(`INSERT INTO votes(nickname, voice) VALUES ($1, $2);`, vote.Nickname, vote.Voice)
	return err
}

func AddVoiceToThread(thread models.Thread, voice int) error {
	_, err := models.DB.Exec(`UPDATE threads SET votes=$1 WHERE slug=$2;`, voice, thread.Slug)
	return err
}

func UpdateVote(vote models.Vote) error {
	_, err := models.DB.Exec(`UPDATE votes SET voice=$1 WHERE nickname=$2;`, vote.Voice, vote.Nickname)
	if err != nil {
		return err
	}
	return nil
}

func SelectVote(nickname string) (models.Vote, error) {
	row := models.DB.QueryRow(`SELECT nickname, voice FROM votes WHERE nickname ILIKE $1;`, nickname)
	var v models.Vote
	err := row.Scan(&v.Nickname, &v.Voice)
	return v, err
}

func LastVote(nickname string) (int, error) {
	var last int
	row := models.DB.QueryRow(`SELECT voice FROM votes WHERE nickname ILIKE $1;`, nickname)
	err := row.Scan(&last)
	return last, err
}
