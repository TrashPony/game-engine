package game_loop_move

import (
	"github.com/TrashPony/game-engine/node/mechanics/actual_target"
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/target"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"sync"
)

func moveGlobalUnit(gameUnit *unit.Unit, mp *_map.Map, wg *sync.WaitGroup) {

	if wg != nil {
		defer wg.Done()
	}

	followTarget, path, currentPoint, _ := gameUnit.GetMovePathState()
	if followTarget == nil && gameUnit.GetID() < 0 {
		gameUnit.GetPhysicalModel().SetWASD(false, false, false, false, false, false, false)
		gameUnit.RemoveMovePath()
		return
	}

	if !actual_target.GetXYZTarget(gameUnit.GetGunner(), followTarget, mp, nil) {
		gameUnit.RemoveMovePath()
		return
	}

	distToTarget := int(game_math.GetBetweenDist(gameUnit.GetPhysicalModel().GetX(), gameUnit.GetPhysicalModel().GetY(), followTarget.GetX(), followTarget.GetY()))
	if followTarget.Radius > distToTarget {
		gameUnit.RemoveMovePath()
		return
	}

	if len(*path) == currentPoint {
		if followTarget.Radius > distToTarget {
			gameUnit.RemoveMovePath()
		} else {
			gameUnit.SetFindMovePath() // todo возмодная утечка ресурсов из за того что путь не находим, а так это костыль
		}
		return
	}

	if len(*path) > currentPoint {
		point := (*path)[currentPoint]

		if len(gameUnit.Body.Weapons) == 0 {
			gameUnit.SetWeaponTarget(&target.Target{Type: "map", X: point.X, Y: point.Y})
		}

		wt := gameUnit.GetWeaponTarget()
		if wt == nil {
			gameUnit.SetWeaponTarget(&target.Target{X: point.X, Y: point.Y, Attack: false})
		}

		if followTarget.Radius > distToTarget {
			if moveAI(gameUnit.GetPhysicalModel(), point.X, point.Y, float64(followTarget.Radius)+100,
				float64(followTarget.Radius), 1, followTarget.Radius) {

				gameUnit.NextMovePoint()
			}
		} else {
			if moveAI(gameUnit.GetPhysicalModel(), point.X, point.Y, 50, 0, 10, _const.CellSize*2) {

				gameUnit.NextMovePoint()
			}
		}
	}
}

func moveAI(obj moveObject, toX, toY int, slowDist, stopDist, kPower float64, radius int) bool {

	dist := game_math.GetBetweenDist(toX, toY, obj.GetX(), obj.GetY())
	needAngle := game_math.GetBetweenAngle(float64(toX), float64(toY), float64(obj.GetX()), float64(obj.GetY()))
	diffAngle := game_math.ShortestBetweenAngle(obj.GetRotate(), needAngle)

	if int(dist) < radius && (diffAngle <= 2 && diffAngle >= -2 || obj.GetChassisType() == "caterpillars") {
		obj.SetWASD(false, false, false, false, false, false, false)
		return true
	}

	if diffAngle <= 2 && diffAngle >= -2 {

		obj.SetAngularVelocity(0)

		if dist < stopDist {
			if obj.GetPowerMove()-obj.GetReverse() > obj.GetPowerFactor()*kPower {
				// притормаживаем перед финишом
				obj.SetWASD(false, false, true, false, false, false, false)
			} else {
				if obj.GetPowerMove()-obj.GetReverse() > 0 {
					obj.SetWASD(false, false, false, false, false, false, false)
				} else {
					obj.SetWASD(true, false, false, false, false, false, false)
				}
			}

		} else {
			obj.SetWASD(true, false, false, false, false, false, false)
		}

		return false
	}

	if obj.GetChassisType() == "caterpillars" || obj.GetChassisType() == "fly" {
		if diffAngle <= 10 && diffAngle >= -10 {
			obj.SetWASD(true, diffAngle < 0, false, diffAngle > 0, false, false, false)
		} else {
			obj.SetWASD(false, diffAngle < 0, false, diffAngle > 0, false, false, false)
		}
	}

	return false
}
