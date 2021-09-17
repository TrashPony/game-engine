package _map

import (
	"errors"
	"fmt"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/coordinate"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/dynamic_map_object"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/obstacle_point"
	"math/rand"
	"strconv"
	"sync"
)

type Map struct {
	Id            int      `json:"id"`
	TypeID        int      `json:"-"`
	Name          string   `json:"name"`
	XSize         int      `json:"x_size"`
	YSize         int      `json:"y_size"`
	DefaultLevel  float64  `json:"default_level"`
	Specification string   `json:"specification"`
	Type          string   `json:"-"`
	Spawns        []*Spawn `json:"spawns"`

	// интерактивные координаты на карте, типо телепорты, выхода с баз, заходы на базы и тд
	OneLayerMap map[int]map[int]*coordinate.Coordinate `json:"one_layer_map"`

	// текстуры земли
	Flore map[int]map[int]*dynamic_map_object.Flore `json:"flore"`
	// все обьекты у которых нет поведения находятся в карте OneLayerMap, дороги, горы, базы.
	// Игрок видит эти обьекты всегда независимо от радара/обзора
	StaticObjects     map[int]*dynamic_map_object.Object `json:"-"`
	StaticObjectsJSON map[int]string                     `json:"static_objects_json"` // статисные обьекты не изменяются поэтому кешируем тут обьекты для фронтенда
	StaticObjectsMX   sync.RWMutex                       `json:"-"`
	// в матрице DynamicObjects находятся обьекты которые могут передвигатся/уничтожатся/рождатся
	// тоесть это обьекты с поведением, игрок их видит и запоминает последнее их состояние у себя.
	// Игрок не видит если с обьектов что либо произошло вне поле его зрения.
	DynamicObjects   map[int]*dynamic_map_object.Object `json:"-"`
	DynamicObjectsMX sync.RWMutex                       `json:"-"`
	// вспомогательная мапа для уменьшение блокировок, сюда падают обьекты только с build=true
	DynamicBuildObjects   map[int]*dynamic_map_object.Object `json:"-"`
	DynamicBuildObjectsMX sync.RWMutex                       `json:"-"`

	GeoData []*obstacle_point.ObstaclePoint `json:"geo_data"`

	// разделяем карту на зоны (DiscreteSize х DiscreteSize) при загрузке сервера,
	// добавляем в зону все поинты которые пересекают данных квадрат и ближайшие к нему
	// когда надо найти колизию с юнитом делем его полизию на 100 и отбрасываем дровь так мы получим зону
	// например положение юнита 55/DiscreteSize:77/DiscreteSize = зона 0:0, 257/DiscreteSize:400/DiscreteSize = 1:1, 1654/DiscreteSize:2340/DiscreteSize = 6:9
	// и смотрим только те точки которые находятся в данной зоне
	GeoZones [][]*Zone `json:"-"`

	LevelMap map[string]*LvlMap `json:"level_map"`
	lvlMX    sync.RWMutex
	LoopInit bool           `json:"-"`
	Time     int64          `json:"-"` // игровое время относительно сервертиков
	Clouds   map[int]*Cloud `json:"-"`
	Exit     bool           `json:"-"`
	mx       sync.RWMutex
}

type Spawn struct {
	ID           int         `json:"-"`
	X            int         `json:"x"`
	Y            int         `json:"y"`
	Name         string      `json:"name"`
	Radius       int         `json:"radius"`
	Rotate       int         `json:"rotate"`
	Type         string      `json:"type"`
	CaptureTeam  string      `json:"capture_team"`
	Capture      int         `json:"capture"`
	TypeGame     string      `json:"type_game"`
	SubTypeGame  int         `json:"sub_type_game"`
	CaptureFact  bool        `json:"capture_fact"`
	ObjectsOwner map[int]int `json:"-"` // [id обьекта в базе] id обьекта когда получен при создание игры
	SpawnOwner   map[int]int `json:"-"` // [id спавна в базе]
	ReloadTime   int         `json:"-"`
}

