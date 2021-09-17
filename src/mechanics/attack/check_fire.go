package attack

import (
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/mechanics/factories/units"
	"github.com/TrashPony/game_engine/src/mechanics/fly_bullets"
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_math/collisions"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/body"
	_map "github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/target"
)

func CheckFire(gunner Gunner, gunnerType string, gunnerID int, mp *_map.Map, wTarget *target.Target, weaponSlot *body.WeaponSlot) bool {

	// по какойто причине цель не достяжима
	if !GetXYZTarget(gunner, wTarget, mp, weaponSlot) {
		return false
	}

	_, _, targetLvl := mp.GetPosLevel(wTarget.GetX(), wTarget.GetY())
	maxRange, _ := gunner.GetWeaponMaxRange(targetLvl, weaponSlot.Number)
	if maxRange < int(game_math.GetBetweenDist(wTarget.GetX(), wTarget.GetY(), gunner.GetX(), gunner.GetY())) {
		// цель дальше чем может выстрелить оружие
		return false
	}

	gunner.SetWeaponTarget(wTarget)
	x, y := collisionWeaponRangeCollision(gunner, mp, weaponSlot, gunnerType, gunnerID)
	if x == -1 && y == -1 || wTarget == nil {
		return false
	}

	radius := weaponSlot.GetAmmo().AreaCovers
	if radius < _const.AmmoRadius {
		radius = _const.AmmoRadius
	} else {
		radius = radius / 2
	}

	if wTarget.Type == "map" {
		return true
	}

	if wTarget.Type == "unit" {
		targetUnit := units.Units.GetUnitByIDAndMapID(wTarget.ID, mp.Id)
		if targetUnit != nil {
			return collisions.CircleUnit(x, y, radius, targetUnit)
		}
	}

	if wTarget.Type == "object" {
		obj := mp.GetDynamicObjectsByID(wTarget.ID)
		if obj != nil {
			return collisions.CircleDynamicObj(x, y, radius, obj, false)
		}
	}

	return false
}

// метод что оружие может стрелять в цель (между оружием и целью нет колизий)
func collisionWeaponRangeCollision(gunner Gunner, mp *_map.Map, weaponSlot *body.WeaponSlot, gunnerType string, gunnerID int) (int, int) {

	bullets, _ := CreateBullets(gunner, mp, gunner.GetWeaponSlot(weaponSlot.Number), gunnerType, gunnerID, false, false)
	if bullets == nil {
		return -1, -1
	}

	bullet := bullets[0]
	InitBullet(bullet, gunner, true, weaponSlot)

	// лазеры особенные :)
	if weaponSlot.Weapon.Type == _const.LaserWeapon {
		x, y, _, _ := fly_bullets.FlyLaser(bullet, mp, true)
		return x, y
	} else {
		return fly_bullets.FlyBullet(bullet, mp, true)
	}
}
