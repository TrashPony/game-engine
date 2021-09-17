package attack

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/bullet"
	"math"
)

func GetBulletFireAngle(bullet *bullet.Bullet, artillery, force bool) {

	if bullet.Ammo.Type == "missile" && artillery && !force {
		// артилерийские ракеты сначало летят на орибату, потому к цели, после падают на цель
		bullet.StartRadian = game_math.DegToRadian(75)
		return
	}

	if bullet.Ammo.Type == "missile" && !force {
		bullet.StartRadian = game_math.DegToRadian(0)
		return
	}

	bullet.StartRadian = game_math.GetReachAngle(bullet.GetX(), bullet.GetY(), bullet.Target.GetX(), bullet.Target.GetY(), bullet.Target.GetZ(),
		bullet.StartZ, float64(bullet.Speed), artillery, bullet.Ammo.Type)

	// недостижимый угол, поэтому стреляем под 45
	if math.IsNaN(bullet.StartRadian) {
		bullet.StartRadian = game_math.DegToRadian(45)
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
