package rpc

import (
	"encoding/gob"
	"fmt"
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/mechanics/factories/nodes"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/rpc_request"
	"github.com/TrashPony/game-engine/router/web_socket"
	uuid "github.com/satori/go.uuid"
	"github.com/valyala/gorpc"
	"log"
)

var rpc *RPC

func GetRPC() *RPC {

	if rpc == nil {
		rpc = initRpc()
	}

	return rpc
}

type RPC struct {
	Server           *gorpc.Server
	VeliriMainClient *gorpc.Client
	Error            bool
}

func GobRegister() {
	gob.Register(rpc_request.Request{})
	gob.Register(map[int]string{})
	gob.Register(map[int][]byte{})
}

func initRpc() *RPC {
	GobRegister()

	rpc := &RPC{
		Server: gorpc.NewTCPServer(_const.Config.GetParams("veliriURL"), router),
	}

	go func() {
		err := rpc.Server.Serve()
		if err != nil {
			log.Fatal(err)
		}

		rpc.Server.LogError = func(format string, args ...interface{}) {
			// TODO
		}
	}()

	fmt.Println("RPC READY")
	return rpc
}

func router(clientAddr string, request interface{}) interface{} {

	req := request.(rpc_request.Request)

	switch req.Event {
	case "InitNode":
		newNodeUUID := uuid.NewV1().String()

		fmt.Println("New node: " + newNodeUUID)
		_, err := nodes.Nodes().AddNode(newNodeUUID, req.UUID, req.Slot)
		if err != nil {
			return rpc_request.Request{OK: false}
		}

		if !_const.MasterInit {
			_const.MasterInit = true
			fmt.Println("Master init: complete")
		}

		return rpc_request.Request{Event: newNodeUUID, OK: true}
	case "SendData":
		web_socket.SendMessage(req.Response)
	}

	return nil
}
