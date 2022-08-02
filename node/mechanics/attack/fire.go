package attack

import (
	"github.com/TrashPony/game-engine/node/mechanics/factories/bullets"
	"github.com/TrashPony/game-engine/node/mechanics/fly_bullets"
	_const "github.com/TrashPony/game-engine/router/const"
	battle2 "github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/body"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/bullet"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/target"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
)

type FireMessage struct {
	TypeID              int `json:"type_id"`
	X                   int `json:"x"`
	Y                   int `json:"y"`
	Z                   int `json:"z"`
	Rotate              int `json:"r"`
	MapID               int `json:"m"`
	AccumulationPercent int `json:"ap"`
}

func Fire(gunner Gunner, weaponSlot *body.WeaponSlot, gunnerType string, gunnerID int, b *battle2.Battle, units []*unit.Unit) []FireMessage {

	fireWeapons := make([]FireMessage, 0)
	weaponTarget := gunner.GetWeaponTarget()
	if weaponTarget == nil {
		return fireWeapons
	}

	if weaponSlot != nil && weaponSlot.GetReload() {
		weaponSlot.SetCurrentReload(weaponSlot.GetCurrentReload() - _const.ServerTick)
		return fireWeapons
	}

	if weaponSlot.Weapon.AccumulationFirePower {
		if !accumulationFirePowerWeapon(gunner, b.Map, weaponSlot, weaponTarget, &fireWeapons) {
			return fireWeapons
		}
	}

	fire(gunner, gunnerType, gunnerID, weaponSlot, weaponTarget, &fireWeapons, b, units)

	return fireWeapons
}

func fire(gunner Gunner, gunnerType string, gunnerID int, weaponSlot *body.WeaponSlot, weaponTarget *target.Target, fireWeapons *[]FireMessage, b *battle2.Battle, units []*unit.Unit) {

	oldBullets, pos := gunner.GetBurstOfShots().GetBullets(weaponSlot.Number)
	if (weaponSlot == nil || weaponSlot.GetReload() || (!weaponTarget.Attack && weaponSlot.AccumulationCurrent == 0)) && oldBullets == nil {
		return
	}

	defer func() {
		weaponSlot.AccumulationTimeOut = 0
		weaponSlot.AccumulationCurrent = 0
	}()

	if oldBullets == nil {

		if CheckFire(gunner, gunnerType, gunnerID, b, weaponTarget, weaponSlot, false, units) || weaponTarget.Force {

			newBullets, fire := CreateBullets(gunner, b.Map, weaponSlot, gunnerType, gunnerID, false, true, weaponTarget)
			if fire {
				gunner.GetBurstOfShots().AddBullets(weaponSlot.Number, newBullets)
				startAttack(gunner, newBullets, 0, weaponSlot, b.Map, fireWeapons, b, units)
			}
		}
	} else {
		startAttack(gunner, oldBullets, pos, weaponSlot, b.Map, fireWeapons, b, units)
	}
}

func startAttack(gunner Gunner, newBullets []*bullet.Bullet, pos int, weaponSlot *body.WeaponSlot, mp *_map.Map, fireWeapons *[]FireMessage, b *battle2.Battle, units []*unit.Unit) {

	timeTick := _const.ServerTick - 1

	for timeTick > 0 && len(newBullets) > pos {

		InitBullet(mp, newBullets[pos], gunner, false, false, weaponSlot)

		if len(*fireWeapons) == 0 {
			// для отыгрыша анимации выстрела
			*fireWeapons = append(*fireWeapons, FireMessage{
				TypeID: weaponSlot.Weapon.ID,
				X:      newBullets[pos].GetX(),
				Y:      newBullets[pos].GetY(),
				Z:      int(newBullets[pos].GetZ()),
				Rotate: int(newBullets[pos].GetRotate()),
				MapID:  mp.Id,
			})
		}

		// лазеры особенные :)
		if weaponSlot.Weapon.Type == _const.LaserWeapon {
			bullets.Bullets.AddBullet(newBullets[pos])
		} else {
			fly_bullets.FlyBullet(newBullets[pos], b, false, units)
		}

		if weaponSlot.Weapon.DelayFollowingFire < _const.ServerTick {
			timeTick -= weaponSlot.Weapon.DelayFollowingFire
		} else {
			weaponSlot.StartReload(weaponSlot.Weapon.DelayFollowingFire, false)
			timeTick = 0
		}

		pos++
		weaponSlot.SetAmmoQuantity(weaponSlot.GetAmmoQuantity() - 1)
	}

	if len(newBullets) == pos {

		if weaponSlot.GetAmmoQuantity() <= 0 {
			weaponSlot.StartReload(gunner.GetWeaponReloadAmmoTime(weaponSlot.Number), true)
			weaponSlot.AmmoQuantity = weaponSlot.Weapon.AmmoCapacity
		} else {
			weaponSlot.StartReload(gunner.GetWeaponReloadTime(weaponSlot.Number), false)
		}

		gunner.GetBurstOfShots().AddBullets(weaponSlot.Number, nil)
	} else {
		gunner.GetBurstOfShots().ChangePos(weaponSlot.Number, pos)
	}
}

func accumulationFirePowerWeapon(gunner Gunner, mp *_map.Map, weaponSlot *body.WeaponSlot, weaponTarget *target.Target, fireWeapons *[]FireMessage) bool {

	if weaponSlot == nil || weaponSlot.GetReload() {
		return false
	}

	if weaponTarget.Attack && (weaponSlot.Weapon.AccumulationFullTimeOut > weaponSlot.AccumulationTimeOut) {
		weaponSlot.AccumulationCurrent += _const.ServerTick

		if weaponSlot.AccumulationCurrent > weaponSlot.Weapon.AccumulationFull {
			weaponSlot.AccumulationCurrent = weaponSlot.Weapon.AccumulationFull
			weaponSlot.AccumulationTimeOut += _const.ServerTick
		}

		x, y := gunner.GetWeaponPosInMap(weaponSlot.Number)
		// сообщение о зарядке для анимации для всех

		_, _, lvl := mp.GetPosLevel(gunner.GetX(), gunner.GetY())
		*fireWeapons = append(*fireWeapons, FireMessage{
			TypeID:              weaponSlot.Weapon.ID,
			X:                   x,
			Y:                   y,
			Z:                   int(lvl + gunner.GetMapHeight()),
			Rotate:              int(weaponSlot.GetGunRotate()),
			AccumulationPercent: int((weaponSlot.AccumulationCurrent/weaponSlot.Weapon.AccumulationFull)*100) + 1,
			MapID:               mp.Id,
		})
		return false
	}

	if weaponSlot.AccumulationCurrent > 0 && weaponSlot.Weapon.AccumulationFullTimeOut <= weaponSlot.AccumulationTimeOut {
		weaponTarget.Force = true
	}

	return true
}
