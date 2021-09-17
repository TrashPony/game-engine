package ai_methods

import (
	"github.com/TrashPony/game_engine/src/mechanics/factories/maps"
	"github.com/TrashPony/game_engine/src/mechanics/game_math/collisions"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/coordinate"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/target"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
	"math/rand"
)

func Scouting(gameUnit *unit.Unit) bool {

	var currentTarget *target.Target
	if gameUnit.MovePath != nil {
		currentTarget = gameUnit.MovePath.FollowTarget
	}

	if currentTarget == nil || currentTarget.Type != "map" {
		mp, _ := maps.Maps.GetByID(gameUnit.MapID)

		border := 100
		x := rand.Intn((mp.XSize-border)-border) + border
		y := rand.Intn((mp.YSize-border)-border) + border

		collision, _, _ := collisions.CircleAllCollisionCheck(x, y, gameUnit.GetPhysicalModel().Radius, 0, mp, 0, gameUnit.GetID(), false, nil)
		if collision {
			return false
		}

		gameUnit.MovePath = &unit.MovePath{
			NeedFindPath: true,
			Path:         &[]*coordinate.Coordinate{{X: x, Y: y}},
			FollowTarget: &target.Target{Type: "map", X: x, Y: y, Radius: 50},
		}
	}

	return false
}
