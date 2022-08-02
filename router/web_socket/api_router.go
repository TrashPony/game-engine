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
	UserName    string `json:"user_name,omitempty"`
	ID          int    `json:"id,omitempty"`
	Type        string `json:"type,omitempty"`
	UUID        string `json:"uuid,omitempty"`
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
	TypeSlot    int    `json:"type_slot"`
	Slot        int    `json:"slot"`
	SrcSlot     int    `json:"src_slot"`
	DstSlot     int    `json:"dst_slot"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Resolution  string `json:"resolution"`
	Name        string `json:"name"`
	Left        int    `json:"left"`
	Top         int    `json:"top"`
	Height      int    `json:"height"`
	Width       int    `json:"width"`
	Open        bool   `json:"open"`
	Ready       bool   `json:"ready"`
	Price       int    `json:"price"`
	Count       int    `json:"count"`
	Data        string `json:"data"`
	Date        string `json:"date"`
	ToDate      string `json:"to_date"`
	ToPage      int    `json:"to_page"`
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
