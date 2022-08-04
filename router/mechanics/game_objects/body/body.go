package body

type Body struct {
	ID            int                 `json:"id"`
	Name          string              `json:"name"`
	Texture       string              `json:"texture"`
	MaxHP         int                 `json:"max_hp"`
	Scale         int                 `json:"scale"`
	Length        int                 `json:"length"`
	Width         int                 `json:"width"`
	Height        int                 `json:"height"`
	Radius        int                 `json:"radius"`
	RangeView     int                 `json:"range_view"`
	RangeRadar    int                 `json:"range_radar"`
	Weapons       map[int]*WeaponSlot `json:"weapons"`
	ChassisType   string              `json:"chassis_type"`
	Speed         float64             `json:"speed"`
	ReverseSpeed  float64             `json:"-"`
	PowerFactor   float64             `json:"power_factor"`
	ReverseFactor float64             `json:"-"`
	TurnSpeed     float64             `json:"turn_speed"`
	MoveDrag      float64             `json:"move_drag"`
	AngularDrag   float64             `json:"-"`
	Weight        float64             `json:"-"`
}
