package server

import (
	"park_2020/api-database/models"
)

func InsertVote(vote models.Vote) error {
	_, err := models.DB.Exec(`INSERT INTO votes(nickname, voice, thread) VALUES ($1, $2, NULLIF($3, 0));`, vote.Nickname, vote.Voice, vote.Thread)
	return err
}

func UpdateVote(vote models.Vote) error {
	_, err := models.DB.Exec(`UPDATE votes SET voice=$1 WHERE nickname=$2 AND thread=$3;`, vote.Voice, vote.Nickname, vote.Thread)
	if err != nil {
		return err
	}
	return nil
}
