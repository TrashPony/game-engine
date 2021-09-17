package ai

import (
	"github.com/TrashPony/game_engine/src/ai/ai_methods"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/behavior_rule"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/player"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
)

func AI(bot *player.Player, aiUnit *unit.Unit, rule *behavior_rule.BehaviorRule) {

	if rule == nil {
		return
	}

	if rule.Action == "find_hostile_in_range_view" {
		AI(bot, aiUnit, rule.StopRule)
	}

	if rule.Action == "follow_attack_target" {
		AI(bot, aiUnit, rule.StopRule)
	}

	if rule.Action == "scouting" {
		if ai_methods.Scouting(aiUnit) {
			AI(bot, aiUnit, rule.PassRule)
		} else {
			AI(bot, aiUnit, rule.StopRule)
		}
	}
}
