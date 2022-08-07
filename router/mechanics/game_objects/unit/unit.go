package unit

import (
	"github.com/TrashPony/game-engine/node/binary_msg"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/behavior_rule"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/body"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/burst_of_shots"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/gunner"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/physical_model"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/target"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/visible_objects"
	"sync"
)

type Unit struct {
	ID             int        `json:"id"`
	OwnerID        int        `json:"owner_id"`
	MapID          int        `json:"map_id"`
	HP             int        `json:"hp"`
	Body           *body.Body `json:"body"`
	TeamID         int        `json:"team_id"`
	movePath       *MovePath  // специальный обьект для пути
	moveMx         sync.Mutex
	BehaviorRules  *behavior_rule.BehaviorRules `json:"-"`
	physicalModel  *physical_model.PhysicalModel
	gunner         *gunner.Gunner
	burstOfShots   *burst_of_shots.BurstOfShots
	weaponTarget   *target.Target
	CacheJson      []byte `json:"-"`
	CreateJsonTime int64  `json:"-"`
	visibleObjects *visible_objects.VisibleObjectsStore
	mx             sync.RWMutex
}

func (u *Unit) GetUpdateData(mapTime int64) []byte {
	// TODO оптимизация и совмещение с src/game_loop/game_loop_view/update_messages.go
	command := []byte{}
	command = append(command, binary_msg.GetIntBytes(u.HP)...)
	command = append(command, binary_msg.GetIntBytes(u.TeamID)...)
	command = append(command, binary_msg.GetIntBytes(u.GetRangeView())...)
	command = append(command, binary_msg.GetIntBytes(u.GetRangeRadar())...)
	return command
}

func (u *Unit) GetTeamID() int {
	return u.TeamID
}

func (u *Unit) GetOwnerID() int {
	return u.OwnerID
}

func (u *Unit) GetWeight() float64 {
	return u.GetPhysicalModel().GetWeight()
}

func (u *Unit) AddVelocity(x float64, y float64) {
	u.GetPhysicalModel().AddVelocity(x, y)
}

// ячейки которые отображены на панеле быстрого доступа игрока
type EquipSell struct {
	Number int    `json:"-"`
	Source string `json:"source"` // squadInventory, Constructor, "empty" - убрать все выделения

	// если body то эти параметры важны
	TypeSlot int `json:"type_slot"` // 0,1,2,3
	Slot     int `json:"slot"`

	// эти параметры берутся каждый тик по состоянию снаряжения дл фронтенда
	StartReload  int64 `json:"sr"`
	EndReload    int64 `json:"er"`
	AmmoQuantity int   `json:"aq"`
	On           bool  `json:"on"`
}

func (u *Unit) GetType() string {
	return "unit"
}

func (u *Unit) GetMapID() int {
	return u.MapID
}

func (u *Unit) SetMapID(mapID int) {
	u.MapID = mapID
}

func (u *Unit) GetX() int {
	return u.GetPhysicalModel().GetX()
}

func (u *Unit) GetY() int {
	return u.GetPhysicalModel().GetY()
}

func (u *Unit) GetID() int {
	return u.ID
}

func (u *Unit) GetBurstOfShots() *burst_of_shots.BurstOfShots {
	if u.burstOfShots == nil {
		u.burstOfShots = &burst_of_shots.BurstOfShots{}
	}

	return u.burstOfShots
}

func (u *Unit) GetGunner() *gunner.Gunner {
	if u.gunner == nil {
		u.gunner = &gunner.Gunner{GunUser: u}
	}

	return u.gunner
}

