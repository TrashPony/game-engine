package collisions

import (
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/dynamic_map_object"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/physical_model"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
)

func BodyCheckCollisionsOnStaticMap(ph *physical_model.PhysicalModel, x, y int, mp *_map.Map, bigGrid bool) bool {

	if x < 0 || y < 0 || x > mp.XSize || y > mp.YSize {
		return false
	}

	if !bigGrid {
		// приравниваем все к сетке из дискретных клеток, а то кеш будет просто невьебенно огромен
		x = x / _const.CellSize
		y = y / _const.CellSize
	}

	x = x*_const.CellSize + _const.CellSize/2
	y = y*_const.CellSize + _const.CellSize/2

	collision, _ := checkCircleObjectsMap(mp, x, y, ph.GetRadius())
	if collision {
		return false
	}

	return true
}

func CheckCollisionsPlayersByRect(mUserRect *game_math.Polygon, mpID, unitID int, units []*unit.Unit) *unit.Unit {
	// TODO data race, но если использовать геттеры и сетеры то производительность падает в десятки раз
	for _, otherUnit := range units {

		if otherUnit == nil || mpID != otherUnit.MapID {
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
	rect := GetBodyRect(checkUnit.GetPhysicalModel(), float64(checkUnit.GetPhysicalModel().GetX()), float64(checkUnit.GetPhysicalModel().GetY()),
		checkUnit.GetPhysicalModel().GetRotate(), false, false)
	if rect.DetectCollisionRectToCircle(&game_math.Point{X: float64(xCenter), Y: float64(yCenter)}, radius) {
		return true
	}

	return false
}

func BodyCheckCollisionDynamicObjectsInMove(rect *game_math.Polygon, mp *_map.Map, x, y int, mapObjects map[int]*dynamic_map_object.Object) *dynamic_map_object.Object {

	obstacles, mx := mp.GetObstaclesByZoneUnsafe(x, y)
	if mx == nil {
		return nil
	}
	defer mx.RUnlock()

	// проверять колизии по зонам, если происходит колизия доставать обьект по типу и ид, возвращает обьект
	for _, obstacle := range obstacles {

		if obstacle.GetMove() || obstacle.GetResource() {
			continue
		}

		if rect.DetectCollisionRectToCircle(&game_math.Point{X: float64(obstacle.GetX()), Y: float64(obstacle.GetY())}, obstacle.GetRadius()) {

			if obstacle.GetParentType() == "object" {
				obj := mapObjects[obstacle.GetParentID()]
				if obj != nil {
					return obj
				}
			}

			fakeObj := dynamic_map_object.Object{Type: "static"}
			fakeObj.GetPhysicalModel().SetPos(float64(obstacle.X), float64(obstacle.Y), 0)
			return &fakeObj
		}
	}

	return nil
}
