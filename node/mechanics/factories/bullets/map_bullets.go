package bullets

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/bullet"
	"sync"
)

type mapBullets struct {
	bullets []*bullet.Bullet
	mx      sync.RWMutex
}

func (s *mapBullets) AddBullet(newBullet *bullet.Bullet) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if s.bullets == nil {
		s.bullets = make([]*bullet.Bullet, 0)
	}

	s.bullets = append(s.bullets, newBullet)
}

func (s *mapBullets) GetCopyArrayBullets(basket []*bullet.Bullet) []*bullet.Bullet {

	s.mx.RLock()
	defer s.mx.RUnlock()

	basket = basket[:0]
	for _, b := range s.bullets {
		basket = append(basket, b)
	}

	return basket
}

func (s *mapBullets) RemoveBullet(bullet *bullet.Bullet) {
	s.mx.Lock()
	defer s.mx.Unlock()

	index := -1
	for i, u := range s.bullets {
		if u.ID == bullet.ID {
			index = i
			break
		}
	}

	if index >= 0 {
		s.bullets[index] = s.bullets[len(s.bullets)-1]
		s.bullets = s.bullets[:len(s.bullets)-1]
	}
}
