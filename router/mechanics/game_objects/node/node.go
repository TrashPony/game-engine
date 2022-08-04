package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/player"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/rpc_request"
	"github.com/valyala/gorpc"
	"time"
)

type Node struct {
	Name        string `json:"name"`
	Url         string `json:"url"`
	MaxSessions int    `json:"max_sessions"`
	Error       bool   `json:"error"`
	rpcClient   *gorpc.Client
}

func (n *Node) Connect() error {

	n.rpcClient = gorpc.NewTCPClient(n.Url)
	n.rpcClient.Start()

	n.rpcClient.LogError = func(format string, args ...interface{}) {
		n.Error = true
	}

	try := 50
	for n.rpcClient.Conns == 0 {

		try--
		if try == 0 {
			return errors.New("timeout")
		}

		if n.Error {
			return errors.New("error connect")
		}

		time.Sleep(time.Millisecond * 32)
	}

	return n.InitNode()
}

func (n *Node) Stop() {

	defer func() {
		recover()
	}()

	if n.rpcClient != nil {
		n.rpcClient.Stop()
	}
}

func (n *Node) InitNode() error {
	fmt.Println("InitNode: start")
	req := rpc_request.Request{Event: "InitNode"}

	_, err := n.rpcClient.Call(req)
	if err != nil {
		return err
	}

	fmt.Println("OK:", n.Name)
	return nil
}

func (n *Node) CreateBattle(startPlayers []*player.Player) string {

	data := map[string]interface{}{
		"start_players": startPlayers,
		"map_id":        1,
	}

	dataJson, err := json.Marshal(data)
	if err != nil {
		fmt.Println("failed create game")
		return ""
	}

	req := rpc_request.Request{Event: "CreateBattle", Data: string(dataJson)}
	response, err := n.rpcClient.Call(req)
	if err != nil {
		return ""
	}

	if response == nil {
		return ""
	}

	resp := response.(rpc_request.Request)
	return resp.UUID
}

func (n *Node) FindBattle(battleUUID string) bool {

	req := rpc_request.Request{Event: "FindBattle", UUID: battleUUID}

	response, err := n.rpcClient.Call(req)
	if err != nil {
		return false
	}

	if response == nil {
		return false
	}

	resp := response.(rpc_request.Request)
	return resp.OK
}

func (n *Node) InitBattle(battleUUID string, playerID int) rpc_request.Request {

	req := rpc_request.Request{Event: "InitBattle", UUID: battleUUID, ID: playerID}

	response, err := n.rpcClient.Call(req)
	if err != nil {
		return rpc_request.Request{}
	}

	if response == nil {
		return rpc_request.Request{}
	}

	resp := response.(rpc_request.Request)
	return resp
}

func (n *Node) StartLoad(battleUUID string, playerID int) rpc_request.Request {

	req := rpc_request.Request{Event: "StartLoad", UUID: battleUUID, ID: playerID}

	response, err := n.rpcClient.Call(req)
	if err != nil {
		return rpc_request.Request{}
	}

	if response == nil {
		return rpc_request.Request{}
	}

	resp := response.(rpc_request.Request)
	return resp
}

func (n *Node) Input(battleUUID string, playerID int, w, a, s, d, sp, st, z bool, x, y int, fire bool) {
	req := rpc_request.Request{Event: "i", UUID: battleUUID, ID: playerID, W: w, A: a, S: s, D: d, Sp: sp, St: st, Z: z, X: x, Y: y, Fire: fire}
	_, _ = n.rpcClient.Call(req)
}

func (n *Node) CreateUnit(battleUUID string, playerID, x, y int) {
	req := rpc_request.Request{Event: "CreateUnit", UUID: battleUUID, ID: playerID, X: x, Y: y}
	n.rpcClient.Call(req)
}

func (n *Node) CreateBot(battleUUID string, playerID, x, y, teamID int) {
	req := rpc_request.Request{Event: "CreateBot", UUID: battleUUID, ID: playerID, X: x, Y: y, TeamID: teamID}
	n.rpcClient.Call(req)
}

func (n *Node) CreateObj(battleUUID string, playerID, objID, x, y, teamID int) {
	req := rpc_request.Request{Event: "CreateObj", UUID: battleUUID, ObjectID: objID, ID: playerID, X: x, Y: y, TeamID: teamID}
	n.rpcClient.Call(req)
}
