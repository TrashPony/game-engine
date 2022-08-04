package player

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/user"
)

type Player struct {
	ID                 int    `json:"id"`
	Login              string `json:"login"`
	GameUUID           string `json:"game_uuid"`
	NodeName           string `json:"node_name"`
	MapID              int    `json:"map_id"`
	Ready              bool   `json:"ready"`
	Bot                bool   `json:"bot"`
	BehaviorController bool   `json:"-"`
	TeamID             int    `json:"team_id"`
	OwnerID            int    `json:"-"`
	CacheSenderCommand []byte `json:"-"`
	gameUnitStore      *userUnitStore
	owner              *user.User
}

func (client *Player) SetOwner(owner *user.User) {
	client.OwnerID = owner.ID
	client.owner = owner
}

func (client *Player) GetLogin() (login string) {
	return client.Login
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
