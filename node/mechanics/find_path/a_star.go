package find_path

import (
	"errors"
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/coordinate"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/physical_model"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"sort"
)

type Points struct {
	open   bool
	points [][]*Point
	count  int
	// небольшая оптимизация поиска минF
	minFIndex  int
	minFArrays []int
	minFKeys   map[int][]*Point
}

func (p *Points) checkCoordinate(x, y, xSize, ySize int) bool {

	// за пределами карты
	if x > xSize-1 || y > ySize-1 || x < 1 || y < 1 {
		return true
	}

	if p.points[x][y] == nil {
		return false
	}

	return true
}

func (p *Points) addPoint(c *Point, minFCalc bool) {

	if minFCalc {
		if p.minFKeys[c.f] == nil {
			p.minFKeys[c.f] = minFArrayHeap.Pop()
			p.addNumberInSortArray(c.f) // вставка в отсортированый масив, нет необходимости его каждый раз сортировать
		}
		p.minFKeys[c.f] = append(p.minFKeys[c.f], c)
	}

	p.points[c.x][c.y] = c // добавляем точку в масив не посещеных
	p.count++
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

func (p *Points) removePoint(c *Point) {

	index := -1
	for i, rC := range p.minFKeys[c.f] {
		if rC.x == c.x && rC.y == c.y {
			index = i
			break
		}
	}

	if index >= 0 {
		p.minFKeys[c.f][index] = p.minFKeys[c.f][len(p.minFKeys[c.f])-1]
		p.minFKeys[c.f] = p.minFKeys[c.f][:len(p.minFKeys[c.f])-1]
	}

	p.points[c.x][c.y] = nil
	p.count--
}

func (p *Points) GetMinF() *Point {

	if len(p.minFArrays) > p.minFIndex && p.minFKeys[p.minFArrays[p.minFIndex]] != nil {

		if len(p.minFKeys[p.minFArrays[p.minFIndex]]) == 0 {
			p.minFIndex++
		}

		if len(p.minFArrays) > p.minFIndex && len(p.minFKeys[p.minFArrays[p.minFIndex]]) > 0 {
			return p.minFKeys[p.minFArrays[p.minFIndex]][0]
		}
	}

	min := &Point{f: 9999}
	for _, yLine := range p.points {
		for _, point := range yLine {
			if point != nil && point.f < min.f {
				min = point
			}
		}
	}

	return min
}

func (p *Points) GetMinH() *Point {

	var min *Point

	for _, yLine := range p.points {
		for _, p := range yLine {
			if p != nil && (min == nil || p.h < min.h) {
				min = p
			}
		}
	}

	return min
}

func MoveUnit(ph *physical_model.PhysicalModel, unitID int, startX, startY, ToX, ToY float64, mp *_map.Map, units []*unit.Unit, revisionEnd bool) (path []*coordinate.Coordinate, err error) {

	start := &coordinate.Coordinate{X: int(startX), Y: int(startY)}
	end := &coordinate.Coordinate{X: int(ToX), Y: int(ToY)}

	err, path = findPath(mp, start, end, ph, unitID, units, revisionEnd)

	return path, err
}

func PrepareInData(mp *_map.Map, start, end *coordinate.Coordinate) (*Point, *Point, int, int, error) {

	xSize, ySize := mp.SetXYSize(_const.CellSize) // расчтиамем высоту и ширину карты в ху

	startPoint := &Point{
		x: start.X / _const.CellSize,
		y: start.Y / _const.CellSize,
	}

	endPoint := &Point{
		x: end.X / _const.CellSize,
		y: end.Y / _const.CellSize,
	}

	if end == nil || start == nil || endPoint.x >= xSize || endPoint.y >= ySize || endPoint.x < 0 || endPoint.y < 0 || startPoint.x >= xSize || startPoint.y >= ySize || startPoint.x < 0 || startPoint.y < 0 {
		return nil, nil, 0, 0, errors.New("start or end out the range")
	}

	return startPoint, endPoint, xSize, ySize, nil
}

func findPath(gameMap *_map.Map, start, end *coordinate.Coordinate, ph *physical_model.PhysicalModel, unitID int, units []*unit.Unit, revisionEnd bool) (error, []*coordinate.Coordinate) {

	startPoint, endPoint, xSize, ySize, err := PrepareInData(gameMap, start, end)
	if err != nil {
		return err, nil
	}

	// создаем 2 карты для посещенных (open) и непосещеных (close) точек
	openPoints, closePoints := initPoints(xSize, ySize)

	startPoint.h = GetH(startPoint, endPoint)
	openPoints.addPoint(startPoint, true) // кладем в карту посещенных точек стартовую точку

	defer func() {
		for _, x := range openPoints.points {
			for _, p := range x {
				if p != nil {
					pointHead.Push(p)
				}
			}
		}

		for _, x := range closePoints.points {
			for _, p := range x {
				if p != nil {
					pointHead.Push(p)
				}
			}
		}

		for _, x := range openPoints.minFKeys {
			minFArrayHeap.Push(x)
		}
	}()

	var path []*coordinate.Coordinate
	var noSortedPath []*Point
	newPointBasket := make([]*Point, 8, 8)

	for {
		if openPoints.count == 0 {
			if revisionEnd {
				// найти самую близжайшую точку к той куда ищем путь, назначить ее концом и найти путь заного
				point := closePoints.GetMinH()
				if point == nil {
					return errors.New("no path"), nil
				}

				endPoint = point
				for !point.Equal(startPoint) { // идем обратно до тех пока пока не дойдем до стартовой точки

					point = point.parent // по родительским точкам

					if !point.Equal(startPoint) { // если текущая точка попрежнему не стартовая то добавляем в путь координату
						noSortedPath = append(noSortedPath, point)
					}
				}
			} else {
				return errors.New("no path"), nil
			}

			break
		}

		current := MinF(openPoints, closePoints) // Берем точку с мин стоимостью пути

		if current.Equal(endPoint) { // если текущая точка и есть конец начинаем генерить путь

			for !current.Equal(startPoint) { // идем обратно до тех пока пока не дойдем до стартовой точки

				current = current.parent // по родительским точкам

				if !current.Equal(startPoint) { // если текущая точка попрежнему не стартовая то добавляем в путь координату
					noSortedPath = append(noSortedPath, current)
				}
			}
			break
		}

		parseNeighbours(current, openPoints, closePoints, gameMap, endPoint, ph, unitID, xSize, ySize, units, newPointBasket)
	}

	if len(noSortedPath) > 0 {
		for i := len(noSortedPath); i > 0; i-- {
			path = append(path, &coordinate.Coordinate{
				X: noSortedPath[i-1].x * _const.CellSize,
				Y: noSortedPath[i-1].y * _const.CellSize,
			})
		}

		end.X = endPoint.x * _const.CellSize
		end.Y = endPoint.y * _const.CellSize

		path = append(path, end)
		return nil, path
	} else {
		return errors.New("no path"), nil
	}
}

func parseNeighbours(curr *Point, open, close *Points, gameMap *_map.Map, end *Point,
	ph *physical_model.PhysicalModel, unitID int, xSize, ySize int, units []*unit.Unit, newPointBasket []*Point) {

	nCoordinate := generateNeighboursCoordinate(curr, gameMap, ph, unitID, xSize, ySize, close, open, units, newPointBasket) // берем всех соседей этой клетки

	for _, c := range nCoordinate {

		if c == nil {
			continue
		}

		c.g = curr.GetG(c) // стоимость клетки
		c.h = GetH(c, end) // приближение от точки до конечной цели.
		c.GetF()           // длина пути до цели
		c.parent = curr

		open.addPoint(c, true)
	}
}

func initPoints(xSize, ySize int) (*Points, *Points) {
	openPoints, closePoints := &Points{
		points:     make([][]*Point, xSize),
		minFKeys:   make(map[int][]*Point),
		minFArrays: make([]int, 0),
		open:       true,
	}, &Points{
		points:     make([][]*Point, xSize),
		minFKeys:   make(map[int][]*Point),
		minFArrays: make([]int, 0),
	}

	for i := 0; i < len(openPoints.points); i++ {
		openPoints.points[i] = make([]*Point, ySize)
	}

	for i := 0; i < len(closePoints.points); i++ {
		closePoints.points[i] = make([]*Point, ySize)
	}

	return openPoints, closePoints
}

func GetH(a, b *Point) int { // эвристическое приближение стоимости пути от v до конечной цели. (длинна пути)
	return int(game_math.GetBetweenDist(b.x, b.y, a.x, a.y) * 10)
}

func MinF(open, close *Points) *Point { // берет точку с минимальной стоимостью пути из масива не посещеных

	min := open.GetMinF()

	open.removePoint(min)
	close.addPoint(min, false) // добавляем в массив посещенные

	return min
}
