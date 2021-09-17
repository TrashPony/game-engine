package player

import (
	"github.com/TrashPony/game_engine/src/dbConnect"
	"log"
)

func init() {

}

func removeUser(id int) {

	_, err := dbConnect.GetDBConnect().Exec(`DELETE FROM player_interface WHERE id_players = $1`, id)
	if err != nil {
		log.Fatal("delete ui" + err.Error())
	}

	_, err = dbConnect.GetDBConnect().Exec(`DELETE FROM players WHERE id = $1`, id)
	if err != nil {
		log.Fatal("delete user" + err.Error())
	}
}
