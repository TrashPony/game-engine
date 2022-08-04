package visible_objects

import (
	"fmt"
	"sync"
)

func (v *VisibleObjectsStore) AddDynamicObject(object *VisibleObject) {
	v.mx.Lock()
	defer v.mx.Unlock()

	if v.visibleObjects == nil {
		v.visibleObjects = make([]*VisibleObject, 0)
	}

	v.visibleObjects = append(v.visibleObjects, object)
}

func (v *VisibleObjectsStore) RemoveDynamicObject(id int) {
	v.mx.Lock()
	defer v.mx.Unlock()

	index := -1
	for i, o := range v.visibleObjects {
		if o.ID == id {
			index = i
			break
		}
	}

	if index >= 0 {
		v.visibleObjects[index] = v.visibleObjects[len(v.visibleObjects)-1]
		v.visibleObjects = v.visibleObjects[:len(v.visibleObjects)-1]
	}
}

func (v *VisibleObjectsStore) GetMapDynamicObjectByID(id int) *VisibleObject {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in GetMapDynamicObjectByUUID squad", r)
		}
	}()

	v.mx.RLock()
	defer v.mx.RUnlock()

	if v == nil || v.visibleObjects == nil {
		return nil
	}

	for _, o := range v.visibleObjects {
		if o.ID == id {
			return o
		}
	}

	return nil
}

func (v *VisibleObjectsStore) GetMapDynamicObjects(mapID int) <-chan *VisibleObject {

	v.mx.RLock()
	objChan := make(chan *VisibleObject, len(v.visibleObjects))

	go func() {

		defer func() {
			v.mx.RUnlock()
			close(objChan)
		}()

		for _, obj := range v.visibleObjects {
			if obj.MapID == mapID {
				objChan <- obj
			}
		}
	}()

	return objChan
}

func (v *VisibleObjectsStore) UnsafeRangeMapDynamicObjects() ([]*VisibleObject, *sync.RWMutex) {
	return v.visibleObjects, &v.mx
}
