package player

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
	"sync"
)

type userUnitStore struct {
	units map[int]*unit.Unit `json:"units"`
	mx    sync.RWMutex
}

func (u *userUnitStore) GetUnit(id int) *unit.Unit {
	u.mx.RLock()
	defer u.mx.RUnlock()
	return u.units[id]
}

func (u *userUnitStore) CountUnits() int {
	u.mx.RLock()
	defer u.mx.RUnlock()
	return len(u.units)
}

func (u *userUnitStore) AddUnit(addUnit *unit.Unit) {
	u.mx.Lock()
	defer u.mx.Unlock()
	u.units[addUnit.ID] = addUnit
}

func (u *userUnitStore) RangeUnits() <-chan *unit.Unit {
	u.mx.RLock()
	units := make(chan *unit.Unit, len(u.units))

	go func() {
		defer func() {
			close(units)
			u.mx.RUnlock()
		}()

		for _, u := range u.units {
			units <- u
		}
	}()

	return units
}
