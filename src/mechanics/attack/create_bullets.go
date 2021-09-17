package attack

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/body"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/bullet"
	_map "github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	uuid "github.com/satori/go.uuid"
)

func CreateBullets(gunner Gunner, mp *_map.Map, weaponSlot *body.WeaponSlot, gunnerType string, gunnerID int, noLvlMap, nextFirePos bool) ([]*bullet.Bullet, bool) {

	bullets := make([]*bullet.Bullet, 0)
	firePos := gunner.GetWeaponFirePos(weaponSlot.Number)

	fireTarget := gunner.GetWeaponTarget()
	if !GetXYZTarget(gunner, fireTarget, mp, weaponSlot) {
		gunner.SetWeaponTarget(nil)
		return nil, false
	}

	//  создаем обьект пули, дать ему направление и начальную позицию
	for i := 0; i < weaponSlot.Weapon.CountFireBullet; i++ {
		// если количество пулей больше чем точек то пули вылетают по кругу

		// для каждой пули своя цель
		bulletTarget := fireTarget.GetCopy()

		damage := gunner.GetDamage(weaponSlot.Number)
		_, _, lvl := mp.GetPosLevel(gunner.GetX(), gunner.GetY())

		maxRange, maxAngle := gunner.GetWeaponMaxRange(lvl, weaponSlot.Number)
		_, minAngle := gunner.GetWeaponMinRange(lvl, weaponSlot.Number)

		speed := weaponSlot.Weapon.BulletSpeed + weaponSlot.GetAmmo().BulletSpeed
		newBullet := &bullet.Bullet{
			UUID:      uuid.NewV1().String(),
			FirePos:   firePos[weaponSlot.GetLastFirePosition()],
			Weapon:    weaponSlot.Weapon,
			Ammo:      weaponSlot.GetAmmo(),
			Speed:     speed,
			Target:    bulletTarget,
			OwnerType: gunnerType,
			OwnerID:   gunnerID,
			Damage:    damage,
			MaxRange:  maxRange,
			Rotate:    gunner.GetGunRotate(weaponSlot.Number),
			MapID:     mp.Id,
			HP:        1,
			MaxAngle:  game_math.DegToRadian(maxAngle),
			MinAngle:  game_math.DegToRadian(minAngle),
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

		bullets = append(bullets, newBullet)
	}

	return bullets, true
}
