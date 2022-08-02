package ai

import (
	"github.com/TrashPony/game-engine/node/ai/ai_methods"
	battle2 "github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/behavior_rule"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/player"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
)

func AI(bot *player.Player, aiUnit *unit.Unit, rule *behavior_rule.BehaviorRule, b *battle2.Battle) {

	if rule == nil {
		return
	}

	if rule.Action == "find_hostile_in_range_view" {
		if ai_methods.FindHostileInRangeView(bot, aiUnit, 0, b, rule.Meta) {
			AI(bot, aiUnit, rule.PassRule, b)
		} else {
			AI(bot, aiUnit, rule.StopRule, b)
		}
	}

	if rule.Action == "follow_attack_target" {
		if ai_methods.FollowAttackTarget(b, aiUnit) {
			AI(bot, aiUnit, rule.PassRule, b)
		} else {
			AI(bot, aiUnit, rule.StopRule, b)
		}
	}

	if rule.Action == "scouting" {
		if ai_methods.Scouting(b, aiUnit) {
			AI(bot, aiUnit, rule.PassRule, b)
		} else {
			AI(bot, aiUnit, rule.StopRule, b)
		}
	}
}
