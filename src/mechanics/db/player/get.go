package player

import (
	"database/sql"
	"github.com/TrashPony/game_engine/src/dbConnect"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/player"
	"log"
)

func Player(id int) *player.Player {
	rows, err := dbConnect.GetDBConnect().Query(`SELECT id, name, avatar, biography, title FROM players WHERE id=$1;`, id)
	if err != nil {
		log.Fatal("get players " + err.Error())
	}
	defer rows.Close()

	newUser := player.Player{}

	for rows.Next() {

		var id int
		var name string
		var avatar, biography, title sql.NullString

		err := rows.Scan(&id, &name, &avatar, &biography, &title)
		if err != nil {
			log.Fatal("get user " + err.Error())
		}

		newUser.SetID(id)
		newUser.SetLogin(name)
		getUI(&newUser)
	}

	return &newUser
}

func GetPlayersIDsByUserID(userID int) map[int]bool {
	rows, err := dbConnect.GetDBConnect().Query(`SELECT id FROM players WHERE user_id=$1 `, userID)
	if err != nil {
		log.Fatal("get players by user id " + err.Error())
	}
	defer rows.Close()

	ids := make(map[int]bool)

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			log.Fatal("scan players by user id " + err.Error())
		}

		ids[id] = true
	}

	return ids
}

func getUI(user *player.Player) {
	rows, err := dbConnect.GetDBConnect().Query(`SELECT data FROM player_interface WHERE id_players=$1`, user.GetID())
	if err != nil {
		log.Fatal("get ui " + err.Error())
	}
	defer rows.Close()

	var uiJson []byte

	for rows.Next() {
		err := rows.Scan(&uiJson)
		if err != nil {
			println("get scan ui " + err.Error())
		}

		user.SetWindowStateFromJSON(uiJson)
	}
}
