package game_loop_move

import (
	collisions2 "github.com/TrashPony/game-engine/node/mechanics/collisions"
	units2 "github.com/TrashPony/game-engine/node/mechanics/factories/units"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	battle2 "github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/dynamic_map_object"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"github.com/TrashPony/game-engine/router/web_socket"
	"math"
)

func checkMoveCollision(typeObj string, id int, moveObject moveObject, x, y, rotate float64, b *battle2.Battle, units []*unit.Unit, ms *web_socket.MessagesStore) bool {
	if moveObject.IsFly() { // летающие обьекты не сталкиваются
		return true
	}

	collision := false

	collisionUnit := collisionToUnit(typeObj, id, moveObject, x, y, rotate, b.Map, units)
	if collisionUnit != nil {

		collisionUnit.SetFindMovePath()
		updatePath(b.Map.Id, id, typeObj)

		if typeObj == "unit" {
			gameUnit := units2.Units.GetUnitByIDAndMapID(id, b.Map.Id)
			if gameUnit != nil {
				gameUnit.SetFindMovePath()
			}
		}

		x2, y2 := collisionUnit.GetPhysicalModel().GetRealPos()
		game_math.CollisionReactionBallBall(
			moveObject, collisionUnit.GetPhysicalModel(), moveObject.GetWeight(),
			collisionUnit.GetPhysicalModel().GetWeight(),
			moveObject.GetPowerFactor(), collisionUnit.GetPhysicalModel().GetPowerFactor(), x2, y2)

		collision = true
	}

	fakeObj, obj := collisionToObject(typeObj, id, moveObject, x, y, rotate, b.Map)
	if obj != nil {

		if typeObj == "unit" {
			gameUnit := units2.Units.GetUnitByIDAndMapID(id, b.Map.Id)
			if gameUnit != nil {
				gameUnit.SetFindMovePath()
			}
		}

		x2, y2 := fakeObj.GetPhysicalModel().GetRealPos()
		game_math.CollisionReactionBallBall(moveObject, obj.GetPhysicalModel(),
			moveObject.GetWeight(), obj.GetPhysicalModel().GetWeight(), moveObject.GetPowerFactor(),
			obj.GetPhysicalModel().GetPowerFactor(), x2, y2)

		collision = true
	}

	border := collisionBorderSector(moveObject, x, y, rotate, b.Map)
	if border != nil {
		x2, y2 := border.GetPhysicalModel().GetRealPos()
		game_math.CollisionReactionBallBall(moveObject, border.GetPhysicalModel(),
			moveObject.GetWeight(), math.MaxInt32, moveObject.GetPowerFactor(),
			border.GetPhysicalModel().GetPowerFactor(), x2, y2)
	}

	return !collision
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
		moveObject.GetPhysicalModel().GetNextPolygon().UpdateCenterRect(x, y, moveObject.GetLength(), moveObject.GetWidth())
		moveObject.GetPhysicalModel().GetNextPolygon().Rotate(rotate)
		return collisions2.CheckCollisionsPlayersByRect(moveObject.GetPhysicalModel().GetNextPolygon(), mp.Id, id, units)
	}

	if typeObj == "object" {

		for _, obstaclePoint := range moveObject.GetGeoData() {
			collision, gameUnit, _ := collisions2.CircleUnits(
				obstaclePoint.GetX(),
				obstaclePoint.GetY(),
				obstaclePoint.GetRadius(),
				nil, mp.Id, 0, 0, units)
			if collision {
				return gameUnit
			}
		}
	}

	return nil
}

func collisionToObject(typeObj string, id int, moveObject moveObject, x, y, rotate float64, mp *_map.Map) (*dynamic_map_object.Object, *dynamic_map_object.Object) {

	if typeObj == "unit" {
		moveObject.GetPhysicalModel().GetNextPolygon().UpdateCenterRect(x, y, moveObject.GetLength(), moveObject.GetWidth())
		moveObject.GetPhysicalModel().GetNextPolygon().Rotate(rotate)
		return collisions2.BodyCheckCollisionDynamicObjectsInMove(moveObject.GetPhysicalModel().GetNextPolygon(), mp, int(x), int(y))
	}

	if typeObj == "object" {

		for _, obstaclePoint := range moveObject.GetGeoData() {

			collision, _, _, obstacle := collisions2.CircleObjectsMap(
				obstaclePoint.GetX(),
				obstaclePoint.GetY(),
				obstaclePoint.GetRadius(),
				[]int{id}, mp, true, 0, 0, true)

			if collision {

				var obj *dynamic_map_object.Object
				if obstacle.ParentType == "object" || obstacle.ParentType == "static_object" {
					obj = mp.GetDynamicObjectsByID(int(obstacle.ParentID))
					if obj == nil {
						obj = mp.StaticObjects[int(obstacle.ParentID)]
					}
				}

				fakeObj := dynamic_map_object.Object{Type: "object", Weight: math.MaxInt}
				fakeObj.GetPhysicalModel().SetPos(float64(obstaclePoint.X), float64(obstaclePoint.Y), 0)
				if obj != nil {
					fakeObj.Weight = obj.Weight
				}

				return &fakeObj, obj
			}
		}
	}

	return nil, nil
}

func updatePath(mpId, id int, typeObj string) {
	if typeObj == "unit" {
		gameUnit := units2.Units.GetUnitByIDAndMapID(id, mpId)
		if gameUnit != nil {
			gameUnit.SetFindMovePath()
		}
	}
}
