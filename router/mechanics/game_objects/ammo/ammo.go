package ammo

type Ammo struct {
	ID                    int     `json:"id"`
	Name                  string  `json:"name"`
	Type                  string  `json:"type"`
	MinDamage             int     `json:"min_damage"`
	MaxDamage             int     `json:"max_damage"`
	AreaCovers            int     `json:"area_covers"`
	ChaseTarget           bool    `json:"chase_target"`
	ChaseOption           string  `json:"chase_option"`
	ChaseCatchDestination int     `json:"chase_catch_destination"`
	Rotate                float64 `json:"rotate"`
	BulletSpeed           int     `json:"bullet_speed"`
	PushingPower          int     `json:"pushing_power"`
	Gravity               float64 `json:"-"`
}
