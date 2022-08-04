package _map

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/dynamic_map_object"
	"math/rand"
	"strconv"
	"sync"
)

type Map struct {
	Id           int     `json:"id"`
	XSize        int     `json:"x_size"`
	YSize        int     `json:"y_size"`
	DefaultLevel float64 `json:"default_level"`

	// текстуры земли
	Flore map[int]map[int]*dynamic_map_object.Flore `json:"flore"`
	// не изменяемые обьекты которые игрок видит эти обьекты всегда независимо от радара/обзора
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

func (mp *Map) GetFlore(x, y int) *dynamic_map_object.Flore {
	flore := mp.Flore[x][y]
	return flore
}

func (mp *Map) SetXYSize(Scale int) (int, int) {
	return mp.XSize / Scale, mp.YSize / Scale
}

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
