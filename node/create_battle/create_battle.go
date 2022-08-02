package create_battle

import (
	"github.com/TrashPony/game-engine/node/mechanics/generators"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/player"
	uuid "github.com/satori/go.uuid"
)

func CreateBattle(startPlayers []*player.Player, mapID int) *battle.Battle {
	newBattle := &battle.Battle{
		UUID:                uuid.NewV4().String(),
		SessionPlayersState: make(map[int]*battle.SessionPlayer),
		Teams: map[int]*battle.Team{
			1: {ID: 1, PlayersIDs: make([]int, 0)},
			2: {ID: 2, PlayersIDs: make([]int, 0)},
		},
	}

	newBattle.Map = generators.CreateBattleMap(mapID)
	for _, p := range startPlayers {
		p.GameUUID = newBattle.UUID
		p.TeamID = 1
		newBattle.Teams[1].PlayersIDs = append(newBattle.Teams[1].PlayersIDs, p.GetID())
		newBattle.AddUser(p)
	}

	return newBattle
}
