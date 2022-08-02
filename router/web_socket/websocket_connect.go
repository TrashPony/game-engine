package web_socket

import (
	"fmt"
	"github.com/TrashPony/game-engine/router/mechanics/factories/players"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: false,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ReadSocket(id int, w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	ws.SetReadLimit(10485760)

	user := players.Users().GetUser(id)
	if user == nil {
		ws.Close()
		return
	}

	clients.AddUser(ws, user)
	fmt.Println("login: " + user.Login)
	Reader(ws, user)
}
