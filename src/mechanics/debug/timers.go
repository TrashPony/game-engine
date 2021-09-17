package debug

import (
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"
)

type timersPool struct {
	timers map[int]map[string]*Timer
	mx     sync.RWMutex
}

type Timer struct {
	UUID      string    `json:"uuid"`
	StartTime time.Time `json:"start_time"`
	WorkMS    int64     `json:"work_ms"`
	FuncName  string    `json:"func_name"`
}

func (t *timersPool) Start(mapID int, funcName string) string {
	timerUUID := uuid.NewV4().String()

	if mapID != -1 {
		return ""
	}

	go func() {
		t.mx.Lock()
		defer t.mx.Unlock()

		if t.timers[mapID] == nil {
			t.timers[mapID] = make(map[string]*Timer)
		}

		t.timers[mapID][timerUUID] = &Timer{UUID: timerUUID, FuncName: funcName, StartTime: time.Now()}
	}()

	return timerUUID
}

func (t *timersPool) Stop(mapID int, uuid string) {

	if uuid == "" {
		return
	}

	go func() {
		t.mx.RLock()
		defer t.mx.RUnlock()

		if t.timers[mapID] == nil {
			t.timers[mapID] = make(map[string]*Timer)
		}

		timer, ok := t.timers[mapID][uuid]
		if ok {
			timer.WorkMS = time.Since(timer.StartTime).Milliseconds()
		}
	}()
}

func (t *timersPool) GetAll(mapID int) []*Timer {
	t.mx.RLock()
	defer t.mx.RUnlock()

	timers := make([]*Timer, 0)

	for _, timer := range t.timers[mapID] {
		if timer.WorkMS > 0 {
			timers = append(timers, timer)
		}
	}

	return timers
}
