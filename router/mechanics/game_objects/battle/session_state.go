package battle

import (
	"sync"
)

type SessionPlayer struct {
	UUID         string `json:"uuid"`
	PlayerID     int    `json:"player_id"`
	Login        string `json:"login"`
	GameType     string `json:"game_type"`
	Live         bool   `json:"live"`
	Leave        bool   `json:"leave"`
	Deaths       int    `json:"deaths"`
	DoRespawn    bool   `json:"-"`
	CountRespawn int    `json:"count_respawn"`
	RespawnTime  int    `json:"respawn_time"`
	mx           sync.RWMutex
}
