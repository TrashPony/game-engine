package web_socket

import (
	"github.com/TrashPony/game_engine/src/mechanics/factories/players"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/user"
	"github.com/gorilla/websocket"
	"sync"
)

var clients = createClientStore()

type clientStore struct {
	clients  map[int]*user.User
	connects map[int]*websocket.Conn
	mx       sync.RWMutex
}

func (s *clientStore) AddUser(ws *websocket.Conn, user *user.User) {

	s.RemoveUser(user.GetID())

	s.mx.Lock()
	defer s.mx.Unlock()

	s.clients[user.GetID()] = user
	s.connects[user.GetID()] = ws
}

func (s *clientStore) RemoveUser(userID int) {
	s.mx.Lock()
	defer s.mx.Unlock()

	for _, p := range players.Users.GetPlayersByUserID(userID) {
		p.SetReady(false)
	}

	delete(s.clients, userID)
	delete(s.connects, userID)
}

func (s *clientStore) GetUsersChan() <-chan *user.User {
	s.mx.RLock()

	users := make(chan *user.User, len(s.clients))
	go func() {
		defer func() {
			s.mx.RUnlock()
			close(users)
		}()

		for _, user := range s.clients {
			users <- user
		}
	}()

	return users
}

func (s *clientStore) GetConnect(userID int) *websocket.Conn {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.connects[userID]
}

func createClientStore() *clientStore {
	return &clientStore{clients: make(map[int]*user.User), connects: make(map[int]*websocket.Conn)}
}
