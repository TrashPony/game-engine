package ai_methods

import (
	"github.com/TrashPony/game-engine/node/mechanics/collisions"
	battle2 "github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/target"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"math/rand"
)

func Scouting(b *battle2.Battle, gameUnit *unit.Unit) bool {

	followTarget := gameUnit.GetFollowTarget()

	if followTarget == nil || followTarget.Type != "map" {
		border := 100
		x := rand.Intn((b.Map.XSize-border)-border) + border
		y := rand.Intn((b.Map.YSize-border)-border) + border

		collision, _, _, _ := collisions.CircleAllCollisionCheck(x, y, gameUnit.GetPhysicalModel().Radius, 0, b.Map, nil, []int{gameUnit.GetID()}, false, nil)
		if collision {
			return false
		}

		gameUnit.SetMovePathTarget(&target.Target{Type: "map", X: x, Y: y, Radius: 50})
		return true
	}

	return false
}
