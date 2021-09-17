package debug

import (
	"sync"
)

var ServerState = newServerStateStore()

// TODO брать все стату только в режиме дебага
type serverState struct {
	bots             int // всего ботов в игре
	calcMove         int // сколько юнитов пытается искать путь
	aiMove           int // сколько ботов пользуется рулем
	meteoriteWaitEnd int
	mapTick          map[int][]int // время каждого тика в gameLoop для каждой карты [map_id][]ms
	messagesCount    map[string]int
	mx               sync.RWMutex
	Timers           *timersPool
}

func newServerStateStore() *serverState {
	return &serverState{
		mapTick:       make(map[int][]int),
		messagesCount: make(map[string]int),
		Timers:        &timersPool{timers: make(map[int]map[string]*Timer)},
	}
}

func (s *serverState) AddMessage(service string) {
	go func() {
		s.mx.Lock()
		defer s.mx.Unlock()
		s.messagesCount[service]++
	}()
}

func (s *serverState) RemoveMessage(service string) {
	go func() {
		s.mx.Lock()
		defer s.mx.Unlock()
		s.messagesCount[service]--
	}()
}

func (s *serverState) GetMessagesQueueCount() map[string]int {

	s.mx.Lock()
	defer s.mx.Unlock()

	queues := make(map[string]int)
	if s.messagesCount == nil {
		s.messagesCount = make(map[string]int)
	}

	for serves, count := range s.messagesCount {
		queues[serves] = count
	}

	return queues
}

func (s *serverState) GetMeteoriteWaitEnd() int {
	return s.meteoriteWaitEnd
}

func (s *serverState) AddMeteoriteWaitEnd() {
	go func() {
		s.mx.Lock()
		defer s.mx.Unlock()
		s.meteoriteWaitEnd++
	}()
}

func (s *serverState) RemoveMeteoriteWaitEnd() {
	go func() {
		s.mx.Lock()
		defer s.mx.Unlock()
		s.meteoriteWaitEnd--
	}()
}

func (s *serverState) GetAiMove() int {
	return s.aiMove
}

func (s *serverState) SetAiMove(count int) {
	go func() {
		s.mx.Lock()
		defer s.mx.Unlock()
		s.aiMove = count
	}()
}

func (s *serverState) GetCountBots() int {
	return s.bots
}

func (s *serverState) SetBots(count int) {
	go func() {
		s.mx.Lock()
		defer s.mx.Unlock()
		s.bots = count
	}()
}

func (s *serverState) GetCountCalcMove() int {
	return s.calcMove
}

func (s *serverState) SetCalcMove(count int) {
	go func() {
		s.mx.Lock()
		defer s.mx.Unlock()
		s.calcMove = count
	}()
}

func (s *serverState) GetMapsTick() map[int][]int {

	s.mx.Lock()
	defer s.mx.Unlock()

	ticks := make(map[int][]int)
	if s.mapTick == nil {
		s.mapTick = make(map[int][]int)
	}

	for mp, mpTicks := range s.mapTick {
		ticks[mp] = mpTicks
	}

	return ticks
}

func (s *serverState) AddMapTick(mapID int, tickMS int) {
	go func() {
		s.mx.Lock()
		defer s.mx.Unlock()

		if s.mapTick[mapID] == nil {
			s.mapTick[mapID] = make([]int, 0)
		}

		s.mapTick[mapID] = append(s.mapTick[mapID], tickMS)
	}()
}
