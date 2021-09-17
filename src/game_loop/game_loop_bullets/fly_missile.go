package game_loop_bullets

import (
	"github.com/TrashPony/game_engine/src/mechanics/attack"
	"github.com/TrashPony/game_engine/src/mechanics/factories/units"
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/bullet"
	_map "github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
)

func Missile(bullet *bullet.Bullet, gameMap *_map.Map) {

	_, _, lvl := gameMap.GetPosLevel(bullet.GetX(), bullet.GetY())
	if bullet.GetZ()-100 < lvl { //что бы ракета не падала на землю
		bullet.StartRadian = game_math.DegToRadian(3)
	} else {
		bullet.StartRadian = game_math.DegToRadian(1) * -1
	}

	bullet.StartZ = bullet.GetZ()
	bullet.StartX = bullet.GetX()
	bullet.StartY = bullet.GetY()

	//attack.GetBulletFireAngle(bullet, false, true)

	// todo 2 - преследует свою цель, игнорируя цель юнита (например если цель это юнит в движение)
	if bullet.Ammo.ChaseTarget {

		// 1 - пуля следует за целью унита
		if bullet.Ammo.ChaseOption == "always_to_target" {
			if bullet.OwnerType == "unit" {
				gameUnit := units.Units.GetUnitByIDAndMapID(bullet.OwnerID, gameMap.Id)
				if gameUnit != nil {
					wTarget := gameUnit.GetWeaponTarget()
					if attack.GetXYZTarget(nil, wTarget, gameMap, nil) {
						bullet.Target.SetX(wTarget.GetX())
						bullet.Target.SetY(wTarget.GetY())
						missileCorrect(bullet)
					}
				}
			}
		}

		// автоматически направляется на врага когда тот оказывается в радиусе захвата (например 200 пикселей)
		if bullet.Ammo.ChaseOption == "distance_chase" {

			if bullet.OwnerType == "unit" {
				gameUnit := units.Units.GetUnitByIDAndMapID(bullet.OwnerID, gameMap.Id)
				if gameUnit != nil {
					wTarget := gameUnit.GetWeaponTarget()
					if attack.GetXYZTarget(nil, wTarget, gameMap, nil) {
						bullet.Target.SetX(wTarget.GetX())
						bullet.Target.SetY(wTarget.GetY())
						distToTarget := game_math.GetBetweenDist(bullet.GetX(), bullet.GetY(), bullet.Target.GetX(), bullet.Target.GetY())
						if distToTarget < 100 {
							missileCorrect(bullet)
						}
					}
				}
			}

		}
	}
}

func missileCorrect(bullet *bullet.Bullet) {
	needRotate := game_math.GetBetweenAngle(float64(bullet.Target.GetX()), float64(bullet.Target.GetY()), float64(bullet.GetX()), float64(bullet.GetY()))

	newRotate := bullet.GetRotate()
	attack.Rotate(&newRotate, &needRotate, bullet.Ammo.Rotate)
	bullet.SetRotate(newRotate)

	bullet.RadRotate = game_math.DegToRadian(bullet.GetRotate())

	bullet.XVelocity, bullet.YVelocity = game_math.SpeedAndAngleToVelocity(bullet.RealSpeed, bullet.RadRotate)
}
