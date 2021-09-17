package user

import (
	"github.com/TrashPony/game_engine/src/dbConnect"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/user"
	"log"
)

func GetUser(id int) *user.User {
	rows, err := dbConnect.GetDBConnect().Query(`SELECT id, name, mail, role, language FROM users WHERE id=$1 `, id)
	if err != nil {
		log.Fatal("get user " + err.Error())
	}
	defer rows.Close()

	newUser := user.User{}

	for rows.Next() {
		var mail string

		err := rows.Scan(&newUser.ID, &newUser.Login, &mail, &newUser.UserRole, &newUser.Language)
		if err != nil {
			log.Fatal("get user " + err.Error())
		}

		newUser.SetEmail(mail)
	}

	return &newUser
}

func UpdateUser(uUser *user.User) {
	_, err := dbConnect.GetDBConnect().Exec(`UPDATE users SET language = $1 where id = $2`,
		uUser.Language, uUser.ID)
	if err != nil {
		println("update user" + err.Error())
	}
}
