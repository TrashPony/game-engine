package attack

import (
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/body"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/target"
	"math"
)

type RotateGunMessage struct {
	ID         int `json:"id"`
	MS         int `json:"ms"`
	Rotate     int `json:"rotate"`
	SlotNumber int `json:"slot_number"`
}

func RotateGun(gunner Gunner, gunnerID int, weaponSlot *body.WeaponSlot) *RotateGunMessage {

	weaponTarget := gunner.GetWeaponTarget()
	if weaponTarget == nil {
		return RotateGunToBody(gunner, gunnerID, weaponSlot)
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
	for weaponSlot := range gunner.RangeWeaponSlots() {
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

		x, y := gunner.GetFirePos(weaponSlotNumber).X, gunner.GetFirePos(weaponSlotNumber).Y
		needRotate := game_math.GetBetweenAngle(float64(weaponTarget.GetX()), float64(weaponTarget.GetY()), float64(x), float64(y))

		diff := Rotate(&currentRotate, &needRotate, 1)
		gunner.SetGunRotate(currentRotate, weaponSlotNumber)
		if diff <= 0.1 {
			break
		}

		countRotateAngle += diff
	}

	// time/realRotateGun = мс за 1 градус,
	// получаем сколько реально времени требуется для поворота
	time = int((float64(time) / realRotateGunSpeed) * countRotateAngle)

	return currentRotate, time, countRotateAngle
}

func Rotate(unitRotate, needRotate *float64, step float64) float64 {

	prepareAngle(unitRotate)
	prepareAngle(needRotate)

	countRotateAngle := 0.0

	for i := 0.0; i < step*10; i++ {

		if math.Round(*unitRotate*10) != math.Round(*needRotate*10) {

			if directionRotate(*unitRotate, *needRotate) {
				*unitRotate += 0.1
				if *unitRotate >= 360 {
					*unitRotate -= 360
				}

				countRotateAngle += 0.1

			} else {
				*unitRotate -= 0.1
				if *unitRotate < 0 {
					*unitRotate += 360
				}

				countRotateAngle += 0.1

			}

		} else {
			return countRotateAngle
		}
	}

	return countRotateAngle
}

func directionRotate(unitAngle, needAngle float64) bool {

	prepareAngle(&unitAngle)
	prepareAngle(&needAngle)

	// true ++
	// false --
	count := 0
	direction := false

	if unitAngle < needAngle {
		for unitAngle < needAngle {
			count++
			direction = true
			unitAngle++
		}
	} else {
		for unitAngle > needAngle {
			count++
			direction = false
			needAngle++
		}
	}

	if direction {
		return count <= 180
	} else {
		return !(count <= 180)
	}
}

func prepareAngle(angle *float64) {
	if *angle < 0 {
		*angle += 360
	}

	if *angle >= 360 {
		*angle -= 360
	}

	if *angle < 0 {
		*angle += 360
	}

	if *angle >= 360 {
		*angle -= 360
	}
}

func RotateGunToBody(gunner Gunner, gunnerID int, weaponSlot *body.WeaponSlot) *RotateGunMessage {

	realRotateGun := float64(gunner.GetWeaponSlot(weaponSlot.Number).Weapon.RotateSpeed) / float64(1000/_const.ServerTick)

	needRotate := gunner.GetRotate()
	rotate := gunner.GetGunRotate(weaponSlot.Number)

	diff := Rotate(&rotate, &needRotate, realRotateGun)
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
