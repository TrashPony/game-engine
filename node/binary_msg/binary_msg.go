package binary_msg

import (
	"fmt"
	_const "github.com/TrashPony/game-engine/router/const"
)

func CreateBinaryUnitMoveMsg(unitID, speed, x, y, z, ms, rotate, angularVelocity int, animate, direction, sky, a, d, w, sp bool) []byte {
	// [1[eventID], 4[unitID], 4[speed], 4[x], 4[y], 4[z], 4[ms], 4[rotate], 4[angularVelocity], 4[mpID] 1[animate], 1[direction], 1[sky]]

	command := make([]byte, 37)

	command[0] = 1
	ReuseByteSlice(&command, 1, GetIntBytes(unitID))
	ReuseByteSlice(&command, 5, GetIntBytes(speed))
	ReuseByteSlice(&command, 9, GetIntBytes(x))
	ReuseByteSlice(&command, 13, GetIntBytes(y))
	ReuseByteSlice(&command, 17, GetIntBytes(z))
	command[21] = byte(ms)
	ReuseByteSlice(&command, 22, GetIntBytes(rotate))
	ReuseByteSlice(&command, 26, GetIntBytes(angularVelocity))
	command[30] = BoolToByte(animate)

	command[31] = BoolToByte(direction)
	command[32] = BoolToByte(sky)
	command[33] = BoolToByte(a)
	command[34] = BoolToByte(d)
	command[35] = BoolToByte(w)
	command[36] = BoolToByte(sp)

	return command
}

func CreateMarkBinaryMove(id, x, y, ms int) []byte {
	// [1[eventID] 4[ID], 4[x], 4[y], 4[ms]]

	command := make([]byte, 14)

	command[0] = 5
	ReuseByteSlice(&command, 1, GetIntBytes(id))
	ReuseByteSlice(&command, 5, GetIntBytes(x))
	ReuseByteSlice(&command, 9, GetIntBytes(y))
	command[13] = byte(ms)

	return command
}

func CreateBulletBinaryFly(typeID, id, x, y, z, ms, rotate int) []byte {
	// [1[eventID] 4[typeID], 4[id], 4[x], 4[y], 4[z], 4[ms], 4[rotate], 4[mpID]]

	command := []byte{6}

	command = append(command, byte(typeID))
	command = append(command, GetIntBytes(id)...)
	command = append(command, GetIntBytes(x)...)
	command = append(command, GetIntBytes(y)...)
	command = append(command, GetIntBytes(z)...)
	command = append(command, byte(ms))
	command = append(command, GetIntBytes(rotate)...)

	return command
}

func CreateBulletLaserFly(typeID, x, y, toX, toY, unitID int) []byte {
	command := []byte{20}

	command = append(command, byte(typeID))
	command = append(command, GetIntBytes(x)...)
	command = append(command, GetIntBytes(y)...)
	command = append(command, GetIntBytes(toX)...)
	command = append(command, GetIntBytes(toY)...)
	command = append(command, GetIntBytes(unitID)...)

	return command
}

func CreateBulletBinaryExplosion(typeID, x, y, z int) []byte {
	// [1[eventID] 4[typeID], 4[x], 4[y], 4[z], 4[mpID]]

	command := []byte{7}

	command = append(command, byte(typeID))
	command = append(command, GetIntBytes(x)...)
	command = append(command, GetIntBytes(y)...)
	command = append(command, GetIntBytes(z)...)

	return command
}

func CreateRotateGunBinaryMsg(id, ms, rotate, slot int) []byte {
	// [1[eventID], 4[ID], 4[ms], 4[rotate]
	command := []byte{8}

	command = append(command, GetIntBytes(id)...)
	command = append(command, byte(ms))
	command = append(command, GetIntBytes(rotate)...)
	command = append(command, byte(slot))

	return command
}

func CreateFireGunBinaryMsg(typeID, x, y, z, rotate, accumulationPercent int) []byte {
	// [1[eventID] 4[typeID], 4[x], 4[y], 4[z], 4[rotate], 4[mpID]]
	command := []byte{9}

	command = append(command, byte(typeID))
	command = append(command, GetIntBytes(x)...)
	command = append(command, GetIntBytes(y)...)
	command = append(command, GetIntBytes(z)...)
	command = append(command, GetIntBytes(rotate)...)
	command = append(command, byte(accumulationPercent))

	return command
}

func RotateTurretGunBinaryMsg(id, rotate, ms int) []byte {
	// [1[eventID], 4[ID], 4[ms], 4[rotate]
	command := []byte{10}

	command = append(command, GetIntBytes(id)...)
	command = append(command, byte(ms))
	command = append(command, GetIntBytes(rotate)...)

	return command
}

