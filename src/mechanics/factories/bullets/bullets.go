package bullets

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/bullet"
	"sync"
	"sync/atomic"
)

var Bullets = initStore()

// todo трефачить как убежусь шо работаед
type store struct {
	mx      sync.RWMutex
	bullets map[int]*mapBullets
	lastID  int64 // типо уникальный ид
}

func initStore() *store {
	return &store{bullets: make(map[int]*mapBullets)}
}

func (s *store) AddBullet(bullet *bullet.Bullet) {

	atomic.AddInt64(&s.lastID, 1)

	s.mx.RLock()
	mapStore := s.bullets[bullet.MapID]
	s.mx.RUnlock()

	bullet.SetID(int(atomic.LoadInt64(&s.lastID)))

	if mapStore == nil {

		mapStore = &mapBullets{}

		s.mx.Lock()
		s.bullets[bullet.MapID] = mapStore
		s.mx.Unlock()
	}

	mapStore.AddBullet(bullet)
}

func (s *store) RemoveBullet(bullet *bullet.Bullet) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	if s.bullets[bullet.MapID] != nil {
		s.bullets[bullet.MapID].RemoveBullet(bullet)
	}
}

func (s *store) GetCount() int {
	s.mx.RLock()
	defer s.mx.RUnlock()
	count := 0

	for _, bullets := range s.bullets {
		count += bullets.GetCountBullets()
	}

	return count
}

func (s *store) GetBulletsByMapID(mapID int) <-chan *bullet.Bullet {

	s.mx.RLock()
	defer s.mx.RUnlock()

	if s.bullets[mapID] != nil {
		return s.bullets[mapID].GetBullets()
	} else {
		bullets := make(chan *bullet.Bullet, 0)
		close(bullets)
		return bullets
	}
}

func (s *store) GetCopyMapBullets(mapID int) map[string]*bullet.Bullet {
	s.mx.RLock()
	defer s.mx.RUnlock()

	if s.bullets[mapID] != nil {
		return s.bullets[mapID].GetCopyMapBullets()
	}

	return nil
}

func (s *store) GetBulletByMapIDAndUUID(mapID int, uuid string) *bullet.Bullet {
	s.mx.RLock()
	defer s.mx.RUnlock()

	if s.bullets[mapID] == nil {
		return nil
	} else {
		return s.bullets[mapID].GetBulletByUUID(uuid)
	}
}

func (s *store) UnsafeRangeByMapID(mapID int) (map[string]*bullet.Bullet, *sync.RWMutex) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	if s.bullets[mapID] == nil {
		return nil, nil
	} else {
		return s.bullets[mapID].UnsafeRange()
	}
}
