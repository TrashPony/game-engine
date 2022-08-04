package main

import (
	"github.com/TrashPony/game-engine/router/rpc"
	"github.com/TrashPony/game-engine/router/web_socket"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"path"
	"strings"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/socket", web_socket.HandleConnections)

	go rpc.GetRPC()
	go web_socket.Sender()
	go web_socket.InitNodeChecker()

	headersOk := handlers.AllowedHeaders([]string{"Access-Control-Allow-Credentials"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:8083"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	router.NotFoundHandler = vueHandler()

	httpPort := "8086"

	log.Println("http server started on :" + httpPort)
	err := http.ListenAndServe(":"+httpPort, handlers.CORS(handlers.AllowCredentials(), originsOk, headersOk, methodsOk)(router))
	if err != nil {
		log.Panic(err)
	}
}

func vueHandler() http.Handler {

	handler := http.FileServer(http.Dir("./static/dist"))

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_path := req.URL.Path

		// static files
		if strings.Contains(_path, ".") || _path == "/" {

			splitUrl := strings.Split(_path, "/")
			if len(splitUrl) > 0 {
				// т.к. все файлы находятся в 1 директории то перенаправляем все запросы к ней
				req.URL.Path = splitUrl[len(splitUrl)-1]
			}

			handler.ServeHTTP(w, req)
			return
		}

		// the all 404 gonna be served as root
		http.ServeFile(w, req, path.Join("./static/dist", "/index.html"))
	})
}
