package ammo

type Ammo struct {
	ID                     int     `json:"id"`
	Name                   string  `json:"name"`
	Specification          string  `json:"specification"`
	Type                   string  `json:"type"`
	MinDamage              int     `json:"min_damage"`
	MaxDamage              int     `json:"max_damage"`
	AreaCovers             int     `json:"area_covers"`
	ChaseTarget            bool    `json:"chase_target"`
	ChaseOption            string  `json:"chase_option"`
	ChaseCatchDestination  int     `json:"chase_catch_destination"`
	Rotate                 float64 `json:"rotate"`
	BulletSpeed            int     `json:"bullet_speed"`
	StandardSize           int     `json:"standard_size"`
	DetonationTimeOut      int     `json:"detonation_time_out"`
	DetonationStartTimeOut int     `json:"detonation_start_time_out"`
	DetonationDistance     int     `json:"detonation_distance"`
	PushingPower           int     `json:"pushing_power"`
	Gravity                float64 `json:"-"`
	ForceAnimate           bool    `json:"force_animate"`
}

func (a *Ammo) GetName() string {
	return a.Name
}

func (a *Ammo) GetStandardSize() int {
	return a.StandardSize
}

func (a *Ammo) GetTypeSlot() int {
	return 0
}
