package player

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"sync"
)

type userUnitStore struct {
	units []*unit.Unit
	mx    sync.RWMutex
}

func (u *userUnitStore) GetUnit(id int) *unit.Unit {
	u.mx.RLock()
	defer u.mx.RUnlock()

	for _, un := range u.units {
		if un.ID == id {
			return un
		}
	}

	return nil
}

func (u *userUnitStore) CountUnits() int {
	u.mx.RLock()
	defer u.mx.RUnlock()
	return len(u.units)
}

func (u *userUnitStore) AddUnit(addUnit *unit.Unit, teamID int, ownerID int) {
	u.mx.Lock()
	defer u.mx.Unlock()

	addUnit.TeamID = teamID
	addUnit.OwnerID = ownerID

	u.units = append(u.units, addUnit)
}

func (u *userUnitStore) RemoveUnit(id int) {
	u.mx.Lock()
	defer u.mx.Unlock()

	index := -1
	for i, un := range u.units {
		if un.ID == id {
			index = i
			break
		}
	}

	if index >= 0 {
		u.units[index] = u.units[len(u.units)-1]
		u.units = u.units[:len(u.units)-1]
	}
}

func (u *userUnitStore) RangeUnits() <-chan *unit.Unit {
	u.mx.RLock()
	units := make(chan *unit.Unit, len(u.units))

	go func() {
		defer func() {
			close(units)
			u.mx.RUnlock()
		}()

		for _, un := range u.units {
			units <- un
		}
	}()

	return units
}

func (u *userUnitStore) UnsafeRangeUnits() ([]*unit.Unit, *sync.RWMutex) {
	u.mx.RLock()
	return u.units, &u.mx
}
