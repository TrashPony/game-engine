package collisions

import (
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/dynamic_map_object"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/physical_model"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"math"
)

func BodyCheckCollisionsOnStaticMap(ph *physical_model.PhysicalModel, x, y int, mp *_map.Map, bigGrid bool) bool {

	if x < 0 || y < 0 || x > mp.XSize || y > mp.YSize || ph.Fly {
		return false
	}

	if !bigGrid {
		// приравниваем все к сетке из дискретных клеток, а то кеш будет просто невьебенно огромен
		x = x / _const.CellSize
		y = y / _const.CellSize
	}

	x = x*_const.CellSize + _const.CellSize/2
	y = y*_const.CellSize + _const.CellSize/2

	collision, _, _, _ := CircleObjectsMap(x, y, ph.GetRadius(), nil, mp, false, 0, 0, true)
	if collision {
		return false
	}

	return true
}

func CheckCollisionsPlayersByRect(mUserRect *game_math.Polygon, mpID, unitID int, units []*unit.Unit) *unit.Unit {
	for _, otherUnit := range units {

		if otherUnit == nil || mpID != otherUnit.MapID || otherUnit.GetPhysicalModel().IsFly() {
			continue
		}

		centerX, centerY := mUserRect.GetCenter()
		if game_math.GetBetweenDist(int(centerX), int(centerY), otherUnit.GetPhysicalModel().GetX(), otherUnit.GetPhysicalModel().GetY()) > 250 {
			continue
		}

		if unitID != otherUnit.GetID() {
			userRect := GetBodyRect(otherUnit.GetPhysicalModel(), float64(otherUnit.GetPhysicalModel().GetX()), float64(otherUnit.GetPhysicalModel().GetY()), otherUnit.GetPhysicalModel().GetRotate(), false, false)
			if mUserRect.DetectCollisionRectToRect(userRect) {
				return otherUnit
			}
		}
	}

	return nil
}

func CircleUnit(xCenter, yCenter, radius int, checkUnit *unit.Unit) bool {

	if checkUnit.GetPhysicalModel().IsFly() {
		return false
	}

	if checkUnit.GetPhysicalModel().GetPolygon().DetectCollisionRectToCircle(&game_math.Point{X: float64(xCenter), Y: float64(yCenter)}, radius) {
		return true
	}

	return false
}

func BodyCheckCollisionDynamicObjectsInMove(rect *game_math.Polygon, mp *_map.Map, x, y int) (*dynamic_map_object.Object, *dynamic_map_object.Object) {

	obstacles, mx := mp.GetObstaclesByZoneUnsafe(x, y)
	if mx == nil {
		return nil, nil
	}
	defer mx.RUnlock()

	// проверять колизии по зонам, если происходит колизия доставать обьект по типу и ид, возвращает обьект
	for _, obstacle := range obstacles {

		if obstacle.GetMove() {
			continue
		}

		if rect.DetectCollisionRectToCircle(&game_math.Point{X: float64(obstacle.GetX()), Y: float64(obstacle.GetY())}, obstacle.GetRadius()) {

			var obj *dynamic_map_object.Object
			if obstacle.ParentType == "object" || obstacle.ParentType == "static_object" {
				obj = mp.GetDynamicObjectsByID(int(obstacle.ParentID))
				if obj == nil {
					obj = mp.StaticObjects[int(obstacle.ParentID)]
				}
			}

			fakeObj := dynamic_map_object.Object{Type: "object", Weight: math.MaxInt}
			fakeObj.GetPhysicalModel().SetPos(float64(obstacle.X), float64(obstacle.Y), 0)
			if obj != nil {
				fakeObj.Weight = obj.Weight
			}

			return &fakeObj, obj
		}
	}

	return nil, nil
}
