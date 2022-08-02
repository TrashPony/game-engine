package web_socket

import (
	"github.com/TrashPony/game-engine/router/mechanics/factories/nodes"
	"github.com/TrashPony/game-engine/router/mechanics/factories/players"
	player2 "github.com/TrashPony/game-engine/router/mechanics/game_objects/player"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/user"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/web_socket_response"
	"github.com/gorilla/websocket"
)

func LobbyService(ws *websocket.Conn, gameUser *user.User, req Request) {

	p := players.Users().GetPlayer(gameUser.GetCurrentPlayerID(), gameUser.GetID())
	if p == nil {
		SendMessage(web_socket_response.Response{Event: req.Event, Error: "no player", UserID: gameUser.ID})
		return
	}

	if req.Event == "gsgs" {
		SendMessage(web_socket_response.Response{
			Event:    "gsgs",
			PlayerID: p.GetID(),
			Data: struct {
				InGame bool `json:"ig"`
			}{
				InGame: p.GameUUID != "",
			},
		})
		return
	}

	if p.GameUUID != "" {
		SendMessage(web_socket_response.Response{Event: req.Event, Error: "player in battle", UserID: gameUser.ID})
		return
	}

	if req.Event == "StartGame" {

		node := nodes.Nodes().GetNode()
		if node == nil {
			return
		}

		newBattleUUID := node.CreateBattle([]*player2.Player{p})

		p.GameUUID = newBattleUUID
		p.NodeName = node.Name
	}
}
