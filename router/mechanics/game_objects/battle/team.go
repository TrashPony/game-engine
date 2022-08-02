package battle

import (
	"github.com/TrashPony/game-engine/router/generate_ids"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/dynamic_map_object"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/visible_objects"
	"sync"
)

type Team struct {
	ID                     int         `json:"id"`
	Hide                   bool        `json:"hide"`
	Spawn                  *_map.Spawn `json:"-"`
	Winner                 bool        `json:"winner"`
	PlayersIDs             []int       `json:"players_ids"`
	Points                 int         `json:"points"`
	TickPoint              int         `json:"tick_point"`
	AI                     bool        `json:"ai"`
	unions                 []int       `json:"unions"`
	memoryDynamicObjects   *visible_objects.VisibleObjectsStore
	visibleObjects         *visible_objects.VisibleObjectsStore
	countUpdateViewObjects int
}

func (t *Team) CheckUnion(id int) bool {
	for _, union_id := range t.unions {
		if union_id == id {
			return true
		}
	}

	return false
}

func (t *Team) AddUnion(id int) {
	if t.unions != nil {
		t.unions = make([]int, 0)
	}
	t.unions = append(t.unions, id)
}

func (t *Team) GetTeamID() int {
	return t.ID
}

func (t *Team) GetID() int {
	return t.ID
}

func (t *Team) GetType() string {
	return "team"
}

func (t *Team) InitVisibleObjects() {
	t.checkVisibleObjectStore()
	t.visibleObjects.InitVisibleObjects()
}

func (t *Team) InitMemmoryObjects() {
	t.checkMemoryObjectStore()
	t.memoryDynamicObjects.InitVisibleObjects()
}

func (t *Team) checkVisibleObjectStore() {
	if t.visibleObjects == nil {
		t.visibleObjects = &visible_objects.VisibleObjectsStore{}
	}
}

func (t *Team) checkMemoryObjectStore() {
	if t.memoryDynamicObjects == nil {
		t.memoryDynamicObjects = &visible_objects.VisibleObjectsStore{}
	}
}

func (t *Team) GetVisibleObjectByTypeAndID(typeObj string, id int) *visible_objects.VisibleObject {
	t.checkVisibleObjectStore()
	return t.visibleObjects.GetVisibleObjectByTypeAndID(typeObj, id)
}

func (t *Team) GetVisibleObjects() <-chan *visible_objects.VisibleObject {
	t.checkVisibleObjectStore()
	return t.visibleObjects.GetVisibleObjects()
}

func (t *Team) UnsafeRangeVisibleObjects() ([]*visible_objects.VisibleObject, *sync.RWMutex) {
	t.checkVisibleObjectStore()
	return t.visibleObjects.UnsafeRangeMapDynamicObjects()
}

func (t *Team) RemoveVisibleObject(removeObj *visible_objects.VisibleObject) {
	t.checkVisibleObjectStore()
	t.visibleObjects.RemoveVisibleObject(removeObj)
}

func (t *Team) AddDynamicObject(object *dynamic_map_object.Object, mapID int, view, radar bool, mapTime int64) {
	t.checkMemoryObjectStore()
	if object.MemoryID == 0 {
		object.MemoryID = generate_ids.GetMarkID()
	}

	vObj := &visible_objects.VisibleObject{
		ID:         object.MemoryID,
		IDObject:   object.ID,
		TypeObject: "object",
		Scale:      object.GetScale(),
		HP:         object.GetHP(),
		MapID:      mapID,
		X:          object.GetPhysicalModel().GetX(),
		Y:          object.GetPhysicalModel().GetY(),
		OwnerID:    object.GetOwnerID(),
		Object:     object,
		ObjectJSON: object.GetBytes(mapTime),
		Work:       object.Work,
		TeamID:     object.TeamID,
	}

	vObj.SetView(view)
	vObj.SetRadar(radar)

	t.memoryDynamicObjects.AddDynamicObject(vObj)
}

func (t *Team) RemoveDynamicObject(id int) {
	t.checkMemoryObjectStore()
	t.memoryDynamicObjects.RemoveDynamicObject(id)
}

func (t *Team) GetMapDynamicObjects(mapID int) <-chan *visible_objects.VisibleObject {
	t.checkMemoryObjectStore()
	return t.memoryDynamicObjects.GetMapDynamicObjects(mapID)
}

func (t *Team) UnsafeRangeMapDynamicObjects() ([]*visible_objects.VisibleObject, *sync.RWMutex) {
	t.checkMemoryObjectStore()
	return t.memoryDynamicObjects.UnsafeRangeMapDynamicObjects()
}

func (t *Team) GetMapDynamicObjectByID(id int) *visible_objects.VisibleObject {
	t.checkMemoryObjectStore()
	return t.memoryDynamicObjects.GetMapDynamicObjectByID(id)
}

func (t *Team) AddVisibleObject(newObj *visible_objects.VisibleObject) {
	t.checkVisibleObjectStore()
	t.visibleObjects.AddVisibleObject(newObj)
}

func (t *Team) GetVisibleObjectStore() *visible_objects.VisibleObjectsStore {
	t.checkVisibleObjectStore()
	return t.visibleObjects
}

// todo мега костыль что бы не обновлять обьекты каждые 32 мс а как укажешь
func (t *Team) GetUpdateViewObjects() bool {
	if t.countUpdateViewObjects == 3 {
		t.countUpdateViewObjects = 0
		return true
	} else {
		t.countUpdateViewObjects++
		return false
	}
}
