package find_path

import (
	"errors"
	"fmt"
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/mechanics/debug"
	"github.com/TrashPony/game_engine/src/mechanics/factories/maps"
	"github.com/TrashPony/game_engine/src/mechanics/game_math/collisions"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/coordinate"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/physical_model"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
	"strconv"
	"time"
)

func LeftHandAlgorithm(moveUnit *unit.Unit, startX, startY, ToX, ToY float64, units []*unit.Unit) ([]*coordinate.Coordinate, error) {
	mp, _ := maps.Maps.GetByID(moveUnit.MapID)

	// создаем фейковую тушку для расчета ее на разных позициях
	physicalModel := moveUnit.GetCopyPhysicalModel()
	physicalModel.Polygon = nil

	// 0 пытаемя проложить путь от начала пути до конечной точки по прямой
	collision := collisions.SearchCollisionInLine(startX, startY, ToX, ToY, mp, physicalModel, moveUnit.GetID(), 5, false, units)
	if !collision {
		return []*coordinate.Coordinate{{X: int(ToX), Y: int(ToY)}}, nil
	} else {
		return startFind(physicalModel, moveUnit.GetID(), startX, startY, ToX, ToY, mp, units)
	}
}

func startFind(ph *physical_model.PhysicalModel, unitID int, x, y, ToX, ToY float64, mp *_map.Map, units []*unit.Unit) ([]*coordinate.Coordinate, error) {

	path := make([]*coordinate.Coordinate, 0)
	last := false

	var points []*coordinate.Coordinate

	for {

		if !last {

			// ищем путь алгоритмом А*
			if points == nil {
				// т.к. он не будет менятся нам нет смысла искать его всегда заного
				var err error
				points, err = MoveUnit(ph, unitID, x, y, ToX, ToY, mp, units)
				if err != nil {
					return nil, errors.New("a start no find path")
				}
			}

			if points == nil {
				return nil, errors.New("a start no find path")
			}

			// находим максимальную отдаленную точку куда может попать юнит
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
					// TODO это не правильно т.к. если нашел путь А* то путь есть 100%
					return path, nil
				}
			}

			if debug.Store.HandAlgorithm && len(path) > 0 {
				debug.Store.AddMessage("CreateLine", "orange",
					path[len(path)-1].X*_const.CellSize+_const.CellSize/2,
					path[len(path)-1].Y*_const.CellSize+_const.CellSize/2,
					int(x),
					int(y),
					_const.CellSize, mp.Id, 20)
			}

			path = append(path, &coordinate.Coordinate{X: int(x), Y: int(y)})
		} else {
			//  2.1.1 если между координатой истиного пути и целью нет препятсвий формируем путь. Выходим из функции.
			// path = append(path, &coordinate.Coordinate{X: int(moveUnit.ToX), Y: int(moveUnit.ToY)})
			break
		}
	}

	return path, nil
}

func SearchPoint(points *[]*coordinate.Coordinate, unitX, unitY float64, mp *_map.Map, ph *physical_model.PhysicalModel, unitID int, units []*unit.Unit) (float64, float64, bool) {

	if debug.Store.Move && unitID > 0 {
		println("-- start search line --")
		startTime := time.Now()
		defer func() {
			elapsed := time.Since(startTime)
			fmt.Println("time search line: " + strconv.FormatFloat(elapsed.Seconds(), 'f', 6, 64))
		}()
	}

	//// ищем самую дальнюю точку до которой можем дойти
	//for i := len(*points) - 1; i >= 0; i-- {
	//
	//	collision := collisions.SearchCollisionInLine(
	//		float64((*points)[i].X)+game_math.CellSize/2,
	//		float64((*points)[i].Y)+game_math.CellSize/2,
	//		float64(unitX),
	//		float64(unitY), mp, gameUnit, game_math.CellSize, false, false, units)
	//	if !collision {
	//
	//		if debug.Store.SearchCollisionLineResult {
	//			debug.Store.AddMessage("CreateLine", "blue", (*points)[i].X, (*points)[i].Y, unitX, unitY, 0, mp.Id, 20)
	//		}
	//
	//		return (*points)[i].X + game_math.CellSize/2, (*points)[i].Y + game_math.CellSize/2, i == len(*points) - 1
	//	}
	//}

	// todo самый дорогой метод на дальних дистациях из за того что он считается много раз, возможно можно просчитать его в 1 фор
	// 	но у меня чет не вышло
	// ищем самую дальнюю точку до которой можем дойти

	// todo вариант оптимизации найти первую точку (искать от юнита, первая точка после которой идет не проходимая)
	// todo отдавать ее сразу и юнит начинает движение, искать остальные точки в фоне и формировать путь в переменную юнита
	// todo от туда он будет уже формировать путь для движения из горутины в пакете сокетов

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

					if debug.Store.SearchCollisionLineStep {
						debug.Store.AddMessage("CreateLine", "white", cell.X+_const.CellSize/2, cell.Y+_const.CellSize/2,
							int(unitX), int(unitY), 0, mp.Id, 20)
					}

					if i > lastIndex {
						x, y, lastIndex = cell.X+_const.CellSize/2, cell.Y+_const.CellSize/2, i
					}

				} else {
					if debug.Store.SearchCollisionLineStep {
						debug.Store.AddMessage("CreateLine", "red", cell.X+_const.CellSize/2, cell.Y+_const.CellSize/2,
							int(unitX), int(unitY), 0, mp.Id, 20)
					}
				}
			}
		}
	}

	for i := 0; i <= lastIndex; i++ {
		(*points)[i] = nil
	}

	if debug.Store.SearchCollisionLineResult {
		debug.Store.AddMessage("CreateLine", "blue", x, y, int(unitX), int(unitY), 0, mp.Id, 20)
	}

	return float64(x), float64(y), len(*points)-1 == lastIndex
}
