package visible_objects

import (
	"sync"
)

type VisibleObjectsStore struct {
	visibleObjects []*VisibleObject // key id_object+type_object
	mx             sync.RWMutex
}

func (v *VisibleObjectsStore) InitVisibleObjects() {
	v.mx.Lock()
	defer v.mx.Unlock()

	if v.visibleObjects != nil {
		for _, v := range v.visibleObjects {
			v.UpdateMsg = nil
			v.Object = nil
			v.ObjectJSON = nil
		}
	}

	v.visibleObjects = make([]*VisibleObject, 0)
}

func (v *VisibleObjectsStore) GetVisibleObjectByTypeAndID(typeObj string, id int) *VisibleObject {
	v.mx.RLock()
	defer v.mx.RUnlock()

	if v.visibleObjects == nil {
		return nil
	}

	for _, vObj := range v.visibleObjects {
		if vObj.TypeObject == typeObj && vObj.IDObject == id {
			return vObj
		}
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

func (v *VisibleObjectsStore) AddVisibleObject(newObj *VisibleObject) {
	v.mx.Lock()
	defer v.mx.Unlock()

	if v.visibleObjects == nil {
		v.visibleObjects = make([]*VisibleObject, 0)
	}

	// TODO := newObj.TypeObject+strconv.Itoa(newObj.IDObject)
	v.visibleObjects = append(v.visibleObjects, newObj)
}

func (v *VisibleObjectsStore) RemoveVisibleObject(removeObj *VisibleObject) {
	v.mx.Lock()
	defer v.mx.Unlock()

	index := -1
	for i, o := range v.visibleObjects {
		if o.ID == removeObj.ID {
			index = i
			break
		}
	}

	if index >= 0 {
		v.visibleObjects[index] = v.visibleObjects[len(v.visibleObjects)-1]
		v.visibleObjects = v.visibleObjects[:len(v.visibleObjects)-1]
	}
}

type VisibleObject struct {
	ID         int    `json:"id_mark"`
	IDObject   int    `json:"id"`
	TypeObject string `json:"to"`
	TeamID     int    `json:"teamID"`
	UUIDObject string `json:"uo"`
	//UUID       string `json:"-"`
	View   bool   `json:"-"`    // в прямой видимости
	Radar  bool   `json:"-"`    // видим только радаром
	Type   string `json:"type"` // fly(летающий), ground(наземный), structure(структура), resource(ресурс)
	update bool

	HP            int         `json:"-"`
	Complete      float64     `json:"-"`
	Scale         int         `json:"-"`
	Energy        int         `json:"-"`
	MapID         int         `json:"mid"`
	X             int         `json:"-"`
	Y             int         `json:"-"`
	OwnerID       int         `json:"-"`
	Object        interface{} `json:"-"`
	ObjectJSON    []byte      `json:"-"`
	UpdateChecker []byte      `json:"-"`
	Work          bool        `json:"-"`

	UpdateMsg *UpdateObjMap `json:"-"`

	mx sync.RWMutex
}

type UpdateObjMap struct {
	UpdateData []byte
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
