package dynamic_map_object

import (
	"encoding/json"
	"fmt"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/body"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/burst_of_shots"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/gunner"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/obstacle_point"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/physical_model"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/target"
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
	UnitOverlap         bool   `json:"unit_overlap"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	MaxHP               int    `json:"max_hp"`
	TypeMaxHP           int    `json:"-"`
	HP                  int    `json:"hp"`
	Scale               int    `json:"scale"`
	MaxScale            int    `json:"-"` // определяется рандомно для растений максимальный размер куста
	Shadow              bool   `json:"shadow"`
	AnimationSpeed      int    `json:"animation_speed"`
	Priority            int    `json:"priority"`

	TypeXShadowOffset int `json:"-"`
	XShadowOffset     int `json:"x_shadow_offset"`
	TypeYShadowOffset int `json:"-"`
	YShadowOffset     int `json:"y_shadow_offset"`
	ShadowIntensity   int `json:"shadow_intensity"`

	TypeGeoData []*obstacle_point.ObstaclePoint `json:"-"`
	HeightType  float64                         `json:"-"`

	Fraction string `json:"fraction"`

	/* постройка */
	OwnerID int  `json:"owner_id"` // ид игрока владельца
	Static  bool `json:"-"`

	NoAnchor bool `json:"-"` // обьект может передвигатся если его например пнуть

	DestroyLeftTimer bool `json:"-"`
	DestroyTimer     int  `json:"-"`

	CacheJson      string `json:"-"`
	CreateJsonTime int64  `json:"-"`

	GrowCycle    int `json:"-"`
	GrowLeftTime int `json:"-"`
	GrowTime     int `json:"grow_time"` // говорит время цикла когда растение росло для гладкой отрисовки

	MemoryUUID string `json:"-"`

	PositionData  interface{}              `json:"position_data"` // описание поции юнита для отображения на фронте
	Weapons       map[int]*body.WeaponSlot `json:"-"`
	physicalModel *physical_model.PhysicalModel
	gunner        *gunner.Gunner
	BurstOfShots  *burst_of_shots.BurstOfShots `json:"-"`
	weaponTarget  *target.Target
	mx            sync.RWMutex
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
	// todo тестовые параметры
	o.physicalModel = &physical_model.PhysicalModel{
		WASD:        physical_model.WASD{},
		MoveDrag:    0.9,
		AngularDrag: 0.9,
		Weight:      20000,
		Height:      o.HeightType,
	}
}

func (o *Object) UpdatePhysicalModel() {
	// todo обонление параметров типо изменилась скорость из за скила например
}

func (o *Object) IsFly() bool {
	return false
}

func (o *Object) CheckGrowthPower() bool {
	return false
}

func (o *Object) CheckGrowthRevers() bool {
	return false
}

func (o *Object) CheckLeftRotate() bool {
	return false
}

func (o *Object) CheckRightRotate() bool {
	return false
}

func (o *Object) GetJSON(mapTime int64) string {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in GetJSON object", r)
		}
	}()

	if o.CreateJsonTime == mapTime && o.CacheJson != "" {
		return o.CacheJson
	}

	o.mx.RLock()
	defer o.mx.RUnlock()

	o.PositionData = o.physicalModel
	jsonObject, err := json.Marshal(struct {
		Obj     *Object  `json:"obj"`
		GeoData []string `json:"geo_data"`
	}{
		Obj:     o,
		GeoData: o.GetGeoDataJSON(),
	})
	if err != nil {
		println("dyn.object to json: ", err.Error())
	}

	o.CacheJson = string(jsonObject)
	o.CreateJsonTime = mapTime

	return o.CacheJson
}

func (o *Object) GetGeoDataJSON() []string {
	geoDataJSON := make([]string, 0)

	for _, geoPoint := range o.physicalModel.GeoData {
		geoDataJSON = append(geoDataJSON, geoPoint.GetJSON())
	}

	return geoDataJSON
}

type Flore struct {
	TextureOverFlore string `json:"texture_over_flore"`
	TexturePriority  int    `json:"texture_priority"`
	X                int    `json:"x"`
	Y                int    `json:"y"`
}

func (o *Object) GetGrowTime() int {
	return o.GrowTime
}
