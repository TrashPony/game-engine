package game_loop_view

import (
	"github.com/TrashPony/game-engine/node/binary_msg"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/dynamic_map_object"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/visible_objects"
)

func updateUnitMessage(unit *unit.Unit) []byte {
	command := make([]byte, 8)

	binary_msg.ReuseByteSlice(&command, 0, binary_msg.GetIntBytes(unit.HP))
	binary_msg.ReuseByteSlice(&command, 4, binary_msg.GetIntBytes(unit.TeamID))

	return command
}

func updateObject(obj *dynamic_map_object.Object, memoryObj *visible_objects.VisibleObject, mp *_map.Map) ([]byte, bool) {

	// т.к. это общий метод для всех игроков то его можно обновлять 1 раз за тик для всех сразу
	// не создавать каждый раз мапу а создать 1 раз и просто обновлять ее и делать json
	if memoryObj == nil {
		return nil, false
	}

	if memoryObj.UpdateMsg == nil {
		memoryObj.UpdateMsg = &visible_objects.UpdateObjMap{}
	}

	memoryObj.UpdateMsg.Mx.Lock()
	defer memoryObj.UpdateMsg.Mx.Unlock()

	if memoryObj.UpdateMsg.ServerTime == mp.Time {
		return memoryObj.UpdateMsg.UpdateData, memoryObj.UpdateMsg.Update
	}

	if memoryObj.UpdateMsg.UpdateData == nil {
		memoryObj.UpdateMsg.UpdateData = make([]byte, 31+len(obj.GetGeoDataBin()))
	}

	memoryObj.UpdateMsg.ServerTime = mp.Time
	memoryObj.UpdateMsg.Update = false

	binary_msg.ReuseByteSlice(&memoryObj.UpdateMsg.UpdateData, 0, []byte{
		binary_msg.BoolToByte(obj.Work),
		byte(obj.Scale),
		byte(obj.XShadowOffset),
		byte(obj.YShadowOffset), // 4
	})

	binary_msg.ReuseByteSlice(&memoryObj.UpdateMsg.UpdateData, 4, binary_msg.GetIntBytes(obj.GetHP()))
	binary_msg.ReuseByteSlice(&memoryObj.UpdateMsg.UpdateData, 8, binary_msg.GetIntBytes(obj.MaxHP))
	binary_msg.ReuseByteSlice(&memoryObj.UpdateMsg.UpdateData, 12, binary_msg.GetIntBytes(obj.CurrentEnergy))
	binary_msg.ReuseByteSlice(&memoryObj.UpdateMsg.UpdateData, 16, binary_msg.GetIntBytes(obj.MaxEnergy))
	binary_msg.ReuseByteSlice(&memoryObj.UpdateMsg.UpdateData, 20, binary_msg.GetIntBytes(obj.OwnerID))
	binary_msg.ReuseByteSlice(&memoryObj.UpdateMsg.UpdateData, 24, binary_msg.GetIntBytes(obj.GrowTime))

	binary_msg.ReuseByteSlice(&memoryObj.UpdateMsg.UpdateData, 28, []byte{
		byte(obj.GetMapHeight()),
	})

	binary_msg.ReuseByteSlice(&memoryObj.UpdateMsg.UpdateData, 29, []byte{
		byte(obj.TeamID),
	})

	if obj.GetScale() != memoryObj.Scale {
		memoryObj.UpdateMsg.Update = true

		binary_msg.ReuseByteSlice(&memoryObj.UpdateMsg.UpdateData, 30, []byte{1})
		binary_msg.ReuseByteSlice(&memoryObj.UpdateMsg.UpdateData, 31, obj.GetGeoDataBin())
	} else {
		binary_msg.ReuseByteSlice(&memoryObj.UpdateMsg.UpdateData, 30, []byte{0})
	}

	if obj.GetHP() != memoryObj.HP {
		memoryObj.UpdateMsg.Update = true
	}

	if obj.OwnerID != memoryObj.OwnerID {
		memoryObj.UpdateMsg.Update = true
	}

	if obj.Work != memoryObj.Work {
		memoryObj.UpdateMsg.Update = true
	}

	if obj.GetPower() != memoryObj.Energy {
		memoryObj.UpdateMsg.Update = true
	}

	if obj.TeamID != memoryObj.TeamID {
		memoryObj.UpdateMsg.Update = true
	}

	return memoryObj.UpdateMsg.UpdateData, memoryObj.UpdateMsg.Update
}
