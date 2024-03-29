package fly_bullets

import (
	"github.com/TrashPony/game-engine/node/mechanics/damage"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/bullet"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
)

type FlyLaserMessage struct {
	TypeID int `json:"type_id"`
	X      int `json:"x"`
	Y      int `json:"y"`
	ToX    int `json:"to_x"`
	ToY    int `json:"to_y"`
	MapID  int `json:"m"`
}

func FlyLaser(bullet *bullet.Bullet, gameMap *_map.Map, noTime bool) (int, int, *FlyLaserMessage, []*damage.Object) {

	// лазер летит со скоростью света, поэтому что все нам надо это отдать стартовое ХУ и конечную ХУ
	// конечная ХУ это координата колизии или карты куда стреляет игрок

	radRotate := game_math.DegToRadian(bullet.GetRotate())
	startX, startY := bullet.GetX(), bullet.GetY()

	toX, toY := game_math.VectorToAngleBySpeed(float64(startX), float64(startY), float64(bullet.MaxRange), bullet.GetRotate())

	// пройденное растояние, что бы лазер не пролетал больше своей дальности
	distanceTraveled := 0.0

	realX, realY := float64(bullet.GetX()), float64(bullet.GetY())
	_, _, _, crater, damageObj, _ := DetailFlyBullet(bullet, &realX, &realY, float64(toX), float64(toY), radRotate,
		gameMap, &distanceTraveled, noTime, true, nil)

	bullet.SetX(int(realX))
	bullet.SetY(int(realY))

	bullet.Target.SetX(int(realX)) // на фронте рисуется линия до Target
	bullet.Target.SetY(int(realY))

	var msg *FlyLaserMessage
	if !noTime {

		msg = &FlyLaserMessage{
			X:      startX,
			Y:      startY,
			ToX:    bullet.Target.GetX(),
			ToY:    bullet.Target.GetY(),
			TypeID: bullet.Ammo.ID,
			MapID:  gameMap.Id,
		}

		gameMap.AddCrater(crater)
		bullet.SetEnd(true)
	}

	return bullet.Target.GetX(), bullet.Target.GetY(), msg, damageObj
}
