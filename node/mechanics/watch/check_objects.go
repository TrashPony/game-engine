package watch

import (
	"bytes"
	"github.com/TrashPony/game-engine/router/generate_ids"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/dynamic_map_object"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/visible_objects"
	"sync"
)

type Watcher interface {
	InitVisibleObjects()
	GetVisibleObjectByTypeAndID(string, int) *visible_objects.VisibleObject
	GetVisibleObjects() <-chan *visible_objects.VisibleObject
	UnsafeRangeVisibleObjects() ([]*visible_objects.VisibleObject, *sync.RWMutex)
	RemoveVisibleObject(removeObj *visible_objects.VisibleObject)
	AddDynamicObject(object *dynamic_map_object.Object, mapID int, view, radar bool, mapTime int64)
	RemoveDynamicObject(id int)
	GetMapDynamicObjects(mapID int) <-chan *visible_objects.VisibleObject
	UnsafeRangeMapDynamicObjects() ([]*visible_objects.VisibleObject, *sync.RWMutex)
	GetMapDynamicObjectByID(id int) *visible_objects.VisibleObject
	AddVisibleObject(newObj *visible_objects.VisibleObject)
	GetTeamID() int
	CheckUnion(id int) bool
}

func CheckObjects(watcher Watcher, oldObj *visible_objects.VisibleObject, id, teamID int,
	typeMark, typeObject string, view, radar bool, uuidObj string, data, updateData []byte) (string, string, *visible_objects.VisibleObject) {
	defer func() {
		if oldObj != nil {
			oldObj.SetUpdate(true)
		}
	}()

	if oldObj == nil && view {
		// мы не видили обьект совсем а теперь видим визуально
		oldObj = &visible_objects.VisibleObject{
			IDObject:   id,
			TypeObject: typeObject,
			ID:         generate_ids.GetMarkID(),
			View:       view, Radar: radar, Type: typeMark,
			UUIDObject: uuidObj,
			ObjectJSON: data,
			TeamID:     teamID,
		}
		watcher.AddVisibleObject(oldObj)

		return "removeRadarMark", "createObj", oldObj
	}

	if oldObj == nil && !view && radar {
		// мы не видили обьект совсем и видим теперь его на радаре

		oldObj = &visible_objects.VisibleObject{
			IDObject:   id,
			TypeObject: typeObject,
			ID:         generate_ids.GetMarkID(),
			View:       view,
			Radar:      radar,
			Type:       typeMark,
			UUIDObject: uuidObj,
			ObjectJSON: data,
			TeamID:     teamID,
		}
		watcher.AddVisibleObject(oldObj)

		return "createRadarMark", "", oldObj
	}

	if oldObj != nil && !oldObj.GetView() && oldObj.GetRadar() && view {
		// мы видили обьект на радаре а теперь видим его визуально
		oldObj.SetView(true)
		oldObj.SetRadar(true)
		return "removeRadarMark", "createObj", oldObj
	}

	if oldObj != nil && oldObj.GetView() && !view && radar {
		// мы видили обьект визуально а теперь видим только на радаре
		oldObj.SetView(false)
		oldObj.SetRadar(true)
		return "createRadarMark", "removeObj", oldObj
	}

	if oldObj != nil && oldObj.GetView() && !view && !radar {
		// мы видили обьект визуально и он пропал
		watcher.RemoveVisibleObject(oldObj)
		oldObj.SetView(false)
		oldObj.SetRadar(false)
		return "removeRadarMark", "removeObj", oldObj
	}

	if oldObj != nil && !oldObj.GetView() && oldObj.GetRadar() && !view && !radar {
		// мы видили обьект на радаре и он пропал
		watcher.RemoveVisibleObject(oldObj)
		oldObj.SetView(false)
		oldObj.SetRadar(false)
		return "removeRadarMark", "", oldObj
	}

	if oldObj != nil && oldObj.GetView() && view && bytes.Compare(oldObj.UpdateChecker, updateData) != 0 {
		// у обьекта что то изменилось
		oldObj.ObjectJSON = data
		oldObj.UpdateChecker = updateData
		return "", "updateObj", oldObj
	}

	return "", "", oldObj
	// во всем остальных случаях изменение состояния не произошло
}
