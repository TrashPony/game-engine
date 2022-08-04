package game_loop_view

import (
	"github.com/TrashPony/game-engine/node/binary_msg"
	"github.com/TrashPony/game-engine/node/mechanics/watch"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/bullet"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/dynamic_map_object"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/visible_objects"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/web_socket_response"
	"github.com/TrashPony/game-engine/router/web_socket"
	"sync"
)

func View(mp *_map.Map, units []*unit.Unit, bullets []*bullet.Bullet, b *battle.Battle, ms *web_socket.MessagesStore) {
	for _, t := range b.Teams {
		worker(t, mp, units, bullets, b, t.GetUpdateViewObjects(), ms)
	}
}

func worker(watcher watch.Watcher, mp *_map.Map, units []*unit.Unit, bullets []*bullet.Bullet, gameBattle *battle.Battle, updateViewObjects bool, ms *web_socket.MessagesStore) {

	// функция должна отслеживать что обьект вышел за пределы радара/обзора и сообщать об этом клиент
	// и наоборот небыл видим стал видим обьект вошел в зону радара/обзора
	//    -- для этого надо хранить предыдущие состояния (в прошлый раз видел, теперь нет - обьект вышел из поля зрения)
	// каждый обьект в зоне радара должен иметь метку например: objectType + id
	// каждой метке радара надо давать uuid что бы можно было ее двигать и удалять
	// так же метод получения UUID метки по objectType + id для фильтра исходящих сообщений в метода CheckView()

	// смотрим на юнитов
	for _, otherUnit := range units {
		check(watcher, mp, "ground", "unit", otherUnit.GetPhysicalModel().GetX(), otherUnit.GetPhysicalModel().GetY(), otherUnit.GetID(), otherUnit.TeamID, updateUnitMessage(otherUnit), otherUnit, gameBattle, ms, units)
	}

	for _, b := range bullets {
		check(watcher, mp, "bullet", "bullet", b.GetX(), b.GetY(), b.ID, 0, nil, b, gameBattle, ms, units)
	}

	viewDynamicObjects(watcher, mp, updateViewObjects, gameBattle, ms, units)
	removeOldDynamicMemoryObject(watcher, mp, updateViewObjects, gameBattle, ms, units)
	removeOutObjects(watcher, mp, ms)
}

type Dater interface {
	GetBytes(mapTime int64) []byte
	GetUpdateData(mapTime int64) []byte
}

func check(watcher watch.Watcher, mp *_map.Map, typeMark, typeObj string, x, y, idObj, teamID int, msgUpdate []byte, dater Dater,
	b *battle.Battle, ms *web_socket.MessagesStore, units []*unit.Unit) {

	oldVisible := watcher.GetVisibleObjectByTypeAndID(typeObj, idObj)
	view, radarVisible := watch.CheckViewCoordinate(watcher, x, y, b, units)

	if typeMark == "bullet" {
		radarVisible = false
	}

	markEvent, objEvent, newMark := watch.CheckObjects(watcher, oldVisible, idObj, teamID, typeMark,
		typeObj, view, radarVisible, "", dater.GetBytes(mp.Time), dater.GetUpdateData(mp.Time))

	appendRadarMessage(watcher, markEvent, objEvent, newMark, dater, msgUpdate, x, y, mp.Time, ms)
}

func removeOutObjects(watcher watch.Watcher, mp *_map.Map, ms *web_socket.MessagesStore) {
	// все не обновленные обьекты считаются потеряными из виду, например телепорт смерть и тд

	removeOut := func(vObj *visible_objects.VisibleObject) {

		if !vObj.GetUpdate() {
			appendRadarMessage(watcher, "removeRadarMark", "removeObj", vObj, nil, nil, 0, 0, mp.Time, ms)
			watcher.RemoveVisibleObject(vObj)
		} else {
			vObj.SetUpdate(false)
		}
	}

	obj, mx := watcher.UnsafeRangeVisibleObjects()
	mx.RLock()
	defer mx.RUnlock()

	for _, vObj := range obj {
		mx.RUnlock()
		removeOut(vObj)
		mx.RLock()
	}
}

