package visible_objects

import (
	"strconv"
	"sync"
)

type VisibleObjectsStore struct {
	visibleObjects map[string]*VisibleObject // key id_object+type_object
	mx             sync.RWMutex
}

func (v *VisibleObjectsStore) InitVisibleObjects() {
	v.mx.Lock()
	defer v.mx.Unlock()
	v.visibleObjects = make(map[string]*VisibleObject)
}

func (v *VisibleObjectsStore) GetVisibleObjectByID(id string) *VisibleObject {
	v.mx.RLock()
	defer v.mx.RUnlock()

	if v.visibleObjects == nil {
		return nil
	}

	object, ok := v.visibleObjects[id]
	if ok {
		return object
	}

	return nil
}

func (v *VisibleObjectsStore) GetVisibleObjects() <-chan *VisibleObject {

	v.mx.RLock()

	objs := make(chan *VisibleObject, len(v.visibleObjects))

	go func() {
		defer func() {
			v.mx.RUnlock()
			close(objs)
		}()

		for _, obj := range v.visibleObjects {
			objs <- obj
		}
	}()

	return objs
}

func (v *VisibleObjectsStore) UnsafeRangeVisibleObjects() (map[string]*VisibleObject, *sync.RWMutex) {
	return v.visibleObjects, &v.mx
}

func (v *VisibleObjectsStore) AddVisibleObject(newObj *VisibleObject) {
	v.mx.Lock()
	defer v.mx.Unlock()

	if v.visibleObjects == nil {
		v.visibleObjects = make(map[string]*VisibleObject)
	}

	v.visibleObjects[newObj.TypeObject+strconv.Itoa(newObj.IDObject)] = newObj
}

func (v *VisibleObjectsStore) RemoveVisibleObject(removeObj *VisibleObject) {
	v.mx.Lock()
	defer v.mx.Unlock()

	delete(v.visibleObjects, removeObj.TypeObject+strconv.Itoa(removeObj.IDObject))
}

type VisibleObject struct {
	ID         int    `json:"id_mark"`
	IDObject   int    `json:"id"`
	TypeObject string `json:"to"`
	UUIDObject string `json:"uo"`
	UUID       string `json:"-"`
	View       bool   `json:"-"`    // в прямой видимости
	Radar      bool   `json:"-"`    // видим только радаром
	Type       string `json:"type"` // fly(летающий), ground(наземный), structure(структура), resource(ресурс)
	update     bool

	HP         int         `json:"-"`
	Complete   float64     `json:"-"`
	Scale      int         `json:"-"`
	Energy     int         `json:"-"`
	MapID      int         `json:"mid"`
	X          int         `json:"-"`
	Y          int         `json:"-"`
	OwnerID    int         `json:"-"`
	Object     interface{} `json:"-"`
	ObjectJSON string      `json:"-"`
	Work       bool        `json:"-"`

	UpdateMsg *UpdateObjMap `json:"-"`

	mx sync.RWMutex
}

type UpdateObjMap struct {
	UpdateMap  map[string]interface{}
	Update     bool
	JSON       string
	ServerTime int64
	Mx         sync.Mutex
}

func (v *VisibleObject) GetUpdate() bool {
	return v.update
}

func (v *VisibleObject) SetUpdate(update bool) {
	v.update = update
}

func (v *VisibleObject) GetView() bool {
	return v.View
}

func (v *VisibleObject) SetView(view bool) {
	v.View = view
}

func (v *VisibleObject) GetRadar() bool {
	return v.Radar
}

func (v *VisibleObject) SetRadar(radar bool) {
	v.Radar = radar
}
