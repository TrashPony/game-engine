package web_socket_response

type Response struct {
	Responses   []*Response `json:"-"`
	Event       string      `json:"event,omitempty"`
	Service     string      `json:"service,omitempty"`
	Error       string      `json:"error,omitempty"`
	Data        interface{} `json:"data,omitempty"`
	AllNoBattle bool        `json:"-"`
	All         bool        `json:"-"`
	ID          int         `json:"-"`
	X           int         `json:"-"`
	Y           int         `json:"-"`
	ToX         int         `json:"-"`
	ToY         int         `json:"-"`
	TeamID      int         `json:"-"`
	UserID      int         `json:"-"`
	PlayerID    int         `json:"-"`
	GameUUID    string      `json:"-"`
	OnlyData    bool        `json:"-"`
	BinaryMsg   []byte      `json:"-"`
	CheckTo     bool        `json:"-"`
	Radius      int         `json:"-"`
}
