package find_path

import (
	"errors"
	"github.com/TrashPony/game-engine/node/mechanics/collisions"
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/coordinate"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/physical_model"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
)

func FindPath(moveUnit *unit.Unit, startX, startY, ToX, ToY float64, units []*unit.Unit, revisionEnd bool, mp *_map.Map) ([]*coordinate.Coordinate, error) {
	// создаем фейковую тушку для расчета ее на разных позициях
	physicalModel := moveUnit.GetCopyPhysicalModel()
	physicalModel.SetPolygon(nil)

	// пытаемя проложить путь от начала пути до конечной точки по прямой, если колизий нет то и искать не надо
	collision := collisions.SearchCollisionInLine(startX, startY, ToX, ToY, mp, physicalModel, moveUnit.GetID(), 5, false, units)
	if !collision {
		return []*coordinate.Coordinate{{X: int(ToX), Y: int(ToY)}}, nil
	} else {
		return startFind(physicalModel, moveUnit.GetID(), startX, startY, ToX, ToY, mp, units, revisionEnd)
	}
}

func startFind(ph *physical_model.PhysicalModel, unitID int, x, y, ToX, ToY float64, mp *_map.Map, units []*unit.Unit, revisionEnd bool) ([]*coordinate.Coordinate, error) {
	// если честно то я уже сам не помню что тут происходит, но это важно
	path := make([]*coordinate.Coordinate, 0)
	last := false

	var points []*coordinate.Coordinate

	for {

		if !last {

			// ищем путь алгоритмом А*
			if points == nil {
				// т.к. он не будет менятся нам нет смысла искать его всегда заного
				var err error
				points, err = MoveUnit(ph, unitID, x, y, ToX, ToY, mp, units, revisionEnd)
				if err != nil {
					return nil, errors.New("a start no find path")
				}
			}

			if points == nil {
				return nil, errors.New("a start no find path")
			}

			// находим максимальную отдаленную точку куда может попасть юнит
			x, y, last = SearchPoint(&points, x, y, mp, ph, unitID, units)
			if x == 0 && y == 0 {

				exit := true
				for i := 0; i < len(points); i++ {
					if points[i] != nil {
						x, y = float64(points[i].X+_const.CellSize/2), float64(points[i].Y+_const.CellSize/2)
						points[i] = nil
						exit = false
						break
					}
				}

				if exit {
					return path, nil
				}
			}

			path = append(path, &coordinate.Coordinate{X: int(x), Y: int(y)})
		} else {
			break
		}
	}

	return path, nil
}

func SearchPoint(points *[]*coordinate.Coordinate, unitX, unitY float64, mp *_map.Map, ph *physical_model.PhysicalModel, unitID int, units []*unit.Unit) (float64, float64, bool) {

	x, y := 0, 0
	lastIndex := 0

	if len(*points) > 0 {
		for i := len(*points) - 1; i >= 0; i-- {

			cell := (*points)[i]

			if cell != nil {
				collision := collisions.SearchCollisionInLine(
					float64(cell.X)+_const.CellSize/2,
					float64(cell.Y)+_const.CellSize/2,
					unitX,
					unitY, mp, ph, unitID, _const.CellSize, false, units)

				if !collision {
					if i > lastIndex {
						x, y, lastIndex = cell.X+_const.CellSize/2, cell.Y+_const.CellSize/2, i
					}
				}
			}
		}
	}

	for i := 0; i <= lastIndex; i++ {
		(*points)[i] = nil
	}

	return float64(x), float64(y), len(*points)-1 == lastIndex
}
