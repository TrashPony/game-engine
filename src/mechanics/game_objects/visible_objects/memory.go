package visible_objects

import (
	"fmt"
	"strconv"
	"sync"
)

func (v *VisibleObjectsStore) AddDynamicObject(object *VisibleObject) {
	v.mx.Lock()
	defer v.mx.Unlock()

	if v.visibleObjects == nil {
		v.visibleObjects = make(map[string]*VisibleObject)
	}

	v.visibleObjects[object.UUID] = object
}

func (v *VisibleObjectsStore) AddDynamicMemoryObject(object *VisibleObject) {
	v.mx.Lock()
	defer v.mx.Unlock()

	if v.visibleObjects == nil {
		v.visibleObjects = make(map[string]*VisibleObject)
	}

	v.visibleObjects[object.UUID] = object
}

func (v *VisibleObjectsStore) RemoveDynamicObject(uuid string) {
	v.mx.Lock()
	defer v.mx.Unlock()
	delete(v.visibleObjects, uuid)
}

func (v *VisibleObjectsStore) GetMapDynamicObjectByID(mapID, id int) *VisibleObject {
	v.mx.RLock()
	defer v.mx.RUnlock()

	if v.visibleObjects == nil {
		return nil
	}

	return v.visibleObjects[strconv.Itoa(mapID)+"object"+strconv.Itoa(id)]
}

func (v *VisibleObjectsStore) GetMapDynamicObjectByUUID(uuid string) *VisibleObject {

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

	return v.visibleObjects[uuid]
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

func (v *VisibleObjectsStore) UnsafeRangeMapDynamicObjects() (map[string]*VisibleObject, *sync.RWMutex) {
	return v.visibleObjects, &v.mx
}
