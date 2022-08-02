package find_path

import (
	collisions2 "github.com/TrashPony/game-engine/node/mechanics/collisions"
	_const "github.com/TrashPony/game-engine/router/const"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/physical_model"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
)

func checkValidForMoveCoordinate(gameMap *_map.Map, x, y, xSize, ySize int, ph *physical_model.PhysicalModel, unitID int, units []*unit.Unit) *Point {

	// за пределами карты
	if x > xSize-1 || y > ySize-1 || x < 1 || y < 1 {
		return nil
	}

	_, collUnit, _ := collisions2.CircleUnits(x*_const.CellSize+_const.CellSize/2, y*_const.CellSize+_const.CellSize/2, ph.GetRadius(), []int{unitID}, gameMap.Id, 0, 0, units)
	if collUnit != nil {
		return nil
	}

	possible := collisions2.BodyCheckCollisionsOnStaticMap(ph, x, y, gameMap, true)
	if possible {
		p := pointHead.Pop()
		p.x = x
		p.y = y
		return p
	} else {
		return nil
	}
}
