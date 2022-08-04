package game_loop_gun

import (
	"github.com/TrashPony/game-engine/node/binary_msg"
	"github.com/TrashPony/game-engine/node/mechanics/attack"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/web_socket_response"
	"github.com/TrashPony/game-engine/router/web_socket"
)

func Unit(units []*unit.Unit, battle *battle.Battle, ms *web_socket.MessagesStore) {

	if battle.WaitReady {
		return
	}

	for _, gunnerUnit := range units {

		if gunnerUnit.GetPhysicalModel().Fly {
			continue
		}

		// метод fire должен возвращать события выстрелов и мы их отправляем
		// метод rotate тоже возаращает нам данные о повороте и мы отправляем сообщения на их основе
		for _, weaponSlot := range gunnerUnit.RangeWeaponSlots() {

			if weaponSlot == nil || weaponSlot.Weapon == nil {
				continue
			}

			rotateMsg := attack.RotateGun(gunnerUnit.GetGunner(), battle.Map, gunnerUnit.GetID(), weaponSlot, true)
			if rotateMsg != nil {
				ms.AddMsg("rotateMsgs", "bin", web_socket_response.Response{
					BinaryMsg: binary_msg.CreateRotateGunBinaryMsg(rotateMsg.ID, rotateMsg.MS, rotateMsg.Rotate, rotateMsg.SlotNumber),
					X:         gunnerUnit.GetPhysicalModel().GetX(),
					Y:         gunnerUnit.GetPhysicalModel().GetY(),
				}, nil)
			}

			fireWeapons := attack.Fire(gunnerUnit.GetGunner(), weaponSlot, "unit", gunnerUnit.GetID(), battle, units)
			if len(fireWeapons) > 0 {
				for _, m := range fireWeapons {
					ms.AddMsg("fireMsgs", "bin", web_socket_response.Response{
						BinaryMsg: binary_msg.CreateFireGunBinaryMsg(m.TypeID, m.X, m.Y, m.Z, m.Rotate, m.AccumulationPercent),
						X:         m.X,
						Y:         m.Y,
					}, nil)
				}
			}

			if gunnerUnit.GetWeaponTarget() != nil {
				weaponTargetMSG := attack.WeaponTarget(gunnerUnit.GetGunner(), weaponSlot, battle, "unit", gunnerUnit.GetID(), units)

				if gunnerUnit != nil {
					ms.AddMsg("targetMsg", "bin", web_socket_response.Response{
						BinaryMsg: binary_msg.WeaponMouseTargetBinary(weaponTargetMSG.X, weaponTargetMSG.Y,
							weaponTargetMSG.Accuracy, weaponTargetMSG.AmmoCount, weaponTargetMSG.AmmoAvailable,
							weaponTargetMSG.AccumulationPercent, weaponTargetMSG.Reload, weaponTargetMSG.Chase,
							weaponTargetMSG.TargetType, weaponTargetMSG.TargetID),
						PlayerID: gunnerUnit.OwnerID,
					}, nil)
				}
			}
		}
	}

	return
}
