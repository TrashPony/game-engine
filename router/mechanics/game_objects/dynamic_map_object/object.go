package dynamic_map_object

import (
	"fmt"
	"github.com/TrashPony/game-engine/node/binary_msg"
	"github.com/TrashPony/game-engine/router/const/game_types"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/body"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/burst_of_shots"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/gunner"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/obstacle_point"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/physical_model"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/target"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/visible_objects"
	"math"
	"sync"
)

type Object struct {
	// TODO Сделать обьект Mover и встраивать его во все обьекты которые могут двигатся
	// везде где есть приставка Type это оригиналдьные данные типа, все остальное сформированые
	ID                  int    `json:"id"`
	TypeID              int    `json:"type_id"`
	Type                string `json:"type"`
	MapID               int    `json:"map_id"`
	Texture             string `json:"texture"`
	AnimateSpriteSheets string `json:"animate_sprite_sheets"`
	AnimateLoop         bool   `json:"animate_loop"`
	MaxHP               int    `json:"max_hp"`
	TypeMaxHP           int    `json:"-"`
	HP                  int    `json:"hp"`
	Scale               int    `json:"scale"`
	Shadow              bool   `json:"shadow"`
	AnimationSpeed      int    `json:"animation_speed"`
	Priority            int    `json:"priority"`
	TeamID              int    `json:"team_id"`

	TypeXShadowOffset int `json:"-"`
	XShadowOffset     int `json:"x_shadow_offset"`
	TypeYShadowOffset int `json:"-"`
	YShadowOffset     int `json:"y_shadow_offset"`
	ShadowIntensity   int `json:"shadow_intensity"`

	TypeGeoData []*obstacle_point.ObstaclePoint `json:"-"`
	HeightType  float64                         `json:"-"`

	Fraction string `json:"fraction"`

	/* постройка */
	OwnerID          int  `json:"owner_id"` // ид игрока владельца
	Static           bool `json:"-"`
	Weight           int  `json:"weight"`
	Immortal         bool `json:"-"`
	Build            bool `json:"build"`
	DestroyLeftTimer bool `json:"-"`
	DestroyTimer     int  `json:"-"`
	Work             bool `json:"work"`

	CacheJson      []byte `json:"-"`
	CacheGeoData   []byte `json:"-"`
	CreateJsonTime int64  `json:"-"`

	GrowTime int `json:"grow_time"` // говорит время цикла когда растение росло для гладкой отрисовки

	MemoryID int `json:"-"`

	WeaponID      int `json:"weapon_id"`
	XAttach       int `json:"x_attach"`
	YAttach       int `json:"y_attach"`
	MaxEnergy     int `json:"max_energy"`
	CurrentEnergy int `json:"current_energy"`
	ViewRange     int `json:"view_range"`

	Weapons map[int]*body.WeaponSlot `json:"weapons"`
	Run     bool                     `json:"-"`

	countUpdateWeaponTarget int
	visibleObjects          *visible_objects.VisibleObjectsStore
	physicalModel           *physical_model.PhysicalModel
	gunner                  *gunner.Gunner
	BurstOfShots            *burst_of_shots.BurstOfShots `json:"-"`
	weaponTarget            *target.Target
	ForceTarget             *target.Target `json:"-"`
	mx                      sync.RWMutex
}

func (o *Object) GetUpdateData(mapTime int64) []byte {
	// TODO оптимизация и совмещение с src/game_loop/game_loop_view/update_messages.go
	command := []byte{}
	command = append(command, binary_msg.GetIntBytes(o.HP)...)
	command = append(command, binary_msg.GetIntBytes(o.TeamID)...)
	return command
}

func (o *Object) CheckViewCoordinate(x, y, radius int) (bool, bool) {
	if o.GetRangeView()+radius >= int(game_math.GetBetweenDist(o.GetPhysicalModel().X, o.GetPhysicalModel().Y, x, y)) {
		return true, true
	}

	if o.GetRadarRange()+radius >= int(game_math.GetBetweenDist(o.GetPhysicalModel().X, o.GetPhysicalModel().Y, x, y)) {
		return false, true
	}

	return false, false
}

func (o *Object) GetRangeView() int {
	if o == nil {
		return 0
	}

	return o.ViewRange
}

func (o *Object) GetRadarRange() int {
	radarRange := 0

	if o == nil {
		return radarRange
	}

	return radarRange
}

func (o *Object) GetPower() int {
	return o.CurrentEnergy
}

func (o *Object) SetPower(power int) {
	o.CurrentEnergy = power
}

func (o *Object) GetGunner() *gunner.Gunner {
	if o.gunner == nil {
		o.gunner = &gunner.Gunner{GunUser: o}
	}

	return o.gunner
}

func (o *Object) GetPhysicalModel() *physical_model.PhysicalModel {
	if o.physicalModel == nil {
		o.initPhysicalModel()
	}

	return o.physicalModel
}

func (o *Object) initPhysicalModel() {

	weight := float64(o.Weight)
	if o.Static {
		weight = float64(math.MaxInt32)
	}

	o.physicalModel = &physical_model.PhysicalModel{
		WASD:        physical_model.WASD{},
		MoveDrag:    0.7,
		AngularDrag: 0.7,
		Type:        "object",
		ID:          o.ID,
		Weight:      weight,
	}
}