func (mp *Map) GetCountDynamicObjects() int {
	mp.DynamicObjectsMX.RLock()
	defer mp.DynamicObjectsMX.RUnlock()

	return len(mp.DynamicObjects)
}

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

func (mp *Map) UnsafeRangeDynamicObjects() (map[int]*dynamic_map_object.Object, *sync.RWMutex) {
	mp.DynamicObjectsMX.RLock()
	return mp.DynamicObjects, &mp.DynamicObjectsMX
}

func (mp *Map) UnsafeRangeBuildDynamicObjects() (map[int]*dynamic_map_object.Object, *sync.RWMutex) {
	mp.DynamicBuildObjectsMX.RLock()
	return mp.DynamicBuildObjects, &mp.DynamicBuildObjectsMX
}

//func (mp *Map) DynamicObjectsUnsafeRange() (map[int]*dynamic_map_object.Object, *sync.RWMutex) {
//	mp.DynamicObjectsMX.RLock()
//	return mp.DynamicObjects, &mp.DynamicObjectsMX
//}

func (mp *Map) GetCopyMapDynamicObjects() (copyMap map[int]*dynamic_map_object.Object, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("error")
			fmt.Println("Recovered in GetCopyMapDynamicObjects", r)
		}
	}()

	copyMap = make(map[int]*dynamic_map_object.Object)

	mp.DynamicObjectsMX.RLock()
	defer mp.DynamicObjectsMX.RUnlock()

	for _, obj := range mp.DynamicObjects {
		if obj != nil {
			copyMap[obj.ID] = obj
		}
	}

	return copyMap, nil
}

func (mp *Map) GetDynamicObjects(x, y int) *dynamic_map_object.Object {
	mp.DynamicObjectsMX.RLock()
	defer mp.DynamicObjectsMX.RUnlock()

	idString := strconv.Itoa(x) + strconv.Itoa(y)
	id, _ := strconv.Atoi(idString)

	obj := mp.DynamicObjects[id]
	return obj
}

func (mp *Map) GetDynamicObjectsByID(id int) *dynamic_map_object.Object {
	mp.DynamicObjectsMX.RLock()
	defer mp.DynamicObjectsMX.RUnlock()
	return mp.DynamicObjects[id]
}

func (mp *Map) AddCrater(crater *dynamic_map_object.Object) {
	if crater != nil {

		crater.DestroyLeftTimer = true
		crater.DestroyTimer = 10000

		mp.AddDynamicObject(crater)
	}
}

func (mp *Map) AddFlore(flore *dynamic_map_object.Flore) {
	// TODO #1
	mp.DynamicObjectsMX.Lock()
	defer mp.DynamicObjectsMX.Unlock()

	_, ok := mp.Flore[flore.X][flore.Y]
	if ok {
		delete(mp.Flore[flore.X], flore.Y)
	}

	if mp.Flore[flore.X] != nil {
		mp.Flore[flore.X][flore.Y] = flore
	} else {
		mp.Flore[flore.X] = make(map[int]*dynamic_map_object.Flore)
		mp.Flore[flore.X][flore.Y] = flore
	}
}

func (mp *Map) AddStaticObject(object *dynamic_map_object.Object) {
	// TODO #1
	mp.DynamicObjectsMX.Lock()
	defer mp.DynamicObjectsMX.Unlock()

	if mp.StaticObjects == nil {
		mp.StaticObjects = make(map[int]*dynamic_map_object.Object)
	}

	idString := strconv.Itoa(object.GetPhysicalModel().GetX()) + strconv.Itoa(object.GetPhysicalModel().GetY())
	object.ID, _ = strconv.Atoi(idString)
	object.MapID = mp.Id

	_, ok := mp.StaticObjects[object.ID]
	if ok {
		delete(mp.StaticObjects, object.ID)
	}

	mp.StaticObjects[object.ID] = object
}

