package attack

import (
	"github.com/TrashPony/game-engine/node/mechanics/actual_target"
	collisions2 "github.com/TrashPony/game-engine/node/mechanics/collisions"
	"github.com/TrashPony/game-engine/node/mechanics/damage"
	"github.com/TrashPony/game-engine/node/mechanics/factories/units"
	"github.com/TrashPony/game-engine/node/mechanics/fly_bullets"
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	battle2 "github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/body"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/target"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
)

type WeaponTargetMsg struct {
	X                   int    `json:"x"`
	Y                   int    `json:"y"`
	Accuracy            int    `json:"a"`
	MapID               int    `json:"m"`
	AmmoCount           int    `json:"ac"`
	AmmoAvailable       int    `json:"aa"`
	Reload              bool   `json:"r"`
	AccumulationPercent int    `json:"ap"`
	Chase               bool   `json:"chase"`
	TargetType          string `json:"target_type"`
	TargetID            int    `json:"target_id"`
}

func CheckFire(gunner Gunner, gunnerType string, gunnerID int, b *battle2.Battle, wTarget *target.Target, weaponSlot *body.WeaponSlot, noRotate bool, unitsBasket []*unit.Unit) bool {

	if weaponSlot == nil || weaponSlot.Weapon == nil || !actual_target.GetXYZTarget(gunner, wTarget, b.Map, weaponSlot) {
		return false
	}

	_, _, targetLvl := b.Map.GetPosLevel(wTarget.GetX(), wTarget.GetY())
	maxRange, _ := gunner.GetWeaponMaxRange(targetLvl, weaponSlot.Number, false)
	if maxRange < int(game_math.GetBetweenDist(wTarget.GetX(), wTarget.GetY(), gunner.GetX(), gunner.GetY())) {
		// цель дальше чем может выстрелить оружие
		return false
	}

	gunner.SetWeaponTarget(wTarget)
	x, y, damageObjects := CollisionWeaponRangeCollision(gunner, b, weaponSlot, gunnerType, gunnerID, wTarget, noRotate, unitsBasket)
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

	for _, do := range damageObjects {
		if do.TypeTarget == wTarget.Type && do.IdTarget == wTarget.ID {
			return true
		}
	}

	if wTarget.Type == "unit" {
		targetUnit := units.Units.GetUnitByIDAndMapID(wTarget.ID, b.Map.Id)
		if targetUnit != nil {
			return collisions2.CircleUnit(x, y, radius, targetUnit)
		}
	}

	if wTarget.Type == "object" {
		obj := b.Map.GetDynamicObjectsByID(wTarget.ID)
		if obj != nil {
			return collisions2.CircleDynamicObj(x, y, radius, obj, false)
		}
	}

	return false
}

func CheckFireToTarget(target *target.Target, slot *body.WeaponSlot, xWeapon, yWeapon int, gunRotate float64, ignoreRotate bool) bool {

	if !ignoreRotate {

		if slot.GetReload() {
			// оружие перезаряжается
			return false
		}

		// смотрим что бы оружие было повернуто в необходимом положение
		needRotate := game_math.GetBetweenAngle(float64(target.GetX()), float64(target.GetY()), float64(xWeapon), float64(yWeapon))
		if needRotate < 0 {
			needRotate += 360
		}
		if gunRotate < 0 {
			gunRotate += 360
		}

		if !(gunRotate >= needRotate-2 && gunRotate <= needRotate+2) {
			return false
		}

	}

	return true
}

