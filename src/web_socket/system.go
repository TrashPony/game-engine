package web_socket

import (
	_const "github.com/TrashPony/game_engine/src/const"
	"github.com/TrashPony/game_engine/src/const/game_types"
	"github.com/TrashPony/game_engine/src/mechanics/factories/players"
	player2 "github.com/TrashPony/game_engine/src/mechanics/game_objects/player"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/user"
	"github.com/gorilla/websocket"
)

// todo вынести в отдельный пакет
func SystemService(ws *websocket.Conn, gameUser *user.User, req Request) {

	if req.Event == "GetPlayers" {
		updatePlayers(gameUser)
	}

	if req.Event == "CreateNewPlayer" {
		if len(players.Users.GetPlayersByUserID(gameUser.ID)) > 0 {
			SendMessage(Response{Event: req.Event, Error: "Already", UserID: gameUser.ID})
		} else {
			newPlayer := &player2.Player{Login: req.UserName}
			gameUser.SetCurrentPlayerID(newPlayer.ID)
			players.Users.AddNewPlayer(newPlayer, gameUser.ID)

			updatePlayers(gameUser)
		}
	}

	if req.Event == "SelectPlayer" {
		gameUser.SetCurrentPlayerID(req.ID)
		player := players.Users.GetPlayer(gameUser.GetCurrentPlayerID(), gameUser.GetID())
		if player == nil {
			ws.Close()
			return
		}

		updatePlayers(gameUser)
	}
}

func updatePlayers(gameUser *user.User) {
	SendMessage(Response{
		Event:  "GetPlayers",
		UserID: gameUser.GetID(),
		Data: struct {
			Player   *player2.Player   `json:"player"`
			Players  []*player2.Player `json:"players"`
			GameUser *user.User        `json:"game_user"`
		}{
			Player:   players.Users.GetPlayer(gameUser.GetCurrentPlayerID(), gameUser.GetID()),
			Players:  players.Users.GetPlayersByUserID(gameUser.GetID()),
			GameUser: gameUser,
		}})

	sendUI(gameUser)
}

func sendUI(gameUser *user.User) {

	p := players.Users.GetPlayer(gameUser.GetCurrentPlayerID(), gameUser.GetID())

	SendMessage(Response{
		Event:   "setWindowsState",
		Service: "init",
		UserID:  gameUser.GetID(),
		Data: struct {
			UserInterface string          `json:"user_interface"`
			AllowWindows  map[string]bool `json:"allow_windows"`
			// TODO DescriptionItems map[string]map[string]map[string]handbook.DescriptionItem `json:"description_items"`
		}{
			UserInterface: p.GetJSONWindowState(),
			AllowWindows:  player2.AllowWindowSave,
			// TODO DescriptionItems: handbook.AllDescription,
		},
	})
}

func sendGameTypes(gameUser *user.User) {
	SendMessage(Response{
		Event:   "setGameTypes",
		Service: "init",
		UserID:  gameUser.GetID(),
		Data: struct {
			Bodies         interface{} `json:"bodies"`
			Weapons        interface{} `json:"weapons"`
			Ammo           interface{} `json:"ammo"`
			UnitSpriteSize int         `json:"unit_sprite_size"`
		}{
			Bodies:         game_types.BodyTypes,
			Weapons:        game_types.WeaponTypes,
			Ammo:           game_types.AmmoTypes,
			UnitSpriteSize: _const.SpriteSize,
		},
	})
}
