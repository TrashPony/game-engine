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
	2: {
		ID:          2,
		Name:        "medium_lens",
		Type:        _const.LaserWeapon,
		MinDamage:   50,
		MaxDamage:   300,
		BulletSpeed: 99999,
	},
	6: {
		ID:                    6,
		Name:                  "aim_small_missile_bullet",
		Type:                  _const.MissileWeapon,
		ChaseTarget:           true,
		ChaseOption:           "distance_chase",
		ChaseCatchDestination: 150,
		Rotate:                5.0,
		AreaCovers:            25,
		MinDamage:             57,
		MaxDamage:             90,
		BulletSpeed:           450,
		PushingPower:          200,
	},
}

func GetNewAmmo(id int) *ammo.Ammo {
	newAmmo := ammoTypes[id]
	return &newAmmo
}

func GetAllAmmo() map[int]ammo.Ammo {
	return ammoTypes
}
