package damage

import (
	collisions2 "github.com/TrashPony/game-engine/node/mechanics/collisions"
	units2 "github.com/TrashPony/game-engine/node/mechanics/factories/units"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/physical_model"
)

func Explosion(mp *_map.Map, typeTarget string, idTarget, damage, areaCovers, x, y int, z float64, startX, startY, pushPower int) []*Object {

	// если произошла колизия с конкретным щитом,
	// то все обьекты под этим щитом будут не тронуты и урон по щиту происходит 1 раз

	damageObjects := make([]*Object, 0)

	// инициатор взрыва если он есть получает 100% урон
	if typeTarget != "" {
		damageObjects = append(damageObjects, &Object{TypeTarget: typeTarget, IdTarget: idTarget, Damage: damage, PushPower: pushPower})
	}

	// тип снаряда имеет зону поражения то дамаг по всему что вокруг
	explosionDamage := func(targetX, targetY, radius, areaCovers int, zTarget float64) (int, int) {
		// чем ближе к эпицентру тем более полный урон

		// по Z мы просто проверяем что обьект находится на уровне, хотя это немного не правильно Х)
		// areaCovers/2 потому что высота за 1Z это как 2X или 2Y
		distZ := z - zTarget
		if int(distZ) > areaCovers/2 {
			return 0, 0
		}

		distXY := game_math.GetBetweenDist(targetX, targetY, x, y)
		distXY -= float64(radius)

		if int(distXY) < areaCovers {
			percentRange := 100 - (distXY*100)/float64(areaCovers)
			push := int(float64(pushPower) * (percentRange / 100))
			d := int(float64(damage) * (percentRange / 100))
			return d, push
		} else {
			// до обьекта снаряд не дастрелил, а получил он в кусочек геодаты поэтому получает половину урона
			return damage / 2, 0
		}
	}

	objects, objectsMX := mp.UnsafeRangeDynamicObjects()
	for _, obj := range objects {

		if typeTarget == "object" && obj.ID == idTarget { // потому что он получает 100% урона по умолчанию
			continue
		}

		if collisions2.CircleDynamicObj(x, y, areaCovers, obj, false) {

			_, _, lvl := mp.GetPosLevel(obj.GetX(), obj.GetY())
			lvl += obj.GetPhysicalModel().Height

			d, _ := explosionDamage(obj.GetX(), obj.GetY(), 0, areaCovers, lvl)

			damageObjects = append(damageObjects, &Object{TypeTarget: "object", IdTarget: obj.GetID(), Damage: d})
		}
	}
	objectsMX.RUnlock()

	units, mxUnits := units2.Units.GetAllUnitsByMapIDUnsafeRange(mp.Id)
	for _, gameUnit := range units {

		if typeTarget == "unit" && gameUnit.GetID() == idTarget {
			continue
		}

		if gameUnit.GetPhysicalModel().Fly {
			continue
		}

		if collisions2.CircleUnit(x, y, areaCovers, gameUnit) {

			_, _, lvl := mp.GetPosLevel(gameUnit.GetX(), gameUnit.GetY())
			lvl += gameUnit.GetPhysicalModel().GetHeight()

			d, p := explosionDamage(gameUnit.GetX(), gameUnit.GetY(), 0, areaCovers, lvl)
			damageObjects = append(damageObjects, &Object{TypeTarget: "unit", IdTarget: gameUnit.GetID(), Damage: d, PushPower: p})
		}
	}
	mxUnits.RUnlock()

	return damageObjects
}

func pushUnit(pushObject *physical_model.PhysicalModel, x, y, p int) {
	if p > 0 {
		radRotate := game_math.DegToRadian(game_math.GetBetweenAngle(
			float64(pushObject.GetX()),
			float64(pushObject.GetY()),
			float64(x),
			float64(y),
		))

		if p > 2000 && pushObject.GetChassisType() == "antigravity" {
			p = 2000 // а то они ваще улетают
		}

		p = int(float64(p) * (2000 / pushObject.GetWeight()))

		pushObject.AddVelocity(
			game_math.Cos(radRotate)*(float64(p)/100),
			game_math.Sin(radRotate)*(float64(p)/100),
		)
	}
}
