package battle

import (
	_const "github.com/TrashPony/game-engine/router/const"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/player"
	uuid "github.com/satori/go.uuid"
	"sync"
)

type Battle struct {
	ID                  int                    `json:"id"`
	UUID                string                 `json:"uuid"`
	Map                 *_map.Map              `json:"-"`
	Teams               map[int]*Team          `json:"teams"`
	WaitReady           bool                   `json:"wait_ready"`
	WaitTimeOut         int                    `json:"wait_time_out"`
	SessionPlayersState map[int]*SessionPlayer `json:"session_players_state"`
	players             []*player.Player
	mx                  sync.RWMutex
}

func (b *Battle) GetPlayerTeams() map[int]*Team {
	b.mx.RLock()
	defer b.mx.RUnlock()

	teams := make(map[int]*Team)

	for key, t := range b.Teams {
		if !t.AI && !t.Hide {
			teams[key] = t
		}
	}

	return teams
}

func (b *Battle) CalcReady(players []*player.Player) {
	allReady := true

	for _, p := range players {
		allReady = allReady && (p.GetReady() || p.Bot)
	}

	if allReady && b.WaitTimeOut > 5*1000 {
		b.WaitTimeOut = 5 * 1000
	}

	if b.WaitTimeOut <= 0 {
		b.WaitReady = false
	}

	b.WaitTimeOut -= _const.ServerTick
}

func (b *Battle) GetPlayerByID(userID int) *player.Player {
	b.mx.Lock()
	defer b.mx.Unlock()

	for _, p := range b.players {
		if p.GetID() == userID {
			return p
		}
	}

	return nil
}

func (b *Battle) AddUser(newPlayer *player.Player) {
	b.mx.Lock()
	defer b.mx.Unlock()

	if b.players == nil {
		b.players = make([]*player.Player, 0)
	}

	_, ok := b.SessionPlayersState[newPlayer.ID]
	if !ok {
		b.SessionPlayersState[newPlayer.ID] = &SessionPlayer{
			UUID:     uuid.NewV4().String(),
			PlayerID: newPlayer.GetID(),
			Login:    newPlayer.GetLogin(),
			Live:     true,
		}
	}

	b.players = append(b.players, newPlayer)
}

func (b *Battle) GetChanPlayers() <-chan *player.Player {
	b.mx.RLock()

	playerChan := make(chan *player.Player, len(b.players))

	go func() {
		defer func() {
			close(playerChan)
			b.mx.RUnlock()
		}()

		for _, p := range b.players {
			playerChan <- p
		}
	}()

	return playerChan
}

func (b *Battle) UnsafeRangePlayers() ([]*player.Player, *sync.RWMutex) {
	b.mx.RLock()
	return b.players, &b.mx
}

func (b *Battle) GetPlayers(basket []*player.Player) []*player.Player {
	b.mx.Lock()
	defer b.mx.Unlock()

	basket = basket[:0]
	for _, p := range b.players {
		if p.GameUUID == b.UUID {
			basket = append(basket, p)
		}
	}

	return basket
}

func (b *Battle) RemovePlayer(removePlayer *player.Player) {
	b.mx.Lock()
	defer b.mx.Unlock()

	index := -1
	for i, p := range b.players {
		if p.GetID() == removePlayer.GetID() {
			index = i
			break
		}
	}

	if index >= 0 {
		b.players[index] = b.players[len(b.players)-1]
		b.players = b.players[:len(b.players)-1]
	}
}