func WeaponMouseTargetBinary(x, y, radius, ammoCount, ammoAvailable, accumulationPercent int, reload, chase bool, targetType string, targetID int) []byte {
	// [1[eventID], 4[id], 4[typeID], 4[x], 4[y], 4[mpID], 4[alpha], 4[rotate]]

	command := make([]byte, 20)

	command[0] = 13
	ReuseByteSlice(&command, 1, GetIntBytes(x))
	ReuseByteSlice(&command, 5, GetIntBytes(y))
	command[9] = byte(radius)
	command[10] = byte(ammoCount)
	command[11] = byte(ammoAvailable)
	command[12] = byte(accumulationPercent)
	command[13] = BoolToByte(reload)
	command[14] = BoolToByte(chase)
	command[15] = byte(_const.MapBinItems[targetType])
	ReuseByteSlice(&command, 16, GetIntBytes(targetID))

	return command
}

func StatusSquadBinaryMsg(hp, energy int) []byte {
	// [1[eventID], 4[hp], 4[energy]
	command := []byte{14}

	command = append(command, GetIntBytes(hp)...)
	command = append(command, GetIntBytes(energy)...)

	return command
}

func DamageTextBinaryMsg(x, y, damage int, typeObject string) []byte {
	// [1[eventID], 4[x], 4[y], 4[d], 4[m], 1[t]]
	command := []byte{17}

	_, ok := _const.MapBinItems[typeObject]
	if !ok {
		fmt.Println("unknown type object: ", typeObject)
	}

	command = append(command, GetIntBytes(x)...)
	command = append(command, GetIntBytes(y)...)
	command = append(command, GetIntBytes(damage)...)
	command = append(command, byte(_const.MapBinItems[typeObject]))

	return command
}

func ObjectDeadBinaryMsg(id, x, y int, typeObject string) []byte {
	// [1[eventID], 4[id], 4[x], 4[y], 4[m], 1[t]]
	command := []byte{18}

	_, ok := _const.MapBinItems[typeObject]
	if !ok {
		fmt.Println("unknown type object: ", typeObject)
	}

	command = append(command, GetIntBytes(id)...)
	command = append(command, GetIntBytes(x)...)
	command = append(command, GetIntBytes(y)...)
	command = append(command, byte(_const.MapBinItems[typeObject]))

	return command
}

var MarksTypes = map[string]byte{
	"fly":       1,
	"ground":    2,
	"structure": 3,
	"resource":  4,
	"bullet":    5,
}

func CreateBinaryCreateRadarMarkMsg(id, x, y int, typeMark string) []byte {
	command := []byte{39}

	command = append(command, GetIntBytes(id)...)
	command = append(command, GetIntBytes(x)...)
	command = append(command, GetIntBytes(y)...)
	command = append(command, MarksTypes[typeMark])

	return command
}

func CreateBinaryRemoveRadarMarkMsg(id int) []byte {
	command := []byte{40}

	command = append(command, GetIntBytes(id)...)

	return command
}

func CreateBinaryRemoveRadarObjectMsg(idObj int, typeObj string) []byte {
	command := []byte{41}

	command = append(command, GetIntBytes(idObj)...)
	command = append(command, byte(_const.MapBinItems[typeObj]))

	return command
}

func CreateBinaryCreateObjMsg(typeObj string, data []byte) []byte {
	command := make([]byte, 6+len(data))

	command[0] = 42
	command[1] = byte(_const.MapBinItems[typeObj])
	ReuseByteSlice(&command, 2, GetIntBytes(len(data)))
	ReuseByteSlice(&command, 6, data)

	return command
}

func CreateBinaryUpdateObjMsg(typeObj string, id int, data []byte) []byte {

	command := make([]byte, 10+len(data))

	command[0] = 43
	command[1] = byte(_const.MapBinItems[typeObj])
	ReuseByteSlice(&command, 2, GetIntBytes(id))
	ReuseByteSlice(&command, 6, GetIntBytes(len(data)))
	ReuseByteSlice(&command, 10, data)

	return command
}

func RadareON() []byte {
	command := []byte{49}
	return command
}

func CreateBinaryObjectMove(id, x, y, ms, rotate int) []byte {
	// [1[eventID], 1[TypeObj] 4[ID], 4[x], 4[y], 4[ms], 4[rotate], 4[mpID]]
	command := []byte{3}

	command = append(command, GetIntBytes(id)...)
	command = append(command, GetIntBytes(x)...)
	command = append(command, GetIntBytes(y)...)
	command = append(command, GetIntBytes(ms)...)
	command = append(command, GetIntBytes(rotate)...)

	return command
}