func (mp *Map) RemoveStaticObject(x, y int) {
	// TODO #2

	mp.DynamicObjectsMX.Lock()
	defer mp.DynamicObjectsMX.Unlock()

	idString := strconv.Itoa(x) + strconv.Itoa(y)
	ID, _ := strconv.Atoi(idString)

	delete(mp.StaticObjects, ID)
}

func (mp *Map) AddDynamicObject(object *dynamic_map_object.Object) {
	// TODO #1
	idString := strconv.Itoa(object.GetPhysicalModel().GetX()) + strconv.Itoa(object.GetPhysicalModel().GetY())
	object.ID, _ = strconv.Atoi(idString)
	object.SetGeoData()
	object.MapID = mp.Id

	//mp.DynamicObjectsMX.Lock()
	//defer mp.DynamicObjectsMX.Unlock()

	if mp.DynamicObjects == nil {
		mp.DynamicObjectsMX.Lock()
		mp.DynamicObjects = make(map[int]*dynamic_map_object.Object)
		mp.DynamicObjectsMX.Unlock()
	}

	mp.DynamicObjectsMX.RLock()
	_, ok := mp.DynamicObjects[object.ID]
	mp.DynamicObjectsMX.RUnlock()
	if ok {
		// нельзя ставить в уже существующией обьект
		return
	}

	// обновлять геодату в зонах
	go mp.AddGeoDataObjectsToZone(object.GetPhysicalModel())
	//go cache.CollisionRectToCircleMap.PurgeCacheByObject(object)
	mp.DynamicObjectsMX.Lock()
	mp.DynamicObjects[object.ID] = object
	mp.DynamicObjectsMX.Unlock()
}

func (mp *Map) RemoveDynamicObject(object *dynamic_map_object.Object) {

	// обновлять геодату в зонах
	go mp.RemoveGeoDataObjectsToZone(object.GetPhysicalModel())

	//go cache.CollisionRectToCircleMap.PurgeCacheByObject(object)
	mp.DynamicObjectsMX.Lock()
	delete(mp.DynamicObjects, object.ID)
	mp.DynamicObjectsMX.Unlock()

	mp.DynamicBuildObjectsMX.Lock()
	delete(mp.DynamicBuildObjects, object.ID)
	mp.DynamicBuildObjectsMX.Unlock()
}

func (mp *Map) RemoveDynamicObjectByXY(x, y int) {

	idString := strconv.Itoa(x) + strconv.Itoa(y)
	id, _ := strconv.Atoi(idString)

	obj := mp.GetDynamicObjectsByID(id)
	if obj != nil {
		mp.RemoveDynamicObject(obj)
	}
}

type ShortInfoMap struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	XSize         int    `json:"x_size"`
	YSize         int    `json:"y_size"`
	Specification string `json:"specification"`
}

func (mp *Map) GetShortInfoMap() *ShortInfoMap {
	return &ShortInfoMap{
		Id:            mp.Id,
		Name:          mp.Name,
		XSize:         mp.XSize,
		YSize:         mp.YSize,
		Specification: mp.Specification,
	}
}

func (mp *Map) GetCoordinate(x, y int) (*coordinate.Coordinate, bool) {
	mapCoordinate, find := mp.OneLayerMap[x][y]
	if !find {
		mapCoordinate = &coordinate.Coordinate{X: x, Y: y}

		if mp.OneLayerMap[x] == nil {
			mp.OneLayerMap[x] = make(map[int]*coordinate.Coordinate)
		}

		mp.OneLayerMap[x][y] = mapCoordinate
	}

	return mapCoordinate, true
}

func (mp *Map) GetStaticObj(x, y int) *dynamic_map_object.Object {
	mp.DynamicObjectsMX.RLock()
	defer mp.DynamicObjectsMX.RUnlock()

	idString := strconv.Itoa(x) + strconv.Itoa(y)
	ID, _ := strconv.Atoi(idString)

	staticObj := mp.StaticObjects[ID]
	return staticObj
}

