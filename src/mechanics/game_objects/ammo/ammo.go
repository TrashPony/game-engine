package ammo

type Ammo struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Specification string  `json:"specification"`
	Type          string  `json:"type"`
	MinDamage     int     `json:"min_damage"`
	MaxDamage     int     `json:"max_damage"`
	AreaCovers    int     `json:"area_covers"`
	ChaseTarget   bool    `json:"chase_target"`
	ChaseOption   string  `json:"chase_option"`
	Rotate        float64 `json:"rotate"`
	BulletSpeed   int     `json:"bullet_speed"`
}

func (a *Ammo) GetName() string {
	return a.Name
}

func (a *Ammo) GetItemType() string {
	return ""
}

func (a *Ammo) GetItemName() string {
	return ""
}

func (a *Ammo) GetInMap() bool {
	return false
}

func (a *Ammo) GetIcon() string {
	return ""
}

func (a *Ammo) GetTypeSlot() int {
	return 0
}
