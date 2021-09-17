package quick_battles

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/battle"
	"sync"
)

var Battles = newStore()

type battles struct {
	battles      map[string]*battle.Battle // [uuid]
	customGameID int
	mx           sync.RWMutex
}

func newStore() *battles {
	battles := &battles{
		battles: make(map[string]*battle.Battle),
	}

	return battles
}

func (b *battles) AddNewGame(newGame *battle.Battle) {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.battles[newGame.UUID] = newGame
}

func (b *battles) GetBattleByUUID(uuid string) *battle.Battle {
	b.mx.Lock()
	defer b.mx.Unlock()

	return b.battles[uuid]
}

func (b *battles) removeBattle(uuid string) {
	b.mx.Lock()
	defer b.mx.Unlock()

	delete(b.battles, uuid)
}

func (b *battles) GetBattleByMapID(id int) *battle.Battle {
	b.mx.Lock()
	defer b.mx.Unlock()

	for _, b := range b.battles {
		if b.Map.Id == id {
			return b
		}
	}

	return nil
}
