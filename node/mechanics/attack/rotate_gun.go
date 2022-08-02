package attack

import (
	"github.com/TrashPony/game-engine/node/mechanics/actual_target"
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/body"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/target"
)

type RotateGunMessage struct {
	ID         int `json:"id"`
	MS         int `json:"ms"`
	Rotate     int `json:"rotate"`
	SlotNumber int `json:"slot_number"`
}

func RotateGun(gunner Gunner, mp *_map.Map, gunnerID int, weaponSlot *body.WeaponSlot, toBody bool) *RotateGunMessage {

	weaponTarget := gunner.GetWeaponTarget()
	if !actual_target.GetXYZTarget(gunner, weaponTarget, mp, weaponSlot) {
		if toBody {
			return RotateGunToBody(gunner, gunnerID, weaponSlot)
		}
		return nil
	}

	currentRotate, time, diffAngle := rotateGunToTarget(weaponTarget, _const.ServerTick, gunner, weaponSlot.Number)
	if diffAngle > 0 {
		gunner.SetGunRotate(currentRotate, weaponSlot.Number)
		return &RotateGunMessage{
			ID:         gunnerID,
			MS:         time,
			Rotate:     int(currentRotate),
			SlotNumber: weaponSlot.Number,
		}
	}

	return nil
}

func SetAllGunRotate(gunner Gunner, addRotate float64) {
	for _, weaponSlot := range gunner.RangeWeaponSlots() {
		if weaponSlot != nil {
			weaponSlot.SetGunRotate(weaponSlot.GetGunRotate() + addRotate)
		}
	}
}

func rotateGunToTarget(weaponTarget *target.Target, time int, gunner Gunner, weaponSlotNumber int) (float64, int, float64) {

	// из за того что пуле вылетают часто не из центра оружия а от смещения, то поворачивать оружие обычным способо некоректно
	// при повороте оружия будет смещатся ХУ позиции вылета пули, а вместе с ним и необходимый угол для поворота
	// следовательно мы должны на каждой итерации проверять позицию, необходимый угол и необходимость поворота

	RotateSpeed := gunner.GetGunRotateSpeed(weaponSlotNumber)
	currentRotate := gunner.GetGunRotate(weaponSlotNumber)

	realRotateGunSpeed := float64(RotateSpeed) / float64(1000/time)
	countRotateAngle := 0.0

	for i := 0.0; i < realRotateGunSpeed; i++ {

		fp := gunner.GetFirePos(weaponSlotNumber)
		x, y := fp.X, fp.Y

		needRotate := game_math.GetBetweenAngle(float64(weaponTarget.GetX()), float64(weaponTarget.GetY()), float64(x), float64(y))

		diff := game_math.Rotate(&currentRotate, &needRotate, 1)
		gunner.SetGunRotate(currentRotate, weaponSlotNumber)
		if diff <= 0.2 {
			break
		}

		countRotateAngle += diff
	}

	// time/realRotateGun = мс за 1 градус,
	// получаем сколько реально времени требуется для поворота
	time = int((float64(time) / realRotateGunSpeed) * countRotateAngle)

	return currentRotate, time, countRotateAngle
}

func RotateGunToBody(gunner Gunner, gunnerID int, weaponSlot *body.WeaponSlot) *RotateGunMessage {

	realRotateGun := float64(weaponSlot.Weapon.RotateSpeed) / float64(1000/_const.ServerTick)

	needRotate := gunner.GetRotate()
	rotate := gunner.GetGunRotate(weaponSlot.Number)

	diff := game_math.Rotate(&rotate, &needRotate, realRotateGun)
	if diff > 1 {
		gunner.SetGunRotate(rotate, weaponSlot.Number)
		return &RotateGunMessage{
			ID:         gunnerID,
			MS:         _const.ServerTick,
			Rotate:     int(rotate),
			SlotNumber: weaponSlot.Number,
		}
	} else {
		return nil
	}
}
