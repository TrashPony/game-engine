package game_loop_move

import (
	"github.com/TrashPony/game_engine/src/web_socket"
)

type PathUnit struct {
	UnitID          int  `json:"id"`
	Speed           int  `json:"s"`
	X               int  `json:"x"`
	Y               int  `json:"y"`
	Z               int  `json:"z"`
	MS              int  `json:"ms"`
	Rotate          int  `json:"r"`
	AngularVelocity int  `json:"av"`
	MapID           int  `json:"m"`
	Animate         bool `json:"a"`
	Direction       bool `json:"d"`
}

func SendMoveUnit(obj moveObject, unitID, x, y, ms, mapID int, rotate, z, speed float64, animate bool, msgs *web_socket.GameLoopMessages) {
	path := &PathUnit{
		UnitID:          unitID,
		Speed:           int(speed * 10),
		X:               x,
		Y:               y,
		Z:               int(z),
		MS:              ms,
		Rotate:          int(rotate),
		AngularVelocity: int(obj.GetAngularVelocity() * 1000),
		MapID:           mapID,
		Animate:         animate,
		Direction:       obj.GetDirection(),
	}

	msgs.AddMessage(path)
}
