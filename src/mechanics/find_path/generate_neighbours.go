package find_path

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/coordinate"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/physical_model"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
)

func generateNeighboursCoordinate(curr *coordinate.Coordinate, gameMap *_map.Map, ph *physical_model.PhysicalModel, unitID int,
	xSize, ySize int, close *Points, open *Points, units []*unit.Unit) []*coordinate.Coordinate {

	// берет все соседние клетки от текущей
	res := make([]*coordinate.Coordinate, 4, 4)

	//строго лево
	if !close.checkCoordinate(curr.X-1, curr.Y, xSize, ySize) && !open.checkCoordinate(curr.X-1, curr.Y, xSize, ySize) {
		res[0] = checkValidForMoveCoordinate(gameMap, curr.X-1, curr.Y, xSize, ySize, ph, unitID, units)
	}

	//строго право
	if !close.checkCoordinate(curr.X+1, curr.Y, xSize, ySize) && !open.checkCoordinate(curr.X+1, curr.Y, xSize, ySize) {
		res[1] = checkValidForMoveCoordinate(gameMap, curr.X+1, curr.Y, xSize, ySize, ph, unitID, units)
	}

	//верх центр
	if !close.checkCoordinate(curr.X, curr.Y-1, xSize, ySize) && !open.checkCoordinate(curr.X, curr.Y-1, xSize, ySize) {
		res[2] = checkValidForMoveCoordinate(gameMap, curr.X, curr.Y-1, xSize, ySize, ph, unitID, units)
	}

	//низ центр
	if !close.checkCoordinate(curr.X, curr.Y+1, xSize, ySize) && !open.checkCoordinate(curr.X, curr.Y+1, xSize, ySize) {
		res[3] = checkValidForMoveCoordinate(gameMap, curr.X, curr.Y+1, xSize, ySize, ph, unitID, units)
	}

	return res
}
