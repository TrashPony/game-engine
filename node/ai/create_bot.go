package ai

import (
	"github.com/TrashPony/game-engine/router/const/game_types"
	"github.com/TrashPony/game-engine/router/generate_ids"
	"github.com/TrashPony/game-engine/router/mechanics/factories/players"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/behavior_rule"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/player"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"math/rand"
	"time"
)

func CreateBot(battleUUID string, teamID int, mp *_map.Map, x, y, rotate int, behavior *behavior_rule.BehaviorRules, addInCommonStore, seed bool, bodyID int) *player.Player {

	generateBotID := generate_ids.GetBotID()

	if seed {
		time.Sleep(time.Nanosecond)
		rand.Seed(time.Now().UnixNano())
	}

	botPlayer := &player.Player{
		ID:                 generateBotID,
		GameUUID:           battleUUID,
		MapID:              mp.Id,
		Ready:              true,
		BehaviorController: false,
		Bot:                true,
		TeamID:             teamID,
	}

	newUnit := &unit.Unit{
		ID:      generateBotID,
		OwnerID: botPlayer.GetID(),
		MapID:   mp.Id,
		HP:      100,
	}

	newUnit.Body = game_types.GetNewBody(bodyID)
	getWeapon(newUnit)

	newUnit.GetPhysicalModel().SetPos(float64(x), float64(y), float64(rotate))

	newUnit.BehaviorRules = behavior
	botPlayer.GetGameUnitsStore().AddUnit(newUnit, botPlayer.TeamID, botPlayer.GetID())

	if addInCommonStore {
		players.Users().AddNewPlayer(botPlayer)
	}

	generateBotID--
	return botPlayer
}

func getWeapon(newUnit *unit.Unit) {
	wSlot := newUnit.GetWeaponSlot(1)
	wSlot.Weapon = game_types.GetRandomWeapon()
	wSlot.Ammo = game_types.GetNewAmmo(wSlot.Weapon.DefaultAmmoTypeID)
	wSlot.SetAnchor()
}
