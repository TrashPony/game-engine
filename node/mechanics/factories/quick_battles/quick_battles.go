package quick_battles

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
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

func (b *battles) RemoveBattle(uuid string) {
	b.mx.Lock()
	defer b.mx.Unlock()

	delete(b.battles, uuid)
}

func (b *battles) GetAll(basket []*battle.Battle) []*battle.Battle {
	b.mx.RLock()
	defer b.mx.RUnlock()

	basket = basket[:0]
	for _, mp := range b.battles {
		basket = append(basket, mp)
	}

	return basket
}
