package player

import (
	"database/sql"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/player"
	"log"
)

func UpdateUser(p *player.Player) {
	// TODO
	//if p.GetID() < 1 {
	//	return
	//}
	//
	//tx, err := dbConnect.GetDBConnect().Begin()
	//if err != nil || tx == nil {
	//	return
	//}
	//defer tx.Rollback()
	//
	//_, err = tx.Exec("UPDATE players "+
	//	"SET "+
	//	"credits = $2, "+
	//	"training = $3, "+
	//	"last_base_id = $4, "+
	//	"fraction = $5, "+
	//	"avatar = $6, "+
	//	"biography = $7, "+
	//	"scientific_points = $8, "+
	//	"attack_points = $9, "+
	//	"production_points = $10, "+
	//	"title = $11, "+
	//	"fraction_points = $12 "+
	//	"WHERE id = $1",
	//	p.GetID(), p.GetCredits(), p.Training, p.LastBaseID, p.Fraction,
	//	p.GetAvatar(), p.Biography, p.ScientificPoints, p.AttackPoints, p.ProductionPoints, p.Title,
	//	p.GetFractionPoints())
	//if err != nil {
	//	log.Fatal("update user " + err.Error())
	//}
	//
	//UpdateUI(p, tx)
	//
	//tx.Commit()
}

func UpdateUI(p *player.Player, tx *sql.Tx) {
	_, err := tx.Exec(`DELETE FROM player_interface WHERE id_players = $1`, p.GetID())
	if err != nil {
		log.Fatal("delete ui" + err.Error())
	}

	_, err = tx.Exec(`INSERT INTO player_interface (data, id_players) VALUES ($1, $2)`,
		p.GetJSONWindowState(), p.GetID())
	if err != nil {
		log.Fatal("add new ui" + err.Error())
	}
}
