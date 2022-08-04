package game_loop_bullets

import (
	"github.com/TrashPony/game-engine/node/binary_msg"
	damage2 "github.com/TrashPony/game-engine/node/mechanics/damage"
	"github.com/TrashPony/game-engine/node/mechanics/factories/bullets"
	"github.com/TrashPony/game-engine/node/mechanics/fly_bullets"
	_const "github.com/TrashPony/game-engine/router/const"
	battle2 "github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/bullet"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/web_socket_response"
	"github.com/TrashPony/game-engine/router/web_socket"
	"sync"
)

func Bullet(mp *_map.Map, b *battle2.Battle, bullets []*bullet.Bullet, units []*unit.Unit, ms *web_socket.MessagesStore) {

	var wg sync.WaitGroup

	for _, flyBullet := range bullets {
		if flyBullet.MapID == mp.Id {
			wg.Add(1)
			fly(mp, b, flyBullet, &wg, ms, units)
		}
	}

	wg.Wait()

	return
}

func fly(mp *_map.Map, b *battle2.Battle, flyBullet *bullet.Bullet, wg *sync.WaitGroup, ms *web_socket.MessagesStore, units []*unit.Unit) {
	defer wg.Done()

	if flyBullet.Ammo.Type == _const.LaserWeapon {
		_, _, flyLaserMsg, damageObject := fly_bullets.FlyLaser(flyBullet, mp, false)
		damage(damageObject, b, flyBullet, ms)

		ms.AddMsg("flyLaserMsgs", "bin", web_socket_response.Response{
			BinaryMsg: binary_msg.CreateBulletLaserFly(flyLaserMsg.TypeID, flyLaserMsg.X, flyLaserMsg.Y, flyLaserMsg.ToX, flyLaserMsg.ToY, flyBullet.OwnerID),
			X:         flyLaserMsg.X,
			Y:         flyLaserMsg.Y,
			ToX:       flyLaserMsg.ToX,
			ToY:       flyLaserMsg.ToY,
			CheckTo:   true,
		}, nil)
	}

	if flyBullet.Ammo.Type == _const.FirearmsWeapon || flyBullet.Ammo.Type == _const.MissileWeapon {

		flyMsg, explosionMsg, damageObject, crater := fly_bullets.BulletFlyTick(flyBullet, b, false, units)
		damage(damageObject, b, flyBullet, ms)

		if flyMsg != nil {
			ms.AddMsg("flyMsgs", "move", web_socket_response.Response{
				ID:        flyBullet.ID,
				BinaryMsg: binary_msg.CreateBulletBinaryFly(flyMsg.TypeID, flyMsg.ID, flyMsg.X, flyMsg.Y, flyMsg.Z, flyMsg.MS, flyMsg.Rotate),
				X:         flyMsg.X,
				Y:         flyMsg.Y,
			}, map[string]string{"type_obj": "bullet"})
		}

		if explosionMsg != nil {
			ms.AddMsg("explosionsMsgs", "bin", web_socket_response.Response{
				BinaryMsg: binary_msg.CreateBulletBinaryExplosion(explosionMsg.TypeID, explosionMsg.X, explosionMsg.Y, explosionMsg.Z),
				X:         explosionMsg.X,
				Y:         explosionMsg.Y,
			}, nil)
		}

		mp.AddCrater(crater)
	}

	if flyBullet.GetEnd() {
		bullets.Bullets.RemoveBullet(flyBullet)
	}
}

func damage(damageObjects []*damage2.Object, b *battle2.Battle, damageBullet *bullet.Bullet, ms *web_socket.MessagesStore) {

	if len(damageObjects) == 0 {
		return
	}

	//наносим урона обьектам
	weaponID := 0
	if damageBullet.Weapon != nil {
		weaponID = damageBullet.Weapon.ID
	}

	for _, do := range damage2.Damage(damageObjects, b, damageBullet.OwnerType, damageBullet.OwnerID,
		damageBullet.GetX(), damageBullet.GetY(), damageBullet.Ammo.ID, weaponID, damageBullet.EquipID) {

		if (do.Obj == nil && do.TypeTarget != "shield") || do.TypeTarget == "flore" || do.TypeTarget == "static_object" {
			continue
		}

		sendDamageMsg(b.Map, do, ms)
	}
}

func sendDamageMsg(mp *_map.Map, do *damage2.Object, ms *web_socket.MessagesStore) {
	// акамулируем сообщения и отдаем на отдачу
	ms.AddMsg("objectDamageMsgs", "bin", web_socket_response.Response{
		X:         do.X,
		Y:         do.Y,
		BinaryMsg: binary_msg.DamageTextBinaryMsg(do.X, do.Y, do.Damage, do.TypeTarget),
	}, nil)

	if do.Dead {
		ms.AddMsg("objectDeadMsgs", "bin", web_socket_response.Response{
			X:         do.X,
			Y:         do.Y,
			BinaryMsg: binary_msg.ObjectDeadBinaryMsg(do.IdTarget, do.X, do.Y, do.TypeTarget),
		}, nil)
	}
}
