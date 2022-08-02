package watch

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
)

func CheckViewCoordinate(watcher Watcher, x, y, mpID int, b *battle.Battle, units []*unit.Unit) (bool, bool) {
	// todo тут могла быть ваша формула расчета видимости
	return true, true
}
