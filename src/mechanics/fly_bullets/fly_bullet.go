package fly_bullets

import (
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/mechanics/factories/bullets"
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/bullet"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/dynamic_map_object"
	_map "github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
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

func FlyBullet(bullet *bullet.Bullet, gameMap *_map.Map, noTime bool) (int, int) {

	if bullet == nil {
		return 0, 0
	}

	bullet.RealSpeed = float64(bullet.Speed) / (1000.0 / _const.ServerBulletTick) * math.Cos(bullet.StartRadian)
	if bullet.RealSpeed < 7 {
		bullet.RealSpeed = 7.0
	}

	bullet.RadRotate = game_math.DegToRadian(bullet.GetRotate())
	// пройденное растояние, что бы наприм ракеты не могли пролететь больше своей дальности
	bullet.DistanceTraveled = 0.0
	bullet.RealX, bullet.RealY = float64(bullet.GetX()), float64(bullet.GetY())

	bullet.XVelocity, bullet.YVelocity = game_math.SpeedAndAngleToVelocity(bullet.RealSpeed, bullet.RadRotate)

	if !noTime {
		bullets.Bullets.AddBullet(bullet)
		// первое появление пули из ствола
		// SendMessage(*createFlyBulletMsg("FlyBullet", bullet, gameMap.Id, 10, bullet.X, bullet.Y))
		return 0, 0
	} else {
		for {
			BulletFlyTick(bullet, gameMap, true, nil)
			if bullet.GetEnd() {
				return bullet.GetX(), bullet.GetY()
			}
		}
	}
}

func BulletFlyTick(bullet *bullet.Bullet, mp *_map.Map, noTime bool, units []*unit.Unit) (fly *FlyBulletMessage, explosion *ExplosionBulletMessage, damageObject []*DamageObject, crater *dynamic_map_object.Object) {

	if bullet.Ammo.Type == "missile" {
		// ракеты ведут себя не так как обычные пули
		// TODO пока не работает из за замыкающего импорта Missile(bullet, mp, mapObjects)
	}

	//deltaTime - время затрачено на проверку колизий, оно существенно поэтому надо учитывать
	var percent int
	var end bool

	percent, end, _, crater, damageObject, _ = DetailFlyBullet(bullet, &bullet.RealX, &bullet.RealY,
		bullet.RealX+bullet.XVelocity, bullet.RealY+bullet.YVelocity, bullet.RadRotate, mp, &bullet.DistanceTraveled, noTime, false, units)

	ms := _const.ServerBulletTick
	if end {
		ms = (_const.ServerBulletTick * percent) / 100
	}

	//if debug.Store.FlyBulletLevel {
	//	x, y, lvl := mp.GetPosLevel(bullet.GetX(), bullet.GetY())
	//	CreateRect("orange", x*16, y*16, 16, mp.Id, 0, fmt.Sprintf("%.2f", lvl))
	//}

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
			MapID:  mp.Id,
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
				MapID:  mp.Id,
			}
		}
	}

	return
}
