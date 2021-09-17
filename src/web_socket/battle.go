package web_socket

import (
	"github.com/TrashPony/game_engine/src/ai"
	"github.com/TrashPony/game_engine/src/mechanics/factories/maps"
	"github.com/TrashPony/game_engine/src/mechanics/factories/players"
	"github.com/TrashPony/game_engine/src/mechanics/factories/quick_battles"
	"github.com/TrashPony/game_engine/src/mechanics/factories/units"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/behavior_rule"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/coordinate"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/player"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/target"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/user"
	"github.com/gorilla/websocket"
	"math/rand"
)

func BattleService(ws *websocket.Conn, gameUser *user.User, req Request) {
	p := players.Users.GetPlayer(gameUser.GetCurrentPlayerID(), gameUser.GetID())
	if p == nil {
		SendMessage(Response{Event: req.Event, Error: "no player", UserID: gameUser.ID})
		return
	}

	battle := quick_battles.Battles.GetBattleByUUID(p.GameUUID)
	if battle == nil {
		SendMessage(Response{Event: "Error", Error: "no battle", UserID: gameUser.ID})
		return
	}

	if req.Event == "InitBattle" {
		p.SetReady(false)
		battle.Map.GetJSONStaticObjects()
		p.MapID = battle.Map.Id

		SendMessage(Response{
			Event:  "InitBattle",
			UserID: gameUser.GetID(),
			Data: struct {
				Maps     map[int]*maps.MapPosition `json:"maps"`
				Player   *player.Player            `json:"player"`
				Spectrum bool                      `json:"spectrum"`
			}{
				// todo небольшой костыль из за клиента который взят из другой игоры
				Maps:     map[int]*maps.MapPosition{battle.Map.Id: {X: 0, Y: 0, Map: battle.Map}},
				Player:   p,
				Spectrum: false,
			}})
	}

	if req.Event == "StartLoad" {

		dynamicObjects := make(map[int]string)
		for memoryObj := range p.GetMapDynamicObjects(battle.Map.Id) {
			dynamicObjects[memoryObj.IDObject] = memoryObj.ObjectJSON
		}

		p.InitVisibleObjects()
		p.SetReady(true)

		SendMessage(Response{
			Event:  "RefreshDynamicObj",
			UserID: gameUser.GetID(),
			Data:   dynamicObjects,
		})
	}

	if req.Event == "CreateBot" {
		bot := ai.CreateBot(battle.UUID, battle.Map, rand.Intn(battle.Map.XSize), rand.Intn(battle.Map.YSize), rand.Intn(360), &behavior_rule.BehaviorRules{
			Rules: []*behavior_rule.BehaviorRule{
				{
					Action: "find_hostile_in_range_view",
					PassRule: &behavior_rule.BehaviorRule{
						Action: "follow_attack_target",
					},
					StopRule: &behavior_rule.BehaviorRule{
						Action: "scouting",
					},
				},
			},
			Meta: &behavior_rule.Meta{},
		})

		for u := range bot.GetUnitsStore().RangeUnits() {
			units.Units.AddUnit(u)
		}

		ai.AddLifeBot(bot)
	}

	if req.Event == "CreateUnit" {
		newUnit := &unit.Unit{ID: units.Units.GetBotID(), OwnerID: p.GetID(), MapID: battle.Map.Id, HP: 100}
		newUnit.GetPhysicalModel().SetPos(float64(rand.Intn(battle.Map.XSize)), float64(rand.Intn(battle.Map.YSize)), float64(rand.Intn(360)))
		p.GetUnitsStore().AddUnit(newUnit)
		units.Units.AddUnit(newUnit)
	}

	if req.Event == "MoveTo" {

		for _, id := range req.SelectUnits {
			gameUnit := units.Units.GetUnitByIDAndMapID(id, battle.Map.Id)
			if gameUnit != nil {
				gameUnit.MovePath = &unit.MovePath{
					NeedFindPath: true,
					Path:         &[]*coordinate.Coordinate{{X: req.X, Y: req.Y}},
					FollowTarget: &target.Target{X: req.X, Y: req.Y, Radius: 50},
				}

				gameUnit.SetWeaponTarget(&target.Target{
					Type: "map",
					X:    req.X,
					Y:    req.Y,
				})
			}
		}
	}

	//if req.Event == "i" {
	//
	//	for _, id := range req.SelectUnits {
	//		gameUnit := units.Units.GetUnitByIDAndMapID(id, battle.Map.Id)
	//		if gameUnit != nil {
	//
	//			// одновременно может быть нажато только 2 клавишы, если зажаты 4 то эфект нигилируется
	//			if req.W && req.S {
	//				req.W = false
	//				req.S = false
	//			}
	//
	//			if req.A && req.D {
	//				req.A = false
	//				req.D = false
	//			}
	//
	//			gameUnit.GetPhysicalModel().SetWASD(req.W, req.A, req.S, req.D)
	//
	//			gameUnit.SetWeaponTarget(&target.Target{
	//				Type:   "map",
	//				X:      req.X,
	//				Y:      req.Y,
	//				Force:  true,
	//				Attack: req.Fire,
	//			})
	//		}
	//	}
	//}
}
