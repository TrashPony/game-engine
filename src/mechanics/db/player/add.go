package player

import (
	"github.com/TrashPony/game_engine/src/dbConnect"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/player"
	"log"
)

func AddNewPlayer(newPlayer *player.Player, userID int) {
	tx, err := dbConnect.GetDBConnect().Begin()
	defer tx.Rollback()

	err = tx.QueryRow(`INSERT INTO players (name, user_id) VALUES ($1, $2) RETURNING id`,
		newPlayer.Login, userID).Scan(&newPlayer.ID)
	if err != nil {
		log.Fatal("registration new player " + err.Error())
	}

	err = tx.Commit()

	UpdateUser(newPlayer) // что бы убрать все null
}
