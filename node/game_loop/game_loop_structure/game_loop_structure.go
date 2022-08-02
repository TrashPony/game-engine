package game_loop_structure

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/dynamic_map_object"
)

func Structure(b *battle.Battle, buildObjects []*dynamic_map_object.Object) {
	for _, obj := range buildObjects {
		if obj.Build && (obj.GetHP() > 0 || obj.Immortal) {
			StructureRouter(obj, b)
		}
	}

	return
}

func StructureRouter(obj *dynamic_map_object.Object, b *battle.Battle) {

	if !obj.Run {
		return
	}

	if obj.Type == "turret" {
		TurretTarget(obj, b)
	}
}
