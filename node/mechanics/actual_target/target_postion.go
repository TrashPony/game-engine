package actual_target

import (
	"github.com/TrashPony/game-engine/node/mechanics/factories/units"
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/body"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/physical_model"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/position"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/target"
	"math"
)

type Gunner interface {
	GetX() int
	GetY() int
	GetFirePos(int) *position.Positions
	GetMapHeight() float64
}

func GetXYZTarget(gunner Gunner, target *target.Target, mp *_map.Map, weaponSlot *body.WeaponSlot) bool {

	if target == nil {
		return false
	}

	target.SetZ(0)
	defer func() {
		// высота цели
		_, _, lvl := mp.GetPosLevel(target.GetX(), target.GetY())
		if target.GetZ() > 0 {
			target.SetZ(target.GetZ()/2 + lvl)
		} else {
			target.SetZ(lvl)
		}
	}()

	if target.Type == "map" {
		return true
	}

	if target.Type == "object" {
		obj := mp.GetDynamicObjectsByID(target.ID)
		if obj == nil || obj.MapID != mp.Id {
			// по той или иной причине юнит перестал видит цель и больше не знает существует оно или нет
			return false
		} else {
			// цель может двигатся значит надо делать поправку на время полета пули до цели
			x, y := GetOffsetSpeedTarget(gunner, obj.GetPhysicalModel(), mp, weaponSlot)

			target.SetX(x)
			target.SetY(y)
			target.SetZ(obj.GetPhysicalModel().GetHeight())
		}
	}

	if target.Type == "unit" {

		targetUnit := units.Units.GetUnitByIDAndMapID(target.ID, mp.Id)
		if targetUnit == nil || targetUnit.MapID != mp.Id {
			// по той или иной причине юнит перестал видит цель и больше не знает существует оно или нет
			return false
		} else {
			// цель может двигатся значит надо делать поправку на время полета пули до цели
			x, y := GetOffsetSpeedTarget(gunner, targetUnit.GetPhysicalModel(), mp, weaponSlot)

			target.SetX(x)
			target.SetY(y)
			target.SetZ(targetUnit.GetPhysicalModel().GetHeight())
		}
	}

	return true
}

func GetOffsetSpeedTarget(gunner Gunner, target *physical_model.PhysicalModel, mp *_map.Map, weaponSlot *body.WeaponSlot) (int, int) {

	xVelocity, yVelocity := target.GetVelocity()
	if gunner == nil || weaponSlot == nil || int(xVelocity) == 0 && int(yVelocity) == 0 {
		// юнит стоит на месте
		return target.GetX(), target.GetY()
	}

	_, _, lvl := mp.GetPosLevel(target.GetX(), target.GetY())
	lvl = lvl + (target.GetHeight() / 2)

	//посчитать время полета пули (расстояние до цели / скорость пули в сек = время полета в секундах)
	bulletSpeed := GetBulletSpeedToTargetByWeaponSlot(target.GetX(), target.GetY(), lvl, mp, gunner, weaponSlot)
	if math.IsNaN(bulletSpeed) {
		// орудие не достреливает до юнита
		return target.GetX(), target.GetY()
	}

	distToTarget := game_math.GetBetweenDist(target.GetX(), target.GetY(), gunner.GetX(), gunner.GetY())

	// посчитать сколько за это время проедет цель ( скорость цели * на время полета пули = расстояние которое преодолеет цель)
	targetSpeed := target.GetCurrentSpeed()
	timeToTarget := distToTarget / bulletSpeed //(sec)
	distPathUnit := timeToTarget * (targetSpeed * (1000 / _const.ServerTick))

	// берём вектор направления цели и прокладываем путь на расстояние которое цель пройдет за время полета пули.
	// результат это предположительно место где окажется юнит когда долетит пуля
	x, y := game_math.VectorToAngleBySpeed(float64(target.GetX()), float64(target.GetY()), distPathUnit, game_math.RadianToDeg(target.GetVelocityRotate()))
	return x, y
}

func GetBulletSpeedToTargetByWeaponSlot(xTarget, yTarget int, zTarget float64, mp *_map.Map, gunner Gunner, weaponSlot *body.WeaponSlot) float64 {
	if weaponSlot == nil || weaponSlot.Weapon == nil || weaponSlot.GetAmmo() == nil {
		return 0
	}

	weaponPos := gunner.GetFirePos(weaponSlot.Number)
	_, _, lvl := mp.GetPosLevel(gunner.GetX(), gunner.GetY())
	lvl = lvl + gunner.GetMapHeight()

	bulletSpeed := float64(weaponSlot.GetAmmo().BulletSpeed + weaponSlot.Weapon.BulletSpeed)
	startRadian := game_math.GetReachAngle(weaponPos.X, weaponPos.Y, xTarget, yTarget, zTarget, lvl, bulletSpeed,
		weaponSlot.Weapon.Artillery, weaponSlot.GetAmmo().Type, weaponSlot.GetAmmo().Gravity)

	// время в секундах
	realSpeed := bulletSpeed / math.Cos(startRadian)

	return realSpeed
}
