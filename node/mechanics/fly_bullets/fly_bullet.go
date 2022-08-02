package fly_bullets

import (
	"github.com/TrashPony/game-engine/node/mechanics/damage"
	"github.com/TrashPony/game-engine/node/mechanics/factories/bullets"
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	battle2 "github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/bullet"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/dynamic_map_object"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"math"
)

type FlyBulletMessage struct {
	ID     int `json:"id"`
	TypeID int `json:"type_id"`
	X      int `json:"x"`
	Y      int `json:"y"`
	Z      int `json:"z"`
	MS     int `json:"ms"`
	Rotate int `json:"r"`
	MapID  int `json:"m"`
}

type ExplosionBulletMessage struct {
	TypeID int `json:"type_id"`
	X      int `json:"x"`
	Y      int `json:"y"`
	Z      int `json:"z"`
	MapID  int `json:"m"`
}

func FlyBullet(bullet *bullet.Bullet, b *battle2.Battle, noTime bool, units []*unit.Unit) (int, int, []*damage.Object) {

	if bullet == nil {
		return 0, 0, nil
	}

	bullet.RealSpeed = float64(bullet.Speed) / (1000.0 / _const.ServerBulletTick) * math.Cos(bullet.StartRadian)
	if bullet.RealSpeed < 7 {
		bullet.RealSpeed = 7.0
	}

	bullet.RadRotate = game_math.DegToRadian(bullet.GetRotate())
	// пройденное растояние
	bullet.DistanceTraveled = 0.0
	bullet.RealX, bullet.RealY = float64(bullet.GetX()), float64(bullet.GetY())

	bullet.XVelocity, bullet.YVelocity = game_math.SpeedAndAngleToVelocity(bullet.RealSpeed, bullet.RadRotate)

	if !noTime {
		bullets.Bullets.AddBullet(bullet)
		return 0, 0, nil
	} else {
		for {
			_, _, do, _ := BulletFlyTick(bullet, b, true, units)
			if bullet.GetEnd() {
				return bullet.GetX(), bullet.GetY(), do
			}
		}
	}
}

func BulletFlyTick(bullet *bullet.Bullet, b *battle2.Battle, noTime bool, units []*unit.Unit) (fly *FlyBulletMessage, explosion *ExplosionBulletMessage, damageObject []*damage.Object, crater *dynamic_map_object.Object) {

	// коректировка ракет
	if bullet.Ammo.Type == _const.MissileWeapon {
		Missile(bullet, b.Map, noTime, units)
	}

	var percent int
	var end bool

	percent, end, _, crater, damageObject, _ = DetailFlyBullet(bullet, &bullet.RealX, &bullet.RealY,
		bullet.RealX+bullet.XVelocity, bullet.RealY+bullet.YVelocity, bullet.RadRotate, b.Map, &bullet.DistanceTraveled, noTime, false, units)

	ms := _const.ServerBulletTick
	if end {
		ms = (_const.ServerBulletTick * percent) / 100
	}

	bullet.SetX(int(bullet.RealX))
	bullet.SetY(int(bullet.RealY))

	if !noTime {
		fly = &FlyBulletMessage{
			ID:     bullet.ID,
			TypeID: bullet.Ammo.ID,
			X:      bullet.GetX(),
			Y:      bullet.GetY(),
			Z:      int(bullet.GetZ()),
			MS:     ms,
			Rotate: int(bullet.GetRotate()),
			MapID:  b.Map.Id,
		}
	}

	if end {
		bullet.SetEnd(end)

		if !noTime {
			explosion = &ExplosionBulletMessage{
				TypeID: bullet.Ammo.ID,
				X:      bullet.GetX(),
				Y:      bullet.GetY(),
				Z:      int(bullet.GetZ()),
				MapID:  b.Map.Id,
			}
		}
	}

	return
}