func (u *Unit) GetBytes(mapTime int64) []byte {

	if u.CreateJsonTime == mapTime && len(u.CacheJson) > 0 {
		return u.CacheJson
	}

	command := []byte{}

	command = append(command, binary_msg.GetIntBytes(u.GetID())...)
	command = append(command, binary_msg.GetIntBytes(u.OwnerID)...)
	command = append(command, binary_msg.GetIntBytes(u.HP)...)
	command = append(command, byte(u.Body.ID))
	command = append(command, byte(u.TeamID))

	// position data
	command = append(command, binary_msg.GetIntBytes(u.GetX())...)
	command = append(command, binary_msg.GetIntBytes(u.GetY())...)
	command = append(command, binary_msg.GetIntBytes(int(u.GetRotate()))...)

	command = append(command, binary_msg.GetIntBytes(u.Body.MaxHP)...)

	command = append(command, byte(len([]byte(u.Body.Texture))))
	command = append(command, []byte(u.Body.Texture)...)

	// weapon data
	for _, unitWeaponSlot := range u.RangeWeaponSlots() {

		command = append(command, byte(unitWeaponSlot.Number))
		command = append(command, binary_msg.GetIntBytes((unitWeaponSlot.GetRealXAttach()))...)
		command = append(command, binary_msg.GetIntBytes((unitWeaponSlot.GetRealYAttach()))...)
		command = append(command, byte(unitWeaponSlot.GetXAnchor()*100))
		command = append(command, byte(unitWeaponSlot.GetYAnchor()*100))
		command = append(command, binary_msg.GetIntBytes(int(unitWeaponSlot.GetGunRotate()))...)

		if unitWeaponSlot.Weapon != nil {
			command = append(command, byte(len([]byte(unitWeaponSlot.Weapon.WeaponTexture))))
			command = append(command, []byte(unitWeaponSlot.Weapon.WeaponTexture)...)
		} else {
			command = append(command, byte(0))
		}
	}

	return command
}

func (u *Unit) GetWeaponSlot(slotNumber int) *body.WeaponSlot {
	return u.Body.Weapons[slotNumber]
}

func (u *Unit) RangeWeaponSlots() map[int]*body.WeaponSlot {
	// мы никогда не пишут в карту слотов оружия поэтому этот метод безопасен (по крайне мере пока)
	return u.Body.Weapons
}

func (u *Unit) GetMapHeight() float64 {
	return u.GetPhysicalModel().GetHeight()
}

func (u *Unit) GetRotate() float64 {
	return u.GetPhysicalModel().GetRotate()
}

func (u *Unit) GetScale() int {
	return u.Body.Scale
}

func (u *Unit) GetWeaponTarget() *target.Target {
	return u.weaponTarget
}

func (u *Unit) SetWeaponTarget(target *target.Target) {
	u.weaponTarget = target
}

func (u *Unit) CheckViewCoordinate(x, y, radius int) (bool, bool) {

	if u.Body == nil || u.HP <= 0 {
		return false, false
	}

	if u.GetRangeView()+radius >= int(game_math.GetBetweenDist(u.GetX(), u.GetY(), x, y)) {
		return true, true
	}

	if u.GetRangeRadar()+radius >= int(game_math.GetBetweenDist(u.GetX(), u.GetY(), x, y)) {
		return false, true
	}

	return false, false
}

func (u *Unit) SetDamage(damage int) int {
	u.HP = u.HP - damage
	return damage
}

func (u *Unit) GetMaxHP() int {
	return u.Body.MaxHP
}

func (u *Unit) GetMoveMaxPower() float64 {
	return u.Body.Speed
}

func (u *Unit) GetMaxReverse() float64 {
	return u.Body.ReverseSpeed
}

func (u *Unit) GetPowerFactor() float64 {
	return u.Body.PowerFactor
}

func (u *Unit) GetReverseFactor() float64 {
	return u.Body.ReverseFactor
}

func (u *Unit) GetTurnSpeed() float64 {
	return u.Body.TurnSpeed
}

func (u *Unit) GetRangeView() int {
	return u.Body.RangeView
}
func (u *Unit) GetRangeRadar() int {
	return u.Body.RangeRadar
}

func (u *Unit) SetVisibleObjectStore(v *visible_objects.VisibleObjectsStore) {
	u.visibleObjects = v
}

func (u *Unit) checkVisibleObjectStore() {
	if u.visibleObjects == nil {
		u.visibleObjects = &visible_objects.VisibleObjectsStore{}
	}
}

func (u *Unit) UnsafeRangeVisibleObjects() ([]*visible_objects.VisibleObject, *sync.RWMutex) {
	u.checkVisibleObjectStore()
	return u.visibleObjects.UnsafeRangeMapDynamicObjects()
}
