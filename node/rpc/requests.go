package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/TrashPony/game-engine/node/ai"
	binary_msg2 "github.com/TrashPony/game-engine/node/binary_msg"
	"github.com/TrashPony/game-engine/node/create_battle"
	"github.com/TrashPony/game-engine/node/mechanics/factories/quick_battles"
	"github.com/TrashPony/game-engine/node/mechanics/factories/units"
	"github.com/TrashPony/game-engine/router/const/game_types"
	"github.com/TrashPony/game-engine/router/generate_ids"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/behavior_rule"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/body"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/dynamic_map_object"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/obstacle_point"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/player"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/rpc_request"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/target"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/web_socket_response"
	"math/rand"
)

type StartGameObject struct {
	StartPlayers []*player.Player `json:"start_players"`
	MapID        int              `json:"map_id"`
}

func createBattle(req *rpc_request.Request) string {
	data := StartGameObject{}

	err := json.Unmarshal([]byte(req.Data), &data)
	if err != nil {
		fmt.Println("failed create game")
		return ""
	}

	b := create_battle.CreateBattle(data.StartPlayers, data.MapID)
	if b != nil {
		quick_battles.Battles.AddNewGame(b)
		return b.UUID
	}

	return ""
}

func i(p *player.Player, req *rpc_request.Request) {
	for gameUnit := range p.GetGameUnitsStore().RangeUnits() {
		if gameUnit != nil {

			// одновременно может быть нажато только 2 клавишы, если зажаты 4 то эфект нигилируется
			if req.W && req.S {
				req.W = false
				req.S = false
			}

			if req.A && req.D {
				req.A = false
				req.D = false
			}

			gameUnit.GetPhysicalModel().SetWASD(req.W, req.A, req.S, req.D, req.Sp, req.St, req.Z)
			gameUnit.SetWeaponTarget(&target.Target{
				Type:   "map",
				X:      req.X,
				Y:      req.Y,
				Force:  true,
				Attack: req.Fire,
			})
		}
	}
}

type MapPosition struct {
	X   int       `json:"x"`
	Y   int       `json:"Y"`
	Map *_map.Map `json:"map"`
}

func initBattle(b *battle.Battle, p *player.Player) web_socket_response.Response {
	p.SetReady(false)
	p.MapID = b.Map.Id

	playerNames := make(map[int]string)
	for p := range b.GetChanPlayers() {
		playerNames[p.GetID()] = p.GetLogin()
	}

	command := []byte{100}
	for u := range p.GetGameUnitsStore().RangeUnits() {
		if u.MapID == b.Map.Id {
			stateMsg := binary_msg2.StatusSquadBinaryMsg(u.HP, 0)
			command = append(command, binary_msg2.GetIntBytes(len(stateMsg))...)
			command = append(command, stateMsg...)
		}
	}

	data, _ := json.Marshal(struct {
		Maps        map[int]*MapPosition `json:"maps"`
		Player      *player.Player       `json:"player"`
		Spectrum    bool                 `json:"spectrum"`
		Battle      *battle.Battle       `json:"battle"`
		PlayerNames map[int]string       `json:"player_names"`
	}{
		// todo небольшой костыль из за клиента который взят из другой игоры
		Maps:        map[int]*MapPosition{b.Map.Id: {X: 0, Y: 0, Map: b.Map}},
		Player:      p,
		Spectrum:    false,
		Battle:      b,
		PlayerNames: playerNames,
	})

	return web_socket_response.Response{
		Event:     "InitBattle",
		PlayerID:  p.GetID(),
		Responses: []*web_socket_response.Response{{BinaryMsg: command}},
		Data:      string(data),
	}
}

