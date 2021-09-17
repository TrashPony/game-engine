package fly_bullets

import (
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_math/collisions"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/bullet"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/dynamic_map_object"
	_map "github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
	"math"
)

func DetailFlyBullet(bullet *bullet.Bullet, realX, realY *float64, toX, toY, radRotate float64, gameMap *_map.Map,
	distanceTraveled *float64, noTime bool, laser bool, units []*unit.Unit) (int, bool, bool, *dynamic_map_object.Object, []*DamageObject, bool) {

	startDist := game_math.GetBetweenDist(bullet.GetX(), bullet.GetY(), int(toX), int(toY))
	minDist := startDist
	dist := startDist

	x, y := *realX, *realY

	xV, yV := bullet.XVelocity*0.1, bullet.YVelocity*0.1
	speed := float64(_const.AmmoRadius)

	for {

		newZ := game_math.GetZBulletByXPath(bullet.StartZ, bullet.StartRadian, float64(bullet.Speed), bullet.StartX, bullet.StartY, int(x), int(y), bullet.Ammo.Type)
		bullet.SetZ(newZ)

		if distanceTraveled != nil {
			if laser {
				*distanceTraveled += speed
			} else {
				*distanceTraveled += math.Sqrt(xV*xV + yV*yV)
			}

			if *distanceTraveled > 100 {
				// AttackUnitID, AttackStructID нужна что бы колизия не сработала сразу после запуска,
				// а так как пуля уже далеко и ее могут например вернуть гравипушкой то тут это
				bullet.IgnoreOwner = true
			}
		}

		percentPath := 100 - (int((dist * 100) / startDist))

		if laser {
			x += speed * game_math.Cos(radRotate) // идем по вектору движения пуди
			y += speed * game_math.Sin(radRotate)
		} else {
			x, y = x+xV, y+yV
		}

		dist = game_math.GetBetweenDist(int(x), int(y), int(toX), int(toY))

		var crater *dynamic_map_object.Object
		var damageObj []*DamageObject
		var end bool

		// пулю сбили в воздухе
		if bullet.HP <= 0 {
			return percentPath, true, false, nil, nil, true
		}

		// колизии с обьектами, юнитами и геодатой по курсу движения
		// проверяем колизии ток для прямо летящих снарядов или низко летящей арты

		// урон по ящикам и транспортам проходит только если стрелять по ним целенаправленно или задело взрывом
		collision, typeCollision, id := collisions.CircleAllCollisionCheck(int(x), int(y), _const.AmmoRadius, bullet.GetZ(), gameMap, 0, 0, true, units)
		if typeCollision == "flore" {
			if !noTime {
				crater = GetCraterByBullet(bullet, int(x), int(y))
			}
		}

		if typeCollision == bullet.OwnerType && id == bullet.OwnerID && !bullet.IgnoreOwner {
			collision = false
		}

		if collision && typeCollision != "" {
			if !noTime {
				// обработка урона если пуля врезалась в обьект
				damageObj = CollisionDamage(typeCollision, id, bullet, gameMap)
			}
		}

		end = checkEndPath(bullet, distanceTraveled)
		if end && !collision && bullet.Ammo.AreaCovers > 0 {
			// обработка урона если пуля долетела до цели и имеет зону поражения
			if !noTime {
				damageObj = Explosion(bullet, gameMap, "", 0)
			}
		}

		end = end || collision

		if dist <= 3 || minDist < dist || end {

			*realX, *realY = x, y

			if collision {
				bullet.Target.SetX(int(x))
				bullet.Target.SetY(int(y))
			}

			return percentPath, end, collision, crater, damageObj, false
		}

		minDist = dist
	}
}

func checkEndPath(bullet *bullet.Bullet, distanceTraveled *float64) bool {

	if distanceTraveled == nil {
		return false
	}

	if bullet.Ammo.Type == "missile" || bullet.Ammo.Type == "laser" {
		// конец пути ракетного снаряда кончается когда он достигает цель
		return *distanceTraveled >= float64(bullet.MaxRange)
	}

	return false
}
