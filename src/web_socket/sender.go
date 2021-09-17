package web_socket

import (
	"github.com/TrashPony/game_engine/src/mechanics/factories/players"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/user"
	"sync"
)

var senderPipe = make(chan Response)

type GameLoopMessages struct {
	Messages []interface{}
	mx       sync.RWMutex
}

func (m *GameLoopMessages) AddMessage(msg interface{}) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if m.Messages == nil {
		m.Messages = make([]interface{}, 0)
	}

	m.Messages = append(m.Messages, msg)
}

type Response struct {
	Event            string      `json:"event,omitempty"`
	Service          string      `json:"service,omitempty"`
	UserID           int         `json:"-"`
	PlayerID         int         `json:"-"`
	MapID            int         `json:"-"`
	LobbySessionUUID string      `json:"-"`
	Data             interface{} `json:"data,omitempty"`
	Error            string      `json:"error,omitempty"`
}

func SendMessage(message Response) {
	go func() {
		senderPipe <- message
	}()
}

func Sender() {
	for {
		resp := <-senderPipe

		for client := range clients.GetUsersChan() {

			if (resp.UserID > 0 && resp.PlayerID > 0) || (resp.UserID > 0 && resp.MapID > 0) || (resp.PlayerID > 0 && resp.MapID > 0) {
				panic("bat priority")
			}

			/** UserID имеет первый приоритет, отправляем только юзеру **/
			if resp.UserID != 0 && client.GetID() == resp.UserID {
				send(client, &resp)
				continue
			}

			/** UserID имеет второй приоритет, отправляем только игроку **/
			if resp.PlayerID != 0 && client.GetCurrentPlayerID() == resp.PlayerID {
				send(client, &resp)
				continue
			}

			/** MapID, отправляется всем кто находится на этой карте **/
			if resp.MapID != 0 {
				player := players.Users.GetPlayer(client.GetCurrentPlayerID(), client.GetID())
				if player.MapID == resp.MapID {
					send(client, &resp)
					continue
				}
			}

			/** всем кто в лобби **/
			if resp.LobbySessionUUID != "" {
				player := players.Users.GetPlayer(client.GetCurrentPlayerID(), client.GetID())
				if player.LobbyUUID == resp.LobbySessionUUID {
					send(client, &resp)
					continue
				}
			}
		}
	}
}

func send(client *user.User, resp *Response) {
	connect := clients.GetConnect(client.GetID())
	if connect == nil {
		return
	}

	err := connect.WriteJSON(resp)
	if err != nil {
		clients.RemoveUser(client.GetID())
	}
}
