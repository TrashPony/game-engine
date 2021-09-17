package units

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
	"sync"
)

var Units = newStore()

type store struct {
	units map[int]*unitsInMap // карта юнитов в игре

	mx sync.RWMutex

	// отрицательный ид, это переменная итерационно уменьшается на каждом сгенерируемом юните
	generateUnitID int
}

func newStore() *store {
	return &store{
		units: make(map[int]*unitsInMap),
	}
}

func (c *store) GetBotID() int {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.generateUnitID -= 1
	return c.generateUnitID
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
	c.units[mapID] = &unitsInMap{units: make(map[int]*unit.Unit)}
}

func (c *store) GetAllUnitsArray(mapID int) []*unit.Unit {
	unitsArray := make([]*unit.Unit, 0)

	units := c.getUnitsInMap(mapID)
	units.mx.RLock()
	for _, gameUnit := range units.units {
		unitsArray = append(unitsArray, gameUnit)
	}
	units.mx.RUnlock()

	return unitsArray
}

func (c *store) AddUnit(unit *unit.Unit) {
	unitsInMap := c.getUnitsInMap(unit.MapID)
	unitsInMap.addUnit(unit)
}

func (c *store) GetAllUnitsByMapIDChan(mapID int) <-chan *unit.Unit {
	units := c.getUnitsInMap(mapID)
	return units.getAllUnitsChan()
}

func (c *store) GetAllUnitsByMapIDUnsafeRange(mapID int) (map[int]*unit.Unit, *sync.RWMutex) {
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
