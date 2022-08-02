package collisions

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/physical_model"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"math"
)

func SearchCollisionInLine(startX, startY, ToX, ToY float64, mp *_map.Map, ph *physical_model.PhysicalModel, unitID int, speed float64, bigGrid bool, units []*unit.Unit) bool {
	// текущее положение курсора
	currentX, currentY := startX, startY

	// угол от старта до конца
	angle := game_math.GetBetweenAngle(ToX, ToY, startX, startY)
	radian := angle * math.Pi / 180

	// перменная для контроля зависаний, если дальность начала возрастает значит алгоритм проебал точку выхода
	minDist := game_math.GetBetweenDist(int(currentX), int(currentY), int(ToX), int(ToY))

	for {
		// находим длинную вектора до цели
		distToEnd := game_math.GetBetweenDist(int(currentX), int(currentY), int(ToX), int(ToY))
		if distToEnd < speed || minDist < distToEnd {
			return false
		}

		possibleMove := BodyCheckCollisionsOnStaticMap(ph, int(currentX), int(currentY), mp, bigGrid)
		if !possibleMove {
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
