package bullets

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/bullet"
	"sync"
)

type mapBullets struct {
	bullets map[string]*bullet.Bullet
	mx      sync.RWMutex
}

func (s *mapBullets) AddBullet(newBullet *bullet.Bullet) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if s.bullets == nil {
		s.bullets = make(map[string]*bullet.Bullet)
	}

	s.bullets[newBullet.UUID] = newBullet
}

func (s *mapBullets) GetBulletByUUID(uuid string) *bullet.Bullet {
	s.mx.RLock()
	defer s.mx.RUnlock()

	if s.bullets == nil {
		return nil
	} else {
		return s.bullets[uuid]
	}
}

func (s *mapBullets) GetCopyMapBullets() map[string]*bullet.Bullet {
	bullets := make(map[string]*bullet.Bullet)

	s.mx.RLock()
	defer s.mx.RUnlock()

	for key, b := range s.bullets {
		bullets[key] = b
	}

	return bullets
}

func (s *mapBullets) RemoveBullet(bullet *bullet.Bullet) {
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.bullets, bullet.UUID)
}

func (s *mapBullets) GetBullets() <-chan *bullet.Bullet {

	s.mx.RLock()
	bullets := make(chan *bullet.Bullet, len(s.bullets))

	go func() {
		defer func() {
			close(bullets)
			s.mx.RUnlock()
		}()

		for _, mpBullet := range s.bullets {
			bullets <- mpBullet
		}
	}()

	return bullets
}

func (s *mapBullets) GetCountBullets() int {

	s.mx.RLock()
	defer s.mx.RUnlock()

	return len(s.bullets)
}

func (s *mapBullets) UnsafeRange() (map[string]*bullet.Bullet, *sync.RWMutex) {
	s.mx.RLock()
	return s.bullets, &s.mx
}
