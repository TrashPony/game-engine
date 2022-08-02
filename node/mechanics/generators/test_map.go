package generators

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/dynamic_map_object"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/obstacle_point"
	"math/rand"
)

var stepX = 200
var stepY = 200

func getTestMap() *_map.Map {
	brMap := &_map.Map{
		XSize:         1024,
		YSize:         1024,
		DefaultLevel:  100,
		Flore:         map[int]map[int]*dynamic_map_object.Flore{},
		StaticObjects: map[int]*dynamic_map_object.Object{},
	}

	rowCount := 0
	terrains := []string{"grass_1", "grass_2", "grass_3"}
	for y := 0; y <= brMap.YSize; y += stepY {
		rowCount++
		for x := 0; x <= brMap.XSize; x += stepX {

			var flore dynamic_map_object.Flore

			flore.TextureOverFlore = terrains[rand.Intn(len(terrains))]
			flore.X = offsetX(x, rowCount)
			flore.Y = y
			flore.TexturePriority = brMap.GetMaxPriorityTexture() + 1

			brMap.AddFlore(&flore)
		}
	}

	obj := &dynamic_map_object.Object{
		Type:              "mountains",
		Texture:           "mountain_5",
		MaxHP:             -2,
		TypeMaxHP:         -2,
		HP:                -2,
		Scale:             40,
		Shadow:            true,
		TypeXShadowOffset: 5,
		XShadowOffset:     5,
		TypeYShadowOffset: 6,
		YShadowOffset:     6,
		ShadowIntensity:   40,
		Immortal:          true,
		TypeGeoData: []*obstacle_point.ObstaclePoint{
			{X: 5, Y: -37, Radius: 350, Move: false},
			{X: 147, Y: 125, Radius: 275, Move: false},
			{X: 395, Y: -23, Radius: 65, Move: false},
			{X: -114, Y: 292, Radius: 65, Move: false},
			{X: -137, Y: 300, Radius: 65, Move: false},
			{X: -304, Y: 245, Radius: 78, Move: false},
			{X: -357, Y: 149, Radius: 63, Move: false},
			{X: -285, Y: 90, Radius: 63, Move: false},
			{X: -234, Y: 279, Radius: 63, Move: false},
			{X: -146, Y: 418, Radius: 47, Move: false},
			{X: -80, Y: 405, Radius: 47, Move: false},
			{X: -116, Y: 358, Radius: 53, Move: false},
			{X: -239, Y: -166, Radius: 177, Move: false},
			{X: -253, Y: -323, Radius: 133, Move: false},
			{X: -369, Y: -302, Radius: 86, Move: false},
			{X: 214, Y: -347, Radius: 86, Move: false},
			{X: 285, Y: -177, Radius: 86, Move: false},
			{X: 266, Y: -212, Radius: 86, Move: false},
			{X: -121, Y: -380, Radius: 86, Move: false},
			{X: 50, Y: -364, Radius: 112, Move: false},
			{X: -289, Y: 170, Radius: 59, Move: false},
		},
		HeightType: 400,
		Static:     true,
	}
	obj.GetPhysicalModel().SetPos(512, 512, 0)
	brMap.AddStaticObject(obj)
	obj.SetGeoData()

	obj = &dynamic_map_object.Object{
		Type:        "mountains",
		Texture:     "mountain_5_2",
		MaxHP:       -2,
		TypeMaxHP:   -2,
		HP:          -2,
		Scale:       39,
		TypeGeoData: []*obstacle_point.ObstaclePoint{},
		Static:      true,
		Immortal:    true,
	}
	obj.GetPhysicalModel().SetPos(511, 512, 0)
	brMap.AddStaticObject(obj)
	obj.SetGeoData()

	for i := 0; i < 50; i++ {
		obj = &dynamic_map_object.Object{
			Type:       "unknown_civilization",
			Texture:    "unknown_civilization_7_4",
			MaxHP:      -2,
			TypeMaxHP:  -2,
			HP:         -2,
			Scale:      25,
			HeightType: 200,
			Weight:     1000,
			TypeGeoData: []*obstacle_point.ObstaclePoint{
				{X: 1, Y: 60, Radius: 57, Move: false},
				{X: -2, Y: 10, Radius: 57, Move: false},
				{X: 0, Y: -56, Radius: 57, Move: false},
			},
			Static:   false,
			Immortal: true,
		}

		obj.GetPhysicalModel().SetPos(float64(game_math.GetRangeRand(50, 950)), float64(game_math.GetRangeRand(50, 950)), float64(rand.Intn(360)))
		brMap.AddDynamicObject(obj)
		obj.SetGeoData()
	}

	brMap.GetJSONStaticObjects()
	return brMap
}

func offsetX(x, rowCount int) int {
	if rowCount%2 == 0 {
		return x + stepX/2
	} else {
		return x
	}
}
