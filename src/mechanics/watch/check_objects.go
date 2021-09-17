package watch

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/dynamic_map_object"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/visible_objects"
	"github.com/satori/go.uuid"
	"sync"
)

type Watcher interface {
	InitVisibleObjects()
	GetVisibleObjectByID(string) *visible_objects.VisibleObject
	GetVisibleObjects() <-chan *visible_objects.VisibleObject
	UnsafeRangeVisibleObjects() (map[string]*visible_objects.VisibleObject, *sync.RWMutex)
	RemoveVisibleObject(removeObj *visible_objects.VisibleObject)
	AddDynamicObject(object *dynamic_map_object.Object, mapID int, view, radar bool, mapTime int64)
	RemoveDynamicObject(uuid string)
	GetMapDynamicObjects(mapID int) <-chan *visible_objects.VisibleObject
	UnsafeRangeMapDynamicObjects() (map[string]*visible_objects.VisibleObject, *sync.RWMutex)
	GetMapDynamicObjectByID(mapID, id int) *visible_objects.VisibleObject
	GetMapDynamicObjectByUUID(uuid string) *visible_objects.VisibleObject
	AddVisibleObject(newObj *visible_objects.VisibleObject)
}

var markIDGenerate = 0

func CheckObjects(watcher Watcher, oldObj *visible_objects.VisibleObject, id, hp int,
	typeMark, typeObject string, view, radar bool, uuidObj string) (string, string, *visible_objects.VisibleObject) {
	defer func() {
		if oldObj != nil {
			oldObj.SetUpdate(true)
		}
	}()

	if oldObj == nil && view {
		// мы не видили обьект совсем а теперь видим визуально
		oldObj = &visible_objects.VisibleObject{
			HP:         hp,
			IDObject:   id,
			TypeObject: typeObject,
			UUID:       uuid.NewV1().String(),
			ID:         markIDGenerate,
			View:       view, Radar: radar, Type: typeMark,
			UUIDObject: uuidObj,
		}
		markIDGenerate++
		watcher.AddVisibleObject(oldObj)

		return "removeRadarMark", "createObj", oldObj
	}

	if oldObj == nil && !view && radar {
		// мы не видили обьект совсем и видим теперь его на радаре

		oldObj = &visible_objects.VisibleObject{
			IDObject:   id,
			TypeObject: typeObject,
			UUID:       uuid.NewV1().String(),
			ID:         markIDGenerate,
			View:       view,
			Radar:      radar,
			Type:       typeMark,
			UUIDObject: uuidObj,
		}
		markIDGenerate++
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

	if oldObj != nil && oldObj.GetView() && view && oldObj.HP != hp {
		// у обьекта сменился уровень здоровья
		oldObj.HP = hp
		return "", "updateObj", oldObj
	}

	return "", "", oldObj
	// во всем остальных случаях изменение состояния не произошло (но это не точно)
}
