package ai_methods

import (
	"github.com/TrashPony/game-engine/node/ai/check_target"
	"github.com/TrashPony/game-engine/node/mechanics/actual_target"
	"github.com/TrashPony/game-engine/node/mechanics/attack"
	"github.com/TrashPony/game-engine/node/mechanics/factories/quick_battles"
	"github.com/TrashPony/game-engine/node/mechanics/factories/units"
	"github.com/TrashPony/game-engine/node/mechanics/watch"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	battle2 "github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/behavior_rule"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/player"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/target"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
)

func FindHostileInRangeView(bot *player.Player, aiUnit *unit.Unit, maxDist int, b *battle2.Battle, meta *behavior_rule.Meta) bool {

	team := b.Teams[bot.TeamID]
	if team == nil {
		return false
	}

	currentTarget := aiUnit.GetWeaponTarget()

	// поиск противника в дальности видимости
	hostiles := check_target.FindHostile(team, aiUnit.GetX(), aiUnit.GetY(), b.Map, bot.TeamID, maxDist, meta)
	typeTarget, idTarget := takePriorityTarget(hostiles, aiUnit, b, meta)

	// проверяем что бы у юнита небыло цели, и что бы текущая цель была все еще в пределах видимости
	if currentTarget != nil && (currentTarget.Type == "unit" || currentTarget.Type == "object") {
		if actual_target.GetXYZTarget(aiUnit.GetGunner(), currentTarget, b.Map, aiUnit.GetWeaponSlot(1)) { // todo хардкод слота

			b := quick_battles.Battles.GetBattleByUUID(bot.GameUUID)

			view, radarView := watch.CheckViewCoordinate(team, currentTarget.GetX(), currentTarget.GetY(), b.Map.Id, b, units.Units.GetAllUnitsArray(b.Map.Id, make([]*unit.Unit, 10)))

			hostile, _, _ := check_target.CheckTarget(team, currentTarget, b.Map, bot.TeamID)
			if (view || radarView) && hostile {
				if currentTarget.Type == typeTarget && currentTarget.ID == idTarget {
					return true
				}
			}
		}
	}

	if typeTarget != "" {
		aiUnit.SetWeaponTarget(&target.Target{Type: typeTarget, ID: idTarget, Attack: true})
		return true
	}
	// если некого не видно то и цели нет
	aiUnit.SetWeaponTarget(nil)
	return false
}

func takePriorityTarget(hostiles []*check_target.Hostile, aiUnit *unit.Unit, b *battle2.Battle, meta *behavior_rule.Meta) (string, int) {
	// тут мы смотрим всех кандидатов и выбираем самую приоритетную цель по каким либо критериям

	if len(hostiles) == 0 {
		return "", 0
	}

	maxHate := &check_target.Hostile{}
	currentTarget := aiUnit.GetWeaponTarget()

	for _, hostile := range hostiles {

		canAttack := attack.CheckFire(aiUnit.GetGunner(), "unit", aiUnit.GetID(), b, &target.Target{Type: hostile.Type, ID: hostile.ID}, aiUnit.GetWeaponSlot(1), true, nil) // todo хардкод слота
		if canAttack {
			hostile.Points = 100
		}

		if meta != nil && meta.Type == "zone" {
			dist := int(game_math.GetBetweenDist(hostile.X, hostile.Y, meta.X, meta.Y))
			if dist <= meta.Radius {
				hostile.Points += hostile.Points * 2
			}
		}

		// todo тут можно докидывать поинты

		if maxHate.Points < hostile.Points || maxHate.Points == 0 {
			maxHate = hostile
		}
	}

	aiUnit.SetWeaponTarget(currentTarget) // todo костыль из за того что указывается цель в методе check_fire.CheckFire
	if currentTarget != nil {
		// на случай если с максимальной ненавистью несколько прентендентов а у бота уже есть цель
		for _, hostile := range hostiles {
			if currentTarget.Type == hostile.Type && currentTarget.ID == hostile.ID && hostile.Points == maxHate.Points {
				return hostile.Type, hostile.ID
			}
		}
	}

	return maxHate.Type, maxHate.ID
}
