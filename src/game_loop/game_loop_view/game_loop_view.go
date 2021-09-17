package game_loop_view

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/bullet"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/dynamic_map_object"
	_map "github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/player"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/visible_objects"
	"github.com/TrashPony/game_engine/src/mechanics/watch"
	"github.com/TrashPony/game_engine/src/web_socket"
	"strconv"
	"sync"
)

type RadarMsg struct {
	RadarMark    *visible_objects.VisibleObject `json:"rm"`
	ActionMark   string                         `json:"am"` // удалить/создать/скрыть/раскрыть метку
	ActionObject string                         `json:"ao"` // удалить/создать обьект
	Object       interface{}                    `json:"o"`  // обьект (ящик, транспорт, юнит и тд)
	X            int                            `json:"x"`
	Y            int                            `json:"y"`
}

type radarMsgs struct {
	Event  string      `json:"e"`
	Events []*RadarMsg `json:"ev"`
}

func View(mp *_map.Map, players map[int]*player.Player, units []*unit.Unit, mapObjects map[int]*dynamic_map_object.Object, bullets map[string]*bullet.Bullet) {
	for _, gamePlayer := range players {
		if gamePlayer.GetReady() {
			worker(gamePlayer, gamePlayer, mp, units, mapObjects, bullets)
		}
	}
}

type radarEvents struct {
	events []*RadarMsg
	mx     sync.Mutex
}

func (r *radarEvents) add(event *RadarMsg) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.events = append(r.events, event)
}

func worker(p *player.Player, watcher watch.Watcher, mp *_map.Map, units []*unit.Unit, mapObjects map[int]*dynamic_map_object.Object, bullets map[string]*bullet.Bullet) []*RadarMsg {

	// функция должна отслеживать что обьект вышел за пределы радара/обзора и сообщать об этом клиент
	// и наоборот небыл видим стал видим обьект вошел в зону радара/обзора
	//    -- для этого надо хранить предыдущие состояния (в прошлый раз видел, теперь нет - обьект вышел из поля зрения)
	// каждый обьект в зоне радара должен иметь метку например: objectType + id
	// каждой метке радара надо давать uuid что бы можно было ее двигать и удалять
	// так же метод получения UUID метки по objectType + id для фильтра исходящих сообщений в метода CheckView()

	events := &radarEvents{events: make([]*RadarMsg, 0)}

	// смотрим на юнитов
	for _, otherUnit := range units {
		check(watcher, mp, "ground", "unit", otherUnit.GetPhysicalModel().GetX(), otherUnit.GetPhysicalModel().GetY(), otherUnit.GetID(), otherUnit.HP, updateUnitMessage(otherUnit), events, otherUnit)
	}

	for _, b := range bullets {
		if b.Ammo.Type == "laser" {
			continue
		}
		// мы не видим пули на радаре, но видим ракеты
		check(watcher, mp, "bullet", "bullet", b.GetX(), b.GetY(), b.ID, 0, nil, events, b)
	}

	viewDynamicObjects(watcher, mp, events, true, mapObjects)
	removeOldDynamicMemoryObject(watcher, mp, events, true, mapObjects)
	removeOutObjects(watcher, mp, events)

	if len(events.events) > 0 {
		web_socket.SendMessage(web_socket.Response{
			Event: "rw",
			Data: radarMsgs{
				Event:  "rw",
				Events: events.events,
			},
			PlayerID: p.GetID(),
		})
	}

	return events.events
}

type JSONer interface {
	GetJSON(mapTime int64) string
}

func check(watcher watch.Watcher, mp *_map.Map, typeMark, typeObj string, x, y, idObj, hp int, msgUpdate interface{}, events *radarEvents, jsoner JSONer) {

	oldVisible := watcher.GetVisibleObjectByID(typeObj + strconv.Itoa(idObj))
	view, _ := watch.CheckViewCoordinate()

	markEvent, objEvent, newMark := watch.CheckObjects(watcher, oldVisible, idObj, hp, typeMark,
		typeObj, view, true, "")

	appendRadarMessage(events, markEvent, objEvent, newMark, jsoner, msgUpdate, x, y, mp.Time)
}

func removeOutObjects(watcher watch.Watcher, mp *_map.Map, events *radarEvents) {
	// все не обновленные обьекты считаются потеряными из виду, например телепорт смерть и тд

	removeOut := func(vObj *visible_objects.VisibleObject) {

		if !vObj.GetUpdate() {
			appendRadarMessage(events, "removeRadarMark", "removeObj", vObj, nil, nil, 0, 0, mp.Time)
			watcher.RemoveVisibleObject(vObj)
		} else {
			vObj.SetUpdate(false)
		}
	}

	for vObj := range watcher.GetVisibleObjects() {
		removeOut(vObj)
	}
}

