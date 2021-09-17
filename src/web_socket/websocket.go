package web_socket

import (
	"github.com/TrashPony/game_engine/src/auth"
	"net/http"
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	var login string
	var id int

	login, id = auth.CheckCookie(w, r) // берем из куки данные по логину и ид пользовтеля

	if login == "" || id == 0 || login == "anonymous" {
		println("Соеденение не разрешено: не авторизован")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403 - No Auth!"))
		return
	}

	ReadSocket(id, w, r)
}