func (mp *Map) GetStaticObjByID(id int) *dynamic_map_object.Object {
	mp.DynamicObjectsMX.RLock()
	defer mp.DynamicObjectsMX.RUnlock()

	staticObj := mp.StaticObjects[id]
	return staticObj
}

func (mp *Map) GetFlore(x, y int) *dynamic_map_object.Flore {
	flore := mp.Flore[x][y]
	return flore
}

func (mp *Map) DeleteCoordinate(x, y int) {
	_, find := mp.OneLayerMap[x][y]
	if find {
		delete(mp.OneLayerMap[x], y)
	}
}

func (mp *Map) AddCoordinate(newCoordinate *coordinate.Coordinate) {
	if mp.OneLayerMap[newCoordinate.X] == nil {
		mp.OneLayerMap[newCoordinate.X] = make(map[int]*coordinate.Coordinate)
	}
	mp.OneLayerMap[newCoordinate.X][newCoordinate.Y] = newCoordinate
}

func (mp *Map) SetXYSize(Scale int) (int, int) {
	return mp.XSize / Scale, mp.YSize / Scale
}

// TODO GetMaxPriorityTexture, GetMaxPriorityObject близнецы
func (mp *Map) GetMaxPriorityTexture() int {
	max := 0

	for _, xLine := range mp.Flore {
		for _, flore := range xLine {
			if max < flore.TexturePriority {
				max = flore.TexturePriority
			}
		}
	}

	return max
}

func (mp *Map) GetMaxPriorityObject() int {
	mp.DynamicObjectsMX.RLock()
	defer mp.DynamicObjectsMX.RUnlock()

	max := 0

	for _, obj := range mp.StaticObjects {
		if max < obj.Priority {
			max = obj.Priority
		}
	}

	for _, obj := range mp.DynamicObjects {
		if max < obj.Priority {
			max = obj.Priority
		}
	}

	return max
}

type LvlMap struct {
	X     int     `json:"x"`
	Y     int     `json:"y"`
	Level float64 `json:"level"`
}

func (mp *Map) GetPosLevel(x, y int) (int, int, float64) {

	// х, у это реальный х/16 и реальный y/16
	mp.lvlMX.RLock()
	defer mp.lvlMX.RUnlock()

	x, y = x/16, y/16
	id := strconv.Itoa(x) + ":" + strconv.Itoa(y)
	lvl, ok := mp.LevelMap[id]
	if ok {
		return x, y, lvl.Level
	} else {
		return x, y, mp.DefaultLevel
	}
}

func (mp *Map) GetAllLvl() <-chan *LvlMap {

	lvlChan := make(chan *LvlMap, len(mp.LevelMap))

	go func() {
		mp.lvlMX.RLock()
		defer func() {
			mp.lvlMX.RUnlock()
			close(lvlChan)
		}()

		for _, lvl := range mp.LevelMap {
			lvlChan <- lvl
		}
	}()

	return lvlChan
}

func (mp *Map) AddLvl(x, y int, lvl float64) {
	// х, у это реальный х/16 и реальный y/16
	mp.lvlMX.Lock()
	defer mp.lvlMX.Unlock()

	id := strconv.Itoa(x) + ":" + strconv.Itoa(y)

	if lvl == mp.DefaultLevel {
		delete(mp.LevelMap, id)
	} else {
		mp.LevelMap[id] = &LvlMap{x, y, lvl}
	}
}

func (mp *Map) GetJSONStaticObjects() {
	if mp.StaticObjectsJSON == nil {
		mp.StaticObjectsJSON = make(map[int]string)

		mp.DynamicObjectsMX.Lock()
		defer mp.DynamicObjectsMX.Unlock()

		for id, obj := range mp.StaticObjects {
			mp.StaticObjectsJSON[id] = obj.GetJSON(int64(rand.Intn(999999999)))
		}
	}
}
