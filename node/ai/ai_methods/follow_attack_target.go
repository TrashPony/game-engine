package ai_methods

import (
	"github.com/TrashPony/game-engine/node/mechanics/actual_target"
	"github.com/TrashPony/game-engine/node/mechanics/attack"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	battle2 "github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
)

func FollowAttackTarget(b *battle2.Battle, aiUnit *unit.Unit) bool {

	weaponTarget := aiUnit.GetWeaponTarget()
	if weaponTarget == nil {
		return false
	}

	globalMeta := aiUnit.BehaviorRules.Meta
	if globalMeta != nil && globalMeta.Patrol != nil && len(globalMeta.Patrol.Path) > 0 {
		currentPatrolTarget := globalMeta.Patrol.Path[globalMeta.Patrol.ToIDIndex]

		dist := int(game_math.GetBetweenDist(weaponTarget.GetX(), weaponTarget.GetY(), currentPatrolTarget.X, currentPatrolTarget.Y))
		if dist > currentPatrolTarget.Radius*3 {
			return false
		}
	}

	if actual_target.GetXYZTarget(aiUnit.GetGunner(), weaponTarget, b.Map, aiUnit.GetWeaponSlot(1)) {

		weaponDist, _ := aiUnit.GetGunner().GetWeaponMaxRange(weaponTarget.Z, 1, false) // todo хардкод слота
		distToTarget := int(game_math.GetBetweenDist(aiUnit.GetX(), aiUnit.GetY(), weaponTarget.GetX(), weaponTarget.GetY()) * 1.3)
		canAttack := attack.CheckFire(aiUnit.GetGunner(), "unit", aiUnit.GetID(), b, weaponTarget, aiUnit.GetWeaponSlot(1), true, nil)

		if weaponDist > distToTarget && canAttack {
			aiUnit.RemoveMovePath()
		} else {
			followTarget := aiUnit.GetFollowTarget()
			if followTarget == nil || (followTarget.Type != weaponTarget.Type && followTarget.ID != weaponTarget.ID) {
				aiUnit.SetMovePathTarget(weaponTarget)
			} else {
				aiUnit.SetFindMovePath()
			}
		}
	}

	return true
}
