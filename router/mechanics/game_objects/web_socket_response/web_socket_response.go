package web_socket_response

type Response struct {
	Responses        []*Response `json:"-"`
	Event            string      `json:"event,omitempty"`
	Service          string      `json:"service,omitempty"`
	Error            string      `json:"error,omitempty"`
	Data             interface{} `json:"data,omitempty"`
	NodeName         string      `json:"-"`
	AllNoBattle      bool        `json:"-"`
	All              bool        `json:"-"`
	ID               int         `json:"-"`
	X                int         `json:"-"`
	Y                int         `json:"-"`
	TeamID           int         `json:"-"`
	UserID           int         `json:"-"`
	PlayerID         int         `json:"-"`
	GameUUID         string      `json:"-"`
	LobbySessionUUID string      `json:"-"`
	OnlyData         bool        `json:"-"`
	BinaryMsg        []byte      `json:"-"`
}
