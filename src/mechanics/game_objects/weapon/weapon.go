package weapon

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/coordinate"
)

type Weapon struct {
	ID                 int                      `json:"id"`
	Name               string                   `json:"name"`
	Specification      string                   `json:"specification"`
	MinRange           int                      `json:"min_range"`
	MaxRange           int                      `json:"max_range"`
	MinAngle           int                      `json:"min_angle"`
	MaxAngle           int                      `json:"max_angle"`
	Accuracy           int                      `json:"accuracy"`
	AmmoCapacity       int                      `json:"ammo_capacity"`
	Type               string                   `json:"type"` /* firearms, missile_weapon, laser_weapon */
	XAttach            int                      `json:"x_attach"`
	YAttach            int                      `json:"y_attach"`
	RotateSpeed        int                      `json:"rotate_speed"`
	CountFireBullet    int                      `json:"count_fire_bullet"`
	BulletSpeed        int                      `json:"bullet_speed"`
	ReloadAmmoTime     int                      `json:"reload_ammo_time"`
	ReloadTime         int                      `json:"reload_time"`
	DelayFollowingFire int                      `json:"delay_following_fire"`
	FirePositions      []*coordinate.Coordinate `json:"fire_positions"`
	DamageModifier     float64                  `json:"damage_modifier"`
	WeaponTexture      string                   `json:"weapon_texture"`
	Artillery          bool                     `json:"artillery"`
}

func (w *Weapon) GetName() string {
	return w.Name
}

func (w *Weapon) GetItemType() string {
	return ""
}

func (w *Weapon) GetItemName() string {
	return ""
}

func (w *Weapon) GetInMap() bool {
	return false
}

func (w *Weapon) GetIcon() string {
	return ""
}

func (w *Weapon) GetTypeSlot() int {
	return 0
}