func viewDynamicObjects(watcher watch.Watcher, mp *_map.Map, updateViewObjects bool, b *battle.Battle, ms *web_socket.MessagesStore, units []*unit.Unit) {
	// смотрим динамические обьекты, в отличие от прошлых обьектов эти можно увидеть только визуально
	// при этом пользователь их запоминает, тоесть он 1 раз увидил куст, ушел в другой конец карты,
	// куст будет виден через туман войны, НО если куст убьют то игрок будет всеравно его видеть
	// в тумане пока не откроет его зону снова.

	// TODO брать обьекты по зонам которые попадают в дальность обзора а не все обьекты на карте
	// проверяем видит ли юзер новые обьекты

	scanObject := func(obj *dynamic_map_object.Object) {

		view, radarVisible := viewDynamicObject(watcher, obj, b, units)

		var memoryObj *visible_objects.VisibleObject
		if obj.MemoryID != 0 {
			memoryObj = watcher.GetMapDynamicObjectByID(obj.MemoryID)
		}

		if memoryObj != nil {
			memoryObj.ObjectJSON = obj.GetBytes(mp.Time)
		}

		if obj.Build && radarVisible {
			oldVisible := watcher.GetVisibleObjectByTypeAndID("object", obj.ID)
			markEvent, objEvent, newMark := watch.CheckObjects(watcher, oldVisible, obj.ID, obj.TeamID, "structure",
				"object", view, radarVisible, "", obj.GetBytes(mp.Time), obj.GetUpdateData(mp.Time))

			appendRadarMessage(watcher, markEvent, objEvent, newMark, obj, nil, obj.GetX(), obj.GetY(), mp.Time, ms)
		}

		if view && memoryObj == nil {
			watcher.AddDynamicObject(obj, mp.Id, view, radarVisible, mp.Time)

			ms.AddMsg("createRadarObjMsg", "bin", web_socket_response.Response{
				BinaryMsg: binary_msg.CreateBinaryCreateObjMsg("dynamic_objects", obj.GetBytes(mp.Time)),
				TeamID:    watcher.GetTeamID(),
			}, nil)
		}

		updateData, update := updateObject(obj, memoryObj, mp)
		if view && memoryObj != nil && update {
			watcher.RemoveDynamicObject(memoryObj.ID)
			watcher.AddDynamicObject(obj, mp.Id, view, radarVisible, mp.Time)

			ms.AddMsg("updateRadarObjMsg", "bin", web_socket_response.Response{
				BinaryMsg: binary_msg.CreateBinaryUpdateObjMsg("dynamic_objects", obj.ID, updateData),
				TeamID:    watcher.GetTeamID(),
			}, nil)
		}
	}

	objects, objectsMX := mp.UnsafeRangeDynamicObjects()
	for _, obj := range objects {
		objectsMX.RUnlock()
		if updateViewObjects || obj.Build {
			scanObject(obj)
		}
		objectsMX.RLock()
	}
	objectsMX.RUnlock()
}

