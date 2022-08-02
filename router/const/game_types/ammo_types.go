package game_types

import (
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/ammo"
)

var ammoTypes = map[int]ammo.Ammo{
	1: {
		ID:           1,
		Name:         "piu-piu_2",
		Type:         _const.FirearmsWeapon,
		MinDamage:    9,
		MaxDamage:    18,
		BulletSpeed:  850,
		PushingPower: 25,
	},
}

func GetNewAmmo(id int) *ammo.Ammo {
	newAmmo := ammoTypes[id]
	return &newAmmo
}

func GetAllAmmo() map[int]ammo.Ammo {
	return ammoTypes
}
