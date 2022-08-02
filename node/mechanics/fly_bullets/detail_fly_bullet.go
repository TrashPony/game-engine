package fly_bullets

import (
	"github.com/TrashPony/game-engine/node/mechanics/collisions"
	"github.com/TrashPony/game-engine/node/mechanics/damage"
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/bullet"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/dynamic_map_object"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"math"
)

func DetailFlyBullet(b *bullet.Bullet, realX, realY *float64, toX, toY, radRotate float64, gameMap *_map.Map,
	distanceTraveled *float64, noTime bool, laser bool, units []*unit.Unit) (int, bool, bool, *dynamic_map_object.Object, []*damage.Object, bool) {

	startDist := game_math.GetBetweenDist(b.GetX(), b.GetY(), int(toX), int(toY))
	minDist := startDist
	dist := startDist

	x, y := *realX, *realY

	xV, yV := b.XVelocity*0.1, b.YVelocity*0.1
	speed := float64(_const.AmmoRadius)

	for {

		newZ := game_math.GetZBulletByXPath(b.StartZ, b.StartRadian, float64(b.Speed), b.StartX, b.StartY, int(x), int(y), b.Ammo.Type, b.Ammo.Gravity)
		b.SetZ(newZ)

		if distanceTraveled != nil {
			if laser {
				*distanceTraveled += speed
			} else {
				*distanceTraveled += math.Sqrt(xV*xV + yV*yV)
			}

			if *distanceTraveled > 100 {
				// AttackUnitID, AttackStructID нужна что бы колизия не сработала сразу после запуска,
				// а так как пуля уже далеко и ее могут например вернуть гравипушкой то тут это
				b.IgnoreOwner = true
			}
		}

		percentPath := 100 - (int((dist * 100) / startDist))

		if laser {
			x += speed * game_math.Cos(radRotate)
			y += speed * game_math.Sin(radRotate)
		} else {
			x, y = x+xV, y+yV
		}

		dist = game_math.GetBetweenDist(int(x), int(y), int(toX), int(toY))

		var crater *dynamic_map_object.Object
		var damageObj []*damage.Object
		var collision bool
		var end bool
		var typeCollision string
		var id int

		// пулю сбили в воздухе
		if b.HP <= 0 {
			return percentPath, true, false, nil, nil, true
		}

		if !collision {
			collision, typeCollision, id, _ = collisions.CircleAllCollisionCheck(int(x), int(y), _const.AmmoRadius, b.GetZ(), gameMap, nil, nil, true, units)
			if typeCollision == b.OwnerType && id == b.OwnerID && !b.IgnoreOwner {
				collision = false
			}

			if collision && typeCollision != "" {
				// обработка урона если пуля врезалась в обьект
				damageObj = damage.CollisionDamage(typeCollision, id, b.Damage, b.Ammo.AreaCovers, gameMap, b.GetX(), b.GetY(), b.GetZ(), b.StartX, b.StartY, b.Ammo.PushingPower)
			}
		}

		end = checkEndPath(b, distanceTraveled) || collision || b.ForceExplosion

		if end && !collision && b.Ammo.AreaCovers > 0 {
			// обработка урона, пуля взорвалась
			damageObj = damage.Explosion(gameMap, "", 0, b.Damage, b.Ammo.AreaCovers, b.GetX(), b.GetY(), b.GetZ(), b.StartX, b.StartY, b.Ammo.PushingPower)
		}

		if dist <= 3 || minDist < dist || end {

			*realX, *realY = x, y

			if collision {
				b.Target.SetX(int(x))
				b.Target.SetY(int(y))
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

	return bullet.MaxRange > 0 && *distanceTraveled >= float64(bullet.MaxRange)
}
