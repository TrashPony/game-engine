package ai

import (
	"github.com/TrashPony/game_engine/src/ai/bot_store"
	"github.com/TrashPony/game_engine/src/mechanics/debug"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/player"
	"time"
)

func AddLifeBot(bot *player.Player) {
	bot_store.Bots.Add(bot)
}

func BotLifeLoop() {

	aiTick := time.Duration(1000)
	for {

		bots, mx := bot_store.Bots.UnsafeRange()

		debug.ServerState.SetBots(len(bots))

		for _, bot := range bots {
			mx.RUnlock()
			behavior(bot)
			mx.RLock()
		}
		mx.RUnlock()

		time.Sleep(aiTick * time.Millisecond)
	}
}

func behavior(bot *player.Player) {
	// у 1го бота есть ток 1 юнит в подчинение и если он умер то бот полностью умирает
	if bot.GetUnitsStore().CountUnits() == 0 {
		go bot_store.Bots.Remove(bot)
		return
	}

	// правила поведения работают на уровне юнитов
	for u := range bot.GetUnitsStore().RangeUnits() {
		AI(bot, u, u.BehaviorRules.Rules[0])
	}
}
