package generators

import (
	_const "github.com/TrashPony/game_engine/src/const"
	dbMap "github.com/TrashPony/game_engine/src/mechanics/db/maps"
	"github.com/TrashPony/game_engine/src/mechanics/factories/maps"
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
)

// создаем сектора с данными о проходимости для быстродействия методов движения
func CreateBattleMap(mpID int) *_map.Map {
	brMap := dbMap.GetMapByID(mpID)
	if brMap != nil {

		brMap.Id = maps.Maps.GetNewMapID()
		for obj := range brMap.GetChanMapDynamicObjects() {
			obj.MapID = brMap.Id
		}

		for _, obj := range brMap.StaticObjects {
			obj.MapID = brMap.Id
		}

		FillMapZone(brMap)
		maps.Maps.AddNewMap(brMap)

		return brMap
	}

	return nil
}

func FillMapZone(mp *_map.Map) {
	mp.GeoZones = make([][]*_map.Zone, mp.XSize/_const.DiscreteSize)

	for x := 0; x < mp.XSize; x += _const.DiscreteSize {

		mp.GeoZones[x/_const.DiscreteSize] = make([]*_map.Zone, mp.YSize/_const.DiscreteSize)

		for y := 0; y < mp.YSize; y += _const.DiscreteSize {

			mp.GeoZones[x/_const.DiscreteSize][y/_const.DiscreteSize] = &_map.Zone{
				Size:      _const.DiscreteSize,
				DiscreteX: x / _const.DiscreteSize,
				DiscreteY: y / _const.DiscreteSize,
			}

			fillMapZone(x, y, mp.GeoZones[x/_const.DiscreteSize][y/_const.DiscreteSize], mp)
		}
	}
}

func fillMapZone(x, y int, zone *_map.Zone, mp *_map.Map) {
	// +50 что бы в область папали пограничные препятсвия
	zone.Rect = game_math.GetRect(float64(x), float64(y), _const.DiscreteSize+50, _const.DiscreteSize+50)

	for i := 0; i < len(mp.GeoData); i++ {
		if zone.Rect.DetectCollisionRectToCircle(&game_math.Point{X: float64(mp.GeoData[i].GetX()), Y: float64(mp.GeoData[i].GetY())}, mp.GeoData[i].GetRadius()) {
			zone.AddObstacle(mp.GeoData[i])
		}
	}

	for _, sObj := range mp.StaticObjects {
		for _, gPoint := range sObj.GetPhysicalModel().GetGeoData() {
			if zone.Rect.DetectCollisionRectToCircle(&game_math.Point{X: float64(gPoint.GetX()), Y: float64(gPoint.GetY())}, gPoint.GetRadius()) {
				zone.AddObstacle(gPoint)
			}
		}
	}

	for obj := range mp.GetChanMapDynamicObjects() {
		for _, gPoint := range obj.GetPhysicalModel().GetGeoData() {
			if zone.Rect.DetectCollisionRectToCircle(&game_math.Point{X: float64(gPoint.GetX()), Y: float64(gPoint.GetY())}, gPoint.GetRadius()) {
				zone.AddObstacle(gPoint)
			}
		}
	}
}
