package maps

import (
	"github.com/TrashPony/game_engine/src/dbConnect"
	"github.com/TrashPony/game_engine/src/mechanics/factories/dynamic_objects"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/dynamic_map_object"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/obstacle_point"
	"log"
	"strconv"
)

func Maps() map[int]*_map.Map {
	rows, err := dbConnect.GetDBConnect().Query(`Select id, name, x_size, y_size, level, specification FROM maps`)
	if err != nil {
		log.Fatal(err.Error() + " get all maps")
	}
	defer rows.Close()

	allMap := make(map[int]*_map.Map)

	for rows.Next() {

		mp := &_map.Map{}

		err := rows.Scan(&mp.Id, &mp.Name, &mp.XSize, &mp.YSize, &mp.DefaultLevel, &mp.Specification)
		if err != nil {
			log.Fatal(err.Error() + " scan all maps")
		}

		GetFlore(mp)
		GeoData(mp)
		GetObjects(mp)
		GetLevelMap(mp)
		GetSpawns(mp)

		allMap[mp.Id] = mp
	}

	return allMap
}

func GetSpawns(mp *_map.Map) {

	mp.Spawns = make([]*_map.Spawn, 0)

	rows, err := dbConnect.GetDBConnect().Query(`SELECT id, x, y, radius, rotate FROM maps_spawn WHERE id_map = $1`,
		strconv.Itoa(mp.Id))
	if err != nil {
		log.Fatal(err.Error() + "get map spawn lvl")
	}
	defer rows.Close()

	for rows.Next() {
		var id, x, y, radius, rotate int

		err := rows.Scan(&id, &x, &y, &radius, &rotate)
		if err != nil {
			log.Fatal(err.Error() + "scan map spawn lvl")
		}

		mp.Spawns = append(mp.Spawns, &_map.Spawn{ID: id, X: x, Y: y, Radius: radius, Rotate: rotate})
	}
}

func GetMapByID(id int) *_map.Map {
	rows, err := dbConnect.GetDBConnect().Query(`Select id, name, x_size, y_size, level, specification FROM maps WHERE id = $1`, id)
	if err != nil {
		log.Fatal(err.Error() + " get map by id")
	}
	defer rows.Close()

	for rows.Next() {

		mp := &_map.Map{}

		err := rows.Scan(&mp.Id, &mp.Name, &mp.XSize, &mp.YSize, &mp.DefaultLevel, &mp.Specification)
		if err != nil {
			log.Fatal(err.Error() + " scan map by id")
		}

		if mp.Id == 3 || mp.Id == 6 {
			continue
		}

		mp.TypeID = mp.Id
		GetFlore(mp)
		GeoData(mp)
		GetObjects(mp)
		GetLevelMap(mp)
		GetSpawns(mp)

		return mp
	}

	return nil
}

func GetLevelMap(mp *_map.Map) {
	mp.LevelMap = make(map[string]*_map.LvlMap)

	rows, err := dbConnect.GetDBConnect().Query(`SELECT x, y, lvl
		FROM map_level WHERE id_map = $1`, strconv.Itoa(mp.Id))
	if err != nil {
		log.Fatal(err.Error() + "get map lvl")
	}
	defer rows.Close()

	for rows.Next() {
		var x, y int
		var lvl float64

		err := rows.Scan(&x, &y, &lvl)
		if err != nil {
			log.Fatal(err.Error() + "scan map lvl")
		}

		mp.LevelMap[strconv.Itoa(x)+":"+strconv.Itoa(y)] = &_map.LvlMap{X: x, Y: y, Level: lvl}
	}
}

func GetFlore(mp *_map.Map) {
	mp.Flore = make(map[int]map[int]*dynamic_map_object.Flore)

	rows, err := dbConnect.GetDBConnect().Query(`SELECT x, y, texture_over_flore, texture_priority
		FROM map_constructor WHERE id_map = $1 AND texture_over_flore != ''`, strconv.Itoa(mp.Id))

	if err != nil {
		log.Fatal(err.Error() + "get map flor")
	}

	defer rows.Close()

	for rows.Next() { // заполняем карту значащами клетками
		var flore dynamic_map_object.Flore

		err := rows.Scan(&flore.X, &flore.Y, &flore.TextureOverFlore, &flore.TexturePriority)
		if err != nil {
			log.Fatal(err.Error() + "scan map flor")
		}

		if mp.Flore[flore.X] != nil {
			mp.Flore[flore.X][flore.Y] = &flore
		} else {
			mp.Flore[flore.X] = make(map[int]*dynamic_map_object.Flore)
			mp.Flore[flore.X][flore.Y] = &flore
		}
	}
}

func GetObjects(mp *_map.Map) {

	rows, err := dbConnect.GetDBConnect().Query(`SELECT id_type, x, y,
		scale, rotate, object_priority, x_shadow_offset, y_shadow_offset
		FROM map_constructor WHERE id_map = $1 `, strconv.Itoa(mp.Id))

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

		obj := dynamic_objects.DynamicObjects.GetDynamicObjectByID(typeId, rotate)

		if obj != nil {
			obj.GetPhysicalModel().SetPos(float64(x), float64(y), rotate)
			obj.SetScale(scale)
			obj.Priority = objectPriority
			obj.TypeXShadowOffset += xShadowOffset
			obj.TypeYShadowOffset += yShadowOffset
			obj.SetHP(obj.MaxHP)

			if obj.MaxHP >= -1 {
				obj.HP = -1
				obj.Static = true
				mp.AddDynamicObject(obj)
			} else {
				obj.Static = true
				mp.AddStaticObject(obj)
			}

			obj.SetGeoData()
		}
	}
}

func GeoData(mp *_map.Map) {
	mp.GeoData = make([]*obstacle_point.ObstaclePoint, 0)
	rows, err := dbConnect.GetDBConnect().Query(""+
		"Select "+
		"id, "+
		"x, "+
		"y, "+
		"radius "+
		""+
		"FROM global_geo_data WHERE id_map = $1", mp.Id)
	if err != nil {
		log.Fatal(err.Error() + "db get geo data")
	}

	for rows.Next() { // заполняем карту значащами клетками
		var obstaclePoint obstacle_point.ObstaclePoint

		var id, x, y, radius int
		err := rows.Scan(&id, &x, &y, &radius)

		obstaclePoint.SetID(id)
		obstaclePoint.SetX(x)
		obstaclePoint.SetY(y)
		obstaclePoint.SetRadius(radius)
		obstaclePoint.SetParentType("static_geo")

		mp.GeoData = append(mp.GeoData, &obstaclePoint)
		if err != nil {
			log.Fatal(err.Error() + "scan geo data")
		}
	}
	defer rows.Close()
}
