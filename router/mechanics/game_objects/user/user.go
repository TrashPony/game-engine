package user

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type User struct {
	ID              int          `json:"id"`
	Login           string       `json:"login"`
	UserRole        string       `json:"user_role"`
	CurrentPlayerID int          `json:"current_player_id"`
	AllowPlayers    map[int]bool `json:"-"`
	Language        string       `json:"language"`
	externalIDs     map[string]string
	email           string
	connect         *websocket.Conn
	mx              sync.RWMutex
	connectMx       sync.Mutex
}

func (u *User) AddExternalID(key, id string) {
	u.mx.Lock()
	defer u.mx.Unlock()

	if u.externalIDs == nil {
		u.externalIDs = make(map[string]string)
	}

	u.externalIDs[key] = id
}

func (u *User) GetExternalID(key string) (string, bool) {
	u.mx.Lock()
	defer u.mx.Unlock()

	id, ok := u.externalIDs[key]
	return id, ok
}

func (u *User) GetID() int {
	return u.ID
}

func (u *User) SetEmail(email string) {
	u.mx.Lock()
	defer u.mx.Unlock()

	u.email = email
}

func (u *User) GetEmail() string {
	u.mx.RLock()
	defer u.mx.RUnlock()

	return u.email
}

func (u *User) SetCurrentPlayerID(id int) {
	u.mx.Lock()
	defer u.mx.Unlock()

	u.CurrentPlayerID = id
}

func (u *User) GetCurrentPlayerID() int {
	u.mx.RLock()
	defer u.mx.RUnlock()

	return u.CurrentPlayerID
}

func (u *User) SetConnect(connect *websocket.Conn) {
	u.mx.Lock()
	defer u.mx.Unlock()
	u.connect = connect
}

func (u *User) Send(binData []byte, data interface{}) error {
	u.mx.Lock()
	defer u.mx.Unlock()

	conn := u.connect
	if conn == nil {
		return errors.New("no connect")
	}

	var err error
	if binData != nil {
		u.connectMx.Lock()
		err = conn.WriteMessage(2, binData)
		u.connectMx.Unlock()
	} else {
		u.connectMx.Lock()
		err = conn.WriteJSON(data)
		u.connectMx.Unlock()
	}

	return err
}