func startLoad(battle *battle.Battle, p *player.Player) web_socket_response.Response {
	dynamicObjects := make(map[int][]byte)

	team := battle.Teams[p.TeamID]
	if team == nil {
		return web_socket_response.Response{}
	}

	for memoryObj := range team.GetMapDynamicObjects(battle.Map.Id) {
		dynamicObjects[memoryObj.IDObject] = memoryObj.ObjectJSON
	}

	p.SetReady(true)

	updater := make([]*web_socket_response.Response, 0)
	updater = append(updater, &web_socket_response.Response{
		BinaryMsg: binary_msg2.RadareON(),
	})

	objects, mx := team.UnsafeRangeVisibleObjects()
	mx.RLock()
	defer mx.RUnlock()

	for _, vObj := range objects {
		if vObj.View {
			updater = append(updater, &web_socket_response.Response{
				BinaryMsg: binary_msg2.CreateBinaryCreateObjMsg(vObj.TypeObject, vObj.ObjectJSON),
			})
			continue
		}

		if vObj.Radar {
			updater = append(updater, &web_socket_response.Response{
				BinaryMsg: binary_msg2.CreateBinaryCreateRadarMarkMsg(vObj.ID, vObj.X, vObj.Y, vObj.Type),
			})
			continue
		}
	}

	return web_socket_response.Response{
		Event:     "RefreshDynamicObj",
		PlayerID:  p.GetID(),
		Data:      dynamicObjects,
		Responses: updater,
	}
}

func createUnit(battle *battle.Battle, p *player.Player, x, y int) {
	newUnit := &unit.Unit{ID: generate_ids.GetPlayerFakeID(), OwnerID: p.GetID(), MapID: battle.Map.Id, HP: 100}
	newUnit.Body = game_types.GetNewBody(1)

	wSlot := newUnit.GetWeaponSlot(1)
	wSlot.Weapon = game_types.GetNewWeapon(2)
	wSlot.Ammo = game_types.GetNewAmmo(wSlot.Weapon.DefaultAmmoTypeID)
	wSlot.SetAnchor()

	newUnit.GetPhysicalModel().SetPos(float64(x), float64(y), float64(rand.Intn(360)))
	p.GetGameUnitsStore().AddUnit(newUnit, p.GetTeamID(), p.GetID())
	units.Units.AddUnit(newUnit)
}

func createBot(battle *battle.Battle, p *player.Player, x, y, teamID int) {
	bodyID := 1
	if teamID == 2 {
		bodyID = 2
	}

	newBot := ai.CreateBot(battle.UUID, teamID, battle.Map, x, y, rand.Intn(360), &behavior_rule.BehaviorRules{
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
	}, true, true, bodyID)

	for u := range newBot.GetGameUnitsStore().RangeUnits() {
		units.Units.AddUnit(u)
	}

	ai.AddLifeBot(newBot)
}

func createObj(battle *battle.Battle, p *player.Player, x, y, teamID, objectID int) {
	fraction := "replics"
	if teamID == 2 {
		fraction = "reverses"
	}

	obj := &dynamic_map_object.Object{
		TypeID:    1,
		Type:      "turret",
		Texture:   "laser_turret_" + fraction,
		MaxHP:     2000,
		TypeMaxHP: 2000,
		HP:        2000,
		Scale:     25,
		Shadow:    false,
		TeamID:    teamID,
		TypeGeoData: []*obstacle_point.ObstaclePoint{
			{X: 0, Y: 0, Radius: 39, Move: false},
			{X: -32, Y: -31, Radius: 14, Move: false},
			{X: 30, Y: -29, Radius: 14, Move: false},
			{X: -30, Y: 28, Radius: 14, Move: false},
			{X: 26, Y: 26, Radius: 14, Move: false},
			{X: -41, Y: -39, Radius: 14, Move: false},
			{X: 41, Y: -40, Radius: 14, Move: false},
			{X: -41, Y: 39, Radius: 14, Move: false},
			{X: 40, Y: 38, Radius: 14, Move: false},
		},
		HeightType: 300,
		Fraction:   fraction,
		Static:     true,
		Build:      true,
		WeaponID:   100,
		XAttach:    64,
		YAttach:    64,
		ViewRange:  250,
	}

	if obj.WeaponID > 0 {
		weapon := game_types.GetNewWeapon(obj.WeaponID)
		obj.Weapons = make(map[int]*body.WeaponSlot)
		obj.Weapons[1] = &body.WeaponSlot{Number: 1, Weapon: weapon, XAttach: obj.XAttach, YAttach: obj.YAttach}
		obj.Weapons[1].Weapon.Name += "_" + fraction
		obj.Weapons[1].SetAnchor()
		obj.GetTurretAmmo()
	}

	if obj.Build {
		obj.SetPower(obj.MaxEnergy)
		obj.Run = true
		obj.Work = true
	}

	obj.GetPhysicalModel().SetPos(float64(x), float64(y), 0)
	battle.Map.AddDynamicObject(obj)
	obj.SetGeoData()
}
