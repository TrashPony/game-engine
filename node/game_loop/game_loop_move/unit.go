package game_loop_move

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"github.com/TrashPony/game-engine/router/web_socket"
)

func Unit(mp *_map.Map, units []*unit.Unit, battle *battle.Battle, ms *web_socket.MessagesStore) {

	if battle.WaitReady {
		return
	}

	moveObjArray := make([]initMoveObj, 0)
	for _, u := range units {
		moveObjArray = append(moveObjArray, u)
	}

	initMove("unit", moveObjArray, battle, units, ms)

	for _, mUnit := range units {
		moveGlobalUnit(mUnit, mp, nil)
	}

	return
}

func SetUnitsPos(units []*unit.Unit) {
	for _, gameUnit := range units {

		posFunc := gameUnit.GetPhysicalModel().GetPosFunc()
		if posFunc != nil {
			posFunc()
			gameUnit.GetPhysicalModel().SetPosFunc(nil)
		}
	}
}
