package unit

import (
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/physical_model"
)

func (u *Unit) GetPhysicalModel() *physical_model.PhysicalModel {
	if u.physicalModel == nil {
		u.initPhysicalModel()
	}

	return u.physicalModel
}

func (u *Unit) GetCopyPhysicalModel() *physical_model.PhysicalModel {
	pm := *u.physicalModel
	return &pm
}

func (u *Unit) initPhysicalModel() {
	u.physicalModel = &physical_model.PhysicalModel{
		Speed:         u.GetMoveMaxPower() / _const.ServerTickSecPart,
		ReverseSpeed:  u.GetMaxReverse() / _const.ServerTickSecPart,
		PowerFactor:   u.GetPowerFactor() / _const.ServerTickSecPart,
		ReverseFactor: u.GetReverseFactor() / _const.ServerTickSecPart,
		TurnSpeed:     u.GetTurnSpeed() / _const.ServerTickSecPart,
		WASD:          physical_model.WASD{},
		MoveDrag:      u.Body.MoveDrag,
		AngularDrag:   u.Body.AngularDrag,
		Weight:        u.Body.Weight,
		ChassisType:   u.Body.ChassisType,
	}

	// применяем настройки размера к обьектам колизий
	sizeOffset := float64(u.Body.Scale) / 100
	u.physicalModel.Height = float64(u.Body.Height) * sizeOffset
	u.physicalModel.Length = float64(u.Body.Length) * sizeOffset
	u.physicalModel.Width = float64(u.Body.Width) * sizeOffset
	u.physicalModel.Radius = int(float64(u.Body.Radius) * sizeOffset)
	u.HP = u.Body.MaxHP
}
