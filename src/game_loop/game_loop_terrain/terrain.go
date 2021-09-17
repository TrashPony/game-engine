package game_loop_terrain

import (
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/mechanics/game_math/collisions"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/dynamic_map_object"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/plant_life_game"
)

/* ФУНКЦИЯ ПО УПРОВЛЕНИЮ ОВОЩАМИ! */

// реализация алгоритма зеро плей игры "Жизнь" для динамиской популяции растений на карте
// https://ru.wikipedia.org/wiki/%D0%98%D0%B3%D1%80%D0%B0_%C2%AB%D0%96%D0%B8%D0%B7%D0%BD%D1%8C%C2%BB

func TerrainLife(mp *_map.Map, mapObjects map[int]*dynamic_map_object.Object) {
	growTerrain(mp, mapObjects)
	plant_life_game.Population(mp, mapObjects)
}

func growTerrain(mp *_map.Map, mapObjects map[int]*dynamic_map_object.Object) {

	allEnd := true
	for _, obj := range mapObjects {

		if obj.DestroyLeftTimer {
			if obj.DestroyTimer <= 0 {
				mp.RemoveDynamicObject(obj)
				continue
			} else {
				obj.DestroyTimer -= _const.ServerTick
			}
		}

		// радар игроков сам обновит обьект
		if obj.GrowCycle <= 0 {
			continue
		}

		allEnd = false

		obj.GrowLeftTime -= _const.ServerTick
		if obj.GrowLeftTime > 0 {
			continue
		}

		obj.GrowLeftTime = obj.GrowTime
		obj.GrowCycle--

		if obj.GetScale() < obj.MaxScale && !collisions.CheckObjectCollision(obj, mp, true) {
			obj.SetScale(obj.GetScale() + 1)
			obj.CalculateScale()
			// удаляем старый кеш обьекта
		}
	}

	if allEnd {
		plant_life_game.Evolution(mp, mapObjects)
	}
}
