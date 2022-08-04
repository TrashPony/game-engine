package game_loop

import (
	"fmt"
	binary_msg2 "github.com/TrashPony/game-engine/node/binary_msg"
	"github.com/TrashPony/game-engine/node/mechanics/watch"
	"github.com/TrashPony/game-engine/node/rpc"
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/player"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/web_socket_response"
	"github.com/TrashPony/game-engine/router/web_socket"
)

func SendMessagesToUsers(b *battle.Battle, players []*player.Player, messagesStore *web_socket.MessagesStore, units []*unit.Unit) {

	for _, p := range players {

		team := b.Teams[p.GetTeamID()]
		if team == nil {
			continue
		}

		// что бы каждый раз не создавать и растить масив заного
		if p.CacheSenderCommand == nil {
			p.CacheSenderCommand = make([]byte, 0, 128)
		}
		p.CacheSenderCommand = p.CacheSenderCommand[:0]
		p.CacheSenderCommand = append(p.CacheSenderCommand, 100)

		radarMsg := make([]byte, 0)

		for _, g := range messagesStore.GetMessageGroups() {

			if g.Type == "move" {

				var moveRadarBinMsg []byte
				var playerMsgLength int

				startIndex := len(p.CacheSenderCommand)                                                          // берем индекс что бы можно было потом вернутся
				p.CacheSenderCommand = append(p.CacheSenderCommand, binary_msg2.GetIntBytes(playerMsgLength)...) // резервируем место для числа длины данных
				p.CacheSenderCommand, moveRadarBinMsg, playerMsgLength = moveMessagesToUser(p.CacheSenderCommand, team, g.Messages, g.Attributes["type_obj"], b.Map.Id, b, units)
				binary_msg2.ReuseByteSlice(&p.CacheSenderCommand, startIndex, binary_msg2.GetIntBytes(playerMsgLength)) // вставляем эту длину в зарезервированое место

				// сообщение меток, если будут другие обьекты их над будет акамулировать - поэтму так
				radarMsg = append(radarMsg, moveRadarBinMsg...)
			}

			if g.Type == "bin" {
				var binDataLength int
				startIndex := len(p.CacheSenderCommand)                                                                            // берем индекс что бы можно было потом вернутся
				p.CacheSenderCommand = append(p.CacheSenderCommand, binary_msg2.GetIntBytes(binDataLength)...)                     // резервируем место для числа длины данных
				p.CacheSenderCommand, binDataLength = ResponseToBin(p.CacheSenderCommand, team, p, g.Messages, b.Map.Id, b, units) // вставляем данные и возвращаем длинну
				binary_msg2.ReuseByteSlice(&p.CacheSenderCommand, startIndex, binary_msg2.GetIntBytes(binDataLength))              // вставляем эту длину в зарезервированое место
			}
		}

		p.CacheSenderCommand = append(p.CacheSenderCommand, binary_msg2.GetIntBytes(len(radarMsg))...)
		p.CacheSenderCommand = append(p.CacheSenderCommand, radarMsg...)

		if len(p.CacheSenderCommand) > 1 && !p.Bot {
			rpc.SendData(&web_socket_response.Response{PlayerID: p.GetID(), BinaryMsg: p.CacheSenderCommand, GameUUID: b.UUID})
		}
	}
}

func moveMessagesToUser(command []byte, team *battle.Team, msgs *web_socket.GameLoopMessages, objType string, mpID int, b *battle.Battle, units []*unit.Unit) ([]byte, []byte, int) {

	//playerMsg := make([]byte, 0)
	playerRadarMsg := make([]byte, 0)
	playerMsgLength := 0

	for i := range msgs.Messages {
		msg, radarMark := filterMoveTo(team, &msgs.Messages[i], objType, mpID, b, units)
		if msg != nil {
			if radarMark {
				playerRadarMsg = append(playerRadarMsg, msg...)
			} else {
				playerMsgLength += len(msg)
				command = append(command, msg...)
			}
		}
	}

	return command, playerRadarMsg, playerMsgLength
}

func filterMoveTo(team *battle.Team, resp *web_socket_response.Response, objType string, mpID int, b *battle.Battle, units []*unit.Unit) ([]byte, bool) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in filterMoveTo", r)
		}
	}()

	view, radarVisible := watch.CheckViewCoordinate(team, resp.X, resp.Y, b, units, 0)
	if view {
		return resp.BinaryMsg, false
	}

	if radarVisible {
		// получаем метку, подменяем обьект на метку затирая методанные, а путь оставляем прежним
		radarMark := team.GetVisibleObjectByTypeAndID(objType, resp.ID)
		// если метки нет значит радар еще не нашел нехера и тупо игнорируем
		if radarMark != nil {
			return binary_msg2.CreateMarkBinaryMove(radarMark.ID, resp.X, resp.Y, _const.ServerTick), true
		}
	}

	return nil, false
}

func ResponseToBin(command []byte, team *battle.Team, p *player.Player, msgs *web_socket.GameLoopMessages, mpID int, b *battle.Battle, units []*unit.Unit) ([]byte, int) {
	//binMsg := make([]byte, 0)
	binMsgLength := 0
	for i := range msgs.Messages {

		if msgs.Messages[i].All ||
			(msgs.Messages[i].PlayerID != 0 && msgs.Messages[i].PlayerID == p.GetID()) ||
			(msgs.Messages[i].TeamID != 0 && msgs.Messages[i].TeamID == p.GetTeamID()) {

			command = append(command, msgs.Messages[i].BinaryMsg...)
			binMsgLength += len(msgs.Messages[i].BinaryMsg)
			continue
		}

		if msgs.Messages[i].PlayerID != 0 && msgs.Messages[i].PlayerID != p.GetID() {
			continue
		}

		if msgs.Messages[i].TeamID > 0 && msgs.Messages[i].TeamID != p.TeamID {
			continue
		}

		view, _ := watch.CheckViewCoordinate(team, msgs.Messages[i].X, msgs.Messages[i].Y, b, units, msgs.Messages[i].Radius)
		if !view && msgs.Messages[i].CheckTo {
			view, _ = watch.CheckViewCoordinate(team, msgs.Messages[i].ToX, msgs.Messages[i].ToY, b, units, msgs.Messages[i].Radius)
		}

		if view {
			command = append(command, msgs.Messages[i].BinaryMsg...)
			binMsgLength += len(msgs.Messages[i].BinaryMsg)
		}
	}

	return command, binMsgLength
}
