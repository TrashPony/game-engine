package attack

import (
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/mechanics/factories/bullets"
	"github.com/TrashPony/game_engine/src/mechanics/fly_bullets"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/body"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/bullet"
	_map "github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/target"
)

type FireMessage struct {
	TypeID int `json:"type_id"`
	X      int `json:"x"`
	Y      int `json:"y"`
	Z      int `json:"z"`
	Rotate int `json:"r"`
	MapID  int `json:"m"`
}

func Fire(gunner Gunner, weaponSlot *body.WeaponSlot, gunnerType string, gunnerID int, mp *_map.Map) []FireMessage {

	fireWeapons := make([]FireMessage, 0)
	weaponTarget := gunner.GetWeaponTarget()
	if weaponTarget == nil {
		return fireWeapons
	}

	if weaponSlot != nil && weaponSlot.GetReload() {
		weaponSlot.SetCurrentReload(weaponSlot.GetCurrentReload() - _const.ServerTick)
		return fireWeapons
	}

	fire(gunner, gunnerType, gunnerID, mp, weaponSlot, weaponTarget, &fireWeapons)

	return fireWeapons
}

func fire(gunner Gunner, gunnerType string, gunnerID int, mp *_map.Map, weaponSlot *body.WeaponSlot, weaponTarget *target.Target, fireWeapons *[]FireMessage) {

	if weaponSlot == nil || weaponSlot.GetReload() || !weaponTarget.Attack {
		return
	}

	bullets, pos := gunner.GetBurstOfShots().GetBullets(weaponSlot.Number)
	if bullets == nil {

		if CheckFire(gunner, gunnerType, gunnerID, mp, weaponTarget, weaponSlot) || weaponTarget.Force {

			bullets, fire := CreateBullets(gunner, mp, weaponSlot, gunnerType, gunnerID, false, true)
			if fire {
				gunner.GetBurstOfShots().AddBullets(weaponSlot.Number, bullets)
				startAttack(gunner, bullets, 0, weaponSlot, mp, fireWeapons)
			}
		}
	} else {
		startAttack(gunner, bullets, pos, weaponSlot, mp, fireWeapons)
	}
}

func startAttack(gunner Gunner, newBullets []*bullet.Bullet, pos int, weaponSlot *body.WeaponSlot, mp *_map.Map, fireWeapons *[]FireMessage) {

	timeTick := _const.ServerTick - 1

	for timeTick > 0 && len(newBullets) > pos {

		InitBullet(newBullets[pos], gunner, false, weaponSlot)

		// для отыгрыша анимации выстрела
		*fireWeapons = append(*fireWeapons, FireMessage{
			TypeID: weaponSlot.Weapon.ID,
			X:      newBullets[pos].GetX(),
			Y:      newBullets[pos].GetY(),
			Z:      int(newBullets[pos].GetZ()),
			Rotate: int(newBullets[pos].GetRotate()),
			MapID:  mp.Id,
		})

		// лазеры особенные :)
		if weaponSlot.Weapon.Type == _const.LaserWeapon {
			bullets.Bullets.AddBullet(newBullets[pos])
		} else {
			fly_bullets.FlyBullet(newBullets[pos], mp, false)
		}

		timeTick -= weaponSlot.Weapon.DelayFollowingFire
		pos++
	}

	if len(newBullets) == pos {
		gunner.GetBurstOfShots().AddBullets(weaponSlot.Number, nil)
		weaponSlot.StartReload(weaponSlot.Weapon.ReloadTime)
	} else {
		gunner.GetBurstOfShots().ChangePos(weaponSlot.Number, pos)
	}
}
