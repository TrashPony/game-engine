package web_socket

import (
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/const/game_types"
	"github.com/TrashPony/game-engine/router/mechanics/factories/players"
	player2 "github.com/TrashPony/game-engine/router/mechanics/game_objects/player"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/user"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/web_socket_response"
	"github.com/gorilla/websocket"
	"time"
)

func SystemService(ws *websocket.Conn, gameUser *user.User, req Request) {
	if req.Event == "GetPlayers" {
		updatePlayers(gameUser)
	}
}

func updatePlayers(gameUser *user.User) {
	SendMessage(web_socket_response.Response{
		Event:  "GetPlayers",
		UserID: gameUser.GetID(),
		Data: struct {
			Player     *player2.Player `json:"player"`
			GameUser   *user.User      `json:"game_user"`
			ServerTime int64           `json:"st"`
		}{
			Player:     players.Users().GetPlayer(gameUser.GetCurrentPlayerID(), gameUser.GetID()),
			GameUser:   gameUser,
			ServerTime: time.Now().UTC().Unix(),
		}})
}

func setCurrentPlayerID(ws *websocket.Conn, gameUser *user.User, playerID int) {
	gameUser.SetCurrentPlayerID(playerID)
	player := players.Users().GetPlayer(gameUser.GetCurrentPlayerID(), gameUser.GetID())
	if player == nil {
		ws.Close()
		return
	}

	player.SetOwner(gameUser)
	updatePlayers(gameUser)
}

func sendGameTypes(gameUser *user.User) {
	SendMessage(web_socket_response.Response{
		Event:   "setGameTypes",
		Service: "init",
		UserID:  gameUser.GetID(),
		Data: struct {
			Bodies         interface{} `json:"bodies"`
			Weapons        interface{} `json:"weapons"`
			Ammo           interface{} `json:"ammo"`
			MapBinItems    interface{} `json:"map_bin_items"`
			UnitSpriteSize int         `json:"unit_sprite_size"`
		}{
			Bodies:         game_types.GetAllBody(),
			Weapons:        game_types.GetAllWeapons(),
			Ammo:           game_types.GetAllAmmo(),
			UnitSpriteSize: _const.SpriteSize,
			MapBinItems:    _const.MapBinItems,
		},
	})
}
