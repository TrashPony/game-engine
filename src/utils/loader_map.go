package utils

import (
	"database/sql"
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/dbConnect"
	"github.com/TrashPony/game_engine/src/mechanics/factories/dynamic_objects"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/dynamic_map_object"
	"log"
)

var oldDB *sql.DB

func init() {
	command := "postgres://" +
		_const.Config.GetParams("dbLogin") + ":" +
		_const.Config.GetParams("dbPass") + "@" +
		_const.Config.GetParams("dbAddress") + "/game?sslmode=disable"

	var err error
	oldDB, err = sql.Open("postgres", command)
	if err != nil {
		log.Fatal(err)
	}

	if err = oldDB.Ping(); err != nil {
		log.Panic(err)
	}
}

func LoadMap(oldID, newID int) {

	object := AllTypeCoordinate()
	mapObjects := make(map[int]*dynamic_map_object.Object)
	for _, obj := range object {
		mapObjects[obj.TypeID] = obj
	}

	rows, err := oldDB.Query("SELECT id_type, x, y, "+
		"scale, rotate, object_priority, x_shadow_offset, y_shadow_offset "+
		"FROM map_constructor "+
		"WHERE id_map = $1 ", oldID)

	if err != nil {
		log.Fatal(err.Error() + "get map obj")
	}

	defer rows.Close()

	for rows.Next() { // заполняем карту значащами клетками

		var typeId, x, y, scale, objectPriority, xShadowOffset, yShadowOffset int
		var rotate float64

		err := rows.Scan(&typeId, &x, &y, &scale, &rotate, &objectPriority, &xShadowOffset, &yShadowOffset)
		if err != nil {
			log.Fatal(err.Error() + "scan map obj")
		}

		if typeId > 0 {
			objName := mapObjects[typeId]
			obj := dynamic_objects.DynamicObjects.GetDynamicObjectByTexture(objName.Texture, 0)

			var id int
			err = dbConnect.GetDBConnect().QueryRow(`
			INSERT INTO map_constructor (id_map, id_type, x, y, scale, rotate, object_priority, x_shadow_offset, y_shadow_offset) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
				newID, obj.TypeID, x, y, scale, rotate, objectPriority, xShadowOffset, yShadowOffset).Scan(&id)
			if err != nil {
				log.Fatal("load map " + err.Error())
			}
		}
	}
}

func AllTypeCoordinate() []*dynamic_map_object.Object {
	rows, err := oldDB.Query("SELECT id, texture_object, animate_sprite_sheets FROM coordinate_type")
	if err != nil {
		log.Fatal(err.Error() + "get all type object")
	}

	objS := make([]*dynamic_map_object.Object, 0)

	for rows.Next() {
		var obj dynamic_map_object.Object

		err := rows.Scan(&obj.TypeID, &obj.Texture, &obj.AnimateSpriteSheets)
		if err != nil {
			log.Fatal(err.Error() + " scan all type coorinate")
		}

		objS = append(objS, &obj)
	}

	return objS
}
