package attack

import (
	"github.com/TrashPony/game-engine/node/mechanics/actual_target"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/body"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/bullet"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/target"
)

func CreateBullets(gunner Gunner, mp *_map.Map, weaponSlot *body.WeaponSlot, gunnerType string, gunnerID int, noLvlMap, nextFirePos bool, fireTarget *target.Target) ([]*bullet.Bullet, bool) {

	_, _, lvl := mp.GetPosLevel(gunner.GetX(), gunner.GetY())
	maxRange, maxAngle := gunner.GetWeaponMaxRange(lvl, weaponSlot.Number, false)
	var cTarget *target.Target

	if chaseTarget(weaponSlot) {
		passObj := getChaseTarget(gunner, weaponSlot, mp, fireTarget)
		if passObj != nil {
			cTarget = &target.Target{Type: passObj.realType, ID: passObj.realID}
			actual_target.GetXYZTarget(nil, cTarget, mp, nil)
		}
	}

	bullets := make([]*bullet.Bullet, weaponSlot.Weapon.CountFireBullet)
	//  создаем обьект пули, дать ему направление и начальную позицию
	for i := 0; i < weaponSlot.Weapon.CountFireBullet; i++ {
		// если количество пулей больше чем точек то пули вылетают по кругу

		// для каждой пули своя цель
		bulletTarget := fireTarget.GetCopy()

		damage := getDamage(gunner, weaponSlot)

		_, minAngle := gunner.GetWeaponMinRange(lvl, weaponSlot.Number)

		speed := weaponSlot.GetAmmo().BulletSpeed + weaponSlot.Weapon.BulletSpeed
		newBullet := &bullet.Bullet{
			FirePos:     weaponSlot.GetLastFirePosition(),
			Weapon:      weaponSlot.Weapon,
			Ammo:        weaponSlot.GetAmmo(),
			Speed:       speed,
			Target:      bulletTarget,
			OwnerTeamID: gunner.GetTeamID(),
			OwnerType:   gunnerType,
			OwnerID:     gunnerID,
			Damage:      damage,
			MaxRange:    maxRange,
			Rotate:      gunner.GetGunRotate(weaponSlot.Number),
			MapID:       mp.Id,
			HP:          1,
			MaxAngle:    game_math.DegToRadian(maxAngle),
			MinAngle:    game_math.DegToRadian(minAngle),
			ChaseTarget: cTarget,
		}

		if nextFirePos {
			weaponSlot.NextLastFirePosition()
		}

		if noLvlMap {
			newBullet.SetZ(gunner.GetMapHeight())
			newBullet.StartZ = gunner.GetMapHeight()
		} else {
			newBullet.SetZ(lvl + gunner.GetMapHeight())
			newBullet.StartZ = lvl + gunner.GetMapHeight()
		}

		bullets[i] = newBullet
	}

	return bullets, true
}

func getDamage(gunner Gunner, weaponSlot *body.WeaponSlot) int {
	var damage int
	// высчитываем урон из уровня зарядки вепона
	if weaponSlot.Weapon.AccumulationFirePower {
		// мин дамаг снаряда это при 0 процентах, а максимальный при 100
		diffMinMax := gunner.GetMaxDamage(weaponSlot.Number) - gunner.GetMinDamage(weaponSlot.Number)
		percent := (weaponSlot.AccumulationCurrent / weaponSlot.Weapon.AccumulationFull) * 100
		damage = gunner.GetMinDamage(weaponSlot.Number) + int(float64(diffMinMax)*(percent/100))
	} else {
		damage = gunner.GetDamage(weaponSlot.Number)
	}

	return damage
}