func viewDynamicObjects(watcher watch.Watcher, mp *_map.Map, events *radarEvents, updateViewObjects bool, mapObjects map[int]*dynamic_map_object.Object) {
	// смотрим динамические обьекты, в отличие от прошлых обьектов эти можно увидеть только визуально
	// при этом пользователь их запоминает, тоесть он 1 раз увидил куст, ушел в другой конец карты,
	// куст будет виден через туман войны, НО если куст убьют то игрок будет всеравно его видеть
	// в тумане пока не откроет его зону снова.

	// TODO брать обьекты по зонам которые попадают в дальность обзора а не все обьекты на карте
	// проверяем видит ли юзер новые обьекты

	scanObject := func(obj *dynamic_map_object.Object) {

		view, radarVisible := viewDynamicObject()

		var memoryObj *visible_objects.VisibleObject
		if obj.MemoryUUID != "" {
			memoryObj = watcher.GetMapDynamicObjectByUUID(obj.MemoryUUID)
		}

		if view && memoryObj == nil {
			watcher.AddDynamicObject(obj, mp.Id, view, radarVisible, mp.Time)

			events.add(&RadarMsg{
				RadarMark:    &visible_objects.VisibleObject{TypeObject: "dynamic_objects", IDObject: obj.ID},
				ActionMark:   "",
				ActionObject: "createObj",
				Object:       obj.GetJSON(mp.Time),
				X:            obj.GetPhysicalModel().GetX(),
				Y:            obj.GetPhysicalModel().GetY(),
			})
		}

		jsonMsg, update := updateObject(obj, memoryObj, mp)
		if view && memoryObj != nil && update {
			watcher.RemoveDynamicObject(memoryObj.UUID)
			watcher.AddDynamicObject(obj, mp.Id, view, radarVisible, mp.Time)

			events.add(&RadarMsg{
				RadarMark:    &visible_objects.VisibleObject{TypeObject: "dynamic_objects", IDObject: obj.ID},
				ActionMark:   "",
				ActionObject: "updateObj",
				Object:       jsonMsg,
				X:            obj.GetPhysicalModel().GetX(),
				Y:            obj.GetPhysicalModel().GetY(),
			})
		}
	}

	for _, obj := range mapObjects {
		scanObject(obj)
	}
}

func removeOldDynamicMemoryObject(watcher watch.Watcher, mp *_map.Map, events *radarEvents, updateViewObjects bool, mapObjects map[int]*dynamic_map_object.Object) {

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

		obj := mapObjects[memoryObj.IDObject]
		if obj != nil {
			return
		}

		view, _ := watch.CheckViewCoordinate()

		if view {

			watcher.RemoveDynamicObject(memoryObj.UUID)

			events.add(&RadarMsg{
				RadarMark:    &visible_objects.VisibleObject{TypeObject: "dynamic_objects", IDObject: memoryObj.IDObject},
				ActionMark:   "",
				ActionObject: "removeObj",
				Object:       memoryObj,
				X:            memoryObj.X,
				Y:            memoryObj.Y,
			})
		}
	}

	var removeOldWait sync.WaitGroup

	objects, mx := watcher.UnsafeRangeMapDynamicObjects()

	mx.RLock()
	for _, memoryObj := range objects {
		removeOldWait.Add(1)
		go removeOld(memoryObj, &removeOldWait)
	}
	mx.RUnlock()

	removeOldWait.Wait()
}

func viewDynamicObject() (bool, bool) {
	return true, true
}

func appendRadarMessage(events *radarEvents, markEvent, objEvent string, newMark *visible_objects.VisibleObject, jsoner JSONer, msgUpdate interface{}, x, y int, mpTime int64) {

	if markEvent != "" || objEvent != "" {

		var msg interface{}

		if objEvent == "createObj" && jsoner != nil {
			msg = jsoner.GetJSON(mpTime)
		}

		if objEvent == "updateObj" {
			msg = msgUpdate
		}

		events.add(&RadarMsg{
			RadarMark:    newMark,
			ActionMark:   markEvent,
			ActionObject: objEvent,
			Object:       msg,
			X:            x,
			Y:            y,
		})
	}
}
