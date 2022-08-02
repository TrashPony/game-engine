package web_socket

import (
	"github.com/TrashPony/game-engine/router/mechanics/factories/nodes"
	"github.com/TrashPony/game-engine/router/mechanics/factories/players"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/player"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/user"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/web_socket_response"
	"github.com/gorilla/websocket"
)

func BattleService(ws *websocket.Conn, gameUser *user.User, req Request) {
	p := players.Users().GetPlayer(gameUser.GetCurrentPlayerID(), gameUser.GetID())
	if p == nil {
		SendMessage(web_socket_response.Response{Event: req.Event, Error: "no player", UserID: gameUser.ID})
		return
	}

	node := nodes.Nodes().GetNodeByName(p.NodeName)
	if node == nil {
		toLobby(p, gameUser)
		return
	}

	if req.Event == "InitBattle" {
		resp := node.InitBattle(p.GameUUID, p.GetID()).Response
		if resp.Event == "error" {
			toLobby(p, gameUser)
		} else {
			SendMessage(resp)
		}
	}

	if req.Event == "StartLoad" {
		SendMessage(node.StartLoad(p.GameUUID, p.GetID()).Response)
	}

	if req.Event == "ExitGame" {
		toLobby(p, gameUser)
	}

	if req.Event == "i" {
		node.Input(p.GameUUID, p.GetID(), req.W, req.A, req.S, req.D, req.Sp, req.St, req.Z, req.X, req.Y, req.Fire)
	}

	if req.Event == "CreateUnit" {
		node.CreateUnit(p.GameUUID, p.GetID(), 50, 50)
	}

	if req.Event == "CreateBot" {
		x, y := 50, 50

		if req.ID == 2 {
			x, y = 975, 975
		}

		node.CreateBot(p.GameUUID, p.GetID(), x, y, req.ID)
	}

	if req.Event == "CreateObj" {
		x, y := game_math.GetRangeRand(50, 250), game_math.GetRangeRand(50, 250)

		if req.ID == 2 {
			x, y = game_math.GetRangeRand(750, 975), game_math.GetRangeRand(750, 975)
		}

		node.CreateObj(p.GameUUID, p.GetID(), 1, x, y, req.ID)
	}
}

func toLobby(p *player.Player, gameUser *user.User) {

	p.TeamID = 0
	p.GameUUID = ""
	p.NodeName = ""

	SendMessage(web_socket_response.Response{
		Event:  "ToLobby",
		UserID: gameUser.GetID(),
	})
}
