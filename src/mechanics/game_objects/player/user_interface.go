package player

import "encoding/json"

// текущие положение интерфейса пользователя
type Window struct {
	Left   int  `json:"left"`
	Top    int  `json:"top"`
	Height int  `json:"height"`
	Width  int  `json:"width"`
	Open   bool `json:"open"`
}

var AllowWindowSave = map[string]bool{
	"CreateGame": true,
	"miniMap":    true,
}

func (client *Player) SetWindowStateFromJSON(jsonInterface []byte) {
	client.mx.Lock()
	defer client.mx.Unlock()

	err := json.Unmarshal(jsonInterface, &client.userInterface)
	if err != nil {
		println("unmarshal ui " + err.Error())
	}
}

func (client *Player) GetJSONWindowState() string {

	if client == nil {
		return ""
	}

	client.mx.Lock()
	defer client.mx.Unlock()

	jsonInterface, err := json.Marshal(client.userInterface)
	if err != nil {
		println("get json user interface", err.Error())
	}

	return string(jsonInterface)
}

func (client *Player) SetWindowState(name, resolution string, open bool, height, width, left, top int) {

	client.mx.Lock()
	defer client.mx.Unlock()

	if client.userInterface == nil {
		// resolution, window_id, state
		client.userInterface = make(map[string]map[string]*Window)
	}

	_, ok := AllowWindowSave[name]
	if !ok {
		return
	}

	setState := func(window *Window) {
		window.Height = height
		window.Width = width
		window.Left = left
		window.Top = top
		window.Open = open
	}

	resol, ok := client.userInterface[resolution]
	if ok {

		_, ok := resol[name]
		if !ok {
			resol[name] = &Window{}
		}
		setState(resol[name])

	} else {
		client.userInterface[resolution] = make(map[string]*Window)
		client.userInterface[resolution][name] = &Window{}
		setState(client.userInterface[resolution][name])
	}
}
