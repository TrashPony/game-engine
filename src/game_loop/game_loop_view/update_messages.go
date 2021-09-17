package game_loop_view

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/dynamic_map_object"
	_map "github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/visible_objects"
)

type updateUnitMsg struct {
	HP int `json:"hp"`
}

func updateUnitMessage(unit *unit.Unit) updateUnitMsg {
	return updateUnitMsg{HP: unit.HP}
}

func updateObject(obj *dynamic_map_object.Object, memoryObj *visible_objects.VisibleObject, mp *_map.Map) (map[string]interface{}, bool) {

	// TODO каждый юзер запомнил обьект по своему поэтому это не работает в полной мере
	// т.к. это общий метод для всех игроков то его можно обновлять 1 раз за тик для всех сразу
	// не создавать каждый раз мапу а создать 1 раз и просто обновлять ее и делать json

	//uMap := make(map[string]interface{})
	//update := false

	if memoryObj == nil {
		return nil, false
	}

	if memoryObj.UpdateMsg == nil {
		memoryObj.UpdateMsg = &visible_objects.UpdateObjMap{UpdateMap: make(map[string]interface{})}
	}

	memoryObj.UpdateMsg.Mx.Lock()
	defer memoryObj.UpdateMsg.Mx.Unlock()

	if memoryObj.UpdateMsg.ServerTime == mp.Time {
		return memoryObj.UpdateMsg.UpdateMap, memoryObj.UpdateMsg.Update
	}

	memoryObj.UpdateMsg.ServerTime = mp.Time
	memoryObj.UpdateMsg.Update = false

	if obj.GetScale() != memoryObj.Scale {
		memoryObj.UpdateMsg.Update = true
		memoryObj.UpdateMsg.UpdateMap["s"] = obj.GetScale()
		memoryObj.UpdateMsg.UpdateMap["mhp"] = obj.MaxHP
		memoryObj.UpdateMsg.UpdateMap["gt"] = obj.GetGrowTime()
		memoryObj.UpdateMsg.UpdateMap["xs"] = obj.XShadowOffset
		memoryObj.UpdateMsg.UpdateMap["ys"] = obj.YShadowOffset
		memoryObj.UpdateMsg.UpdateMap["gd"] = obj.GetGeoDataJSON()
	} else {
		delete(memoryObj.UpdateMsg.UpdateMap, "s")
		delete(memoryObj.UpdateMsg.UpdateMap, "mhp")
		delete(memoryObj.UpdateMsg.UpdateMap, "gt")
		delete(memoryObj.UpdateMsg.UpdateMap, "xs")
		delete(memoryObj.UpdateMsg.UpdateMap, "ys")
		delete(memoryObj.UpdateMsg.UpdateMap, "gd")
	}

	if obj.GetHP() != memoryObj.HP {
		memoryObj.UpdateMsg.Update = true
		memoryObj.UpdateMsg.UpdateMap["hp"] = obj.GetHP()
	} else {
		delete(memoryObj.UpdateMsg.UpdateMap, "hp")
	}

	if obj.OwnerID != memoryObj.OwnerID {
		memoryObj.UpdateMsg.Update = true
		memoryObj.UpdateMsg.UpdateMap["ow"] = obj.GetOwnerID()
	} else {
		delete(memoryObj.UpdateMsg.UpdateMap, "ow")
	}

	return memoryObj.UpdateMsg.UpdateMap, memoryObj.UpdateMsg.Update
}
