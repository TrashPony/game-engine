package find_path

import (
	"errors"
	"fmt"
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/mechanics/debug"
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/coordinate"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/physical_model"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
	"sort"
	"strconv"
	"time"
)

type Points struct {
	points map[string]*coordinate.Coordinate

	// небольшая оптимизация поиска минF
	minFIndex  int
	minFArrays []int
	minFKeys   map[int]map[string]bool
}

func (p *Points) checkCoordinate(x, y, xSize, ySize int) bool {

	// за пределами карты
	if x > xSize-1 || y > ySize-1 || x < 1 || y < 1 {
		return true
	}

	if p.points[strconv.Itoa(x)+":"+strconv.Itoa(y)] == nil {
		return false
	}

	return true
}

func (p *Points) addPoint(c *coordinate.Coordinate) {
	if p.minFKeys[c.F] == nil {
		p.minFKeys[c.F] = make(map[string]bool)
		//p.minFArrays = append(p.minFArrays, c.F)
		p.addNumberInSortArray(c.F) // todo вставка в отсортированый масив, нет необходимости его каждый раз сортировать
	}

	key := c.Key()

	p.minFKeys[c.F][key] = false
	p.points[key] = c // добавляем точку в масив не посещеных
}

func (p *Points) addNumberInSortArray(number int) {
	// ищем индекс куда вставить элемент c помощью бинарного поиска
	index := binarySearch(p.minFArrays, number)
	if index == -1 || len(p.minFArrays) <= index+1 {

		p.minFArrays = append(p.minFArrays, number)

		if index == -1 {
			sort.Ints(p.minFArrays)
		}
	} else {
		index++
		p.minFArrays = append(p.minFArrays[:index+1], p.minFArrays[index:]...)
		p.minFArrays[index] = number
	}
}

func binarySearch(a []int, search int) (result int) {
	mid := len(a) / 2
	switch {
	case len(a) == 0:
		result = -1 // not found
	case a[mid] > search:
		result = binarySearch(a[:mid], search)
	case a[mid] < search:
		result = binarySearch(a[mid+1:], search)
		result += mid + 1
	default: // a[mid] == search
		result = mid // found
	}

	return
}

func (p *Points) removePoint(c *coordinate.Coordinate) {

	key := c.Key()

	delete(p.minFKeys[c.F], key)
	delete(p.points, key)
}

func (p *Points) GetMinF() *coordinate.Coordinate {

	if len(p.minFArrays) > p.minFIndex && p.minFKeys[p.minFArrays[p.minFIndex]] != nil {

		if len(p.minFKeys[p.minFArrays[p.minFIndex]]) == 0 {
			p.minFIndex++
		}

		if len(p.minFArrays) > p.minFIndex {
			for key := range p.minFKeys[p.minFArrays[p.minFIndex]] {
				return p.points[key]
			}
		}
	}

	min := &coordinate.Coordinate{F: 9999}

	for _, p := range p.points {
		if p.F < min.F {
			min = p
		}
	}

	return min
}

func (p *Points) GetMinH() *coordinate.Coordinate {

	var min *coordinate.Coordinate

	for _, p := range p.points {
		if min == nil || p.H < min.H {
			min = p
		}
	}

	return min
}

func MoveUnit(ph *physical_model.PhysicalModel, unitID int, startX, startY, ToX, ToY float64, mp *_map.Map, units []*unit.Unit) (path []*coordinate.Coordinate, err error) {

	defer func() {
		if r := recover(); r != nil {
			// это тут из за очень странного поведения, когда в todo #30, возможно уже не актуально
			fmt.Println("Recovered in f", r)

			// подменяем результат функции
			path = nil
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
		}
	}()

	start := &coordinate.Coordinate{X: int(startX), Y: int(startY)}
	end := &coordinate.Coordinate{X: int(ToX), Y: int(ToY)}

	err, path = FindPath(mp, start, end, ph, unitID, units)
	if debug.Store.AStartResult {
		for _, cell := range path {
			debug.Store.AddMessage("CreateRect", "green", cell.X, cell.Y, 0, 0, _const.CellSize, mp.Id, 0)
		}
	}

	return path, err
}

func PrepareInData(mp *_map.Map, start, end *coordinate.Coordinate) (*coordinate.Coordinate, *coordinate.Coordinate, int, int, error) {

	xSize, ySize := mp.SetXYSize(_const.CellSize) // расчтиамем высоту и ширину карты в ху

	start.X, start.Y = start.X/_const.CellSize, start.Y/_const.CellSize
	start.State = 1

	end.X, end.Y = end.X/_const.CellSize, end.Y/_const.CellSize

	if end == nil || start == nil || end.X >= xSize || end.Y >= ySize || end.X < 0 || end.Y < 0 || start.X >= xSize || start.Y >= ySize || start.X < 0 || start.Y < 0 {
		return nil, nil, 0, 0, errors.New("start or end out the range")
	}

	return start, end, xSize, ySize, nil
}

