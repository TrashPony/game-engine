package find_path

import (
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/mechanics/game_math/collisions"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/coordinate"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/physical_model"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
)

func checkValidForMoveCoordinate(gameMap *_map.Map, x, y, xSize, ySize int, ph *physical_model.PhysicalModel, unitID int, units []*unit.Unit) *coordinate.Coordinate {

	// за пределами карты
	if x > xSize-1 || y > ySize-1 || x < 1 || y < 1 {
		return nil
	}

	userRect := collisions.GetBodyRect(ph, float64(x*_const.CellSize+_const.CellSize/2), float64(y*_const.CellSize+_const.CellSize/2), 0, true, false)
	collUnit := collisions.CheckCollisionsPlayersByRect(userRect, gameMap.Id, unitID, units)
	if collUnit != nil {
		return nil
	}

	possible := collisions.BodyCheckCollisionsOnStaticMap(ph, x, y, gameMap, true)
	if possible {
		return &coordinate.Coordinate{X: x, Y: y}
	}

	return nil
}
