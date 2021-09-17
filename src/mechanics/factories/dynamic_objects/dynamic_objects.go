package dynamic_objects

import (
	"github.com/TrashPony/game_engine/src/mechanics/db/objects"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/dynamic_map_object"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/obstacle_point"
	"github.com/getlantern/deepcopy"
	"math/rand"
	"sync"
)

var DynamicObjects = newObjectsStore()

type store struct {
	mx      sync.RWMutex
	objects map[int]*dynamic_map_object.Object
	flores  map[int]*dynamic_map_object.Flore
}

func newObjectsStore() *store {
	object := objects.AllTypeCoordinate()

	mapObjects := make(map[int]*dynamic_map_object.Object)
	mapFlores := make(map[int]*dynamic_map_object.Flore)

	for _, obj := range object {
		mapObjects[obj.TypeID] = obj
	}

	return &store{
		objects: mapObjects,
		flores:  mapFlores,
	}
}

func (d *store) GetDynamicObjectByID(id int, rotate float64) *dynamic_map_object.Object {

	d.mx.RLock()
	factoryObj, ok := d.objects[id]
	d.mx.RUnlock()

	if ok {
		return getCopyObject(factoryObj, rotate)
	} else {
		return nil
	}
}

func (d *store) GetDynamicObjectByType(typeObj string) *dynamic_map_object.Object {
	d.mx.RLock()
	objTypeID := 0
	for _, factoryObj := range d.objects {
		if factoryObj.Type == typeObj {
			objTypeID = factoryObj.TypeID
		}
	}
	d.mx.RUnlock()

	return d.GetDynamicObjectByID(objTypeID, 0)
}

func (d *store) GetDynamicObjectByTexture(name string, rotate float64) *dynamic_map_object.Object {
	d.mx.RLock()
	defer d.mx.RUnlock()

	for _, factoryObj := range d.objects {
		if factoryObj.Texture == name {
			return getCopyObject(factoryObj, rotate)
		}
	}

	return nil
}

func (d *store) AddNewGeoData(x, y, radius, idType int, move bool) {
	obj, ok := d.objects[idType]
	if !ok {
		return
	}

	obj.TypeGeoData = append(obj.TypeGeoData, &obstacle_point.ObstaclePoint{X: int32(x), Y: int32(y), Radius: int32(radius), Move: move})
}

func (d *store) GetFloreByID(id int) *dynamic_map_object.Flore {
	//todo уровень костылезации - МАКСИМАЛЬНЫЙ

	var flores = map[int]string{
		0:  "arctic",
		1:  "arctic_2",
		2:  "desertDunes",
		3:  "desertDunes_2",
		4:  "grass_1",
		5:  "grass_2",
		6:  "grass_3",
		7:  "soil",
		8:  "soil_2",
		9:  "tundra",
		10: "tundra_2",
		11: "xenos",
		12: "xenos_2",
		13: "water_1",
	}

	return &dynamic_map_object.Flore{TextureOverFlore: flores[id]}
}

func (d *store) GetFloreByName(name string) *dynamic_map_object.Flore {
	return &dynamic_map_object.Flore{TextureOverFlore: name}
}

func (d *store) GetRandomWreckage(typeWreckage string) *dynamic_map_object.Object {

	d.mx.RLock()
	defer d.mx.RUnlock()

	unitWreckage := make([]*dynamic_map_object.Object, 0)

	for _, obj := range d.objects {
		if obj.Type == typeWreckage {
			unitWreckage = append(unitWreckage, obj)
		}
	}

	wreckage := unitWreckage[rand.Intn(len(unitWreckage))]

	return getCopyObject(wreckage, float64(rand.Intn(360)))
}

func (d *store) GetObjectWreckage(typeWreckage string, objectID int) *dynamic_map_object.Object {
	d.mx.RLock()
	defer d.mx.RUnlock()

	// TODO костыль

	var wreckages = map[int]string{
		173: "extractor_damaged",
		174: "repair_station_damaged",
		97:  "laser_turret_damaged",
		98:  "shield_generator_damaged",
		99:  "tank_turret_damaged",
		100: "artillery_turret_damaged",
		101: "energy_generator_damaged",
		102: "jammer_generator_damaged",
		103: "missile_defense_damaged",
		105: "radar_damaged",
		107: "beacon_damaged",
		106: "storage_damaged",
	}

	textureName := wreckages[objectID]

	for _, obj := range d.objects {
		if obj.Type == typeWreckage && obj.Texture == textureName {
			return getCopyObject(obj, float64(rand.Intn(360)))
		}
	}

	return nil
}

func getCopyObject(factoryObj *dynamic_map_object.Object, rotate float64) *dynamic_map_object.Object {
	var newObj dynamic_map_object.Object

	err := deepcopy.Copy(&newObj, &factoryObj)
	if err != nil {
		println(err.Error())
	}

	newObj.GetPhysicalModel().Rotate = rotate
	newObj.GetStartScale()

	newObj.TypeYShadowOffset = factoryObj.TypeYShadowOffset
	newObj.TypeXShadowOffset = factoryObj.TypeXShadowOffset
	newObj.TypeMaxHP = factoryObj.TypeMaxHP
	newObj.HeightType = factoryObj.HeightType

	if len(factoryObj.TypeGeoData) > 0 {
		err = deepcopy.Copy(&newObj.TypeGeoData, &factoryObj.TypeGeoData)
		if err != nil {
			newObj.TypeGeoData = make([]*obstacle_point.ObstaclePoint, 0)
		}
	}

	return &newObj
}
