package players

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/player"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/user"
	uuid "github.com/satori/go.uuid"
	"sync"
)

type usersStore struct {
	mx      sync.RWMutex
	users   map[int]*user.User
	players map[int]*player.Player
}

var users *usersStore

func Users() *usersStore {
	if users == nil {
		users = newUsersStore()
	}

	return users
}

func newUsersStore() *usersStore {

	store := &usersStore{
		users:   make(map[int]*user.User),
		players: make(map[int]*player.Player),
	}

	return store
}

func (usersStore *usersStore) getUser(id int) *user.User {
	usersStore.mx.RLock()
	defer usersStore.mx.RUnlock()

	val, ok := usersStore.users[id]
	if !ok {
		val = &user.User{
			ID:    id,
			Login: uuid.NewV4().String(),
		}

		usersStore.users[id] = val
	}

	return val
}

func (usersStore *usersStore) GetUser(id int) *user.User {
	if id <= 0 {
		return nil
	}

	return usersStore.getUser(id)
}

func (usersStore *usersStore) getPlayer(id int) *player.Player {
	usersStore.mx.RLock()
	defer usersStore.mx.RUnlock()
	return usersStore.players[id]
}

func (usersStore *usersStore) AddNewPlayer(newPlayer *player.Player) {
	usersStore.mx.Lock()
	defer usersStore.mx.Unlock()
	usersStore.players[newPlayer.ID] = newPlayer
}

func (usersStore *usersStore) GetPlayer(id, userID int) *player.Player {
	return usersStore.getPlayer(id)
}
