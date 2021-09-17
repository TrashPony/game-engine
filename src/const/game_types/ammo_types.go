package game_types

import (
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/ammo"
)

// TODO тестовые данные
var AmmoTypes = map[int]ammo.Ammo{
	1: {
		ID:          1,
		Name:        "piu-piu_2",
		Type:        _const.FirearmsWeapon,
		MinDamage:   1,
		MaxDamage:   2,
		BulletSpeed: 700,
	},
	2: {
		ID:          2,
		Name:        "medium_lens",
		Type:        _const.LaserWeapon,
		MinDamage:   1,
		MaxDamage:   2,
		BulletSpeed: 99999,
	},
	3: {
		ID:          3,
		Name:        "ballistics_artillery_bullet",
		Type:        _const.FirearmsWeapon,
		MinDamage:   1,
		MaxDamage:   2,
		BulletSpeed: 700,
	},
	4: {
		ID:          4,
		Name:        "small_lens",
		Type:        _const.LaserWeapon,
		MinDamage:   1,
		MaxDamage:   2,
		BulletSpeed: 99999,
	},
	5: {
		ID:          5,
		Name:        "small_missile_bullet",
		Type:        _const.MissileWeapon,
		ChaseTarget: true,
		ChaseOption: "distance_chase",
		MinDamage:   1,
		MaxDamage:   2,
		Rotate:      12.0,
		BulletSpeed: 450,
	},
	6: {
		ID:          6,
		Name:        "small_missile_bullet",
		Type:        _const.MissileWeapon,
		ChaseTarget: true,
		ChaseOption: "always_to_target",
		MinDamage:   1,
		MaxDamage:   2,
		Rotate:      3.0,
		BulletSpeed: 450,
	},
}
