package battle

import (
	_map "github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/player"
	"sync"
)

type Battle struct {
	ID      int       `json:"id"`
	UUID    string    `json:"uuid"`
	Map     *_map.Map `json:"-"`
	End     bool      `json:"end"`
	players map[int]*player.Player
	mx      sync.RWMutex
}

func (b *Battle) GetUserByID(userID int) *player.Player {
	b.mx.Lock()
	defer b.mx.Unlock()

	return b.players[userID]
}

func (b *Battle) AddUser(newPlayer *player.Player) {
	b.mx.Lock()
	defer b.mx.Unlock()

	if b.players == nil {
		b.players = make(map[int]*player.Player)
	}

	b.players[newPlayer.GetID()] = newPlayer
}

func (b *Battle) GetChanPlayers() <-chan *player.Player {
	b.mx.Lock()

	playerChan := make(chan *player.Player, len(b.players))

	go func() {
		defer func() {
			close(playerChan)
			b.mx.Unlock()
		}()

		for _, p := range b.players {
			playerChan <- p
		}
	}()

	return playerChan
}

func (b *Battle) GetPlayers() map[int]*player.Player {
	players := make(map[int]*player.Player)

	for p := range b.GetChanPlayers() {
		players[p.ID] = p
	}

	return players
}

func (b *Battle) SetPlayers(players map[int]*player.Player) {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.players = players
}
