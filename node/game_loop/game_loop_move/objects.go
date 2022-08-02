package game_loop_move

import (
	battle2 "github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"github.com/TrashPony/game-engine/router/web_socket"
)

func Objects(b *battle2.Battle, units []*unit.Unit, ms *web_socket.MessagesStore) {

	if b.WaitReady {
		return
	}

	moveObjArray := make([]initMoveObj, 0)

	objects, objectsMX := b.Map.UnsafeRangeDynamicObjects()
	for _, obj := range objects {
		if !obj.Static {
			objectsMX.RUnlock()

			if obj.GetPhysicalModel().GetID() == 0 {
				obj.GetPhysicalModel().ID = obj.GetID()
			}

			moveObjArray = append(moveObjArray, obj)

			objectsMX.RLock()
		}
	}
	objectsMX.RUnlock()

	initMove("object", moveObjArray, b, units, ms)
}

func SetObjectsPos(mp *_map.Map) {
	objects, objectsMX := mp.UnsafeRangeDynamicObjects()
	for _, obj := range objects {
		objectsMX.RUnlock()

		if obj.GetPhysicalModel().PosFunc != nil {
			obj.GetPhysicalModel().PosFunc()
			obj.GetPhysicalModel().PosFunc = nil

			mp.RemoveGeoDataObjectsToZone(obj.GetPhysicalModel())
			obj.GetPhysicalModel().GeoData = nil
			obj.SetGeoData()
			mp.AddGeoDataObjectsToZone(obj.GetPhysicalModel())
		}

		objectsMX.RLock()
	}
	objectsMX.RUnlock()
}
