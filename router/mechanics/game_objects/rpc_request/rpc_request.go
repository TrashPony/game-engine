package rpc_request

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/web_socket_response"
)

type Request struct {
	Event    string                       `json:"-"`
	OK       bool                         `json:"-"`
	UUID     string                       `json:"-"`
	ID       int                          `json:"-"`
	Response web_socket_response.Response `json:"-"`
	Slot     int                          `json:"-"`
	W        bool                         `json:"-"`
	A        bool                         `json:"-"`
	S        bool                         `json:"-"`
	D        bool                         `json:"-"`
	Z        bool                         `json:"-"`
	Sp       bool                         `json:"-"`
	St       bool                         `json:"-"`
	Fire     bool                         `json:"-"`
	X        int                          `json:"-"`
	Y        int                          `json:"-"`
	Data     string                       `json:"-"`
	TeamID   int                          `json:"-"`
	ObjectID int                          `json:"-"`
}
