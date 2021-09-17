package lobby_sessions

import (
	"errors"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/player"
	uuid "github.com/satori/go.uuid"
	"sync"
)

var Store = newStore()

type store struct {
	sessions map[string]*LobbySession // [uuid]
	mx       sync.RWMutex
}

type LobbySession struct {
	UUID     string           `json:"uuid"`
	LeaderID int              `json:"leader_id"`
	Players  []*player.Player `json:"players"`
	MapID    int              `json:"map_id"`
}

func newStore() *store {
	return &store{
		sessions: make(map[string]*LobbySession),
	}
}

func (s *store) CreateSession(leader *player.Player) *LobbySession {
	s.mx.Lock()
	defer s.mx.Unlock()

	newSession := &LobbySession{
		UUID:     uuid.NewV4().String(),
		LeaderID: leader.GetID(),
		Players:  make([]*player.Player, 0),
		MapID:    0,
	}

	newSession.Players = append(newSession.Players, leader)
	s.sessions[newSession.UUID] = newSession
	leader.LobbyUUID = newSession.UUID

	return newSession
}

func (s *store) JoinSession(p *player.Player, sessionUUID string) (error, *LobbySession) {

	session := s.GetSessionByUUID(sessionUUID)
	if session == nil {
		return errors.New("no lobby session"), nil
	}

	// todo удалять его из других сессий
	session.Players = append(session.Players, p)
	p.LobbyUUID = session.UUID

	return nil, session
}

func (s *store) GetSessionByUUID(uuid string) *LobbySession {
	s.mx.Lock()
	defer s.mx.Unlock()

	return s.sessions[uuid]
}
