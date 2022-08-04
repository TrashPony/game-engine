package game_types

import (
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/coordinate"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/weapon"
	"math/rand"
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
	8: {
		ID:                      8,
		Name:                    "explores_weapon_1",
		MinAngle:                -25,
		MaxAngle:                25,
		Accuracy:                4,
		AmmoCapacity:            5,
		Type:                    _const.LaserWeapon,
		XAttach:                 44,
		YAttach:                 64,
		RotateSpeed:             120,
		CountFireBullet:         1,
		ReloadAmmoTime:          750,
		ReloadTime:              750,
		DelayFollowingFire:      200,
		MaxRange:                450,
		DefaultAmmoTypeID:       2,
		Unit:                    true,
		AccumulationFirePower:   true,
		AccumulationFull:        2000, // сколько времени надо до 100% зарядки
		AccumulationFullTimeOut: 2000, // сколько оружие может ждать после полной зарядки
		FirePositions: []*coordinate.Coordinate{
			{
				X: 61,
				Y: 64,
			},
		},
		DamageModifier: 1,
		PowerPoints:    200,
	},
	9: {
		ID:                 9,
		Name:               "reverses_weapon_1",
		MinAngle:           -25,
		MaxAngle:           25,
		Accuracy:           13,
		AmmoCapacity:       6,
		Type:               _const.MissileWeapon,
		XAttach:            18,
		YAttach:            64,
		RotateSpeed:        180,
		CountFireBullet:    1,
		ReloadAmmoTime:     128,
		ReloadTime:         128,
		DelayFollowingFire: 0,
		MaxRange:           500,
		DefaultAmmoTypeID:  6,
		Unit:               true,
		FirePositions: []*coordinate.Coordinate{
			{
				X: 42,
				Y: 53,
			}, {
				X: 42,
				Y: 75,
			},
		},
		DamageModifier: 0.75,
		PowerPoints:    150,
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

func GetRandomWeapon() *weapon.Weapon {
	unitWeapon := make([]int, 0)

	for id, w := range weaponTypes {
		if w.Unit {
			unitWeapon = append(unitWeapon, id)
		}
	}

	return GetNewWeapon(unitWeapon[rand.Intn(len(unitWeapon))])
}

func GetAllWeapons() map[int]*weapon.Weapon {
	return weaponTypes
}
