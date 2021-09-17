package game_loop

import (
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/game_loop/game_loop_bullets"
	"github.com/TrashPony/game_engine/src/game_loop/game_loop_gun"
	"github.com/TrashPony/game_engine/src/game_loop/game_loop_move"
	"github.com/TrashPony/game_engine/src/game_loop/game_loop_terrain"
	"github.com/TrashPony/game_engine/src/game_loop/game_loop_view"
	"github.com/TrashPony/game_engine/src/mechanics/factories/bullets"
	"github.com/TrashPony/game_engine/src/mechanics/factories/maps"
	"github.com/TrashPony/game_engine/src/mechanics/factories/quick_battles"
	units2 "github.com/TrashPony/game_engine/src/mechanics/factories/units"
	_map "github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/web_socket"
	"time"
)

func GameLoopInit() {
	for {

		for mp := range maps.Maps.GetAllMap() {

			if mp.LoopInit {
				continue
			}

			if mp.Id < 0 {
				time.Sleep(time.Millisecond) // что бы сектора работали не одновременно
				go GameLoop(mp)
			}
		}

		time.Sleep(time.Second)
	}
}

// каждая карта являет отдельным миром, поэтому на каждую карту свой гейм лооп, возможно для бесшовного мира это будет не так
func GameLoop(mp *_map.Map) {

	if mp.LoopInit {
		return
	}

	mp.LoopInit = true
	defer func() {
		mp.LoopInit = false
	}()

	gameBattle := quick_battles.Battles.GetBattleByMapID(mp.Id)
	if gameBattle == nil {
		return
	}

	for {

		if mp.Exit {
			return
		}

		startTime := time.Now()

		units := units2.Units.GetAllUnitsArray(mp.Id)
		dynObjects, _ := mp.GetCopyMapDynamicObjects()
		mapBullets := bullets.Bullets.GetCopyMapBullets(mp.Id)
		// карта привязана к определенной катке ищем катку по карте и отуда берем игроков
		players := gameBattle.GetPlayers()

		game_loop_move.SetUnitsPos(units)
		unitMoveMsg := game_loop_move.Unit(mp, units, dynObjects)
		FlyBulletsMsgs, FlyLaserMsgs, ExplosionsMsgs := game_loop_bullets.Bullet(mp, mapBullets, units)

		game_loop_terrain.TerrainLife(mp, dynObjects)                   // не отправляем изменения потому что они отслеживаются в пакете game_loop_view
		game_loop_view.View(mp, players, units, dynObjects, mapBullets) // тут мы не возвращаем сообщения т.к. внутри делается для каждого юзера индивидуально

		fireMsgs, rotateMsgs := game_loop_gun.Unit(mp, units)

		SendMessagesToUsers(mp, unitMoveMsg, FlyBulletsMsgs, FlyLaserMsgs, ExplosionsMsgs, fireMsgs, rotateMsgs)

		deltaTime := int(time.Since(startTime).Nanoseconds() / int64(time.Millisecond))
		if deltaTime > _const.ServerTick {
			println("ID: ", mp.Id, " deltaTime: ", deltaTime)
		}

		if deltaTime < _const.ServerTick {
			time.Sleep(time.Millisecond * time.Duration(_const.ServerTick-deltaTime))
		}

		mp.Time += _const.ServerTick
	}
}

func SendMessagesToUsers(mp *_map.Map, unitMoveMsg, FlyBulletsMsgs, FlyLaserMsgs, ExplosionsMsgs, fireMsgs, rotateMsgs *web_socket.GameLoopMessages) {
	if len(unitMoveMsg.Messages) > 0 {
		web_socket.SendMessage(web_socket.Response{
			Event: "um", // unit move
			Data:  unitMoveMsg.Messages,
			MapID: mp.Id,
		})
	}

	if len(FlyBulletsMsgs.Messages) > 0 {
		web_socket.SendMessage(web_socket.Response{
			Event: "fb", // fly bullets
			Data:  FlyBulletsMsgs.Messages,
			MapID: mp.Id,
		})
	}

	if len(FlyLaserMsgs.Messages) > 0 {
		web_socket.SendMessage(web_socket.Response{
			Event: "fl", // fly laser
			Data:  FlyLaserMsgs.Messages,
			MapID: mp.Id,
		})
	}

	if len(ExplosionsMsgs.Messages) > 0 {
		web_socket.SendMessage(web_socket.Response{
			Event: "eb", // explosion bullets
			Data:  ExplosionsMsgs.Messages,
			MapID: mp.Id,
		})
	}

	if len(fireMsgs.Messages) > 0 {
		web_socket.SendMessage(web_socket.Response{
			Event: "ufw", // unit fire weapon
			Data:  fireMsgs.Messages,
			MapID: mp.Id,
		})
	}

	if len(rotateMsgs.Messages) > 0 {
		web_socket.SendMessage(web_socket.Response{
			Event: "uwr", // unit weapon rotate
			Data:  rotateMsgs.Messages,
			MapID: mp.Id,
		})
	}
}
