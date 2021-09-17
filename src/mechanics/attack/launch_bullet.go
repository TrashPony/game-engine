package attack

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/body"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/bullet"
)

func InitBullet(bullet *bullet.Bullet, gunner Gunner, absoluteAccuracy bool, weaponSlot *body.WeaponSlot) {

	// если юнит находится в движение то из за DelayFollowingFire, позиция сместиться
	bullet.SetX(bullet.FirePos.X)
	bullet.SetY(bullet.FirePos.Y)
	bullet.StartX, bullet.StartY = bullet.FirePos.X, bullet.FirePos.Y

	GetBulletFireAngle(bullet, false, false)
	// влияние точности оружия на выстрел
	AccuracyWeapon(bullet, gunner.GetWeaponAccuracy(weaponSlot.Number), absoluteAccuracy)
}
