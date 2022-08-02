package collisions

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/dynamic_map_object"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
)

func checkCircleObjectsMap(mp *_map.Map, x, y, radius int) (bool, int) {
	collision, typeCollision, idCollision, _ := CircleObjectsMap(x, y, radius, nil, mp, false, 0, 0, true)
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
