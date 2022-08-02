package body

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/coordinate"
)

type Body struct {
	ID                        int                              `json:"id"`
	Name                      string                           `json:"name"`
	OnlyConfig                bool                             `json:"only_config"`
	Texture                   string                           `json:"texture"`
	MaxHP                     int                              `json:"max_hp"`
	Scale                     int                              `json:"scale"`
	Length                    int                              `json:"length"`
	Width                     int                              `json:"width"`
	Height                    int                              `json:"height"`
	Radius                    int                              `json:"radius"`
	RangeView                 int                              `json:"range_view"`
	RangeRadar                int                              `json:"range_radar"`
	RecoveryPower             int                              `json:"recovery_power"`
	RecoveryPowerCycle        int                              `json:"recovery_power_cycle"`
	CurrentRecoveryPowerCycle int                              `json:"-"`
	MaxPower                  int                              `json:"max_power"`
	CurrentPower              int                              `json:"current_power"`
	Weapons                   map[int]*WeaponSlot              `json:"weapons"`
	StandardSize              int                              `json:"standard_size"`
	ChassisType               string                           `json:"chassis_type"`
	WheelsPosNoScale          map[string]coordinate.Coordinate `json:"-"`
	WheelAnchors              map[string]Anchor                `json:"wheel_anchors"`
	Speed                     float64                          `json:"speed"`
	ReverseSpeed              float64                          `json:"-"`
	PowerFactor               float64                          `json:"power_factor"`
	ReverseFactor             float64                          `json:"-"`
	TurnSpeed                 float64                          `json:"turn_speed"`
	MoveDrag                  float64                          `json:"move_drag"`
	AngularDrag               float64                          `json:"-"`
	Weight                    float64                          `json:"-"`
	Attributes                map[string]int                   `json:"attributes"`
	PowerPoints               int                              `json:"power_points"`
}

type Anchor struct {
	XAnchor     float64 `json:"x_anchor"`
	YAnchor     float64 `json:"y_anchor"`
	RealXAttach int     `json:"real_x_attach"`
	RealYAttach int     `json:"real_y_attach"`
}

func (body *Body) GetName() string {
	return body.Name
}

func (body *Body) GetStandardSize() int {
	return body.StandardSize
}

func (body *Body) GetTypeSlot() int {
	return 0
}

func (body *Body) GetUsePower() int {
	var allPower int

	//for _, slot := range body.Equips {
	//	if slot.Equip != nil {
	//		allPower = allPower + slot.Equip.Power
	//	}
	//}

	//for _, slot := range body.Weapons {
	//	if slot.Weapon != nil {
	//		allPower = allPower + slot.Weapon.Power
	//	}
	//}

	return allPower
}

func (body *Body) SetWheelsPositions() {

	body.WheelAnchors = make(map[string]Anchor)

	for key, pos := range body.WheelsPosNoScale {
		xAnchor, yAnchor, realXAttach, realYAttach := game_math.GetAnchorWeapon(64, 64, int(pos.X), int(pos.Y))
		body.WheelAnchors[key] = Anchor{
			XAnchor:     xAnchor,
			YAnchor:     yAnchor,
			RealXAttach: realXAttach,
			RealYAttach: realYAttach,
		}
	}
}
