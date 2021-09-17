package attack

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/bullet"
	"math"
)

func AccuracyWeapon(bullet *bullet.Bullet, accuracy int, absoluteAccuracy bool) {
	if bullet.Ammo.Type == "firearms" || bullet.Ammo.Type == "laser" || (bullet.Ammo.Type == "missile" && !bullet.Ammo.ChaseTarget) {
		firearmsAccuracy(bullet, float64(accuracy), absoluteAccuracy)
	}
}

func firearmsAccuracy(bullet *bullet.Bullet, accuracy float64, absoluteAccuracy bool) {
	// чем дальше атака тем она менее точнее от 0 до 10ти
	oldAccuracy := accuracy
	if !absoluteAccuracy {
		accuracy = GetGunAccuracy(bullet.GetX(), bullet.GetY(), bullet.Target.GetX(), bullet.Target.GetY(), bullet.MaxRange, float64(bullet.Speed),
			bullet.StartZ, bullet.Target.GetZ(), 45.0, accuracy, bullet.Weapon.Type)

		if accuracy == math.NaN() || int(accuracy) < 0 {
			accuracy = oldAccuracy
		}

		accuracy = float64(game_math.GetRangeRand(0, int(accuracy)))
		accuracyRotate := game_math.DegToRadian(float64(game_math.GetRangeRand(0, 360)))

		bullet.Target.SetX(bullet.Target.GetX() + int(accuracy*game_math.Cos(accuracyRotate)))
		bullet.Target.SetY(bullet.Target.GetY() + int(accuracy*game_math.Sin(accuracyRotate)))

	}

	GetBulletFireAngle(bullet, false, false)
}

func GetGunAccuracy(x, y, xTarget, yTarget, maxWeaponRange int, bulletSpeed, z, zTarget, gunZAngle, accuracy float64, weaponType string) float64 {
	// чем дальше атака тем она менее точнее от 0 до 10ти
	dist := game_math.GetBetweenDist(x, y, xTarget, yTarget)

	maxRange := 0.0
	if weaponType == "missile" || weaponType == "laser" {
		maxRange = float64(maxWeaponRange)
	} else {
		maxRange = game_math.GetMaxDistBallisticWeapon(bulletSpeed, zTarget,
			z, game_math.DegToRadian(gunZAngle), weaponType)
	}

	percentRange := (dist * 100) / maxRange

	randomK := percentRange / 10

	return randomK * accuracy
}
