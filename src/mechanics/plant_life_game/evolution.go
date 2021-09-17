package plant_life_game

import (
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/coordinate"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/dynamic_map_object"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"strconv"
)

func Evolution(mp *_map.Map, mapObjects map[int]*dynamic_map_object.Object) {

	checkCoordinate := func(x, y int) bool {

		idString := strconv.Itoa(x) + strconv.Itoa(y)
		id, _ := strconv.Atoi(idString)

		obj := mapObjects[id]
		if obj != nil && obj.Texture == "plant_4" {
			return true
		} else {
			return false
		}
	}

	checkCountNeighbors := func(x, y int, life bool) *lifeObject {
		count := 0
		neighbors := make([]*coordinate.Coordinate, 0)

		// если вышли за пределы карты то появляемся в противоположеной стороны
		// todo но чет это не работает
		if mp.XSize < x {
			x = x - mp.XSize
		}
		if mp.YSize < y {
			y = y - mp.YSize
		}

		//строго лево
		if checkCoordinate(x-_const.CellSize, y) {
			count++
		} else {
			neighbors = append(neighbors, &coordinate.Coordinate{X: x - _const.CellSize, Y: y})
		}

		//строго право
		if checkCoordinate(x+_const.CellSize, y) {
			count++
		} else {
			neighbors = append(neighbors, &coordinate.Coordinate{X: x + _const.CellSize, Y: y})
		}

		//верх центр
		if checkCoordinate(x, y-_const.CellSize) {
			count++
		} else {
			neighbors = append(neighbors, &coordinate.Coordinate{X: x, Y: y - _const.CellSize})
		}

		//низ центр
		if checkCoordinate(x, y+_const.CellSize) {
			count++
		} else {
			neighbors = append(neighbors, &coordinate.Coordinate{X: x, Y: y + _const.CellSize})
		}

		//верх лево
		if checkCoordinate(x-_const.CellSize, y-_const.CellSize) {
			count++
		} else {
			neighbors = append(neighbors, &coordinate.Coordinate{X: x - _const.CellSize, Y: y - _const.CellSize})
		}

		//верх право
		if checkCoordinate(x+_const.CellSize, y-_const.CellSize) {
			count++
		} else {
			neighbors = append(neighbors, &coordinate.Coordinate{X: x + _const.CellSize, Y: y - _const.CellSize})
		}

		//низ лево
		if checkCoordinate(x-_const.CellSize, y+_const.CellSize) {
			count++
		} else {
			neighbors = append(neighbors, &coordinate.Coordinate{X: x - _const.CellSize, Y: y + _const.CellSize})
		}

		//низ право
		if checkCoordinate(x+_const.CellSize, y+_const.CellSize) {
			count++
		} else {
			neighbors = append(neighbors, &coordinate.Coordinate{X: x + _const.CellSize, Y: y + _const.CellSize})
		}

		return &lifeObject{X: x, Y: y, life: life, Neighbors: neighbors, countLifeNeighbors: count}
	}

	lifeObjects := make([]*lifeObject, 0)

	// смотрим все живые клетки на наличие соседей, и будущего состояния жизни
	for _, obj := range mapObjects {
		if obj.Texture == "plant_4" {

			lifeObj := checkCountNeighbors(obj.GetPhysicalModel().GetX(), obj.GetPhysicalModel().GetY(), true)

			if lifeObj.countLifeNeighbors == 3 || lifeObj.countLifeNeighbors == 2 {
				lifeObj.futureLife = true
			}

			lifeObjects = append(lifeObjects, lifeObj)
		}
	}

	// смотрим все мертвые клетки будущее состояния жизни
	for _, lifeUnit := range lifeObjects {
		for _, deadUnit := range lifeUnit.Neighbors {

			deadObj := checkCountNeighbors(deadUnit.X, deadUnit.Y, false)

			if deadObj.countLifeNeighbors == 3 {
				deadObj.futureLife = true
				lifeObjects = append(lifeObjects, deadObj)
			}
		}
	}

	// мы применяем действия к клеткам
	for _, terrainUnit := range lifeObjects {
		if !terrainUnit.life && terrainUnit.futureLife {
			createPlant(terrainUnit.X, terrainUnit.Y, "plant_4", mp)
		}

		if terrainUnit.life && !terrainUnit.futureLife {
			mp.RemoveDynamicObjectByXY(terrainUnit.X, terrainUnit.Y)
		}
	}
}
