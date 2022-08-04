package _map

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/dynamic_map_object"
	"strconv"
	"sync"
)

func (mp *Map) GetChanMapDynamicObjects() <-chan *dynamic_map_object.Object {

	if mp == nil {
		objChan := make(chan *dynamic_map_object.Object, 0)
		close(objChan)
		return objChan
	}

	mp.DynamicObjectsMX.RLock()
	objChan := make(chan *dynamic_map_object.Object, len(mp.DynamicObjects))

	go func() {
		defer func() {
			mp.DynamicObjectsMX.RUnlock()
			close(objChan)
		}()

		for _, obj := range mp.DynamicObjects {
			objChan <- obj
		}
	}()

	return objChan
}

func (mp *Map) UnsafeRangeDynamicObjects() ([]*dynamic_map_object.Object, *sync.RWMutex) {
	mp.DynamicObjectsMX.RLock()
	return mp.DynamicObjects, &mp.DynamicObjectsMX
}

func (mp *Map) UnsafeRangeBuildDynamicObjects() ([]*dynamic_map_object.Object, *sync.RWMutex) {
	mp.DynamicBuildObjectsMX.RLock()
	return mp.DynamicBuildObjects, &mp.DynamicBuildObjectsMX
}

func (mp *Map) GetCopyArrayBuildDynamicObjects(basket []*dynamic_map_object.Object) []*dynamic_map_object.Object {

	mp.DynamicBuildObjectsMX.RLock()
	defer mp.DynamicBuildObjectsMX.RUnlock()

	basket = basket[:0]
	for _, obj := range mp.DynamicBuildObjects {
		basket = append(basket, obj)
	}

	return basket
}

func (mp *Map) GetDynamicObjectsByID(id int) *dynamic_map_object.Object {
	mp.DynamicObjectsMX.RLock()
	defer mp.DynamicObjectsMX.RUnlock()

	for _, obj := range mp.DynamicObjects {
		if obj.ID == id {
			return obj
		}
	}

	return nil
}

func (mp *Map) AddDynamicObject(object *dynamic_map_object.Object) {
	// TODO #1
	idString := strconv.Itoa(object.GetX()) + strconv.Itoa(object.GetY())
	object.ID, _ = strconv.Atoi(idString)
	object.SetGeoData()
	object.MapID = mp.Id

	if mp.DynamicObjects == nil {
		mp.DynamicObjectsMX.Lock()
		mp.DynamicObjects = make([]*dynamic_map_object.Object, 0)
		mp.DynamicObjectsMX.Unlock()
	}

	obj := mp.GetDynamicObjectsByID(object.ID)
	if obj != nil {
		// нельзя ставить в уже существующией обьект
		return
	}

	// обновлять геодату в зонах
	go mp.AddGeoDataObjectsToZone(object.GetPhysicalModel())

	mp.DynamicObjectsMX.Lock()
	mp.DynamicObjects = append(mp.DynamicObjects, object)
	mp.DynamicObjectsMX.Unlock()

	if object.Build {
		mp.DynamicBuildObjectsMX.Lock()
		if mp.DynamicBuildObjects == nil {
			mp.DynamicBuildObjects = make([]*dynamic_map_object.Object, 0)
		}

		mp.DynamicBuildObjects = append(mp.DynamicBuildObjects, object)
		mp.DynamicBuildObjectsMX.Unlock()
	}
}

func (mp *Map) RemoveDynamicObject(object *dynamic_map_object.Object) {

	// обновлять геодату в зонах
	go mp.RemoveGeoDataObjectsToZone(object.GetPhysicalModel())

	mp.DynamicObjectsMX.Lock()
	indexD := -1
	for i, o := range mp.DynamicObjects {
		if o.ID == object.ID {
			indexD = i
			break
		}
	}

	if indexD >= 0 {
		mp.DynamicObjects = append(mp.DynamicObjects[:indexD], mp.DynamicObjects[indexD+1:]...)
	}
	mp.DynamicObjectsMX.Unlock()

	mp.DynamicBuildObjectsMX.Lock()
	indexB := -1
	for i, o := range mp.DynamicBuildObjects {
		if o.ID == object.ID {
			indexB = i
			break
		}
	}

	if indexB >= 0 {
		mp.DynamicBuildObjects = append(mp.DynamicBuildObjects[:indexB], mp.DynamicBuildObjects[indexB+1:]...)
	}
	mp.DynamicBuildObjectsMX.Unlock()
}
