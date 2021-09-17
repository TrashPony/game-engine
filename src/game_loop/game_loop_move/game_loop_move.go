package game_loop_move

import (
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/mechanics/attack"
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/dynamic_map_object"
	_map "github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/obstacle_point"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
	"github.com/TrashPony/game_engine/src/web_socket"
	"math"
)

type moveObject interface {
	GetX() int
	GetY() int
	GetVelocity() (float64, float64)
	SetVelocity(float64, float64)
	MultiplyVelocity(float64, float64)
	AddVelocity(float64, float64)
	GetRotate() float64
	SetPos(float64, float64, float64)
	GetRealPos() (float64, float64)
	GetWeight() float64
	GetDirection() bool
	GetCurrentSpeed() float64
	SetPowerMove(float64)
	GetPowerMove() float64
	GetWidth() float64
	GetLength() float64
	GetVelocityRotate() float64

	CheckGrowthPower() bool
	CheckGrowthRevers() bool
	CheckLeftRotate() bool
	CheckRightRotate() bool

	GetAngularVelocity() float64
	SetAngularVelocity(float64)

	GetReverse() float64
	SetReverse(float64)

	GetMoveMaxPower() float64
	GetMaxReverse() float64

	GetPowerFactor() float64
	GetReverseFactor() float64
	GetTurnSpeed() float64

	GetMoveDrag() float64
	GetAngularDrag() float64

	SetWASD(bool, bool, bool, bool)
	GetZ() float64
	SetZ(float64)

	GetGeoData() []*obstacle_point.ObstaclePoint
	SetPosFunc(func())
}

func WasdMove(obj moveObject, gunner attack.Gunner, typeObj string, id int, mp *_map.Map, msgs *web_socket.GameLoopMessages, units []*unit.Unit, mapObjects map[int]*dynamic_map_object.Object) {

	if obj.CheckGrowthPower() {
		obj.SetPowerMove(obj.GetPowerMove() + obj.GetPowerFactor())
	} else {
		obj.SetPowerMove(obj.GetPowerMove() - obj.GetPowerFactor())
	}

	if obj.CheckGrowthRevers() {
		obj.SetReverse(obj.GetReverse() + obj.GetReverseFactor())
	} else {
		obj.SetReverse(obj.GetReverse() - obj.GetReverseFactor())
	}

	obj.SetPowerMove(math.Max(0, math.Min(obj.GetMoveMaxPower(), obj.GetPowerMove())))
	obj.SetReverse(math.Max(0, math.Min(obj.GetMaxReverse(), obj.GetReverse())))

	direction := 1.0
	if obj.GetPowerMove() < obj.GetReverse() {
		direction = -1
	}

	if obj.CheckLeftRotate() {
		obj.SetAngularVelocity(obj.GetAngularVelocity() - direction*obj.GetTurnSpeed())
	}

	if obj.CheckRightRotate() {
		obj.SetAngularVelocity(obj.GetAngularVelocity() + direction*obj.GetTurnSpeed())
	}

	radRotate := game_math.DegToRadian(obj.GetRotate())
	obj.AddVelocity(
		game_math.Cos(radRotate)*(obj.GetPowerMove()-obj.GetReverse()),
		game_math.Sin(radRotate)*(obj.GetPowerMove()-obj.GetReverse()),
	)

	xV, yV := obj.GetVelocity()
	xR, yR := obj.GetRealPos()
	xNext, yNext := xR+xV, yR+yV

	// проверяем колизии
	checkMoveCollision(typeObj, id, obj, xNext, yNext, obj.GetRotate()+game_math.RadianToDeg(obj.GetAngularVelocity()), mp, units, mapObjects)

	// заного перепроверяем поицию и скорость т.к. за время колизии могло произойти столкновение
	xV, yV = obj.GetVelocity()
	xNext, yNext = xR+xV, yR+yV

	// оповещаем мир как двигается обьект
	if typeObj == "unit" {
		_, _, lvl := mp.GetPosLevel(obj.GetX(), obj.GetY())
		SendMoveUnit(obj, id, int(xNext), int(yNext), _const.ServerTick, mp.Id, obj.GetRotate()+game_math.RadianToDeg(obj.GetAngularVelocity()), lvl, obj.GetCurrentSpeed(), true, msgs)
	}

	// принимаем позицию, эта функция отрботает в начале следующего тика
	obj.SetPosFunc(func() {

		obj.MultiplyVelocity(obj.GetMoveDrag(), obj.GetMoveDrag())
		obj.SetPos(xNext, yNext, obj.GetRotate()+game_math.RadianToDeg(obj.GetAngularVelocity()))
		obj.SetAngularVelocity(obj.GetAngularVelocity() * obj.GetAngularDrag())

		// тело поворачивается, и пуха тоже
		if gunner != nil {
			attack.SetAllGunRotate(gunner, game_math.RadianToDeg(obj.GetAngularVelocity()))
		}
	})
}
