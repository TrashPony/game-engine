package plant_life_game

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/dynamic_map_object"
)

func GetCountPlantIMap(textureName string, mapObjects map[int]*dynamic_map_object.Object) int {
	count := 0

	for _, obj := range mapObjects {
		if obj.Texture == textureName {
			count++
		}
	}

	return count
}
