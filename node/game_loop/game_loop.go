package game_loop

import (
	"github.com/TrashPony/game-engine/node/game_loop/game_loop_bullets"
	"github.com/TrashPony/game-engine/node/game_loop/game_loop_gun"
	"github.com/TrashPony/game-engine/node/game_loop/game_loop_move"
	"github.com/TrashPony/game-engine/node/game_loop/game_loop_structure"
	"github.com/TrashPony/game-engine/node/game_loop/game_loop_view"
	"github.com/TrashPony/game-engine/node/mechanics/factories/bullets"
	"github.com/TrashPony/game-engine/node/mechanics/factories/quick_battles"
	units2 "github.com/TrashPony/game-engine/node/mechanics/factories/units"
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/bullet"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/dynamic_map_object"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/player"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"github.com/TrashPony/game-engine/router/web_socket"
	"time"
)

func GameLoopInit() {
	for {
		tick()
		time.Sleep(time.Second)
	}
}

func tick() {
	for _, b := range quick_battles.Battles.GetAll(make([]*battle.Battle, 0)) {

		if b.Map.LoopInit {
			continue
		}

		time.Sleep(time.Millisecond) // что бы сектора работали не одновременно
		go GameLoop(b)
	}
}

// каждая карта являет отдельным миром, поэтому на каждую карту свой GameLoop
func GameLoop(gameBattle *battle.Battle) {

	if gameBattle.Map.LoopInit {
		return
	}

	gameBattle.Map.LoopInit = true
	defer func() {
		gameBattle.Map.LoopInit = false
	}()

	messagesStore := &web_socket.MessagesStore{}
	unitBasket := make([]*unit.Unit, 0)
	mapBulletsBasket := make([]*bullet.Bullet, 0)
	playersBasket := make([]*player.Player, 0)
	buildObjectsBasket := make([]*dynamic_map_object.Object, 0)

	for {

		if gameBattle.Map.Exit {
			return
		}

		startTime := time.Now()
		messagesStore.Clear()

		unitBasket = units2.Units.GetAllUnitsArray(gameBattle.Map.Id, unitBasket)
		mapBulletsBasket = bullets.Bullets.GetCopyArrayBullets(gameBattle.Map.Id, mapBulletsBasket)
		playersBasket = gameBattle.GetPlayers(playersBasket) // карта привязана к определенной катке ищем катку по карте и отуда берем игроков
		buildObjectsBasket = gameBattle.Map.GetCopyArrayBuildDynamicObjects(buildObjectsBasket)

		game_loop_move.SetUnitsPos(unitBasket)
		game_loop_move.SetObjectsPos(gameBattle.Map)

		game_loop_move.Unit(gameBattle.Map, unitBasket, gameBattle, messagesStore)
		game_loop_move.Objects(gameBattle, unitBasket, messagesStore)

		game_loop_bullets.Bullet(gameBattle.Map, gameBattle, mapBulletsBasket, unitBasket, messagesStore)

		game_loop_structure.Structure(gameBattle, buildObjectsBasket)

		game_loop_view.View(gameBattle.Map, unitBasket, mapBulletsBasket, gameBattle, messagesStore)

		game_loop_gun.Object(messagesStore, gameBattle, unitBasket, buildObjectsBasket)
		game_loop_gun.Unit(unitBasket, gameBattle, messagesStore)

		SendMessagesToUsers(gameBattle, playersBasket, messagesStore, unitBasket)

		deltaTime := int(time.Since(startTime).Nanoseconds() / int64(time.Millisecond))
		if deltaTime > _const.ServerTick {
			println("ID: ", gameBattle.Map.Id, " deltaTime: ", deltaTime)
		}

		if deltaTime < _const.ServerTick {
			time.Sleep(time.Millisecond * time.Duration(_const.ServerTick-deltaTime))
		}

		gameBattle.Map.Time += _const.ServerTick
	}
}
