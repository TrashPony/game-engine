package objects

import (
	"encoding/json"
	"github.com/TrashPony/game_engine/src/dbConnect"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/dynamic_map_object"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/obstacle_point"
	"log"
)

func init() {

}

func AllTypeCoordinate() []*dynamic_map_object.Object {
	rows, err := dbConnect.GetDBConnect().Query(`SELECT id, type, texture_object,
		 animate_sprite_sheets, animate_loop, unit_overlap, object_name, object_description,
		 object_hp, shadow_intensity, animate_speed, shadow, geo_data, x_shadow_offset, y_shadow_offset, height 
		FROM coordinate_type`)
	if err != nil {
		log.Fatal(err.Error() + "get all type object")
	}

	objS := make([]*dynamic_map_object.Object, 0)

	for rows.Next() {
		var obj dynamic_map_object.Object
		var geoData []byte

		err := rows.Scan(&obj.TypeID, &obj.Type, &obj.Texture, &obj.AnimateSpriteSheets, &obj.AnimateLoop,
			&obj.UnitOverlap, &obj.Name, &obj.Description, &obj.TypeMaxHP, &obj.ShadowIntensity,
			&obj.AnimationSpeed, &obj.Shadow, &geoData, &obj.TypeXShadowOffset, &obj.TypeYShadowOffset,
			&obj.HeightType)

		if err != nil {
			log.Fatal(err.Error() + " scan all type coorinate")
		}

		obj.MaxHP = obj.TypeMaxHP
		obj.SetHP(obj.TypeMaxHP)
		obj.XShadowOffset = obj.TypeXShadowOffset
		obj.YShadowOffset = obj.TypeYShadowOffset

		err = json.Unmarshal(geoData, &obj.TypeGeoData)
		if err != nil {
			obj.TypeGeoData = make([]*obstacle_point.ObstaclePoint, 0)
		}

		objS = append(objS, &obj)
	}

	return objS
}

func AddGeoData(x, y, radius int, move bool, obj *dynamic_map_object.Object) {

	if obj.TypeGeoData == nil {
		obj.TypeGeoData = make([]*obstacle_point.ObstaclePoint, 0)
	}
	obj.TypeGeoData = append(obj.TypeGeoData, &obstacle_point.ObstaclePoint{X: int32(x), Y: int32(y), Radius: int32(radius), Move: move})

	jsonString, err := json.Marshal(obj.TypeGeoData)
	if err != nil {
		log.Fatal("geoData to json" + err.Error())
	}

	_, err = dbConnect.GetDBConnect().Exec("UPDATE coordinate_type SET geo_data=$1 WHERE id=$2", jsonString, obj.TypeID)
	if err != nil {
		log.Fatal("update geoData item" + err.Error())
	}
}
