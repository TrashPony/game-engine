package web_socket

import (
	"github.com/TrashPony/game-engine/router/generate_ids"
	"github.com/TrashPony/game-engine/router/mechanics/factories/players"
	player2 "github.com/TrashPony/game-engine/router/mechanics/game_objects/player"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/user"
	"github.com/gorilla/websocket"
)

type Request struct {
	Event       string `json:"event,omitempty"`
	Service     string `json:"service,omitempty"`
	ID          int    `json:"id,omitempty"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
	SelectUnits []int  `json:"select_units,omitempty"`
	W           bool   `json:"w"`
	A           bool   `json:"a"`
	S           bool   `json:"s"`
	D           bool   `json:"d"`
	Z           bool   `json:"z"`
	Sp          bool   `json:"sp"`
	St          bool   `json:"st"`
	Fire        bool   `json:"fire"`
}

func Reader(ws *websocket.Conn, gameUser *user.User) {

	defer func() {
		clients.RemoveUser(gameUser)
	}()

	p := players.Users().GetPlayer(gameUser.GetCurrentPlayerID(), gameUser.GetID())
	if p == nil {
		newPlayer := &player2.Player{ID: generate_ids.GetPlayerFakeID(), Login: "Гость"}
		players.Users().AddNewPlayer(newPlayer)
		setCurrentPlayerID(ws, gameUser, newPlayer.ID)
	}

	sendGameTypes(gameUser)
	updatePlayers(gameUser)

	for {
		var msg Request

		err := ws.ReadJSON(&msg)
		if err != nil {
			println(err.Error())
			break
		}

		if msg.Service == "system" {
			SystemService(ws, gameUser, msg)
		}

		if msg.Service == "lobby" {
			LobbyService(ws, gameUser, msg)
		}

		if msg.Service == "battle" {
			BattleService(ws, gameUser, msg)
		}
	}
}
