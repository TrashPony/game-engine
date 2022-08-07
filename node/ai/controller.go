package ai

import (
	"github.com/TrashPony/game-engine/node/ai/bot_store"
	"github.com/TrashPony/game-engine/node/mechanics/factories/quick_battles"
	"github.com/TrashPony/game-engine/node/mechanics/find_path"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/player"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"time"
)

func AddLifeBot(bot *player.Player) {
	bot_store.Bots.Add(bot)
}

func BotLifeLoop() {

	aiTick := time.Duration(1000)
	for {

		bots, mx := bot_store.Bots.UnsafeRange()

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
	if bot.GetGameUnitsStore().CountUnits() == 0 {
		go bot_store.Bots.Remove(bot)
		return
	}

	b := quick_battles.Battles.GetBattleByUUID(bot.GameUUID)
	if b == nil {
		return
	}

	// правила поведения работают на уровне юнитов
	gamePlayerUnits, mx := bot.GetGameUnitsStore().UnsafeRangeUnits()
	defer mx.RUnlock()
	for _, u := range gamePlayerUnits {
		if u.BehaviorRules == nil || u.BehaviorRules.Rules == nil || len(u.BehaviorRules.Rules) == 0 {
			continue
		}

		//continue
		AI(bot, u, u.BehaviorRules.Rules[0], b)
		findPathChecker(b, u)
	}
}

func findPathChecker(b *battle.Battle, u *unit.Unit) {
	if u == nil {
		return
	}

	followTarget, _, _, needCalc := u.GetMovePathState()
	if followTarget == nil {
		u.RemoveMovePath()
		return
	}

	if needCalc {
		// PATH FINDER
		path, err := find_path.FindPath(u,
			float64(u.GetPhysicalModel().GetX()), float64(u.GetPhysicalModel().GetY()),
			float64(followTarget.GetX()), float64(followTarget.GetY()),
			nil, true, b.Map)

		if err != nil {
			return
		}

		// что бы юнита не вьежал в стройку
		followTarget.X = path[len(path)-1].X
		followTarget.Y = path[len(path)-1].Y

		u.SetMovePath(&path)
	}
}
