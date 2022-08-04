package battle

type SessionPlayer struct {
	UUID     string `json:"uuid"`
	PlayerID int    `json:"player_id"`
	Login    string `json:"login"`
	Live     bool   `json:"live"`
}
