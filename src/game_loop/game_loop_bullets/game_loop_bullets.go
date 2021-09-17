package game_loop_bullets

import (
	"fmt"
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/mechanics/factories/bullets"
	"github.com/TrashPony/game_engine/src/mechanics/fly_bullets"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/bullet"
	_map "github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
	"github.com/TrashPony/game_engine/src/web_socket"
	"sync"
)

func Bullet(mp *_map.Map, bullets map[string]*bullet.Bullet, units []*unit.Unit) (*web_socket.GameLoopMessages, *web_socket.GameLoopMessages, *web_socket.GameLoopMessages) {

	FlyMsgs := &web_socket.GameLoopMessages{}
	FlyLaserMsgs := &web_socket.GameLoopMessages{}
	ExplosionsMsgs := &web_socket.GameLoopMessages{}

	var wg sync.WaitGroup

	for _, flyBullet := range bullets {
		if flyBullet.MapID == mp.Id {
			wg.Add(1)
			fly(mp, flyBullet, &wg, FlyMsgs, FlyLaserMsgs, ExplosionsMsgs, units)
		}
	}

	wg.Wait()

	return FlyMsgs, FlyLaserMsgs, ExplosionsMsgs
}

func fly(mp *_map.Map, flyBullet *bullet.Bullet, wg *sync.WaitGroup, FlyMsgs, FlyLaserMsgs, ExplosionsMsgs *web_socket.GameLoopMessages, units []*unit.Unit) {
	defer wg.Done()

	if flyBullet.Ammo.Type == _const.LaserWeapon {
		_, _, flyLaserMsg, damageObject := fly_bullets.FlyLaser(flyBullet, mp, false)
		damage(damageObject)

		FlyLaserMsgs.AddMessage(flyLaserMsg)
	}

	// коректировка ракет
	if flyBullet.Ammo.Type == _const.MissileWeapon {
		Missile(flyBullet, mp)
	}

	if flyBullet.Ammo.Type == _const.FirearmsWeapon || flyBullet.Ammo.Type == _const.MissileWeapon {

		flyMsg, explosionMsg, damageObject, crater := fly_bullets.BulletFlyTick(flyBullet, mp, false, units)
		damage(damageObject)

		if flyMsg != nil {
			FlyMsgs.AddMessage(flyMsg)
		}

		if explosionMsg != nil {
			ExplosionsMsgs.AddMessage(explosionMsg)
		}

		mp.AddCrater(crater)
	}

	if flyBullet.GetEnd() {
		// обязательно в отдельном потоке
		bullets.Bullets.RemoveBullet(flyBullet)
	}
}

func damage(damageObject []*fly_bullets.DamageObject) {
	//наносим урона обьектам
	// TODO
	fmt.Println(damageObject)
	// TODO global.SendDeadObject(attack.DamageObjects(damageObject, mp, flyBullet), mp, flyBullet.GetX(), flyBullet.GetY())
	//создаем кратер после взрыва если он образовался на земле
}
