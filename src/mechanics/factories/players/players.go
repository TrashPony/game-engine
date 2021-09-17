package players

import (
	dbPlayer "github.com/TrashPony/game_engine/src/mechanics/db/player"
	dbUser "github.com/TrashPony/game_engine/src/mechanics/db/user"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/player"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/user"
	"sync"
)

type usersStore struct {
	mx      sync.RWMutex
	users   map[int]*user.User
	players map[int]*player.Player
}

var Users = newUsersStore()

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

	val := usersStore.users[id]
	return val
}

func (usersStore *usersStore) GetUser(id int) *user.User {
	if id <= 0 {
		return nil
	}

	u := usersStore.getUser(id)
	if u != nil {
		return u
	}

	usersStore.mx.Lock()
	defer usersStore.mx.Unlock()

	u = dbUser.GetUser(id)
	u.AllowPlayers = dbPlayer.GetPlayersIDsByUserID(u.ID)

	usersStore.users[id] = u

	return u
}

func (usersStore *usersStore) getPlayer(id int) *player.Player {
	usersStore.mx.RLock()
	defer usersStore.mx.RUnlock()

	val := usersStore.players[id]
	return val
}

func (usersStore *usersStore) AddNewPlayer(newPlayer *player.Player, userID int) {

	defer usersStore.GetPlayersByUserID(userID)

	dbPlayer.AddNewPlayer(newPlayer, userID)
}

func (usersStore *usersStore) GetPlayersByUserID(userID int) []*player.Player {

	players := make([]*player.Player, 0)

	ids := dbPlayer.GetPlayersIDsByUserID(userID)

	u := usersStore.getUser(userID)
	if u != nil {
		u.AllowPlayers = ids
	}

	for id := range ids {
		p := usersStore.GetPlayer(id, userID)
		if p != nil && p.ID > 0 {
			players = append(players, p)
		}
	}

	return players
}

func (usersStore *usersStore) GetPlayer(id, userID int) *player.Player {

	if id <= 0 {
		return nil
	}

	if userID != 0 {
		u := usersStore.getUser(userID)
		if u == nil || !u.AllowPlayers[id] {
			return nil
		}
	}

	p := usersStore.getPlayer(id)
	if p != nil {
		return p
	}

	usersStore.mx.Lock()
	defer usersStore.mx.Unlock()

	p = dbPlayer.Player(id)
	if p.ID > 0 {
		usersStore.players[id] = p
	}

	return p
}
