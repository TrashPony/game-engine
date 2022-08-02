package generators

import (
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/generate_ids"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
)

func CreateBattleMap(mpID int) *_map.Map {

	brMap := getTestMap()
	brMap.Id = generate_ids.GetMapID()
	for obj := range brMap.GetChanMapDynamicObjects() {
		obj.MapID = brMap.Id
	}

	for _, obj := range brMap.StaticObjects {
		obj.MapID = brMap.Id
	}

	FillMapZone(brMap)
	return brMap
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
