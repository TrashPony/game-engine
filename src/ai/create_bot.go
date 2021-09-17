package ai

import (
	"github.com/TrashPony/game_engine/src/mechanics/factories/units"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/behavior_rule"
	_map "github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/player"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
	"strconv"
	"sync"
)

var generateBotID = -1
var mx sync.Mutex

func CreateBot(battleUUID string, mp *_map.Map, x, y, rotate int, behavior *behavior_rule.BehaviorRules) *player.Player {
	mx.Lock()
	defer mx.Unlock()

	botPlayer := &player.Player{
		ID:                 generateBotID,
		Login:              "bot" + strconv.Itoa(generateBotID),
		GameUUID:           battleUUID,
		MapID:              mp.Id,
		Ready:              true,
		BehaviorController: false,
	}
	newUnit := &unit.Unit{ID: units.Units.GetBotID(), OwnerID: botPlayer.GetID(), MapID: mp.Id, HP: 100}
	newUnit.GetPhysicalModel().SetPos(float64(x), float64(y), float64(rotate))

	newUnit.BehaviorRules = behavior

	botPlayer.GetUnitsStore().AddUnit(newUnit)
	//units.Units.AddUnit(newUnit)
	//AddLifeBot(botPlayer)

	generateBotID--
	return botPlayer
}
