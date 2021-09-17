package units

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
	"sync"
)

type unitsInMap struct {
	units map[int]*unit.Unit
	mx    sync.RWMutex
}

func (units *unitsInMap) getAllUnitsChan() <-chan *unit.Unit {
	units.mx.RLock()
	unitChan := make(chan *unit.Unit, len(units.units))

	go func() {
		defer func() {
			close(unitChan)
			units.mx.RUnlock()
		}()

		for _, gameUnit := range units.units {
			unitChan <- gameUnit
		}
	}()

	return unitChan
}

func (units *unitsInMap) unsafeRange() (map[int]*unit.Unit, *sync.RWMutex) {
	units.mx.RLock()
	return units.units, &units.mx
}

func (units *unitsInMap) getByID(id int) *unit.Unit {
	units.mx.RLock()
	defer units.mx.RUnlock()
	return units.units[id]
}

func (units *unitsInMap) addUnit(unit *unit.Unit) {

	if unit == nil {
		return
	}

	units.mx.Lock()
	defer units.mx.Unlock()

	units.units[unit.GetID()] = unit
}

func (units *unitsInMap) removeByID(id int) {
	units.mx.Lock()
	defer units.mx.Unlock()
	delete(units.units, id)
}
