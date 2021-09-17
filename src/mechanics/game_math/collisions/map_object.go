package collisions

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/dynamic_map_object"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
)

func CheckObjectCollision(obj *dynamic_map_object.Object, mp *_map.Map, structCheck bool) bool {
	// это супер мега дорогой метод, однако он используется только для популяции растений и вызывается редко

	// -- не за пределами карты
	if obj.GetPhysicalModel().GetX() > mp.XSize || obj.GetPhysicalModel().GetY() > mp.YSize || obj.GetPhysicalModel().GetX() < 0 || obj.GetPhysicalModel().GetY() < 0 {
		return true
	}

	for _, GeoPoint := range obj.GetPhysicalModel().GeoData {
		_, _, lvl := mp.GetPosLevel(obj.GetPhysicalModel().GetX(), obj.GetPhysicalModel().GetY())
		collision, _, _ := CircleAllCollisionCheck(GeoPoint.GetX(), GeoPoint.GetY(), GeoPoint.GetRadius(), lvl, mp, obj.ID, 0, true, nil)
		if collision {
			return true
		}
	}

	// и 100px от дорог, баз, точек выхода из телепорта и самих телепортов
	// TODO эта опция исключительно для растений
	if structCheck {
		for _, sObj := range mp.StaticObjects {
			if sObj.Type == "roads" {
				if game_math.GetBetweenDist(sObj.GetPhysicalModel().GetX(), sObj.GetPhysicalModel().GetY(), obj.GetPhysicalModel().GetX(), obj.GetPhysicalModel().GetY()) < 90 {
					return true
				}
			}
		}
	}

	return false
}
