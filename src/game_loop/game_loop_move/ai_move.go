package game_loop_move

import (
	"github.com/TrashPony/game_engine/src/mechanics/find_path"
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
	"sync"
)

func moveGlobalUnit(gameUnit *unit.Unit, units []*unit.Unit, wg *sync.WaitGroup) {
	defer wg.Done()

	if gameUnit.MovePath == nil {
		return
	}

	if gameUnit.MovePath.FollowTarget == nil {
		gameUnit.MovePath = nil
		return
	}

	distToTarget := int(game_math.GetBetweenDist(gameUnit.GetPhysicalModel().GetX(), gameUnit.GetPhysicalModel().GetY(), gameUnit.MovePath.FollowTarget.GetX(), gameUnit.MovePath.FollowTarget.GetY()))
	if gameUnit.MovePath.FollowTarget.Radius > distToTarget {
		gameUnit.MovePath = nil
		return
	}

	if gameUnit.MovePath.NeedFindPath {
		// PATH FINDER
		path, err := find_path.LeftHandAlgorithm(
			gameUnit,
			float64(gameUnit.GetPhysicalModel().GetX()), float64(gameUnit.GetPhysicalModel().GetY()),
			float64(gameUnit.MovePath.FollowTarget.GetX()), float64(gameUnit.MovePath.FollowTarget.GetY()),
			units)

		if err != nil {
			// пытаемся крутиться
			moveAI(gameUnit.GetPhysicalModel(), gameUnit.MovePath.FollowTarget.GetX(), gameUnit.MovePath.FollowTarget.GetY(), 50, 0, 2, 20)
			return
		}

		// что бы юнита не вьежал в стройку
		gameUnit.MovePath.FollowTarget.X = path[len(path)-1].X
		gameUnit.MovePath.FollowTarget.Y = path[len(path)-1].Y

		gameUnit.MovePath.NeedFindPath = false
		gameUnit.MovePath.Path = &path
		gameUnit.MovePath.CurrentPoint = 0
	}

	if len(*gameUnit.MovePath.Path) == gameUnit.MovePath.CurrentPoint {
		if gameUnit.MovePath.FollowTarget.Radius > distToTarget {
			gameUnit.MovePath = nil
		} else {
			gameUnit.MovePath.NeedFindPath = true // todo возмодная утечка ресурсов из за того что путь не находим, а так это костыль
		}
		return
	}

	point := (*gameUnit.MovePath.Path)[gameUnit.MovePath.CurrentPoint]
	if gameUnit.MovePath.FollowTarget.Radius > distToTarget {
		if moveAI(gameUnit.GetPhysicalModel(), point.X, point.Y, float64(gameUnit.MovePath.FollowTarget.Radius)+100,
			float64(gameUnit.MovePath.FollowTarget.Radius), 1, gameUnit.MovePath.FollowTarget.Radius) {

			gameUnit.MovePath.CurrentPoint++
		}
	} else {
		if moveAI(gameUnit.GetPhysicalModel(), point.X, point.Y, 50, 0, 2, 20) {

			gameUnit.MovePath.CurrentPoint++
		}
	}
}

func moveAI(obj moveObject, toX, toY int, slowDist, stopDist, kPower float64, radius int) bool {

	dist := game_math.GetBetweenDist(toX, toY, obj.GetX(), obj.GetY())
	if int(dist) < radius {
		obj.SetWASD(false, false, false, false)
		return true
	}

	needAngle := game_math.GetBetweenAngle(float64(toX), float64(toY), float64(obj.GetX()), float64(obj.GetY()))
	diffAngle := game_math.ShortestBetweenAngle(obj.GetRotate(), needAngle)

	if diffAngle <= 2 && diffAngle >= -2 {

		obj.SetAngularVelocity(0)

		// todo ваще тут нуждны супер вычесления но мне как то лень
		if dist < stopDist {
			if dist < slowDist && obj.GetPowerMove()-obj.GetReverse() > 0 {
				obj.SetWASD(false, false, true, false)
			} else {
				obj.SetWASD(false, false, false, false)
			}
		} else {
			if dist < slowDist && obj.GetPowerMove()-obj.GetReverse() > obj.GetPowerFactor()*kPower {
				// притормаживаем перед финишом
				obj.SetWASD(false, false, true, false)
			} else {
				// идем вперед
				obj.SetWASD(true, false, false, false)
			}
		}
	} else {
		if diffAngle <= 10 && diffAngle >= -10 {
			if diffAngle > 0 {
				// идем вперед и поворачиваем направо
				obj.SetWASD(true, false, false, true)
			} else {
				// идем вперед и поворачиваем налево
				obj.SetWASD(true, true, false, false)
			}
		} else {
			if diffAngle > 0 {
				// стоим на месте и крутимся в право
				obj.SetWASD(false, false, false, true)
			} else {
				// стоим на месте и крутимся в лево
				obj.SetWASD(false, true, false, false)
			}
		}
	}

	return false
}
