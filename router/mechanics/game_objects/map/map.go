package _map

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/dynamic_map_object"
	"math/rand"
	"strconv"
	"sync"
)

type Map struct {
	Id            int     `json:"id"`
	TypeID        int     `json:"-"`
	Name          string  `json:"name"`
	XSize         int     `json:"x_size"`
	YSize         int     `json:"y_size"`
	DefaultLevel  float64 `json:"default_level"`
	Specification string  `json:"specification"`
	Type          string  `json:"-"`

	// текстуры земли
	Flore map[int]map[int]*dynamic_map_object.Flore `json:"flore"`
	// Игрок видит эти обьекты всегда независимо от радара/обзора
	StaticObjects   map[int]*dynamic_map_object.Object `json:"-"`
	StaticObjectsMX sync.RWMutex                       `json:"-"`
	// статисные обьекты не изменяются поэтому кешируем тут обьекты для фронтенда
	StaticObjectsJSON map[int][]byte `json:"static_objects_json"`
	// в матрице DynamicObjects находятся обьекты которые могут передвигатся/уничтожатся/рождатся
	// тоесть это обьекты с поведением, игрок их видит и запоминает последнее их состояние у себя.
	// Игрок не видит если с обьектом что либо произошло вне поле его зрения.
	DynamicObjects   []*dynamic_map_object.Object `json:"-"`
	DynamicObjectsMX sync.RWMutex                 `json:"-"`
	// вспомогательная мапа для уменьшение блокировок, сюда падают обьекты только с build=true
	DynamicBuildObjects   []*dynamic_map_object.Object `json:"-"`
	DynamicBuildObjectsMX sync.RWMutex                 `json:"-"`

	// разделяем карту на зоны (DiscreteSize х DiscreteSize) при загрузке сервера,
	// добавляем в зону все поинты которые пересекают данных квадрат и ближайшие к нему
	// когда надо найти колизию с юнитом делем его полизию на 100 и отбрасываем дровь так мы получим зону
	// например положение юнита 55/DiscreteSize:77/DiscreteSize = зона 0:0, 257/DiscreteSize:400/DiscreteSize = 1:1, 1654/DiscreteSize:2340/DiscreteSize = 6:9
	// и смотрим только те точки которые находятся в данной зоне
	GeoZones [][]*Zone `json:"-"`

	LevelMap [][]*LvlMap `json:"level_map"`
	lvlMX    sync.RWMutex
	LoopInit bool  `json:"-"`
	Time     int64 `json:"-"` // игровое время относительно сервер-тиков, с запуска карты
	Exit     bool  `json:"-"`
	mx       sync.RWMutex
}

type Spawn struct {
	ID                     int     `json:"id"`
	X                      int     `json:"x"`
	Y                      int     `json:"y"`
	Name                   string  `json:"name"`
	Radius                 int     `json:"radius"`
	Rotate                 int     `json:"-"`
	Type                   string  `json:"-"`
	CaptureTeam            int     `json:"capture_team"`
	Capture                float64 `json:"capture"`
	CaptureFact            bool    `json:"capture_fact"`
	Close                  bool    `json:"close"` // закрытую базу нельзя захватить
	IgnoreTeamID           int     `json:"ignore_team_id"`
	FireproofAmount        []int   `json:"-"`
	CaptureCountSingleUnit float64 `json:"-"`
	SingleUnit             bool    `json:"single_unit"`
	ReloadTime             int     `json:"-"`
	Owner                  string  `json:"-"`
	Transport              bool    `json:"-"`
	TransportX             int     `json:"-"`
	TransportY             int     `json:"-"`
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

func (mp *Map) GetChanBuildDynamicObjects() <-chan *dynamic_map_object.Object {

	mp.DynamicBuildObjectsMX.RLock()
	objChan := make(chan *dynamic_map_object.Object, len(mp.DynamicBuildObjects))

	go func() {
		defer func() {
			mp.DynamicBuildObjectsMX.RUnlock()
			close(objChan)
		}()

		for _, obj := range mp.DynamicBuildObjects {
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

//func (mp *Map) DynamicObjectsUnsafeRange() (map[int]*dynamic_map_object.Object, *sync.RWMutex) {
//	mp.DynamicObjectsMX.RLock()
//	return mp.DynamicObjects, &mp.DynamicObjectsMX
//}

func (mp *Map) GetCopyArrayBuildDynamicObjects(basket []*dynamic_map_object.Object) []*dynamic_map_object.Object {

	mp.DynamicBuildObjectsMX.RLock()
	defer mp.DynamicBuildObjectsMX.RUnlock()

	basket = basket[:0]
	for _, obj := range mp.DynamicBuildObjects {
		basket = append(basket, obj)
	}

	return basket
}

func (mp *Map) GetDynamicObjects(x, y int) *dynamic_map_object.Object {
	idString := strconv.Itoa(x) + strconv.Itoa(y)
	id, _ := strconv.Atoi(idString)
	return mp.GetDynamicObjectsByID(id)
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

	if x < 0 || y < 0 || len(mp.LevelMap) <= x || len(mp.LevelMap[x]) <= y {
		return x, y, mp.DefaultLevel
	}

	lvl := mp.LevelMap[x][y]
	if lvl != nil {
		return x, y, lvl.Level
	} else {
		return x, y, mp.DefaultLevel
	}
}

func (mp *Map) AddLvl(x, y int, lvl float64) {
	// х, у это реальный х/16 и реальный y/16
	mp.lvlMX.Lock()
	defer mp.lvlMX.Unlock()

	if x < 0 || y < 0 || len(mp.LevelMap) <= x || len(mp.LevelMap[x]) <= y {
		return
	}

	mp.LevelMap[x][y] = &LvlMap{x, y, lvl}
}

func (mp *Map) GetCopyArrayDynamicObjects() (copyArray []*dynamic_map_object.Object) {

	copyArray = make([]*dynamic_map_object.Object, 0, len(mp.DynamicObjects))

	mp.DynamicObjectsMX.RLock()
	defer mp.DynamicObjectsMX.RUnlock()

	for _, obj := range mp.DynamicObjects {
		if obj != nil {
			copyArray = append(copyArray, obj)
		}
	}

	return copyArray
}

func (mp *Map) GetJSONStaticObjects() {
	if mp.StaticObjectsJSON == nil {
		mp.StaticObjectsJSON = make(map[int][]byte)

		mp.DynamicObjectsMX.Lock()
		defer mp.DynamicObjectsMX.Unlock()

		for id, obj := range mp.StaticObjects {
			mp.StaticObjectsJSON[id] = obj.GetBytes(int64(rand.Intn(999999999)))
		}
	}
}
