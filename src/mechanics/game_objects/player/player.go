package player

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/dynamic_map_object"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/visible_objects"
	"strconv"
	"sync"
)

type Player struct {
	ID                   int                           `json:"id"`
	Login                string                        `json:"login"`
	GameUUID             string                        `json:"game_uuid"`
	LobbyUUID            string                        `json:"lobby_uuid"`
	MapID                int                           `json:"map_id"`
	Ready                bool                          `json:"ready"`
	BehaviorController   bool                          `json:"-"`
	userInterface        map[string]map[string]*Window // resolution, window_id, state
	memoryDynamicObjects *visible_objects.VisibleObjectsStore
	visibleObjects       *visible_objects.VisibleObjectsStore
	userUnitStore        *userUnitStore
	mx                   sync.RWMutex
}

func (client *Player) GetUnitsStore() *userUnitStore {
	if client.userUnitStore == nil {
		client.userUnitStore = &userUnitStore{units: make(map[int]*unit.Unit)}
	}

	return client.userUnitStore
}

func (client *Player) InitVisibleObjects() {
	client.checkVisibleObjectStore()
	client.visibleObjects.InitVisibleObjects()
}

func (client *Player) checkVisibleObjectStore() {
	if client.visibleObjects == nil {
		client.visibleObjects = &visible_objects.VisibleObjectsStore{}
	}
}

func (client *Player) checkMemoryObjectStore() {
	if client.memoryDynamicObjects == nil {
		client.memoryDynamicObjects = &visible_objects.VisibleObjectsStore{}
	}
}

func (client *Player) GetVisibleObjectByID(id string) *visible_objects.VisibleObject {
	client.checkVisibleObjectStore()
	return client.visibleObjects.GetVisibleObjectByID(id)
}

func (client *Player) GetVisibleObjects() <-chan *visible_objects.VisibleObject {
	client.checkVisibleObjectStore()
	return client.visibleObjects.GetVisibleObjects()
}

func (client *Player) UnsafeRangeVisibleObjects() (map[string]*visible_objects.VisibleObject, *sync.RWMutex) {
	client.checkVisibleObjectStore()
	return client.visibleObjects.UnsafeRangeMapDynamicObjects()
}

func (client *Player) RemoveVisibleObject(removeObj *visible_objects.VisibleObject) {
	client.checkVisibleObjectStore()
	client.visibleObjects.RemoveVisibleObject(removeObj)
}

func (client *Player) AddDynamicObject(object *dynamic_map_object.Object, mapID int, view, radar bool, mapTime int64) {
	client.checkMemoryObjectStore()
	if object.MemoryUUID == "" {
		object.MemoryUUID = strconv.Itoa(mapID) + "object" + strconv.Itoa(object.ID)
	}

	vObj := &visible_objects.VisibleObject{
		UUID:       object.MemoryUUID,
		IDObject:   object.ID,
		TypeObject: "object",
		Scale:      object.GetScale(),
		HP:         object.GetHP(),
		MapID:      mapID,
		X:          object.GetPhysicalModel().GetX(),
		Y:          object.GetPhysicalModel().GetY(),
		OwnerID:    object.GetOwnerID(),
		Object:     object,
		ObjectJSON: object.GetJSON(mapTime),
	}

	vObj.SetView(view)
	vObj.SetRadar(radar)

	client.memoryDynamicObjects.AddDynamicObject(vObj)
}

func (client *Player) RemoveDynamicObject(uuid string) {
	client.checkMemoryObjectStore()
	client.memoryDynamicObjects.RemoveDynamicObject(uuid)
}

func (client *Player) GetMapDynamicObjects(mapID int) <-chan *visible_objects.VisibleObject {
	client.checkMemoryObjectStore()
	return client.memoryDynamicObjects.GetMapDynamicObjects(mapID)
}

func (client *Player) UnsafeRangeMapDynamicObjects() (map[string]*visible_objects.VisibleObject, *sync.RWMutex) {
	client.checkMemoryObjectStore()
	return client.memoryDynamicObjects.UnsafeRangeMapDynamicObjects()
}

func (client *Player) GetMapDynamicObjectByID(mapID, id int) *visible_objects.VisibleObject {
	client.checkMemoryObjectStore()
	return client.memoryDynamicObjects.GetMapDynamicObjectByID(mapID, id)
}

func (client *Player) GetMapDynamicObjectByUUID(uuid string) *visible_objects.VisibleObject {
	client.checkMemoryObjectStore()
	return client.memoryDynamicObjects.GetMapDynamicObjectByUUID(uuid)
}

func (client *Player) AddVisibleObject(newObj *visible_objects.VisibleObject) {
	client.checkVisibleObjectStore()
	client.visibleObjects.AddVisibleObject(newObj)
}

func (client *Player) SetLogin(login string) {
	client.Login = login
}

func (client *Player) GetLogin() (login string) {
	return client.Login
}

func (client *Player) SetID(id int) {
	client.ID = id
}

func (client *Player) GetID() (id int) {
	return client.ID
}

func (client *Player) GetReady() bool {
	return client.Ready
}

func (client *Player) SetReady(ready bool) {
	client.Ready = ready
}
