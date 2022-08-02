package web_socket

import (
	"fmt"
	"github.com/TrashPony/game-engine/router/generate_ids"
	"net/http"
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {

	tokenString := r.URL.Query().Get("token")
	fmt.Println(tokenString)

	//if err != nil {
	//	println("Соеденение не разрешено: не авторизован: ", err.Error())
	//	w.WriteHeader(code)
	//	w.Write([]byte(err.Error()))
	//	return
	//}

	ReadSocket(generate_ids.GetPlayerFakeID(), w, r)
}
