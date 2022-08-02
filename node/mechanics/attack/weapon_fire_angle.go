package attack

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/bullet"
	"math"
)

func GetBulletFireAngle(bullet *bullet.Bullet, artillery, force bool) {

	if bullet.Ammo.Type == "missile" && artillery && !force {
		bullet.StartRadian = game_math.DegToRadian(75)
		return
	}

	if bullet.Ammo.Type == "missile" && !force {
		bullet.StartRadian = game_math.DegToRadian(0)
		return
	}

	bullet.StartRadian = game_math.GetReachAngle(bullet.GetX(), bullet.GetY(), bullet.Target.GetX(), bullet.Target.GetY(), bullet.Target.GetZ(),
		bullet.StartZ, float64(bullet.Speed), artillery, bullet.Ammo.Type, bullet.Ammo.Gravity)

	// недостижимый угол, поэтому стреляем под 45
	if math.IsNaN(bullet.StartRadian) {
		bullet.StartRadian = game_math.DegToRadian(float64(game_math.GetRangeRand(43, 47)))
	}

	if bullet.Ammo.Type != "missile" {
		if artillery {
			if bullet.StartRadian < bullet.MaxAngle {
				bullet.StartRadian = bullet.MaxAngle
			}

			if bullet.StartRadian > bullet.MinAngle {
				bullet.StartRadian = bullet.MinAngle
			}
		} else {
			if bullet.StartRadian > bullet.MaxAngle {
				bullet.StartRadian = bullet.MaxAngle
			}

			if bullet.StartRadian < bullet.MinAngle {
				bullet.StartRadian = bullet.MinAngle
			}
		}
	}
}
