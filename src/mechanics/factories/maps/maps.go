package maps

import (
	dbMap "github.com/TrashPony/game_engine/src/mechanics/db/maps"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"math/rand"
	"sync"
)

type mapStore struct {
	maps            map[int]*_map.Map
	mapsMX          sync.RWMutex
	generatedMapsID int
}

var Maps = newMapStore()

func newMapStore() *mapStore {
	m := &mapStore{
		maps: dbMap.Maps(),
	}

	return m
}

func (m *mapStore) GetByID(id int) (*_map.Map, bool) {

	m.mapsMX.RLock()
	defer m.mapsMX.RUnlock()

	newMap, ok := m.maps[id]
	return newMap, ok
}

type MapPosition struct {
	X   int       `json:"x"`
	Y   int       `json:"y"`
	Map *_map.Map `json:"map"`
}

func (m *mapStore) GetAllMap() <-chan *_map.Map {
	m.mapsMX.RLock()

	mapChan := make(chan *_map.Map, len(m.maps))

	go func() {
		defer func() {
			m.mapsMX.RUnlock()
			close(mapChan)
		}()

		for _, mp := range m.maps {
			mapChan <- mp
		}
	}()

	return mapChan
}

func (m *mapStore) GetAllShortInfoMap() map[int]*_map.ShortInfoMap {

	m.mapsMX.Lock()
	defer m.mapsMX.Unlock()

	shortMap := make(map[int]*_map.ShortInfoMap)
	for _, mp := range m.maps {
		if mp.Id > 0 {
			shortMap[mp.Id] = mp.GetShortInfoMap()
		}
	}
	return shortMap
}

func (m *mapStore) GetRandomMap() *_map.Map {
	m.mapsMX.Lock()
	defer m.mapsMX.Unlock()

	count := 0
	count2 := rand.Intn(len(m.maps))
	for _, mp := range m.maps {
		if count == count2 {
			return mp
		}
		count++
	}
	return nil
}

func (m *mapStore) GetNewMapID() int {
	m.generatedMapsID--
	return m.generatedMapsID
}

func (m *mapStore) AddNewMap(mp *_map.Map) {
	m.mapsMX.Lock()
	defer m.mapsMX.Unlock()
	m.maps[mp.Id] = mp
}

func (m *mapStore) RemoveMap(mp *_map.Map) {
	m.mapsMX.Lock()
	defer m.mapsMX.Unlock()
	delete(m.maps, mp.Id)
}
