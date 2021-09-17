package game_loop_move

import (
	"fmt"
	units2 "github.com/TrashPony/game_engine/src/mechanics/factories/units"
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_math/collisions"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/dynamic_map_object"
	_map "github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
)

func checkMoveCollision(typeObj string, id int, moveObject moveObject, x, y, rotate float64, mp *_map.Map, units []*unit.Unit, mapObjects map[int]*dynamic_map_object.Object) bool {

	collisionUnit := collisionToUnit(typeObj, id, moveObject, x, y, rotate, mp, units)
	if collisionUnit != nil {

		if collisionUnit.MovePath != nil {
			collisionUnit.MovePath.NeedFindPath = true
		}

		if typeObj == "unit" {
			gameUnit := units2.Units.GetUnitByIDAndMapID(id, mp.Id)
			if gameUnit != nil && gameUnit.MovePath != nil {
				gameUnit.MovePath.NeedFindPath = true
			}
		}

		damage1, damage2 := game_math.CollisionReactionBallBall(
			moveObject, collisionUnit.GetPhysicalModel(), false,
			moveObject.GetWeight(), collisionUnit.GetPhysicalModel().GetWeight(),
			moveObject.GetPowerFactor(), collisionUnit.GetPhysicalModel().GetPowerFactor())

		// TODO melleDamage(damage1, damage2, moveObject, collisionUnit, mp)
		fmt.Println(damage1, damage2)

		return false
	}

	obj := collisionToObject(typeObj, id, moveObject, x, y, rotate, mp, mapObjects)
	if obj != nil {

		if typeObj == "unit" {
			gameUnit := units2.Units.GetUnitByIDAndMapID(id, mp.Id)
			if gameUnit != nil && gameUnit.MovePath != nil {
				gameUnit.MovePath.NeedFindPath = true
			}
		}

		damage1, damage2 := game_math.CollisionReactionBallBall(moveObject, obj.GetPhysicalModel(), false,
			moveObject.GetWeight(), obj.GetPhysicalModel().GetWeight(), moveObject.GetPowerFactor(), obj.GetPhysicalModel().GetPowerFactor())

		if obj.Type == "static" && typeObj == "unit" {
		} else {
			// TODO melleDamage(damage1, damage2, moveObject, obj, mp)
			fmt.Println(damage1, damage2)
		}

		return false
	}

	border := collisionBorderSector(moveObject, x, y, rotate, mp)
	if border != nil {
		game_math.CollisionReactionBallBall(moveObject, border.GetPhysicalModel(), false,
			moveObject.GetWeight(), border.GetPhysicalModel().GetWeight(), moveObject.GetPowerFactor(), border.GetPhysicalModel().GetPowerFactor())
	}

	return true
}

func collisionBorderSector(moveObject moveObject, x, y, rotate float64, mp *_map.Map) *dynamic_map_object.Object {

	fakeObj := &dynamic_map_object.Object{}

	if x < 25 {
		fakeObj.GetPhysicalModel().SetPos(-100, y, 0)
		return fakeObj
	}

	if x > float64(mp.XSize-25) {
		fakeObj.GetPhysicalModel().SetPos(float64(mp.XSize+100), y, 0)
		return fakeObj
	}

	if y < 25 {
		fakeObj.GetPhysicalModel().SetPos(x, -100, 0)
		return fakeObj
	}

	if y > float64(mp.YSize-25) {
		fakeObj.GetPhysicalModel().SetPos(x, float64(mp.YSize+100), 0)
		return fakeObj
	}

	return nil
}

func collisionToUnit(typeObj string, id int, moveObject moveObject, x, y, rotate float64, mp *_map.Map, units []*unit.Unit) *unit.Unit {

	if typeObj == "unit" {
		rect := game_math.GetCenterRect(x, y, moveObject.GetLength(), moveObject.GetWidth())
		rect.Rotate(rotate)

		excludeUnitID := 0
		excludeUnitID = id

		return collisions.CheckCollisionsPlayersByRect(rect, mp.Id, excludeUnitID, units)
	}

	if typeObj == "object" {

		for _, obstaclePoint := range moveObject.GetGeoData() {
			collision, gameUnit := collisions.CircleUnits(
				obstaclePoint.GetX(),
				obstaclePoint.GetY(),
				obstaclePoint.GetRadius(),
				0, mp.Id, 0, 0, units)
			if collision {
				return gameUnit
			}
		}
	}

	return nil
}

func collisionToObject(typeObj string, id int, moveObject moveObject, x, y, rotate float64, mp *_map.Map, mapObjects map[int]*dynamic_map_object.Object) *dynamic_map_object.Object {

	if typeObj == "unit" {
		rect := game_math.GetCenterRect(x, y, moveObject.GetLength(), moveObject.GetWidth())
		rect.Rotate(rotate)
		return collisions.BodyCheckCollisionDynamicObjectsInMove(rect, mp, int(x), int(y), mapObjects)
	}

	if typeObj == "object" {

		for _, obstaclePoint := range moveObject.GetGeoData() {

			collision, typeCollision, idCollision := collisions.CircleObjectsMap(
				obstaclePoint.GetX(),
				obstaclePoint.GetY(),
				obstaclePoint.GetRadius(),
				id, mp, true, 0, 0, false)

			if collision {

				var obj *dynamic_map_object.Object

				if typeCollision == "object" {
					obj = mp.GetDynamicObjectsByID(idCollision)
					return obj
				}

				if typeCollision == "static_object" {
					obj = mp.StaticObjects[idCollision]
					return obj
				}

				// TODO typeCollision == "static_geo"
			}
		}
	}

	return nil
}
