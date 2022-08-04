package watch

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
)

func CheckViewCoordinate(watcher Watcher, x, y int, b *battle.Battle, units []*unit.Unit, radius int) (bool, bool) {
	var view, radarVisible bool

	if b != nil {

		for _, u := range units {
			if u.TeamID == watcher.GetTeamID() {
				unitView, unitRadarVisible := u.CheckViewCoordinate(x, y, radius)
				if unitView {
					return true, true
				}

				radarVisible = radarVisible || unitRadarVisible
			}
		}

		objects, objectsMX := b.Map.UnsafeRangeBuildDynamicObjects()
		defer objectsMX.RUnlock()

		for _, obj := range objects {

			if obj.TeamID != watcher.GetTeamID() {
				continue
			}

			objView, objRadarVisible := obj.CheckViewCoordinate(x, y, radius)
			if objView {
				return true, true
			}

			radarVisible = radarVisible || objRadarVisible
		}
	}

	return view, radarVisible
}
