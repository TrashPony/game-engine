package collisions

import (
	"github.com/TrashPony/game_engine/src/mechanics/debug"
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/physical_model"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/target"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
	"math"
)

func SearchCollisionInLine(startX, startY, ToX, ToY float64, mp *_map.Map, ph *physical_model.PhysicalModel, unitID int, speed float64, bigGrid bool, units []*unit.Unit) bool {
	// текущее положение курсора
	currentX, currentY := startX, startY

	// угол от старта до конца
	angle := game_math.GetBetweenAngle(ToX, ToY, startX, startY)
	radian := angle * math.Pi / 180

	// перменная для контроля зависаний, если дальность начала возрастать значит алгоритм проебал точку выхода
	minDist := game_math.GetBetweenDist(int(currentX), int(currentY), int(ToX), int(ToY))

	for {
		// находим длинную вектора до цели
		distToEnd := game_math.GetBetweenDist(int(currentX), int(currentY), int(ToX), int(ToY))
		if distToEnd < speed || minDist < distToEnd {
			return false
		}

		if debug.Store.SearchCollisionLine {
			debug.Store.AddMessage("CreateRect", "orange", int(currentX), int(currentY), 0, 0, 5, mp.Id, 20)
		}

		possibleMove := BodyCheckCollisionsOnStaticMap(ph, int(currentX), int(currentY), mp, bigGrid)
		if !possibleMove {
			if debug.Store.SearchCollisionLine {
				debug.Store.AddMessage("CreateRect", "red", int(currentX), int(currentY), 0, 0, 5, mp.Id, 20)
			}
			return true
		}

		userRect := GetBodyRect(ph, currentX, currentY, 0, true, false)
		collUser := CheckCollisionsPlayersByRect(userRect, mp.Id, unitID, units)
		if collUser != nil {
			return true
		}

		stopX, stopY := speed*game_math.Cos(radian), speed*game_math.Sin(radian)
		currentX += stopX
		currentY += stopY

		minDist = distToEnd
	}
}

func SearchCircleCollisionInLine(startX, startY, ToX, ToY float64, mp *_map.Map, radius int,
	target *target.Target, excludeObjID, excludeUnitID int) bool {
	// текущее положение курсора
	currentX, currentY := startX, startY

	// угол от старта до конца
	angle := game_math.GetBetweenAngle(ToX, ToY, startX, startY)
	radian := angle * math.Pi / 180

	// перменная для контроля зависаний, если дальность начала возрастать значит алгоритм проебал точку выхода
	minDist := game_math.GetBetweenDist(int(currentX), int(currentY), int(ToX), int(ToY))

	speed := 1.0

	for {
		// находим длинную вектора до цели
		distToEnd := game_math.GetBetweenDist(int(currentX), int(currentY), int(ToX), int(ToY))
		if distToEnd < speed || minDist < distToEnd {
			return false
		}

		_, _, mpLvl := mp.GetPosLevel(int(currentX), int(currentY))
		mpLvl += 0.1

		collision, typeCollision, id := CircleAllCollisionCheck(int(currentX), int(currentY), radius, mpLvl, mp, excludeObjID, excludeUnitID, true, nil)
		if collision {
			// если колизия наша цель то считай что до цели колизий нет
			if target != nil && target.Type == typeCollision && target.ID == id {
				return false
			} else {
				return true
			}
		}

		stopX, stopY := speed*game_math.Cos(radian), speed*game_math.Sin(radian)
		currentX += stopX
		currentY += stopY

		minDist = distToEnd
	}
}
