package units

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"sync"
)

var Units = newStore()

type store struct {
	units map[int]*unitsInMap // карта юнитов в игре
	mx    sync.RWMutex
}

func newStore() *store {
	return &store{
		units: make(map[int]*unitsInMap),
	}
}

func (c *store) getUnitsInMap(mapID int) *unitsInMap {
	c.mx.RLock()

	if c.units[mapID] == nil {
		c.mx.RUnlock()
		c.UnitsInInit(mapID)
		c.mx.RLock()
	}

	defer c.mx.RUnlock()

	return c.units[mapID]
}

func (c *store) UnitsInInit(mapID int) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.units[mapID] = &unitsInMap{units: make([]*unit.Unit, 0)}
}

func (c *store) GetAllUnitsArray(mapID int, basket []*unit.Unit) []*unit.Unit {
	units := c.getUnitsInMap(mapID)

	units.mx.RLock()
	basket = basket[:0]

	for _, gameUnit := range units.units {
		basket = append(basket, gameUnit)
	}
	units.mx.RUnlock()

	return basket
}

func (c *store) AddUnit(unit *unit.Unit) {
	unitsInMap := c.getUnitsInMap(unit.MapID)
	unitsInMap.addUnit(unit)
}

func (c *store) GetAllUnitsByMapIDUnsafeRange(mapID int) ([]*unit.Unit, *sync.RWMutex) {
	units := c.getUnitsInMap(mapID)
	return units.unsafeRange()
}

func (c *store) GetUnitByIDAndMapID(id, mapID int) *unit.Unit {
	unitsInMap := c.getUnitsInMap(mapID)
	return unitsInMap.getByID(id)
}

func (c *store) RemoveUnitByID(id, mapID int) {

	unitsInMap := c.getUnitsInMap(mapID)
	removeUnit := unitsInMap.getByID(id)

	if removeUnit != nil {
		unitsInMap.removeByID(removeUnit.GetID())
	}
}
