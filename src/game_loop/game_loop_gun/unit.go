package game_loop_gun

import (
	"github.com/TrashPony/game_engine/src/mechanics/attack"
	_map "github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
	"github.com/TrashPony/game_engine/src/web_socket"
)

func Unit(mp *_map.Map, units []*unit.Unit) (*web_socket.GameLoopMessages, *web_socket.GameLoopMessages) {

	fireMsgs := &web_socket.GameLoopMessages{}
	rotateMsgs := &web_socket.GameLoopMessages{}

	for _, gunnerUnit := range units {
		// метод fire должен возвращать события выстрелов и мы их отправляем, если выстрелов не произошло то можно крутить турельку
		// метод rotate тоже возаращает нам данные о повороте и мы отправляем сообщения на их основе
		for weaponSlot := range gunnerUnit.RangeWeaponSlots() {
			fireWeapons := attack.Fire(gunnerUnit.GetGunner(), weaponSlot, "unit", gunnerUnit.GetID(), mp)
			if len(fireWeapons) == 0 {
				rotateMsg := attack.RotateGun(gunnerUnit.GetGunner(), gunnerUnit.GetID(), weaponSlot)
				if rotateMsg != nil {
					rotateMsgs.AddMessage(rotateMsg)
				}
			} else {
				for _, m := range fireWeapons {
					fireMsgs.AddMessage(m)
				}
			}
		}
	}

	return fireMsgs, rotateMsgs
}
