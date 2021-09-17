package collisions

import (
	units2 "github.com/TrashPony/game_engine/src/mechanics/factories/units"
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/dynamic_map_object"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
)

func checkCircleObjectsMap(mp *_map.Map, x, y, radius int) (bool, int) {
	collision, typeCollision, idCollision := CircleObjectsMap(x, y, radius, 0, mp, false, 0, 0, true)
	if collision {
		var obj *dynamic_map_object.Object

		if typeCollision == "object" {
			obj = mp.GetDynamicObjectsByID(idCollision)
			if obj != nil {
				return true, obj.ID
			} else {
				return true, -1
			}
		}

		return true, -1
	}

	return false, -1
}

func checkCollisionsUnits(rect *game_math.Polygon, mapID int, excludeUnitID int) bool {
	for otherUnit := range units2.Units.GetAllUnitsByMapIDChan(mapID) {
		if otherUnit == nil {
			continue
		}

		if mapID != otherUnit.MapID {
			continue
		}

		if excludeUnitID == otherUnit.GetID() {
			continue
		}

		userRect := GetBodyRect(otherUnit.GetPhysicalModel(), float64(otherUnit.GetPhysicalModel().GetX()), float64(otherUnit.GetPhysicalModel().GetY()),
			otherUnit.GetPhysicalModel().GetRotate(), false, false)

		if rect.DetectCollisionRectToRect(userRect) {
			return false
		}
	}

	return true
}
