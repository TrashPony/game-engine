package game_types

import (
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/coordinate"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/weapon"
)

var weaponTypes = map[int]*weapon.Weapon{
	2: {
		ID:                 2,
		Name:               "replic_weapon_2",
		MinAngle:           -25,
		MaxAngle:           10,
		Accuracy:           15,
		AmmoCapacity:       60,
		Type:               _const.FirearmsWeapon,
		XAttach:            25,
		YAttach:            64,
		RotateSpeed:        120,
		CountFireBullet:    2,
		ReloadAmmoTime:     64,
		ReloadTime:         64,
		DelayFollowingFire: 32,
		DefaultAmmoTypeID:  1,
		Unit:               true,
		FirePositions: []*coordinate.Coordinate{
			{
				X: 63,
				Y: 52,
			}, {
				X: 63,
				Y: 75,
			},
		},
		DamageModifier: 0.75,
		PowerPoints:    100,
	},
	100: {
		ID:                 100,
		Name:               "laser_turret_weapon",
		MinAngle:           -25,
		MaxAngle:           10,
		Accuracy:           15,
		AmmoCapacity:       60,
		Type:               _const.FirearmsWeapon,
		XAttach:            68,
		YAttach:            64,
		RotateSpeed:        90,
		CountFireBullet:    1,
		ReloadAmmoTime:     64,
		ReloadTime:         64,
		DelayFollowingFire: 32,
		DefaultAmmoTypeID:  1,
		FirePositions: []*coordinate.Coordinate{
			{
				X: 91,
				Y: 36,
			},
		},
		DamageModifier: 1,
	},
}

func GetNewWeapon(id int) *weapon.Weapon {
	sWeapon, ok := weaponTypes[id]
	if !ok {
		return nil
	}
	newWeapon := *sWeapon
	newWeapon.FirePositions = make([]*coordinate.Coordinate, 0)

	for _, c := range weaponTypes[id].FirePositions {
		newWeapon.FirePositions = append(newWeapon.FirePositions, c)
	}

	if newWeapon.Unit {
		newWeapon.WeaponTexture = newWeapon.Name
	}

	return &newWeapon
}

func GetAllWeapons() map[int]*weapon.Weapon {
	return weaponTypes
}
