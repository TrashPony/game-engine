package fly_bullets

import (
	"github.com/TrashPony/game-engine/node/mechanics/actual_target"
	"github.com/TrashPony/game-engine/node/mechanics/factories/units"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/bullet"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
)

func Missile(bullet *bullet.Bullet, gameMap *_map.Map, noTime bool, mapUnits []*unit.Unit) {

	distToTarget := game_math.GetBetweenDist(bullet.GetX(), bullet.GetY(), bullet.Target.GetX(), bullet.Target.GetY())

	_, _, lvl := gameMap.GetPosLevel(bullet.GetX(), bullet.GetY())
	if bullet.GetZ()-25 < lvl { //что бы ракета не падала на землю
		bullet.StartRadian = game_math.DegToRadian(2)
	} else {
		bullet.StartRadian = game_math.DegToRadian(1) * -1
	}

	if distToTarget < 15 {
		bullet.SetZ(lvl - 15)
	}

	bullet.StartZ = bullet.GetZ()
	bullet.StartX = bullet.GetX()
	bullet.StartY = bullet.GetY()

	if bullet.Ammo.ChaseTarget {

		// always_to_mouse_target - пуля следует за целью унита
		if bullet.Ammo.ChaseOption == "always_to_mouse_target" {
			if bullet.OwnerType == "unit" {
				gameUnit := units.Units.GetUnitByIDAndMapID(bullet.OwnerID, gameMap.Id)
				if gameUnit != nil {
					wTarget := gameUnit.GetWeaponTarget()
					if actual_target.GetXYZTarget(nil, wTarget, gameMap, nil) {
						bullet.Target.SetX(wTarget.GetX())
						bullet.Target.SetY(wTarget.GetY())
						missileCorrect(bullet)
					}
				}
			}
		}

		// distance_chase - автоматически направляется на врага когда тот оказывается в радиусе захвата (например 200 пикселей)
		if bullet.Ammo.ChaseOption == "distance_chase" && !noTime {

			type missileTarget struct {
				x, y, dist int
			}

			var passObj *missileTarget
			for _, u := range mapUnits {

				if u.TeamID == bullet.OwnerTeamID {
					continue
				}

				dist := int(game_math.GetBetweenDist(bullet.GetX(), bullet.GetY(), u.GetX(), u.GetY()))
				if (dist < bullet.Ammo.ChaseCatchDestination) && (passObj == nil || passObj.dist > dist) {
					passObj = &missileTarget{x: u.GetX(), y: u.GetY(), dist: dist}
				}
			}

			objects, objectsMX := gameMap.UnsafeRangeBuildDynamicObjects()
			defer objectsMX.RUnlock()

			for _, o := range objects {

				if o.TeamID == bullet.OwnerTeamID {
					continue
				}

				dist := int(game_math.GetBetweenDist(bullet.GetX(), bullet.GetY(), o.GetX(), o.GetY()))
				if (dist < bullet.Ammo.ChaseCatchDestination) && (passObj == nil || passObj.dist > dist) {
					passObj = &missileTarget{x: o.GetX(), y: o.GetY(), dist: dist}
				}
			}

			if passObj != nil {
				bullet.Target.SetX(passObj.x)
				bullet.Target.SetY(passObj.y)
				missileCorrect(bullet)
			}
		}

		// always_to_target - преследует цель за которой закреплен
		if bullet.Ammo.ChaseOption == "always_to_target" {
			if actual_target.GetXYZTarget(nil, bullet.ChaseTarget, gameMap, nil) {
				bullet.Target.SetX(bullet.ChaseTarget.X)
				bullet.Target.SetY(bullet.ChaseTarget.Y)
				missileCorrect(bullet)
			}
		}
	}
}

func missileCorrect(bullet *bullet.Bullet) {
	needRotate := game_math.GetBetweenAngle(float64(bullet.Target.GetX()), float64(bullet.Target.GetY()), float64(bullet.GetX()), float64(bullet.GetY()))

	newRotate := bullet.GetRotate()
	game_math.Rotate(&newRotate, &needRotate, bullet.Ammo.Rotate)
	bullet.SetRotate(newRotate)

	bullet.RadRotate = game_math.DegToRadian(bullet.GetRotate())

	bullet.XVelocity, bullet.YVelocity = game_math.SpeedAndAngleToVelocity(bullet.RealSpeed, bullet.RadRotate)
}
