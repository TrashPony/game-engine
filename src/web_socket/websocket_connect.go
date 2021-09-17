package web_socket

import (
	"fmt"
	"github.com/TrashPony/game_engine/src/mechanics/factories/players"
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

	ws, err := upgrader.Upgrade(w, r, nil) // запрос GET для перехода на протокол websocket
	if err != nil {
		log.Fatal(err)
	}

	ws.SetReadLimit(10485760)

	user := players.Users.GetUser(id)
	if user == nil {
		ws.Close()
		return
	}

	//player := players.Users.GetPlayer(user.GetCurrentPlayerID(), 0)
	//if player == nil {
	//	ws.Close()
	//	return
	//}

	clients.AddUser(ws, user)
	fmt.Println("login: " + user.Login)
	Reader(ws, user)
}
