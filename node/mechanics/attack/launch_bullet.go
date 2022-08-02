package attack

import (
	"github.com/TrashPony/game-engine/node/mechanics/actual_target"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/body"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/bullet"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/target"
)

func InitBullet(mp *_map.Map, bullet *bullet.Bullet, gunner Gunner, absoluteAccuracy, noRotate bool, weaponSlot *body.WeaponSlot) {

	// если юнит находится в движение то из за DelayFollowingFire, позиция сместиться
	//firePos := gunner.GetWeaponFirePos(weaponSlot.Number)
	firePos := gunner.GetWeaponFirePosOne(weaponSlot.Number, bullet.FirePos)
	bullet.SetX(firePos.X)
	bullet.SetY(firePos.Y)
	bullet.StartX, bullet.StartY = firePos.X, firePos.Y

	if !noRotate {

		actual_target.GetXYZTarget(gunner, bullet.Target, mp, weaponSlot)
		_, _, lvl := mp.GetPosLevel(gunner.GetX(), gunner.GetY())
		maxRange, _ := gunner.GetWeaponMaxRange(lvl, weaponSlot.Number, true)

		// т.к. пушка может быть повернута в друном направление высчитываем его
		xWeapon, yWeapon := gunner.GetWeaponPosInMap(weaponSlot.Number)
		dist := game_math.GetBetweenDist(xWeapon, yWeapon, bullet.Target.X, bullet.Target.Y)

		if int(dist) > maxRange {
			dist = float64(maxRange) - 1
		}

		x, y := game_math.VectorToAngleBySpeed(float64(xWeapon), float64(yWeapon), dist, weaponSlot.GetGunRotate())
		_, _, targetLvl := mp.GetPosLevel(x, y)
		bullet.Target = &target.Target{Type: "map", X: x, Y: y, Z: targetLvl}
		bullet.Rotate = gunner.GetGunRotate(weaponSlot.Number)
	}

	GetBulletFireAngle(bullet, bullet.Weapon.Artillery, false)
	// влияние точности оружия на выстрел
	AccuracyWeapon(bullet, gunner.GetWeaponAccuracy(weaponSlot.Number), absoluteAccuracy, weaponSlot)
}
