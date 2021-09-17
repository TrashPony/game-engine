package fly_bullets

import (
	"github.com/TrashPony/game_engine/src/mechanics/factories/dynamic_objects"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/bullet"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/dynamic_map_object"
	"math/rand"
)

func GetCraterByBullet(bullet *bullet.Bullet, x, y int) *dynamic_map_object.Object {

	var crater *dynamic_map_object.Object

	if bullet.Ammo.Name == "small_missile_bullet" || bullet.Ammo.Name == "aim_small_missile_bullet" {

		crater = dynamic_objects.DynamicObjects.GetDynamicObjectByTexture("small_crater_1", bullet.GetRotate())
		crater.SetScale(rand.Intn(6-4) + 4) // подобрано эксперементальным путем
	}

	if bullet.Ammo.Name == "ballistics_artillery_bullet" {
		crater = dynamic_objects.DynamicObjects.GetDynamicObjectByTexture("mid_crater_1", bullet.GetRotate())
		crater.SetScale(rand.Intn(15-5) + 5) // подобрано эксперементальным путем
	}

	if bullet.Ammo.Name == "small_lens" || bullet.Ammo.Name == "medium_lens" || bullet.Ammo.Name == "piu-piu" {
		crater = dynamic_objects.DynamicObjects.GetDynamicObjectByTexture("tilted_mid_crater_1", bullet.GetRotate())
		crater.SetScale(rand.Intn(10-5) + 5) // подобрано эксперементальным путем
	}

	if bullet.Ammo.Name == "piu-piu_2" {
		crater = dynamic_objects.DynamicObjects.GetDynamicObjectByTexture("tilted_mid_crater_1", bullet.GetRotate())
		crater.SetScale(rand.Intn(5-2) + 2) // подобрано эксперементальным путем
	}

	if crater != nil {
		crater.GetPhysicalModel().SetPos(float64(x), float64(y), bullet.GetRotate())
	}

	return crater
}
