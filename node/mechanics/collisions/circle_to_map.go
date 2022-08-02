package collisions

import (
	"fmt"
	units2 "github.com/TrashPony/game-engine/node/mechanics/factories/units"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/dynamic_map_object"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/obstacle_point"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
)

func CircleAllCollisionCheck(xCenter, yCenter, radius int, zCenter float64, mp *_map.Map, excludeObjID, excludeUnitID []int, flore bool,
	units []*unit.Unit) (bool, string, int, *obstacle_point.ObstaclePoint) {

	_, _, lvl := mp.GetPosLevel(xCenter, yCenter)
	if flore {
		if zCenter < lvl {
			// ¯\_(ツ)_/¯
			return true, "flore", 0, nil // удар об землю
		}
	}

	// обьекты на карте
	collision, typeCollision, idCollision, obstacle := CircleObjectsMap(xCenter, yCenter, radius, excludeObjID, mp, true, zCenter, lvl, false)
	if collision {
		return true, typeCollision, idCollision, obstacle
	}

	// все юниты
	collision, gameUnit, obstacle := CircleUnits(xCenter, yCenter, radius, excludeUnitID, mp.Id, zCenter, lvl, units)
	if collision {
		return true, "unit", gameUnit.GetID(), obstacle
	}

	return false, "", 0, nil
}

func CircleObjectsMap(xCenter, yCenter, radius int, excludeObjID []int, mp *_map.Map, moveDestroy bool, zCenter, mpLvl float64, onlyNotMove bool) (bool, string, int, *obstacle_point.ObstaclePoint) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in CircleObjectsMap ", r)
		}
	}()

	obstacles, mx := mp.GetObstaclesByZoneUnsafe(xCenter, yCenter)
	if mx == nil {
		return false, "", 0, nil
	}
	defer mx.RUnlock()

	// динамические обьекты
	for _, obstacle := range obstacles {
		if CircleObject(obstacle, xCenter, yCenter, radius, excludeObjID, mp, zCenter, mpLvl, onlyNotMove) {
			return true, obstacle.GetParentType(), obstacle.GetParentID(), obstacle
		}
	}

	return false, "", 0, nil
}

func CircleObject(obstacle *obstacle_point.ObstaclePoint, xCenter, yCenter, radius int, excludeObjID []int, mp *_map.Map, zCenter, mpLvl float64, onlyNotMove bool) bool {

	if obstacle == nil {
		return false
	}

	if (obstacle.GetParentType() == "object" || obstacle.GetParentType() == "static_object") && contains(excludeObjID, obstacle.GetParentID()) {
		return false
	}

	if zCenter > (mpLvl + obstacle.GetHeight()) {
		// круг выше обьекта, поэтому колизи нет
		return false
	}

	if onlyNotMove && obstacle.Move {
		// если поин проходим и мы ищем только непроходимые
		return false
	}

	distToObstacle := game_math.GetBetweenDist(xCenter, yCenter, obstacle.GetX(), obstacle.GetY())

	if int(distToObstacle) < obstacle.GetRadius()+radius {

		//if obstacle.GetParentType() == "object" {
		//	obj := mp.GetDynamicObjectsByID(obstacle.GetParentID())
		//	if obj == nil {
		//		return false
		//	}
		//}

		return true
	}

	return false
}

func CircleDynamicObj(xCenter, yCenter, radius int, object *dynamic_map_object.Object, onlyNotMove bool) bool {
	for _, sGeoPoint := range object.GetPhysicalModel().GeoData {

		if onlyNotMove && sGeoPoint.GetMove() {
			// если поин проходим и мы ищем только непроходимые
			continue
		}

		distToObstacle := game_math.GetBetweenDist(xCenter, yCenter, sGeoPoint.GetX(), sGeoPoint.GetY())
		if int(distToObstacle) < sGeoPoint.GetRadius()+radius { // если растония меньше чем обра радиуса значит окружности пересекается
			return true
		}
	}

	return false
}
func CircleUnits(xCenter, yCenter, radius int, excludeUnitID []int, mapID int, zCenter, mpLvl float64, units []*unit.Unit) (bool, *unit.Unit, *obstacle_point.ObstaclePoint) {

	if units != nil && len(units) > 0 {

		for _, gameUnit := range units {

			if contains(excludeUnitID, gameUnit.GetID()) {
				continue
			}

			if gameUnit.GetPhysicalModel().Fly {
				continue
			}

			if zCenter > mpLvl+gameUnit.GetPhysicalModel().GetHeight() {
				// ¯\_(ツ)_/¯
				continue
			}

			if int(game_math.GetBetweenDist(xCenter, yCenter, gameUnit.GetPhysicalModel().GetX(), gameUnit.GetPhysicalModel().GetY())) > 150+radius {
				continue
			}

			if CircleUnit(xCenter, yCenter, radius, gameUnit) {
				return true, gameUnit, &obstacle_point.ObstaclePoint{
					X:      int32(gameUnit.GetX()),
					Y:      int32(gameUnit.GetY()),
					Radius: int32(gameUnit.GetPhysicalModel().Radius),
				}
			}
		}

	} else {

		units, mx := units2.Units.GetAllUnitsByMapIDUnsafeRange(mapID)
		defer mx.RUnlock()

		for _, gameUnit := range units {

			if contains(excludeUnitID, gameUnit.GetID()) {
				continue
			}

			if gameUnit.GetPhysicalModel().Fly {
				continue
			}

			if zCenter > mpLvl+gameUnit.GetPhysicalModel().GetHeight() {
				// ¯\_(ツ)_/¯
				continue
			}

			if int(game_math.GetBetweenDist(xCenter, yCenter, gameUnit.GetPhysicalModel().GetX(), gameUnit.GetPhysicalModel().GetY())) > 150+radius {
				continue
			}

			if CircleUnit(xCenter, yCenter, radius, gameUnit) {
				return true, gameUnit, &obstacle_point.ObstaclePoint{
					X:      int32(gameUnit.GetX()),
					Y:      int32(gameUnit.GetY()),
					Radius: int32(gameUnit.GetPhysicalModel().Radius),
				}
			}
		}
	}

	return false, nil, nil
}

func contains(s []int, e int) bool {
	if s == nil {
		return false
	}

	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
