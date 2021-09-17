package game_types

import (
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/coordinate"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/weapon"
)

// TODO тестовые данные
var WeaponTypes = map[int]weapon.Weapon{
	1: {
		ID:                 1,
		Name:               "replic_weapon_1",
		WeaponTexture:      "replic_weapon_1_skin_replics",
		MinAngle:           -25,
		MaxAngle:           25,
		Accuracy:           4,
		AmmoCapacity:       5,
		Type:               _const.FirearmsWeapon,
		XAttach:            14,
		YAttach:            64,
		RotateSpeed:        45,
		CountFireBullet:    1,
		BulletSpeed:        -50,
		ReloadAmmoTime:     3000,
		ReloadTime:         1000,
		DelayFollowingFire: 250,
		FirePositions: []*coordinate.Coordinate{
			{
				X: 61,
				Y: 64,
			},
		},
		DamageModifier: 1,
	},
	2: {
		ID:                 2,
		Name:               "explores_weapon_1",
		WeaponTexture:      "explores_weapon_1_skin_replics",
		MinAngle:           -25,
		MaxAngle:           25,
		Accuracy:           4,
		AmmoCapacity:       5,
		Type:               _const.LaserWeapon,
		XAttach:            44,
		YAttach:            64,
		RotateSpeed:        90,
		CountFireBullet:    1,
		ReloadAmmoTime:     3000,
		ReloadTime:         1000,
		DelayFollowingFire: 200,
		MaxRange:           350,
		FirePositions: []*coordinate.Coordinate{
			{
				X: 61,
				Y: 64,
			},
		},
		DamageModifier: 1,
	},
	3: {
		ID:                 3,
		Name:               "reverses_weapon_1",
		WeaponTexture:      "reverses_weapon_1_skin_replics",
		MinAngle:           -25,
		MaxAngle:           25,
		Accuracy:           3,
		AmmoCapacity:       8,
		Type:               _const.MissileWeapon,
		XAttach:            18,
		YAttach:            64,
		RotateSpeed:        90,
		CountFireBullet:    1,
		ReloadAmmoTime:     180,
		ReloadTime:         1000,
		DelayFollowingFire: 200,
		MaxRange:           400,
		FirePositions: []*coordinate.Coordinate{
			{
				X: 42,
				Y: 53,
			}, {
				X: 42,
				Y: 75,
			},
		},
		DamageModifier: 1,
	},
}
