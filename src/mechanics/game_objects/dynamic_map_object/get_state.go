package dynamic_map_object

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/body"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/burst_of_shots"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/target"
)

func (o *Object) GetOwnerID() int {
	return o.OwnerID
}

func (o *Object) GetHP() int {
	return o.HP
}

func (o *Object) GetScale() int {
	return o.Scale
}

// у некоторых обьектов например стуктур для строительства статичный размер
func (o *Object) GetStartScale() {

	if o.Texture == "explores_antenna" {
		o.SetScale(50)
		return
	}

	if o.Texture == "explores_observatory" {
		o.SetScale(50)
		return
	}

	if o.Texture == "unknown_civilization_jammer" {
		o.SetScale(150)
		return
	}

	if o.Texture == "replic_gauss_gun" {
		o.SetScale(70)
		return
	}

	if o.Type == "turret" {
		o.SetScale(25)
	}

	if o.Texture == "shield_generator" {
		o.SetScale(25)
	}

	if o.Texture == "extractor" {
		o.SetScale(25)
	}

	if o.Texture == "energy_generator" {
		o.SetScale(37)
	}

	if o.Texture == "jammer_generator" {
		o.SetScale(30)
	}

	if o.Texture == "missile_defense" {
		o.SetScale(17)
	}

	if o.Texture == "meteorite_defense" {
		o.SetScale(25)
	}

	if o.Texture == "radar" {
		o.SetScale(37)
	}

	if o.Texture == "storage" {
		o.SetScale(37)
	}

	if o.Texture == "beacon" {
		o.SetScale(37)
	}

	if o.Texture == "repair_station" {
		o.SetScale(25)
	}
}

func (o *Object) GetID() int {
	return o.ID
}

func (o *Object) GetType() string {
	return "object"
}

func (o *Object) GetMapID() int {
	return o.MapID
}

func (o *Object) GetWeaponSlot(slotNumber int) *body.WeaponSlot {
	return o.Weapons[slotNumber]
}

func (o *Object) RangeWeaponSlots() <-chan *body.WeaponSlot {
	// мы никогда не пишут в карту слотов оружия поэтому этот метод безопасен (по крайне мере пока)
	slots := make(chan *body.WeaponSlot, len(o.Weapons))

	go func() {
		defer func() {
			close(slots)
		}()

		for _, weaponSlot := range o.Weapons {
			slots <- weaponSlot
		}
	}()

	return slots
}

func (o *Object) GetMapHeight() float64 {
	return o.GetPhysicalModel().GetHeight()
}

func (o *Object) GetRotate() float64 {
	return o.GetPhysicalModel().GetRotate()
}

func (o *Object) GetWeaponTarget() *target.Target {
	return o.weaponTarget
}

func (o *Object) SetWeaponTarget(target *target.Target) {
	o.weaponTarget = target
}

func (o *Object) GetX() int {
	return o.GetPhysicalModel().GetX()
}

func (o *Object) GetY() int {
	return o.GetPhysicalModel().GetY()
}

func (o *Object) GetBurstOfShots() *burst_of_shots.BurstOfShots {
	if o.BurstOfShots == nil {
		o.BurstOfShots = &burst_of_shots.BurstOfShots{}
	}

	return o.BurstOfShots
}
