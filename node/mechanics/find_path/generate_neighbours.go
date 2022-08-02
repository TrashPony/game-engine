package find_path

import (
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/physical_model"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
)

func generateNeighboursCoordinate(curr *Point, gameMap *_map.Map, ph *physical_model.PhysicalModel, unitID int,
	xSize, ySize int, close *Points, open *Points, units []*unit.Unit, res []*Point) []*Point {

	// берет все соседние клетки от текущей
	for i := range res {
		res[i] = nil
	}

	//строго лево
	if !close.checkCoordinate(curr.x-1, curr.y, xSize, ySize) && !open.checkCoordinate(curr.x-1, curr.y, xSize, ySize) {
		res[0] = checkValidForMoveCoordinate(gameMap, curr.x-1, curr.y, xSize, ySize, ph, unitID, units)
	}

	//строго право
	if !close.checkCoordinate(curr.x+1, curr.y, xSize, ySize) && !open.checkCoordinate(curr.x+1, curr.y, xSize, ySize) {
		res[1] = checkValidForMoveCoordinate(gameMap, curr.x+1, curr.y, xSize, ySize, ph, unitID, units)
	}

	//верх центр
	if !close.checkCoordinate(curr.x, curr.y-1, xSize, ySize) && !open.checkCoordinate(curr.x, curr.y-1, xSize, ySize) {
		res[2] = checkValidForMoveCoordinate(gameMap, curr.x, curr.y-1, xSize, ySize, ph, unitID, units)
	}

	//низ центр
	if !close.checkCoordinate(curr.x, curr.y+1, xSize, ySize) && !open.checkCoordinate(curr.x, curr.y+1, xSize, ySize) {
		res[3] = checkValidForMoveCoordinate(gameMap, curr.x, curr.y+1, xSize, ySize, ph, unitID, units)
	}

	//сверху лево
	if !close.checkCoordinate(curr.x-1, curr.y-1, xSize, ySize) && !open.checkCoordinate(curr.x-1, curr.y-1, xSize, ySize) {
		res[4] = checkValidForMoveCoordinate(gameMap, curr.x-1, curr.y-1, xSize, ySize, ph, unitID, units)
	}

	//сверху право
	if !close.checkCoordinate(curr.x+1, curr.y-1, xSize, ySize) && !open.checkCoordinate(curr.x+1, curr.y-1, xSize, ySize) {
		res[5] = checkValidForMoveCoordinate(gameMap, curr.x+1, curr.y-1, xSize, ySize, ph, unitID, units)
	}

	//низ лево
	if !close.checkCoordinate(curr.x-1, curr.y+1, xSize, ySize) && !open.checkCoordinate(curr.x-1, curr.y+1, xSize, ySize) {
		res[6] = checkValidForMoveCoordinate(gameMap, curr.x-1, curr.y+1, xSize, ySize, ph, unitID, units)
	}

	//низ право
	if !close.checkCoordinate(curr.x+1, curr.y+1, xSize, ySize) && !open.checkCoordinate(curr.x+1, curr.y+1, xSize, ySize) {
		res[7] = checkValidForMoveCoordinate(gameMap, curr.x+1, curr.y+1, xSize, ySize, ph, unitID, units)
	}

	return res
}
