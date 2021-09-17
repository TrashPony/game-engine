package fly_bullets

import (
	units2 "github.com/TrashPony/game_engine/src/mechanics/factories/units"
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_math/collisions"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/bullet"
	_map "github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
)

type DamageObject struct {
	TypeTarget string      `json:"type_target"`
	IdTarget   int         `json:"id_target"`
	Damage     int         `json:"damage"`
	TypeDamage string      `json:"type_damage"`
	Dead       bool        `json:"dead"`
	Obj        interface{} `json:"obj"`
	X          int         `json:"x"`
	Y          int         `json:"y"`
}

func CollisionDamage(typeTarget string, idTarget int, bullet *bullet.Bullet, mp *_map.Map) []*DamageObject {
	if bullet.Ammo.AreaCovers > 0 {
		return Explosion(bullet, mp, typeTarget, idTarget)
	} else {
		return []*DamageObject{{TypeTarget: typeTarget, IdTarget: idTarget, Damage: bullet.Damage}}
	}
}

func Explosion(bullet *bullet.Bullet, mp *_map.Map, typeTarget string, idTarget int) []*DamageObject {

	// если произошла колизия с конкретным щитом,
	// то все обьекты под этим щитом будут не тронуты и урон по щиту происходит 1 раз

	damageObjects := make([]*DamageObject, 0)

	damage := bullet.Damage

	// инициатор взрыва если он есть получает 100% урон
	if typeTarget != "" {
		damageObjects = append(damageObjects, &DamageObject{TypeTarget: typeTarget, IdTarget: idTarget,
			Damage: damage})
	}

	// тип снаряда имеет зону поражения то дамаг по всему что вокруг
	explosionDamage := func(x, y, radius, areaCovers int, zTarget float64) int {
		// todo это работает не правильно
		// чем ближе к эпицентру тем более полный урон

		// по Z мы просто проверяем что обьект находится на уровне, хотя это немного не правильно Х)
		// areaCovers/2 потому что высота за 1Z это как 2X или 2Y
		distZ := bullet.GetZ() - zTarget
		if int(distZ) > areaCovers/2 {
			return 0
		}

		distXY := game_math.GetBetweenDist(bullet.GetX(), bullet.GetY(), x, y)
		distXY -= float64(radius)

		if int(distXY) < areaCovers {
			percentRange := (distXY * 100) / float64(bullet.Ammo.AreaCovers)
			return damage/2 + int(float64(damage/2)*(percentRange/100))

		} else {

			// до обьекта снаряд не дастрелил, а получил он в кусочек геодаты поэтому получает половину урона
			return damage / 2
		}
	}

	units, mxUnits := units2.Units.GetAllUnitsByMapIDUnsafeRange(mp.Id)
	for _, gameUnit := range units {

		if typeTarget == "unit" && gameUnit.GetID() == idTarget {
			continue
		}

		if collisions.CircleUnit(bullet.GetX(), bullet.GetY(), bullet.Ammo.AreaCovers, gameUnit) {

			_, _, lvl := mp.GetPosLevel(gameUnit.GetX(), gameUnit.GetY())
			lvl += gameUnit.GetPhysicalModel().GetHeight()

			damageObjects = append(damageObjects, &DamageObject{TypeTarget: "unit", IdTarget: gameUnit.GetID(),
				Damage: explosionDamage(gameUnit.GetX(), gameUnit.GetY(), 0, bullet.Ammo.AreaCovers, lvl)})
		}
	}
	mxUnits.RUnlock()

	return damageObjects
}
