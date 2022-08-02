package game_loop_move

import (
	"github.com/TrashPony/game-engine/node/binary_msg"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/web_socket_response"
	"github.com/TrashPony/game-engine/router/web_socket"
)

func SendMoveUnit(obj moveObject, unitID, x, y, ms, mapID int, rotate, z, speed float64, animate bool, messagesStore *web_socket.MessagesStore) {

	if obj.CheckHandBrake() || obj.GetPhysicalModel().IsFly() {
		speed = 0
	}

	messagesStore.AddMsg("unitMoveMsg", "move", web_socket_response.Response{
		BinaryMsg: binary_msg.CreateBinaryUnitMoveMsg(unitID, int(speed*10), x, y, int(z), ms, int(rotate),
			int(obj.GetAngularVelocity()*1000), animate, obj.GetDirection(), obj.GetPhysicalModel().IsFly(),
			obj.GetPhysicalModel().WASD.GetA(), obj.GetPhysicalModel().WASD.GetD(), obj.GetPhysicalModel().WASD.GetW(),
			obj.CheckHandBrake()),
		ID: unitID,
		X:  x,
		Y:  y,
	}, map[string]string{"type_obj": "unit"})
}

func SendMoveObject(objID, x, y, ms int, rotate float64, messagesStore *web_socket.MessagesStore) {
	messagesStore.AddMsg("objectMoveMsg", "move", web_socket_response.Response{
		BinaryMsg: binary_msg.CreateBinaryObjectMove(objID, x, y, ms, int(rotate)),
		ID:        objID,
		X:         x,
		Y:         y,
	}, map[string]string{"type_obj": "object"})
}
