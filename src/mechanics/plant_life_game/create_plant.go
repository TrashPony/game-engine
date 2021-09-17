package plant_life_game

import (
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/mechanics/factories/dynamic_objects"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"math/rand"
)

func createPlant(x, y int, texturePlant string, mp *_map.Map) {
	plant := dynamic_objects.DynamicObjects.GetDynamicObjectByTexture(texturePlant, float64(rand.Intn(360)))
	plant.GetPhysicalModel().SetPos(float64(x), float64(y), float64(rand.Intn(360)))
	plant.SetScale(1)
	plant.GrowCycle = 100
	plant.GrowTime = _const.TerrainGrowTime
	plant.GrowLeftTime = _const.TerrainGrowTime

	minSize := 2
	// что бы шанс получить большого был не велик
	if rand.Intn(4) == 1 {
		plant.MaxScale = rand.Intn(50-minSize) + minSize
	} else {
		plant.MaxScale = rand.Intn(20-minSize) + minSize
	}

	plant.SetHP(plant.MaxHP)
	plant.CalculateScale()

	// todo	if !collisions.CheckObjectCollision(plant, mp, true) {
	mp.AddDynamicObject(plant)
	//	}
}
