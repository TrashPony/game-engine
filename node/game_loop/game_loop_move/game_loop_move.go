package game_loop_move

import (
	"fmt"
	"github.com/TrashPony/game-engine/node/mechanics/attack"
	"github.com/TrashPony/game-engine/node/mechanics/collisions"
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	battle2 "github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/gunner"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/obstacle_point"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/physical_model"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"github.com/TrashPony/game-engine/router/web_socket"
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
	CheckGrowthPower() bool
	CheckGrowthRevers() bool
	CheckLeftRotate() bool
	CheckRightRotate() bool
	CheckHandBrake() bool

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

	SetWASD(bool, bool, bool, bool, bool, bool, bool)
	GetZ() float64
	SetZ(float64)
	GetPhysicalModel() *physical_model.PhysicalModel
	GetGeoData() []*obstacle_point.ObstaclePoint
	SetPosFunc(func())
	SetNextPos(x, y float64)
	GetNextPos() (float64, float64)
	GetNeedZ() float64
	IsFly() bool
	GetChassisType() string
	GetType() string
}

type initMoveObj interface {
	GetID() int
	GetPhysicalModel() *physical_model.PhysicalModel
	GetGunner() *gunner.Gunner
}

func initMove(typeObj string, moveObjects []initMoveObj, b *battle2.Battle, units []*unit.Unit, ms *web_socket.MessagesStore) {
	for _, mo := range moveObjects {
		// расчитываем слудующие положение предметов
		WasdMove(mo.GetPhysicalModel())
	}

	for _, mo := range moveObjects {
		// проверяем колизии, и если есть столкновнеия перерасчитываем вектора движения
		xNext, yNext := mo.GetPhysicalModel().GetNextPos()
		checkMoveCollision(typeObj, mo.GetID(), mo.GetPhysicalModel(), xNext, yNext, mo.GetPhysicalModel().GetRotate(), b, units, ms)
	}

	for _, mo := range moveObjects {
		// оповещаем фронт о движение
		mo2 := mo // что бы в замыканеие не переопределился ганер
		obj := mo2.GetPhysicalModel()
		xNext, yNext := obj.GetNextPos()

		// само новое положение примется в начале следующего тика при инициализации
		obj.SetPosFunc(func() {
			obj.MultiplyVelocity(obj.GetMoveDrag(), obj.GetMoveDrag())
			obj.SetAngularVelocity(obj.GetAngularVelocity() * obj.GetAngularDrag())

			if xNext == math.NaN() {
				fmt.Println("math.NaN")
			}

			obj.SetPos(xNext, yNext, obj.GetRotate()+game_math.RadianToDeg(obj.GetAngularVelocity()))

			// принимаем позицию
			// тело поворачивается, и пуха тоже
			collisions.GetBodyRect(obj.GetPhysicalModel(), xNext, yNext, obj.GetRotate(), false, false)
			if mo2.GetGunner() != nil {
				attack.SetAllGunRotate(mo2.GetGunner(), game_math.RadianToDeg(obj.GetAngularVelocity()))
			}
		})

		if typeObj == "unit" {
			_, _, lvl := b.Map.GetPosLevel(obj.GetX(), obj.GetY())
			SendMoveUnit(obj, mo.GetID(), int(xNext), int(yNext), _const.ServerTick, b.Map.Id, obj.GetRotate()+game_math.RadianToDeg(obj.GetAngularVelocity()), lvl+obj.GetZ(), obj.GetCurrentSpeed(), true, ms)
		}

		if typeObj == "object" && (obj.GetX() != int(xNext) || obj.GetY() != int(yNext) || obj.GetAngularVelocity() > 0.001) {
			SendMoveObject(obj.GetID(), int(xNext), int(yNext), _const.ServerTick, obj.GetRotate()+game_math.RadianToDeg(obj.GetAngularVelocity()), ms)
		}
	}
}

func WasdMove(obj moveObject) {
	if obj.GetChassisType() == "" || obj.GetChassisType() == "caterpillars" || obj.GetChassisType() == "wheels" || obj.GetChassisType() == "fly" {
		wheel(obj)
	}
}

func wheel(obj moveObject) {

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

	// ручной тормаз
	if obj.CheckHandBrake() {

		if obj.GetPowerMove() > obj.GetReverse() {
			obj.SetPowerMove(obj.GetPowerMove() - obj.GetPowerMove()/8)
		} else {
			obj.SetReverse(obj.GetReverse() - obj.GetReverse()/8)
		}

		if obj.GetPowerMove() < obj.GetPowerFactor()*10 {
			obj.SetPowerMove(0)
		}

		if obj.GetReverse() < obj.GetReverseFactor()*10 {
			obj.SetReverse(0)
		}
	}

	direction := 1.0
	if obj.GetPowerMove() < obj.GetReverse() {
		direction = -1
	}

	move := obj.GetPowerMove() > 0 || obj.GetReverse() > 0

	if obj.GetChassisType() == "caterpillars" || obj.GetChassisType() == "fly" || (obj.GetChassisType() == "wheels" && move) {
		if obj.CheckLeftRotate() {
			obj.SetAngularVelocity(obj.GetAngularVelocity() - (direction * obj.GetTurnSpeed()))
		}

		if obj.CheckRightRotate() {
			obj.SetAngularVelocity(obj.GetAngularVelocity() + (direction * obj.GetTurnSpeed()))
		}
	}

	radRotate := game_math.DegToRadian(obj.GetRotate())
	obj.AddVelocity(
		game_math.Cos(radRotate)*(obj.GetPowerMove()-obj.GetReverse()),
		game_math.Sin(radRotate)*(obj.GetPowerMove()-obj.GetReverse()),
	)

	xV, yV := obj.GetVelocity()
	xR, yR := obj.GetRealPos()

	obj.SetNextPos(xR+xV, yR+yV)
}
