package rpc

import (
	"fmt"
	"github.com/TrashPony/game-engine/node/mechanics/factories/quick_battles"
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/rpc_request"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/web_socket_response"
	"github.com/TrashPony/game-engine/router/rpc"
	"github.com/valyala/gorpc"
	"time"
)

var rpcNode *rpc.RPC
var nodeName string

func Checker() {
	for {
		time.Sleep(time.Second)
		checker()
	}
}

func checker() {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println("RECOVER: ", r)
		}
	}()

	if rpcNode == nil {
		fmt.Println("InitNode: wait_connect")
		rpcNode = initRpc()

		for rpcNode.VeliriMainClient.Conns == 0 {
			time.Sleep(time.Second)
		}

		InitNode()
		return
	}

	if rpcNode.Error {
		fmt.Println("RpcNode: error")
		rpcNode.VeliriMainClient.Stop()
		rpcNode.Server.Stop()
		rpcNode = nil
	}
}

func GetRPC() *rpc.RPC {
	return rpcNode
}

func initRpc() *rpc.RPC {
	rpc.GobRegister()

	r := &rpc.RPC{
		Server:           gorpc.NewTCPServer(_const.Config.GetParams("nodeUrl"), router),
		VeliriMainClient: gorpc.NewTCPClient(_const.Config.GetParams("veliriURL")),
	}

	r.VeliriMainClient.LogError = func(format string, args ...interface{}) {
		r.Error = true
		fmt.Println(args)
	}

	r.Server.LogError = func(format string, args ...interface{}) {
		r.Error = true
		fmt.Println(args)
	}

	go func() {
		r.VeliriMainClient.Start()
		err := r.Server.Serve()
		if err != nil {
			r.Error = true
			fmt.Println("Init rpc server error:", err)
		}
	}()

	return r
}

func router(clientAddr string, request interface{}) interface{} {

	req := request.(rpc_request.Request)

	if req.Event == "CreateBattle" {
		uuid := createBattle(&req)
		return rpc_request.Request{UUID: uuid}
	}

	if req.Event == "FindBattle" {
		return rpc_request.Request{OK: quick_battles.Battles.GetBattleByUUID(req.UUID) != nil}
	}

	if req.Event == "InitNode" {
		fmt.Println("Init veliri server: ok")
		return rpc_request.Request{}
	}

	battle := quick_battles.Battles.GetBattleByUUID(req.UUID)
	if battle == nil {
		return rpc_request.Request{Response: web_socket_response.Response{
			Event: "error", Error: "no_battle",
		}}
	}

	p := battle.GetPlayerByID(req.ID)
	if p == nil {
		return rpc_request.Request{Response: web_socket_response.Response{
			Event: "error", Error: "no_player",
		}}
	}

	switch req.Event {
	case "i":
		i(p, &req)
	case "InitBattle":
		return rpc_request.Request{Response: initBattle(battle, p)}
	case "StartLoad":
		return rpc_request.Request{Response: startLoad(battle, p)}
	case "CreateUnit":
		createUnit(battle, p, req.X, req.Y)
	case "CreateBot":
		createBot(battle, p, req.X, req.Y, req.TeamID)
	case "CreateObj":
		createObj(battle, p, req.X, req.Y, req.TeamID, req.ObjectID)
	}

	return nil
}

func InitNode() {
	fmt.Println("InitNode: start")

	r := GetRPC()
	if r == nil {
		return
	}

	req := rpc_request.Request{Event: "InitNode", UUID: _const.Config.GetParams("nodeUrl"), Slot: _const.Config.GetIntParams("maxSessions")}
	data, err := r.VeliriMainClient.Call(req)
	if err != nil {
		fmt.Println("InitNode: failed")
		r.Error = true
		return
	}

	resp := data.(rpc_request.Request)
	if !resp.OK {
		fmt.Println("InitNode: failed")
		r.Error = true
		return
	}

	nodeName = resp.Event
	fmt.Println("InitNode: end, MyName:", resp.Event)
}

func SendData(data *web_socket_response.Response) {

	r := GetRPC()
	if r == nil {
		return
	}

	data.NodeName = nodeName
	req := rpc_request.Request{Event: "SendData", Response: *data}
	_, err := r.VeliriMainClient.Call(req)
	if err != nil {
		r.Error = true
	}
}