func WeaponTarget(gunner Gunner, weaponSlot *body.WeaponSlot, b *battle2.Battle, gunnerType string, gunnerID int, units []*unit.Unit) *WeaponTargetMsg {

	wt := gunner.GetWeaponTarget()
	if weaponSlot == nil || weaponSlot.Weapon == nil || weaponSlot.GetAmmo() == nil || wt == nil {
		return nil
	}

	var accuracy float64
	var targetType string
	var x, y, targetID int

	chase := chaseTarget(weaponSlot)
	if !chase {
		x, y, _ = CollisionWeaponRangeCollision(gunner, b, weaponSlot, gunnerType, gunnerID, wt, false, units)
		xWeapon, yWeapon := gunner.GetWeaponPosInMap(weaponSlot.Number)
		_, _, lvl := b.Map.GetPosLevel(xWeapon, yWeapon)
		startBulletLvl := lvl + gunner.GetMapHeight()
		maxRange, maxAngle := gunner.GetWeaponMaxRange(lvl, 1, false)

		accuracy = GetGunAccuracy(xWeapon, yWeapon, x, y, maxRange, float64(weaponSlot.GetAmmo().BulletSpeed+weaponSlot.Weapon.BulletSpeed),
			startBulletLvl, lvl, maxAngle, float64(gunner.GetWeaponAccuracy(weaponSlot.Number)), weaponSlot.Weapon.Type, weaponSlot.GetAmmo().Gravity)

		if accuracy <= 10 {
			accuracy = 10
		}
	} else {
		accuracy = float64(weaponSlot.Ammo.ChaseCatchDestination) * 2
		x, y = wt.GetX(), wt.GetY()
		passObj := getChaseTarget(gunner, weaponSlot, b.Map, wt)
		if passObj != nil {
			targetType = passObj.typeObj
			targetID = passObj.id
		}
	}

	accumulationPercent := 0
	if weaponSlot.Weapon.AccumulationFirePower {
		accumulationPercent = int((weaponSlot.AccumulationCurrent / weaponSlot.Weapon.AccumulationFull) * 100)
	}

	targetMsg := &WeaponTargetMsg{
		X:                   x,
		Y:                   y,
		Accuracy:            int(accuracy),
		MapID:               b.Map.Id,
		AmmoCount:           weaponSlot.AmmoQuantity,
		Reload:              weaponSlot.AmmoReload,
		AmmoAvailable:       weaponSlot.Weapon.AmmoCapacity,
		AccumulationPercent: accumulationPercent,
		Chase:               chase,
		TargetType:          targetType,
		TargetID:            targetID,
	}

	return targetMsg
}

type missileTarget struct {
	id       int
	realID   int
	typeObj  string
	realType string
	dist     int
}

func getChaseTarget(gunner Gunner, weaponSlot *body.WeaponSlot, mp *_map.Map, wt *target.Target) *missileTarget {
	var passObj *missileTarget

	x, y := wt.GetX(), wt.GetY()

	obj, mx := gunner.UnsafeRangeVisibleObjects()
	mx.RLock()
	defer mx.RUnlock()

	for _, vObj := range obj {
		typeObj := vObj.TypeObject
		id := vObj.IDObject
		if !vObj.View && vObj.Radar {
			typeObj = "mark"
			id = vObj.ID
		}

		if vObj.TypeObject == "unit" {
			u := units.Units.GetUnitByIDAndMapID(vObj.IDObject, mp.Id)
			if u != nil && u.TeamID != gunner.GetTeamID() {
				dist := int(game_math.GetBetweenDist(x, y, u.GetX(), u.GetY()))
				if (dist < weaponSlot.Ammo.ChaseCatchDestination) && (passObj == nil || passObj.dist > dist) {
					passObj = &missileTarget{id: id, realID: u.GetID(), typeObj: typeObj, realType: vObj.TypeObject, dist: dist}
				}
			}
		}

		if vObj.TypeObject == "object" {
			o := mp.GetDynamicObjectsByID(vObj.IDObject)
			if o != nil && o.TeamID != gunner.GetTeamID() {
				dist := int(game_math.GetBetweenDist(x, y, o.GetX(), o.GetY()))
				if (dist < weaponSlot.Ammo.ChaseCatchDestination) && (passObj == nil || passObj.dist > dist) {
					passObj = &missileTarget{id: id, realID: o.GetID(), typeObj: typeObj, realType: vObj.TypeObject, dist: dist}
				}
			}
		}
	}

	return passObj
}

func chaseTarget(weaponSlot *body.WeaponSlot) bool {
	if weaponSlot.Ammo.ChaseTarget && (weaponSlot.Ammo.ChaseOption == "always_to_target" || weaponSlot.Ammo.ChaseOption == "always_to_mouse_target") {
		return true
	}

	return false
}

// метод что оружие может стрелять в цель (между оружием и целью нет колизий)
func CollisionWeaponRangeCollision(gunner Gunner, b *battle2.Battle, weaponSlot *body.WeaponSlot, gunnerType string, gunnerID int, target *target.Target, noRotate bool, units []*unit.Unit) (int, int, []*damage.Object) {

	bullets, _ := CreateBullets(gunner, b.Map, gunner.GetWeaponSlot(weaponSlot.Number), gunnerType, gunnerID, false, false, target)
	if bullets == nil {
		return -1, -1, nil
	}

	bullet := bullets[0]
	InitBullet(b.Map, bullet, gunner, true, noRotate, weaponSlot)

	// лазеры особенные :)
	if weaponSlot.Weapon.Type == _const.LaserWeapon {
		x, y, _, damageObject := fly_bullets.FlyLaser(bullet, b.Map, true)
		return x, y, damageObject
	} else {
		return fly_bullets.FlyBullet(bullet, b, true, units)
	}
}