func removeOldDynamicMemoryObject(watcher watch.Watcher, mp *_map.Map, updateViewObjects bool, b *battle.Battle, ms *web_socket.MessagesStore, units []*unit.Unit) {

	if !updateViewObjects {
		return
	}
	// проверяем видит ли место где были старые обьекты но их уже нет

	removeOld := func(memoryObj *visible_objects.VisibleObject, wg *sync.WaitGroup) {

		defer wg.Done()

		// если обьект на другой карте или мы обновили этот обьект то он существует.
		if memoryObj.MapID != mp.Id || memoryObj.GetUpdate() {
			return
		}

		obj := mp.GetDynamicObjectsByID(memoryObj.IDObject)
		if obj != nil {
			return
		}

		view, _ := watch.CheckViewCoordinate(watcher, memoryObj.X, memoryObj.Y, b, units)
		if memoryObj.TeamID == watcher.GetTeamID() || view {

			watcher.RemoveDynamicObject(memoryObj.ID)

			ms.AddMsg("removeRadarObjMsg", "bin", web_socket_response.Response{
				BinaryMsg: binary_msg.CreateBinaryRemoveRadarObjectMsg(memoryObj.IDObject, "dynamic_objects"),
				TeamID:    watcher.GetTeamID(),
			}, nil)
		}
	}

	var removeOldWait sync.WaitGroup

	objects, mx := watcher.UnsafeRangeMapDynamicObjects()

	mx.RLock()
	for _, memoryObj := range objects {
		removeOldWait.Add(1)
		mx.RUnlock()
		removeOld(memoryObj, &removeOldWait)
		mx.RLock()
	}
	mx.RUnlock()

	removeOldWait.Wait()
}

func viewDynamicObject(watcher watch.Watcher, obj *dynamic_map_object.Object, b *battle.Battle, units []*unit.Unit) (bool, bool) {
	if len(obj.GetPhysicalModel().GeoData) == 0 || len(obj.GetPhysicalModel().GeoData) == 1 {
		return watch.CheckViewCoordinate(watcher, obj.GetX(), obj.GetY(), b, units)
	} else {
		return detailCheckObj(watcher, obj, b, units)
	}
}

func detailCheckObj(watcher watch.Watcher, obj *dynamic_map_object.Object, b *battle.Battle, units []*unit.Unit) (bool, bool) {

	radarView := false
	// если игрок видит хотя бы кусок обьекта то игрок видит обьект, а не по центру
	for _, geoPoint := range obj.GetPhysicalModel().GetGeoData() {
		view, radarV := watch.CheckViewCoordinate(watcher, geoPoint.GetX(), geoPoint.GetY(), b, units)
		if view {
			return true, true
		}

		radarView = radarView || radarV
	}

	return false, radarView
}

func appendRadarMessage(watcher watch.Watcher, markEvent, objEvent string, newMark *visible_objects.VisibleObject, dater Dater, msgUpdate []byte, x, y int, mpTime int64, ms *web_socket.MessagesStore) {

	if markEvent != "" || objEvent != "" {
		if markEvent == "createRadarMark" && objEvent != "createObj" {
			ms.AddMsg("createRadarMarkMsg", "bin", web_socket_response.Response{
				BinaryMsg: binary_msg.CreateBinaryCreateRadarMarkMsg(newMark.ID, x, y, newMark.Type),
				TeamID:    watcher.GetTeamID(),
			}, nil)
		}

		if markEvent == "removeRadarMark" {
			ms.AddMsg("removeRadarMarkMsg", "bin", web_socket_response.Response{
				BinaryMsg: binary_msg.CreateBinaryRemoveRadarMarkMsg(newMark.ID),
				TeamID:    watcher.GetTeamID(),
			}, nil)
		}

		if objEvent == "createObj" {
			ms.AddMsg("createRadarObjMsg", "bin", web_socket_response.Response{
				BinaryMsg: binary_msg.CreateBinaryCreateObjMsg(newMark.TypeObject, dater.GetBytes(mpTime)),
				TeamID:    watcher.GetTeamID(),
			}, nil)
		}

		if objEvent == "updateObj" {
			ms.AddMsg("updateRadarObjMsg", "bin", web_socket_response.Response{
				BinaryMsg: binary_msg.CreateBinaryUpdateObjMsg(newMark.TypeObject, newMark.IDObject, msgUpdate),
				TeamID:    watcher.GetTeamID(),
			}, nil)
		}

		if objEvent == "removeObj" {
			ms.AddMsg("removeRadarObjMsg", "bin", web_socket_response.Response{
				BinaryMsg: binary_msg.CreateBinaryRemoveRadarObjectMsg(newMark.IDObject, newMark.TypeObject),
				TeamID:    watcher.GetTeamID(),
			}, nil)
		}
	}
}
