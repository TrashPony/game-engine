package weapon

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/ammo"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/coordinate"
)

type Weapon struct {
	ID                     int                      `json:"id"`
	Name                   string                   `json:"name"`
	OnlyConfig             bool                     `json:"only_config"`
	Specification          string                   `json:"specification"`
	MinRange               int                      `json:"min_range"`
	MaxRange               int                      `json:"max_range"`
	ApproximateFiringRange int                      `json:"approximate_firing_range"`
	MinAngle               int                      `json:"min_angle"`
	MaxAngle               int                      `json:"max_angle"`
	Accuracy               int                      `json:"accuracy"`
	AmmoCapacity           int                      `json:"ammo_capacity"`
	Type                   string                   `json:"type"` /* firearms, missile_weapon, laser_weapon */
	XAttach                int                      `json:"x_attach"`
	YAttach                int                      `json:"y_attach"`
	RotateSpeed            int                      `json:"rotate_speed"`
	CountFireBullet        int                      `json:"count_fire_bullet"`
	BulletSpeed            int                      `json:"bullet_speed"`
	ReloadAmmoTime         int                      `json:"reload_ammo_time"`
	ReloadTime             int                      `json:"reload_time"`
	DelayFollowingFire     int                      `json:"delay_following_fire"`
	FirePositions          []*coordinate.Coordinate `json:"fire_positions"`
	DamageModifier         float64                  `json:"damage_modifier"`
	WeaponTexture          string                   `json:"weapon_texture"`
	Artillery              bool                     `json:"artillery"`
	StandardSize           int                      `json:"standard_size"`
	DefaultAmmoTypeID      int                      `json:"default_ammo_type_id"`
	Unit                   bool                     `json:"unit"`
	//Power              int                      `json:"power"`

	// оружие накопительного типа (зажатая мышка коит энергию, отпускание выпускает снаряд)
	AccumulationFirePower   bool    `json:"accumulation_fire_power"`
	AccumulationFull        float64 `json:"accumulation_full"`
	AccumulationFullTimeOut int     `json:"accumulation_full_time_out"` // время в мс когда произойдет перегрев или автовыстрел
	PowerPoints             int     `json:"power_points"`
}

func (w *Weapon) GetName() string {
	return w.Name
}

func (w *Weapon) GetStandardSize() int {
	return w.StandardSize
}

func (w *Weapon) GetTypeSlot() int {
	return 0
}

func (w *Weapon) GetWeaponMaxRange(ammo *ammo.Ammo, lvlMap, mapHeight float64, realBallistic bool) (int, float64) {
	if (w.Type == "laser" || w.Type == "missile") || (w.MaxRange > 0 && !realBallistic) {
		return w.MaxRange, float64(w.MaxAngle)
	}

	bulletSpeed := ammo.BulletSpeed + w.BulletSpeed

	angle := 0.0

	if w.Artillery {
		angle = float64(w.MinAngle)
	} else {
		angle = float64(w.MaxAngle)
	}

	maxRange := game_math.GetMaxDistBallisticWeapon(
		float64(bulletSpeed),
		lvlMap,
		lvlMap+mapHeight,
		game_math.DegToRadian(angle),
		w.Type,
		ammo.Gravity,
	)

	return int(maxRange), angle
}
