package user

import "sync"

type User struct {
	ID              int          `json:"id"`
	Login           string       `json:"login"`
	email           string       `json:"-"`
	UserRole        string       `json:"user_role"`
	CurrentPlayerID int          `json:"current_player_id"`
	AllowPlayers    map[int]bool `json:"-"`
	Language        string       `json:"language"`
	mx              sync.RWMutex
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
