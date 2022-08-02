package check_target

import (
	"github.com/TrashPony/game-engine/node/mechanics/factories/units"
	"github.com/TrashPony/game-engine/node/mechanics/watch"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/behavior_rule"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/target"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/visible_objects"
	"sync"
)

type Hostile struct {
	Type   string `json:"type"`
	ID     int    `json:"id"`
	Points int    `json:"points"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

func CheckTarget(watcher watch.Watcher, currTarget *target.Target, mp *_map.Map, teamID int) (bool, int, int) {
	if currTarget.Type == "unit" {
		return CheckUnit(watcher, currTarget.ID, mp, teamID)
	}

	if currTarget.Type == "object" {
		return CheckObj(watcher, currTarget.ID, mp, teamID)
	}

	return false, 0, 0
}

func FindHostile(watcher watch.Watcher, watcherX, watcherY int, mp *_map.Map, teamID, maxDist int, meta *behavior_rule.Meta) []*Hostile {

	hostiles := make([]*Hostile, 0)
	if watcher == nil {
		return hostiles
	}

	findTarget := func(objs []*visible_objects.VisibleObject, mx *sync.RWMutex) {

		if mx == nil {
			return
		}

		mx.RLock()
		defer mx.RUnlock()

		for _, vObj := range objs {
			mx.RUnlock()

			if vObj.TypeObject == "unit" {
				hostile, x, y := CheckUnit(watcher, vObj.IDObject, mp, teamID)

				if meta != nil && meta.Type == "zone" {
					dist := int(game_math.GetBetweenDist(x, y, meta.X, meta.Y))
					if dist > meta.Radius {
						mx.RLock()
						continue
					}
				}

				dist := int(game_math.GetBetweenDist(x, y, watcherX, watcherY))
				if hostile && (dist < maxDist || maxDist == 0) {
					hostiles = append(hostiles, &Hostile{Type: "unit", ID: vObj.IDObject, X: x, Y: y})
				}
			}

			if vObj.TypeObject == "object" {
				hostile, x, y := CheckObj(watcher, vObj.IDObject, mp, teamID)

				if meta != nil && meta.Type == "zone" {
					dist := int(game_math.GetBetweenDist(x, y, meta.X, meta.Y))
					if dist > meta.Radius {
						mx.RLock()
						continue
					}
				}

				dist := int(game_math.GetBetweenDist(x, y, watcherX, watcherY))
				if hostile && (dist < maxDist || maxDist == 0) {
					hostiles = append(hostiles, &Hostile{Type: "object", ID: vObj.IDObject, X: x, Y: y})
				}
			}

			mx.RLock()
		}
	}

	findTarget(watcher.UnsafeRangeVisibleObjects())
	return hostiles
}

func CheckUnit(watcher watch.Watcher, id int, mp *_map.Map, teamID int) (bool, int, int) {
	unit := units.Units.GetUnitByIDAndMapID(id, mp.Id)
	if unit != nil {

		if watcher != nil && watcher.CheckUnion(unit.TeamID) {
			return false, 0, 0
		}

		if unit.TeamID > 0 && unit.TeamID != teamID && !unit.GetPhysicalModel().Fly {
			return true, unit.GetX(), unit.GetY()
		} else {
			return false, 0, 0
		}
	}

	return false, 0, 0
}

func CheckObj(watcher watch.Watcher, id int, mp *_map.Map, teamID int) (bool, int, int) {

	obj := mp.GetDynamicObjectsByID(id)
	if obj != nil {

		if watcher != nil && watcher.CheckUnion(obj.TeamID) {
			return false, 0, 0
		}

		if obj.Immortal || obj.Fraction == "ALL" {
			return false, 0, 0
		}

		if obj.TeamID > 0 && obj.TeamID != teamID {
			return true, obj.GetX(), obj.GetY()
		} else {
			return false, 0, 0
		}
	}

	return false, 0, 0
}
