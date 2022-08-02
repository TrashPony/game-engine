package web_socket

import (
	"fmt"
	"github.com/TrashPony/game-engine/router/mechanics/factories/nodes"
	"github.com/TrashPony/game-engine/router/mechanics/factories/players"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/web_socket_response"
	"time"
)

func InitNodeChecker() {
	go func() {
		for {

			for n := range nodes.Nodes().RangeNodes() {
				if n.Error {

					n.Stop()
					nodes.Nodes().RemoveNode(n.Name)
					fmt.Println("Remove node, name:", n.Name)

					for client := range clients.GetUsersChan() {
						player := players.Users().GetPlayer(client.GetCurrentPlayerID(), client.GetID())
						if player != nil && player.NodeName == n.Name {

							toLobby(player, client)
							SendMessage(web_socket_response.Response{
								Event:  "NodeFailed",
								Error:  "node_failed",
								UserID: client.GetID(),
							})
						}
					}
				}
			}

			time.Sleep(time.Second)
		}
	}()
}
