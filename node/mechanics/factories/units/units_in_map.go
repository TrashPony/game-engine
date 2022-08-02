package units

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"sync"
)

type unitsInMap struct {
	units []*unit.Unit
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

func (units *unitsInMap) unsafeRange() ([]*unit.Unit, *sync.RWMutex) {
	units.mx.RLock()
	return units.units, &units.mx
}

func (units *unitsInMap) getByID(id int) *unit.Unit {
	units.mx.RLock()
	defer units.mx.RUnlock()

	for _, u := range units.units {
		if u.ID == id {
			return u
		}
	}

	return nil
}

func (units *unitsInMap) addUnit(unit *unit.Unit) {

	if unit == nil {
		return
	}

	units.mx.Lock()
	defer units.mx.Unlock()

	// todo костыль
	for _, u := range units.units {
		if u.ID == unit.ID {
			return
		}
	}

	units.units = append(units.units, unit)
}

func (units *unitsInMap) removeByID(id int) {
	units.mx.Lock()
	defer units.mx.Unlock()

	index := -1
	for i, u := range units.units {
		if u.ID == id {
			index = i
			break
		}
	}

	if index >= 0 {
		units.units[index] = units.units[len(units.units)-1]
		units.units = units.units[:len(units.units)-1]
	}
}
