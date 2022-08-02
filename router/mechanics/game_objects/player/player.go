package player

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/user"
	"sync"
)

type Player struct {
	ID                 int         `json:"id"`
	Login              string      `json:"login"`
	GameUUID           string      `json:"game_uuid"`
	LobbyUUID          string      `json:"lobby_uuid"`
	NodeName           string      `json:"node_name"`
	MapID              int         `json:"map_id"`
	Ready              bool        `json:"ready"`
	Bot                bool        `json:"bot"`
	BehaviorController bool        `json:"-"`
	TeamID             int         `json:"team_id"`
	BattleRank         *BattleRank `json:"battle_rank"`
	GroupID            int         `json:"group_id"`
	OwnerID            int         `json:"-"`
	CacheSenderCommand []byte      `json:"-"`
	gameUnitStore      *userUnitStore
	mx                 sync.RWMutex
	owner              *user.User
}

func (client *Player) SetOwner(owner *user.User) {
	client.OwnerID = owner.ID
	client.owner = owner
}

func (client *Player) GetOwner() *user.User {
	return client.owner
}

func (client *Player) GetType() string {
	return "player"
}

type Contact struct {
	PlayerID int    `json:"player_id"`
	Type     string `json:"type"` // friend/block
}

type BattleRank struct {
	Lvl            int `json:"lvl"`
	Points         int `json:"points"`
	NeedPointsToUp int `json:"need_points_to_up"`
	mx             sync.Mutex
}

func (client *Player) CheckViewCoordinate(x, y int) (bool, bool) {
	return true, true
}

func (client *Player) SetLogin(login string) {
	client.Login = login
}

func (client *Player) GetLogin() (login string) {
	return client.Login
}

func (client *Player) SetID(id int) {
	client.ID = id
}

func (client *Player) GetID() (id int) {
	return client.ID
}

func (client *Player) GetReady() bool {
	return client.Ready
}

func (client *Player) SetReady(ready bool) {
	client.Ready = ready
}

func (client *Player) GetTeamID() int {
	return client.TeamID
}

func (client *Player) GetGameUnitsStore() *userUnitStore {
	if client.gameUnitStore == nil {
		client.gameUnitStore = &userUnitStore{units: make([]*unit.Unit, 0)}
	}

	return client.gameUnitStore
}

func (client *Player) ClearUnitsStore() {
	client.gameUnitStore = &userUnitStore{units: make([]*unit.Unit, 0)}
}
