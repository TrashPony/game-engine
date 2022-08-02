package game_math

import (
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/coordinate"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/position"
)

func GetAnchorWeapon(weaponXAttach, weaponYAttach, slotXAttach, slotYAttach int) (float64, float64, int, int) {

	xAnchor := (float64(weaponXAttach*100) / _const.SpriteSize) / 100
	yAnchor := (float64(weaponYAttach*100) / _const.SpriteSize) / 100

	// при отрисовке спрайта всегда берется 100% скейл поэтому нельзя считать обычным методом GetWeaponSlotAttachPoint()
	realAttachX := slotXAttach - _const.SpriteSize2
	realAttachY := slotYAttach - _const.SpriteSize2

	return xAnchor, yAnchor, realAttachX, realAttachY
}

func GetWeaponSlotAttachPoint(slotXAttach, slotYAttach, rotate float64, ownerScale int) (float64, float64) {

	sizeOffset := float64(ownerScale) / 100

	// к полизиции оружия прибавляем позицию атаки, получаем точку вылета пули при скейле 1 и угле 0
	x := (slotXAttach - _const.SpriteSize2) * sizeOffset
	y := (slotYAttach - _const.SpriteSize2) * sizeOffset

	newX, newY := RotatePoint(x, y, 0, 0, rotate)

	return newX, newY
}

func GetWeaponPosInMap(ownerX, ownerY, ownerScale int, slotXAttach, slotYAttach, rotate float64) (int, int) {
	x, y := GetWeaponSlotAttachPoint(slotXAttach, slotYAttach, rotate, ownerScale)
	return int(x) + ownerX, int(y) + ownerY
}

func GetWeaponFirePositions(ownerX, ownerY, ownerScale int, ownerRotate float64, gunRotate float64, weaponXAttach, weaponYAttach int,
	weaponFirePositions []*coordinate.Coordinate, slotXAttach, slotYAttach float64) []*position.Positions {

	// смещение размера за счет размера спрайта
	sizeOffset := float64(ownerScale) / 100

	// берем реальную позицию оружия на карте
	xWeapon, yWeapon := GetWeaponPosInMap(ownerX, ownerY, ownerScale, slotXAttach, slotYAttach, ownerRotate)

	realPos := make([]*position.Positions, len(weaponFirePositions), len(weaponFirePositions))

	for i, pos := range weaponFirePositions {
		// TODO повторяющийся код
		// тупо прибовляем позицию атаки, незабывая отнять смещение оружия от 0
		x := float64(xWeapon) + (float64(pos.X-weaponXAttach) * sizeOffset)
		y := float64(yWeapon) + (float64(pos.Y-weaponYAttach) * sizeOffset)

		// поворачиваем получившееся точки на угол оружия
		newX, newY := RotatePoint(x, y, float64(xWeapon), float64(yWeapon), gunRotate)

		// прибавляем точку обьекта к точке оружия
		realPos[i] = &position.Positions{X: int(newX), Y: int(newY)}
	}

	return realPos
}

func GetWeaponFirePosition(ownerX, ownerY, ownerScale int, ownerRotate float64, gunRotate float64, weaponXAttach, weaponYAttach int,
	weaponFirePositions []*coordinate.Coordinate, slotXAttach, slotYAttach float64, firePosition int) *position.Positions {

	// TODO повторяющийся код
	// смещение размера за счет размера спрайта
	sizeOffset := float64(ownerScale) / 100

	// берем реальную позицию оружия на карте
	xWeapon, yWeapon := GetWeaponPosInMap(ownerX, ownerY, ownerScale, slotXAttach, slotYAttach, ownerRotate)

	pos := weaponFirePositions[firePosition]

	// тупо прибовляем позицию атаки, незабывая отнять смещение оружия от 0
	x := float64(xWeapon) + (float64(pos.X-weaponXAttach) * sizeOffset)
	y := float64(yWeapon) + (float64(pos.Y-weaponYAttach) * sizeOffset)

	// поворачиваем получившееся точки на угол оружия
	newX, newY := RotatePoint(x, y, float64(xWeapon), float64(yWeapon), gunRotate)

	// прибавляем точку обьекта к точке оружия
	return &position.Positions{X: int(newX), Y: int(newY)}
}
