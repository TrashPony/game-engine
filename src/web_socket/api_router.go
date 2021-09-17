package web_socket

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/user"
	"github.com/gorilla/websocket"
)

type Request struct {
	Event       string `json:"event,omitempty"`
	Service     string `json:"service,omitempty"`
	UserName    string `json:"user_name,omitempty"`
	ID          int    `json:"id,omitempty"`
	UUID        string `json:"uuid,omitempty"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
	SelectUnits []int  `json:"select_units,omitempty"`
	W           bool   `json:"w"`
	A           bool   `json:"a"`
	S           bool   `json:"s"`
	D           bool   `json:"d"`
	Fire        bool   `json:"fire"`
}

func Reader(ws *websocket.Conn, gameUser *user.User) {

	sendUI(gameUser)
	sendGameTypes(gameUser)

	for {
		var msg Request

		err := ws.ReadJSON(&msg)
		if err != nil {
			println(err.Error())
			clients.RemoveUser(gameUser.GetID())
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
