package burst_of_shots

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/bullet"
	"sync"
)

type BurstOfShots struct {
	bullets    map[int][]*bullet.Bullet //[weapon_slot_number] bullets
	currentPos map[int]int              //[weapon_slot_number] текущая позиция в масиве пуль
	mx         sync.RWMutex
}

func (b *BurstOfShots) AddBullets(weaponSlot int, bullets []*bullet.Bullet) {
	b.mx.Lock()
	defer b.mx.Unlock()

	if b.bullets == nil {
		b.bullets = make(map[int][]*bullet.Bullet)
		b.currentPos = make(map[int]int)
	}

	b.bullets[weaponSlot] = bullets
}

func (b *BurstOfShots) GetBullets(weaponSlot int) ([]*bullet.Bullet, int) {
	b.mx.RLock()
	defer b.mx.RUnlock()

	if b.bullets == nil {
		return nil, 0
	}

	return b.bullets[weaponSlot], b.currentPos[weaponSlot]
}

func (b *BurstOfShots) ChangePos(weaponSlot, newPos int) {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.currentPos[weaponSlot] = newPos
}
