package web_socket

import (
	"github.com/TrashPony/game-engine/router/mechanics/factories/players"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/user"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/web_socket_response"
	"sync"
)

var senderPipe = make(chan web_socket_response.Response)

type MessagesStore struct {
	messagesGroup []*MessagesGroup
	mx            sync.RWMutex
}

type MessagesGroup struct {
	Key        string            `json:"key"`
	Messages   *GameLoopMessages `json:"messages"`
	Type       string            `json:"type"` // move - проверяем функцией "moveMessagesToUser", view - "ResponseToBin"
	Attributes map[string]string `json:"attributes"`
}

type GameLoopMessages struct {
	Messages []web_socket_response.Response
	mx       sync.RWMutex
}

func (ms *MessagesStore) AddMsg(typeMsg, typeCheck string, msg web_socket_response.Response, attributes map[string]string) {
	ms.mx.Lock()
	defer ms.mx.Unlock()

	ms.getMessageGroup(typeMsg, typeCheck, attributes).Messages.AddMessage(msg)
}

func (ms *MessagesStore) GetMessageGroups() []*MessagesGroup {
	// todo скорее всего этот метод безопасен т.к. он испольняется когда уже все добавления происзошли, но это не точно)
	return ms.messagesGroup
}

func (ms *MessagesStore) Clear() {
	ms.messagesGroup = nil
}

func (ms *MessagesStore) getMessageGroup(typeMsg, typeCheck string, attributes map[string]string) *MessagesGroup {
	for _, m := range ms.messagesGroup {
		if m.Key == typeMsg {
			return m
		}
	}

	newGroup := &MessagesGroup{
		Key:        typeMsg,
		Messages:   &GameLoopMessages{Messages: make([]web_socket_response.Response, 0)},
		Type:       typeCheck,
		Attributes: attributes,
	}

	ms.messagesGroup = append(ms.messagesGroup, newGroup)
	return newGroup
}

func (m *GameLoopMessages) AddMessage(msg web_socket_response.Response) {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.Messages = append(m.Messages, msg)
}

func SendMessage(message web_socket_response.Response) {
	go func() {
		senderPipe <- message
	}()
}

func Sender() {
	for {
		resp := <-senderPipe

		go func() {
			//  если сразу указаному кому отправлять сообщение, не искать его через общий пул, а сразу кидать ему
			if resp.UserID != 0 {
				u := clients.GetUser(resp.UserID)
				if u != nil {
					sendWrapper(u, &resp)
				}

				return
			}

			if resp.PlayerID != 0 {
				p := players.Users().GetPlayer(resp.PlayerID, 0)
				if p != nil && (p.GameUUID == resp.GameUUID || resp.GameUUID == "") {
					u := clients.GetUser(p.OwnerID)
					if u != nil {
						sendWrapper(u, &resp)
					}
				}

				return
			}

			for client := range clients.GetUsersChan() {

				if resp.All {
					sendWrapper(client, &resp)
					continue
				}

				/** всем в лобби **/
				if resp.AllNoBattle {
					player := players.Users().GetPlayer(client.GetCurrentPlayerID(), client.GetID())
					if player != nil && player.GameUUID == "" {
						sendWrapper(client, &resp)
						continue
					}
				}

				/** всем в конкретной игре **/
				if resp.GameUUID != "" {
					player := players.Users().GetPlayer(client.GetCurrentPlayerID(), client.GetID())
					if player != nil && player.GameUUID == resp.GameUUID {
						sendWrapper(client, &resp)
						continue
					}
				}
			}
		}()
	}
}

func sendWrapper(client *user.User, resp *web_socket_response.Response) {

	send(client, resp)

	if len(resp.Responses) > 0 {
		for _, r := range resp.Responses {
			send(client, r)
		}
	}
}

func send(client *user.User, resp *web_socket_response.Response) {
	var err error

	if resp.BinaryMsg != nil {
		err = client.Send(resp.BinaryMsg, nil)
	} else {
		if resp.OnlyData {
			err = client.Send(nil, resp.Data)
		} else {
			err = client.Send(nil, resp)
		}
	}

	if err != nil {
		clients.RemoveUser(client)
	}
}