func FindPath(gameMap *_map.Map, start, end *coordinate.Coordinate, ph *physical_model.PhysicalModel, unitID int, units []*unit.Unit) (error, []*coordinate.Coordinate) {

	if debug.Store.Move && unitID > 0 {
		println("-- start find aStar --")
		startTime := time.Now()
		defer func() {
			if debug.Store.Move && unitID > 0 {
				elapsed := time.Since(startTime)
				fmt.Println("time aStar path: " + strconv.FormatFloat(elapsed.Seconds(), 'f', 6, 64))
			}
		}()
	}

	start, end, xSize, ySize, err := PrepareInData(gameMap, start, end)
	if err != nil {
		return err, nil
	}

	if debug.Store.AStartNeighbours {
		debug.Store.AddMessage("CreateRect", "blue", end.X*_const.CellSize, end.Y*_const.CellSize, 0, 0, _const.CellSize, gameMap.Id, 0)
	}

	// создаем 2 карты для посещенных (open) и непосещеных (close) точек
	openPoints, closePoints := Points{
		points:     make(map[string]*coordinate.Coordinate),
		minFKeys:   make(map[int]map[string]bool),
		minFArrays: make([]int, 0),
	}, Points{
		points:     make(map[string]*coordinate.Coordinate),
		minFKeys:   make(map[int]map[string]bool),
		minFArrays: make([]int, 0),
	}

	start.H = GetH(start, end)
	openPoints.points[start.Key()] = start // кладем в карту посещенных точек стартовую точку

	var path []*coordinate.Coordinate
	var noSortedPath []*coordinate.Coordinate

	for {
		if len(openPoints.points) == 0 {

			// найти самую близжайшую точку к результату, назначить ее концом и найти путь заного
			point := closePoints.GetMinH()
			if debug.Store.AStartResult {
				debug.Store.AddMessage("CreateRect", "green", point.X*_const.CellSize, point.Y*_const.CellSize,
					0, 0, _const.CellSize, gameMap.Id, 0)
			}

			end = point
			for !point.Equal(start) { // идем обратно до тех пока пока не дойдем до стартовой точки

				point = point.Parent // по родительским точкам

				if !point.Equal(start) { // если текущая точка попрежнему не стартовая то добавляем в путь координату
					noSortedPath = append(noSortedPath, point)
				}
			}
			break
		}

		current := MinF(&openPoints, &closePoints) // Берем точку с мин стоимостью пути

		if current.Equal(end) { // если текущая точка и есть конец начинаем генерить путь

			for !current.Equal(start) { // идем обратно до тех пока пока не дойдем до стартовой точки

				current = current.Parent // по родительским точкам

				if !current.Equal(start) { // если текущая точка попрежнему не стартовая то добавляем в путь координату
					noSortedPath = append(noSortedPath, current)
				}
			}
			break
		}

		parseNeighbours(current, &openPoints, &closePoints, gameMap, end, ph, unitID, xSize, ySize, units)
	}

	if len(noSortedPath) > 0 {
		for i := len(noSortedPath); i > 0; i-- {
			noSortedPath[i-1].X *= _const.CellSize
			noSortedPath[i-1].Y *= _const.CellSize
			path = append(path, noSortedPath[i-1])
		}

		end.X *= _const.CellSize
		end.Y *= _const.CellSize

		path = append(path, end)
		return nil, path
	} else {
		return errors.New("no path"), nil
	}
}

func parseNeighbours(curr *coordinate.Coordinate, open, close *Points, gameMap *_map.Map, end *coordinate.Coordinate,
	ph *physical_model.PhysicalModel, unitID int, xSize, ySize int, units []*unit.Unit) {

	nCoordinate := generateNeighboursCoordinate(curr, gameMap, ph, unitID, xSize, ySize, close, open, units) // берем всех соседей этой клетки

	for _, c := range nCoordinate {

		if c == nil {
			continue // если ячейка является блокированой или находиться в масиве посещенных то пропускаем ее
		}

		// считаем для поинта значения пути
		c.G = curr.GetG(c) // стоимость клетки
		c.H = GetH(c, end) // приближение от точки до конечной цели.
		c.GetF()           // длина пути до цели
		c.Parent = curr

		open.addPoint(c) // добавляем точку в масив не посещеных

		if debug.Store.AStartNeighbours {
			debug.Store.AddMessage("CreateRect", "orange", c.X*_const.CellSize, c.Y*_const.CellSize, 0, 0, _const.CellSize, gameMap.Id, 0)
		}
	}
}

func GetH(a, b *coordinate.Coordinate) int { // эвристическое приближение стоимости пути от v до конечной цели. (длинна пути)
	return int(game_math.GetBetweenDist(b.X, b.Y, a.X, a.Y) * 10)
}

func MinF(open, close *Points) *coordinate.Coordinate { // берет точку с минимальной стоимостью пути из масива не посещеных

	min := open.GetMinF()

	open.removePoint(min)
	close.addPoint(min) // добавляем в массив посещенные

	return min
}
