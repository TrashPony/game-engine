package web_socket

import (
	"github.com/TrashPony/game_engine/src/mechanics/create_battle"
	"github.com/TrashPony/game_engine/src/mechanics/factories/lobby_sessions"
	"github.com/TrashPony/game_engine/src/mechanics/factories/maps"
	"github.com/TrashPony/game_engine/src/mechanics/factories/players"
	"github.com/TrashPony/game_engine/src/mechanics/factories/quick_battles"
	_map "github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/user"
	"github.com/gorilla/websocket"
)

// todo вынести в отдельный пакет
func LobbyService(ws *websocket.Conn, gameUser *user.User, req Request) {

	p := players.Users.GetPlayer(gameUser.GetCurrentPlayerID(), gameUser.GetID())
	if p == nil {
		SendMessage(Response{Event: req.Event, Error: "no player", UserID: gameUser.ID})
		return
	}

	if req.Event == "CreateLobbySession" {
		newLobbySession := lobby_sessions.Store.CreateSession(p)
		updateLobbySession(newLobbySession)
	}

	if req.Event == "JoinSession" {
		// TODO
	}

	if req.Event == "RemoveSession" {
		// TODO
	}

	if req.Event == "SelectLobbyMap" {
		lobbySession := lobby_sessions.Store.GetSessionByUUID(req.UUID)
		if lobbySession == nil {
			SendMessage(Response{Event: "Error", Error: "no session", UserID: gameUser.ID})
			return
		}

		if lobbySession.LeaderID != p.GetID() {
			SendMessage(Response{Event: "Error", Error: "no leader", UserID: gameUser.ID})
			return
		}

		mp, ok := maps.Maps.GetByID(req.ID)
		if !ok {
			SendMessage(Response{Event: "Error", Error: "wrong map", UserID: gameUser.ID})
			return
		}

		lobbySession.MapID = mp.Id
		updateLobbySession(lobbySession)
	}

	if req.Event == "StartGame" {
		lobbySession := lobby_sessions.Store.GetSessionByUUID(req.UUID)
		if lobbySession == nil {
			SendMessage(Response{Event: "Error", Error: "no session", UserID: gameUser.ID})
			return
		}

		if lobbySession.LeaderID != p.GetID() {
			SendMessage(Response{Event: "Error", Error: "no leader", UserID: gameUser.ID})
			return
		}

		mp, ok := maps.Maps.GetByID(lobbySession.MapID)
		if !ok {
			SendMessage(Response{Event: "Error", Error: "wrong map", UserID: gameUser.ID})
			return
		}

		newBattle := create_battle.CreateBattle(lobbySession.Players, mp.Id)
		quick_battles.Battles.AddNewGame(newBattle)

		SendMessage(Response{
			Event:            "ToBattle",
			LobbySessionUUID: lobbySession.UUID,
		})

		// TODO удаляем лобби сессию
	}
}

func updateLobbySession(lobbySession *lobby_sessions.LobbySession) {
	SendMessage(Response{
		Event:            "CreateLobbySession",
		LobbySessionUUID: lobbySession.UUID,
		Data: struct {
			LobbySession *lobby_sessions.LobbySession `json:"lobby_session"`
			Maps         map[int]*_map.ShortInfoMap   `json:"maps"`
		}{
			LobbySession: lobbySession,
			Maps:         maps.Maps.GetAllShortInfoMap(),
		}})
}
