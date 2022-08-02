package main

import (
	"github.com/TrashPony/game-engine/node/ai"
	"github.com/TrashPony/game-engine/node/game_loop"
	"github.com/TrashPony/game-engine/node/rpc"
	//_ "net/http/pprof" // todo включает веб профайлер
)

func main() {

	//runtime.SetMutexProfileFraction(1)
	//runtime.SetBlockProfileRate(1)
	//runtime.MemProfileRate = 1

	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()

	go ai.BotLifeLoop()
	go rpc.Checker()
	game_loop.GameLoopInit()
}
