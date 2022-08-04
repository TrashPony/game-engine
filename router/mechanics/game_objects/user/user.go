package user

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type User struct {
	ID              int    `json:"id"`
	Login           string `json:"login"`
	CurrentPlayerID int    `json:"current_player_id"`
	connect         *websocket.Conn
	mx              sync.RWMutex
	connectMx       sync.Mutex
}

func (u *User) GetID() int {
	return u.ID
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
