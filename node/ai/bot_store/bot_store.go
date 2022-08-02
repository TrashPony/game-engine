package bot_store

import (
	"github.com/TrashPony/game-engine/node/mechanics/factories/units"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/player"
	"sync"
)

var Bots = botStore{bots: make(map[int]*player.Player)}

type botStore struct {
	bots map[int]*player.Player
	mx   sync.RWMutex
}

func (b *botStore) Add(bot *player.Player) {

	// todo меж галактический костыль из за того что счетчик на роутере и ноде разный
	oldBot := b.GetByID(bot.GetID())
	if oldBot != nil {
		bot.ID += 1000
		b.Add(bot)
	}

	b.mx.Lock()
	defer b.mx.Unlock()

	b.bots[bot.GetID()] = bot
	bot.BehaviorController = true
}

func (b *botStore) Remove(bot *player.Player) {
	b.mx.Lock()
	defer b.mx.Unlock()

	for u := range bot.GetGameUnitsStore().RangeUnits() {
		units.Units.RemoveUnitByID(u.ID, u.MapID)
	}
	bot.BehaviorController = false

	delete(b.bots, bot.GetID())
}

func (b *botStore) UnsafeRange() (map[int]*player.Player, *sync.RWMutex) {
	b.mx.RLock()
	return b.bots, &b.mx
}

func (b *botStore) GetByID(id int) *player.Player {
	b.mx.RLock()
	defer b.mx.RUnlock()
	return b.bots[id]
}
