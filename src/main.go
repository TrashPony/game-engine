package main

import (
	"fmt"
	"github.com/TrashPony/game_engine/src/ai"
	"github.com/TrashPony/game_engine/src/auth"
	"github.com/TrashPony/game_engine/src/game_loop"
	"github.com/TrashPony/game_engine/src/web_socket"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	_ "net/http/pprof" // todo включает веб профайлер
	"path"
	"runtime"
	"strings"
)

func main() {
	router := mux.NewRouter()
	// auth
	router.HandleFunc("/api/login", auth.Login)
	router.HandleFunc("/api/registration", auth.Registration)
	// socket
	router.HandleFunc("/socket", web_socket.HandleConnections)

	// todo vk auth
	//router.HandleFunc("/api/vk-get-oauth-url", auth.VkGetUrlToAuth)
	//router.HandleFunc("/api/vk-oauth", auth.VkAuth)
	//router.HandleFunc("/api/vk-app-login", auth.VkAppLogin)

	go web_socket.Sender()
	go game_loop.GameLoopInit()
	go ai.BotLifeLoop()

	headersOk := handlers.AllowedHeaders([]string{"Access-Control-Allow-Credentials"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:8082"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	router.NotFoundHandler = vueHandler()

	httpPort := "8080"
	httpsPort := "8081"

	// pprof http://localhost:6060/debug/pprof/
	// go tool pprof --trim=false выводит все без обрезки в свг
	// go tool pprof -svg ./ localhost:6060/debug/pprof/allocs?debug=1 > ./allocs.svg
	// go tool pprof -svg ./ localhost:6060/debug/pprof/heap?debug=1 > ./heap.svg
	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	go func() {
		err := http.ListenAndServeTLS(":"+httpsPort, "../ssl/certificate.pem", "../ssl/key", handlers.CORS(handlers.AllowCredentials(), originsOk, headersOk, methodsOk)(router))
		if err != nil {
			fmt.Println("no ssl")
		}

		log.Println("https server started on :" + httpsPort)
	}()

	log.Println("http server started on :" + httpPort)
	err := http.ListenAndServe(":"+httpPort, handlers.CORS(handlers.AllowCredentials(), originsOk, headersOk, methodsOk)(router)) // запускает веб сервер на 8080 порту
	if err != nil {
		log.Panic(err)
	}
}

func vueHandler() http.Handler {

	handler := http.FileServer(http.Dir("../static/dist"))

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
		http.ServeFile(w, req, path.Join("../static/dist", "/index.html"))
	})
}
