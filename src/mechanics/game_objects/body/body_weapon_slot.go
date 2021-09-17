package body

import (
	"encoding/json"
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/ammo"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/weapon"
)

type WeaponSlot struct {
	Number           int            `json:"number"`
	Weapon           *weapon.Weapon `json:"weapon"`
	Ammo             *ammo.Ammo     `json:"ammo"`
	AmmoQuantity     int            `json:"ammo_quantity"`
	XAttach          int            `json:"x_attach"`
	YAttach          int            `json:"y_attach"`
	RealXAttach      int            `json:"real_x_attach"`
	RealYAttach      int            `json:"real_y_attach"`
	XAnchor          float64        `json:"x_anchor"`
	YAnchor          float64        `json:"y_anchor"`
	Reload           bool           `json:"reload"`
	GunRotate        float64        `json:"gun_rotate"`
	TimeReload       int            `json:"time_reload"`
	CurrentReload    int            `json:"current_reload"`
	lastFirePosition int
}

func (s *WeaponSlot) StartReload(reloadTime int) {
	s.SetReload(true)
	s.setTimeReload(reloadTime)
	s.SetCurrentReload(reloadTime)
}

func (s *WeaponSlot) StartAmmoReload() {
	s.SetReload(true)
	s.setTimeReload(s.Weapon.ReloadAmmoTime)
	s.SetCurrentReload(s.Weapon.ReloadAmmoTime)
}

func (s *WeaponSlot) GetReload() bool {
	return s.Reload
}

func (s *WeaponSlot) SetReload(reload bool) {
	s.Reload = reload
}

func (s *WeaponSlot) GetAmmoQuantity() int {
	return s.AmmoQuantity
}

func (s *WeaponSlot) SetAmmoQuantity(quantity int) {
	s.AmmoQuantity = quantity
}

func (s *WeaponSlot) SetAnchor() {
	if s == nil || s.Weapon == nil {
		return
	}

	s.XAnchor, s.YAnchor, s.RealXAttach, s.RealYAttach = game_math.GetAnchorWeapon(s.Weapon.XAttach, s.Weapon.YAttach, s.XAttach, s.YAttach)
}

func (s *WeaponSlot) GetAmmo() *ammo.Ammo {
	return s.Ammo
}

func (s *WeaponSlot) SetAmmo(ammo *ammo.Ammo) {
	s.Ammo = ammo
}

func (s *WeaponSlot) GetGunRotate() float64 {
	return s.GunRotate
}

func (s *WeaponSlot) SetGunRotate(rotate float64) {
	s.GunRotate = rotate
}

func (s *WeaponSlot) getTimeReload() int {
	return s.TimeReload
}

func (s *WeaponSlot) setTimeReload(time int) {
	s.TimeReload = time
}

func (s *WeaponSlot) GetCurrentReload() int {
	return s.CurrentReload
}

func (s *WeaponSlot) SetCurrentReload(time int) {
	if time <= 0 {
		s.SetReload(false)
	}
	s.CurrentReload = time
}

func (s *WeaponSlot) GetXAnchor() float64 {
	return s.XAnchor
}

func (s *WeaponSlot) GetYAnchor() float64 {
	return s.YAnchor
}

func (s *WeaponSlot) GetRealXAttach() int {
	return s.RealXAttach
}

func (s *WeaponSlot) GetRealYAttach() int {
	return s.RealYAttach
}

func (s *WeaponSlot) GetJSON() string {
	jsonSlot, err := json.Marshal(s)
	if err != nil {
		println("weapon Slot to json: ", err.Error())
	}

	return string(jsonSlot)
}

func (s *WeaponSlot) GetCopy() *WeaponSlot {
	copySlot := *s
	return &copySlot
}

func (s *WeaponSlot) NextLastFirePosition() {
	if s.Weapon == nil {
		return
	}

	s.lastFirePosition++
	if len(s.Weapon.FirePositions) == s.lastFirePosition {
		s.lastFirePosition = 0
	}
}

func (s *WeaponSlot) GetLastFirePosition() int {
	return s.lastFirePosition
}
