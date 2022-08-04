package dynamic_map_object

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/body"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/burst_of_shots"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/target"
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

func (o *Object) GetID() int {
	return o.ID
}

func (o *Object) GetWeaponSlot(slotNumber int) *body.WeaponSlot {
	return o.Weapons[slotNumber]
}

func (o *Object) RangeWeaponSlots() map[int]*body.WeaponSlot {
	// мы никогда не пишут в карту слотов оружия поэтому этот метод безопасен (по крайне мере пока)
	return o.Weapons
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
	if o.burstOfShots == nil {
		o.burstOfShots = &burst_of_shots.BurstOfShots{}
	}

	return o.burstOfShots
}
