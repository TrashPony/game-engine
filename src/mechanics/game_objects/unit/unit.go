package unit

import (
	"encoding/json"
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/const/game_types"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/behavior_rule"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/body"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/burst_of_shots"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/coordinate"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/gunner"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/physical_model"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/target"
)

type Unit struct {
	ID              int                          `json:"id"`
	OwnerID         int                          `json:"owner_id"`
	MapID           int                          `json:"map_id"`
	PositionData    interface{}                  `json:"position_data"` // описание поции юнита для отображения на фронте
	WeaponSlotsData interface{}                  `json:"weapon_slots_data"`
	HP              int                          `json:"hp"`
	Body            *body.Body                   `json:"body"`
	MovePath        *MovePath                    `json:"-"` // специальный обьект для пути
	BehaviorRules   *behavior_rule.BehaviorRules `json:"-"`
	physicalModel   *physical_model.PhysicalModel
	gunner          *gunner.Gunner
	BurstOfShots    *burst_of_shots.BurstOfShots `json:"-"`
	weaponTarget    *target.Target
}

func (u *Unit) GetX() int {
	return u.GetPhysicalModel().GetX()
}

func (u *Unit) GetY() int {
	return u.GetPhysicalModel().GetY()
}

type MovePath struct {
	Path         *[]*coordinate.Coordinate
	FollowTarget *target.Target
	CurrentPoint int
	NeedFindPath bool
}

func (u *Unit) GetID() int {
	return u.ID
}

func (u *Unit) GetBurstOfShots() *burst_of_shots.BurstOfShots {
	if u.BurstOfShots == nil {
		u.BurstOfShots = &burst_of_shots.BurstOfShots{}
	}

	return u.BurstOfShots
}

func (u *Unit) GetGunner() *gunner.Gunner {
	if u.gunner == nil {
		u.gunner = &gunner.Gunner{GunUser: u}
	}

	return u.gunner
}

func (u *Unit) GetPhysicalModel() *physical_model.PhysicalModel {
	if u.physicalModel == nil {
		u.initPhysicalModel()
	}

	return u.physicalModel
}

func (u *Unit) initPhysicalModel() {
	// todo тестовые параметры
	u.physicalModel = &physical_model.PhysicalModel{
		X:             100,
		Y:             100,
		Rotate:        0,
		Speed:         90 / _const.ServerTickSecPart,
		ReverseSpeed:  22.5 / _const.ServerTickSecPart,
		PowerFactor:   15 / _const.ServerTickSecPart,
		ReverseFactor: 1.5 / _const.ServerTickSecPart,
		TurnSpeed:     0.3 / _const.ServerTickSecPart,
		WASD:          physical_model.WASD{},
		MoveDrag:      0.65,
		AngularDrag:   0.70,
		Weight:        20000,
	}

	// todo тестовые параметры
	u.Body = game_types.GetNewBody(1)
	newWeapon1 := game_types.WeaponTypes[1]
	newAmmo1 := game_types.AmmoTypes[3]
	u.Body.Weapons[1].Weapon = &newWeapon1
	u.Body.Weapons[1].Ammo = &newAmmo1

	newWeapon2 := game_types.WeaponTypes[2]
	newAmmo2 := game_types.AmmoTypes[2]
	u.Body.Weapons[2].Weapon = &newWeapon2
	u.Body.Weapons[2].Ammo = &newAmmo2

	// распологаем оружие
	u.GetWeaponSlot(1).SetAnchor()
	u.GetWeaponSlot(2).SetAnchor()

	// применяем настройки размера к обьектам колизий
	sizeOffset := float64(u.Body.Scale) / 100
	u.physicalModel.Height = float64(u.Body.Height) * sizeOffset
	u.physicalModel.Length = float64(u.Body.Length) * sizeOffset
	u.physicalModel.Width = float64(u.Body.Width) * sizeOffset
	u.physicalModel.Radius = int(float64(u.Body.Radius) * sizeOffset)
}

func (u *Unit) UpdatePhysicalModel() {
	// todo обонление параметров типо изменилась скорость из за скила/повреждений например
}

func (u *Unit) GetCopyPhysicalModel() *physical_model.PhysicalModel {
	pm := *u.physicalModel
	return &pm
}

func (u *Unit) GetJSON(mapTime int64) string {

	// todo наверн не очень производительно это но пох
	weaponData := make([]interface{}, 0)
	for unitWeaponSlot := range u.RangeWeaponSlots() {
		if unitWeaponSlot != nil {
			weaponData = append(weaponData, body.WeaponSlot{
				Number:      unitWeaponSlot.Number,
				Weapon:      unitWeaponSlot.Weapon,
				XAttach:     unitWeaponSlot.XAttach,
				YAttach:     unitWeaponSlot.YAttach,
				RealXAttach: unitWeaponSlot.GetRealXAttach(),
				RealYAttach: unitWeaponSlot.GetRealYAttach(),
				XAnchor:     unitWeaponSlot.GetXAnchor(),
				YAnchor:     unitWeaponSlot.GetYAnchor(),
				GunRotate:   unitWeaponSlot.GetGunRotate(),
			})
		}
	}

	u.PositionData = u.physicalModel
	u.WeaponSlotsData = weaponData
	jsonShortUnit, err := json.Marshal(u)
	if err != nil {
		println("unit to json: ", err.Error())
	}

	return string(jsonShortUnit)
}

func (u *Unit) GetWeaponSlot(slotNumber int) *body.WeaponSlot {
	return u.Body.Weapons[slotNumber]
}

func (u *Unit) RangeWeaponSlots() <-chan *body.WeaponSlot {
	// мы никогда не пишут в карту слотов оружия поэтому этот метод безопасен (по крайне мере пока)
	slots := make(chan *body.WeaponSlot, len(u.Body.Weapons))

	go func() {
		defer func() {
			close(slots)
		}()

		for _, weaponSlot := range u.Body.Weapons {
			slots <- weaponSlot
		}
	}()

	return slots
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
