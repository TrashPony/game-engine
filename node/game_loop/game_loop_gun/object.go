package game_loop_gun

import (
	"github.com/TrashPony/game-engine/node/binary_msg"
	"github.com/TrashPony/game-engine/node/mechanics/attack"
	battle2 "github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/dynamic_map_object"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/web_socket_response"
	"github.com/TrashPony/game-engine/router/web_socket"
	"math/rand"
)

func Object(ms *web_socket.MessagesStore, b *battle2.Battle, units []*unit.Unit, buildObjects []*dynamic_map_object.Object) {

	for _, obj := range buildObjects {

		// метод fire должен возвращать события выстрелов и мы их отправляем, если выстрелов не произошло то можно крутить турельку
		// метод rotate тоже возаращает нам данные о повороте и мы отправляем сообщения на их основе
		for _, weaponSlot := range obj.RangeWeaponSlots() {

			fireWeapons := attack.Fire(obj.GetGunner(), weaponSlot, "object", obj.GetID(), b, units)

			if len(fireWeapons) == 0 {
				rotateMsg := attack.RotateGun(obj.GetGunner(), b.Map, obj.GetID(), weaponSlot, false)
				if rotateMsg != nil {
					ms.AddMsg("rotateTurretMsgs", "bin", web_socket_response.Response{
						BinaryMsg: binary_msg.RotateTurretGunBinaryMsg(rotateMsg.ID, rotateMsg.Rotate, rotateMsg.MS),
						X:         obj.GetPhysicalModel().GetX(),
						Y:         obj.GetPhysicalModel().GetY(),
					}, nil)
				} else {
					turretTarget := obj.GetWeaponTarget()
					if turretTarget != nil && ((turretTarget.Type == "map" || !turretTarget.Attack) && rand.Intn(25) == 0) {
						// если цели нет то смотрим загадочно в даль, и не часто что бы турель не колбасило
						obj.SetWeaponTarget(nil)
					}
				}
			} else {
				for _, m := range fireWeapons {
					ms.AddMsg("fireMsgs", "bin", web_socket_response.Response{
						BinaryMsg: binary_msg.CreateFireGunBinaryMsg(m.TypeID, m.X, m.Y, m.Z, m.Rotate, m.AccumulationPercent),
						X:         m.X,
						Y:         m.Y,
					}, nil)
				}
			}
		}
	}

	return
}
