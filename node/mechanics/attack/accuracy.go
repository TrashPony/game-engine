package attack

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/body"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/bullet"
	"math"
)

func AccuracyWeapon(bullet *bullet.Bullet, accuracy int, absoluteAccuracy bool, weaponSlot *body.WeaponSlot) {
	if bullet.Ammo.Type == "firearms" || bullet.Ammo.Type == "laser" || (bullet.Ammo.Type == "missile" && !chaseTarget(weaponSlot)) {
		firearmsAccuracy(bullet, float64(accuracy), absoluteAccuracy)
	}
}

func firearmsAccuracy(bullet *bullet.Bullet, accuracy float64, absoluteAccuracy bool) {
	// чем дальше атака тем она менее точнее от 0 до 10ти
	oldAccuracy := accuracy
	if !absoluteAccuracy {
		accuracy = GetGunAccuracy(bullet.GetX(), bullet.GetY(), bullet.Target.GetX(), bullet.Target.GetY(), bullet.MaxRange, float64(bullet.Speed),
			bullet.StartZ, bullet.Target.GetZ(), 45.0, accuracy, bullet.Weapon.Type, bullet.Ammo.Gravity)

		if accuracy == math.NaN() || int(accuracy) < 0 {
			accuracy = oldAccuracy
		}

		accuracy = float64(game_math.GetRangeRand(0, int(accuracy)))
		accuracyRotate := game_math.DegToRadian(float64(game_math.GetRangeRand(0, 360)))

		accuracy = accuracy / 2

		bullet.Target.SetX(bullet.Target.GetX() + int(accuracy*game_math.Cos(accuracyRotate)))
		bullet.Target.SetY(bullet.Target.GetY() + int(accuracy*game_math.Sin(accuracyRotate)))
	}

	bullet.SetRotate(game_math.GetBetweenAngle(float64(bullet.Target.GetX()), float64(bullet.Target.GetY()), float64(bullet.GetX()), float64(bullet.GetY())))
	GetBulletFireAngle(bullet, bullet.Weapon.Artillery, false)
}

func GetGunAccuracy(x, y, xTarget, yTarget, maxWeaponRange int, bulletSpeed, z, zTarget, gunZAngle, accuracy float64, weaponType string, gravity float64) float64 {
	// чем дальше атака тем она менее точнее от 0 до 10ти
	dist := game_math.GetBetweenDist(x, y, xTarget, yTarget)

	maxRange := 0.0
	if maxWeaponRange > 0 {
		maxRange = float64(maxWeaponRange)
	} else {
		maxRange = game_math.GetMaxDistBallisticWeapon(bulletSpeed, zTarget, z, game_math.DegToRadian(gunZAngle), weaponType, gravity)
	}

	percentRange := (dist * 100) / maxRange

	randomK := percentRange / 10

	return randomK * accuracy
}