func (o *Object) UpdatePhysicalModel() {
	// todo обонление параметров типо изменилась скорость из за скила например
}

func (o *Object) GetBytes(mapTime int64) []byte {

	if o.CreateJsonTime == mapTime && len(o.CacheJson) > 0 {
		return o.CacheJson
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in GetJSON object", r)
		}
	}()

	command := []byte{}

	command = append(command, binary_msg.GetIntBytes(o.ID)...)
	command = append(command, binary_msg.GetIntBytes(o.GetX())...)
	command = append(command, binary_msg.GetIntBytes(o.GetY())...)
	command = append(command, binary_msg.GetIntBytes(int(o.GetRotate()))...)
	command = append(command, binary_msg.GetIntBytes(int(o.GetMapHeight()))...)

	command = append(command, binary_msg.GetIntBytes(o.HP)...)
	command = append(command, binary_msg.GetIntBytes(o.MaxHP)...)
	command = append(command, binary_msg.GetIntBytes(o.CurrentEnergy)...)
	command = append(command, binary_msg.GetIntBytes(o.MaxEnergy)...)
	command = append(command, binary_msg.GetIntBytes(o.ViewRange)...)

	command = append(command, binary_msg.GetIntBytes(o.XShadowOffset)...)
	command = append(command, binary_msg.GetIntBytes(o.YShadowOffset)...)
	command = append(command, binary_msg.GetIntBytes(o.ShadowIntensity)...)
	command = append(command, binary_msg.GetIntBytes(o.OwnerID)...)
	command = append(command, binary_msg.GetIntBytes(o.Priority)...)

	command = append(command, byte(o.TeamID))
	command = append(command, binary_msg.BoolToByte(o.Work))
	command = append(command, binary_msg.BoolToByte(o.Build))
	command = append(command, byte(o.Scale))
	command = append(command, binary_msg.BoolToByte(o.Shadow))
	command = append(command, binary_msg.BoolToByte(o.AnimateSpriteSheets != ""))
	command = append(command, byte(o.AnimationSpeed))
	command = append(command, binary_msg.BoolToByte(o.AnimateLoop))
	command = append(command, binary_msg.BoolToByte(o.Static))

	command = append(command, byte(len(o.Type)))
	command = append(command, []byte(o.Type)...)

	if o.AnimateSpriteSheets != "" {
		command = append(command, byte(len(o.AnimateSpriteSheets)))
		command = append(command, []byte(o.AnimateSpriteSheets)...)
	} else {
		command = append(command, byte(len(o.Texture)))
		command = append(command, []byte(o.Texture)...)
	}

	// weapons
	command = append(command, byte(len(o.Weapons)))
	for _, weaponSlot := range o.RangeWeaponSlots() {
		command = append(command, binary_msg.GetIntBytes(int(weaponSlot.GetGunRotate()))...)
		command = append(command, binary_msg.GetIntBytes(weaponSlot.GetRealXAttach())...)
		command = append(command, binary_msg.GetIntBytes(weaponSlot.GetRealYAttach())...)
		command = append(command, byte(weaponSlot.Number))
		command = append(command, byte(weaponSlot.GetXAnchor()*100))
		command = append(command, byte(weaponSlot.GetYAnchor()*100))

		if weaponSlot.Weapon != nil {
			command = append(command, byte(len([]byte(weaponSlot.Weapon.Name))))
			command = append(command, []byte(weaponSlot.Weapon.Name)...)
		} else {
			command = append(command, byte(0))
		}
	}

	// geo data
	command = append(command, o.GetGeoDataBin()...)

	o.CacheJson = command
	o.CreateJsonTime = mapTime

	return o.CacheJson
}

func (o *Object) GetGeoDataBin() []byte {
	return o.CacheGeoData
}

type Flore struct {
	TextureOverFlore string `json:"texture_over_flore"`
	TexturePriority  int    `json:"texture_priority"`
	X                int    `json:"x"`
	Y                int    `json:"y"`
}

func (o *Object) GetTurretAmmo() {
	for _, w := range o.RangeWeaponSlots() {
		if w != nil && w.Weapon != nil && w.Ammo == nil {
			o.GetWeaponSlot(1).SetAmmo(game_types.GetNewAmmo(w.Weapon.DefaultAmmoTypeID))
		}
	}
}

func (o *Object) GetTeamID() int {
	return o.TeamID
}

func (o *Object) CheckVisibleObjectStore() bool {
	return o.visibleObjects != nil
}

func (o *Object) SetVisibleObjectStore(v *visible_objects.VisibleObjectsStore) {
	o.visibleObjects = v
}

func (o *Object) checkVisibleObjectStore() {
	if o.visibleObjects == nil {
		o.visibleObjects = &visible_objects.VisibleObjectsStore{}
	}
}

func (o *Object) UnsafeRangeVisibleObjects() ([]*visible_objects.VisibleObject, *sync.RWMutex) {
	o.checkVisibleObjectStore()
	return o.visibleObjects.UnsafeRangeMapDynamicObjects()
}

func (o *Object) GetUpdateWeaponTarget() bool {
	if o.countUpdateWeaponTarget == 30 {
		o.countUpdateWeaponTarget = 0
		return true
	} else {
		o.countUpdateWeaponTarget++
		return false
	}
}
