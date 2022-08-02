package web_socket

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/user"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

var clients = createClientStore()

type clientStore struct {
	clients  map[int]*user.User
	connects map[int]*websocket.Conn
	mx       sync.RWMutex
}

func (s *clientStore) AddUser(ws *websocket.Conn, user *user.User) {

	// защита от двойного входа
	oldWS := s.GetConnect(user.GetID())
	if oldWS != nil {
		oldWS.Close()
	}

	try := 0
	for try < 5 && oldWS != nil {
		try++
		time.Sleep(time.Second)
		oldWS = s.GetConnect(user.GetID())
	}

	s.mx.Lock()
	s.clients[user.GetID()] = user
	s.connects[user.GetID()] = ws
	s.mx.Unlock()

	user.SetConnect(ws)
}

func (s *clientStore) RemoveUser(user *user.User) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if user == nil {
		return
	}

	user.SetConnect(nil)

	ws := s.connects[user.GetID()]
	if ws != nil {
		ws.Close()
	}

	delete(s.clients, user.GetID())
	delete(s.connects, user.GetID())
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

func (s *clientStore) GetUser(userID int) *user.User {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.clients[userID]
}

func (s *clientStore) GetConnect(userID int) *websocket.Conn {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.connects[userID]
}

func createClientStore() *clientStore {
	return &clientStore{clients: make(map[int]*user.User), connects: make(map[int]*websocket.Conn)}
}
