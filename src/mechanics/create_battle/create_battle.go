package create_battle

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/battle"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/player"
	"github.com/TrashPony/game_engine/src/mechanics/generators"
	uuid "github.com/satori/go.uuid"
)

func CreateBattle(startPlayers []*player.Player, mapID int) *battle.Battle {
	newBattle := &battle.Battle{
		UUID: uuid.NewV4().String(),
	}

	newBattle.Map = generators.CreateBattleMap(mapID)

	for _, p := range startPlayers {
		p.GameUUID = newBattle.UUID
		newBattle.AddUser(p)
	}

	return newBattle
}
