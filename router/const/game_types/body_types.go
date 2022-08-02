package game_types

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/body"
	"math/rand"
)

var bodyTypes = map[int]body.Body{
	1: {
		ID:     1,
		Name:   "replic_start_body",
		Scale:  22,
		Length: 58,
		Height: 40,
		Width:  62,
		Radius: 88,
		// -- balance state
		MaxHP:              600,
		RangeRadar:         600,
		RangeView:          260, // 260
		RecoveryPower:      28,
		RecoveryPowerCycle: 128,
		MaxPower:           2500,
		// -- move balance state
		Speed:         44,
		ReverseSpeed:  44,
		PowerFactor:   0.4,
		ReverseFactor: 0.4,
		TurnSpeed:     0.9,
		MoveDrag:      0.65,
		AngularDrag:   0.70,
		// --
		Weight:      2000,
		ChassisType: "caterpillars",
		Weapons: map[int]*body.WeaponSlot{
			1: {
				Number:  1,
				XAttach: 30,
				YAttach: 64,
			},
		},
		PowerPoints: 100,
	},
	2: {
		ID:     2,
		Name:   "reverses_start_body",
		Scale:  28,
		Length: 58,
		Height: 40,
		Width:  62,
		Radius: 88,
		// -- balance state
		MaxHP:              600,
		RangeRadar:         600,
		RangeView:          260, // 260
		RecoveryPower:      28,
		RecoveryPowerCycle: 128,
		MaxPower:           2500,
		// -- move balance state
		Speed:         44,
		ReverseSpeed:  44,
		PowerFactor:   0.4,
		ReverseFactor: 0.4,
		TurnSpeed:     0.9,
		MoveDrag:      0.65,
		AngularDrag:   0.70,
		// --
		Weight:      2000,
		ChassisType: "caterpillars",
		Weapons: map[int]*body.WeaponSlot{
			1: {
				Number:  1,
				XAttach: 31,
				YAttach: 64,
			},
		},
		PowerPoints: 100,
	},
}

func GetNewBody(id int) *body.Body {
	newBody, ok := bodyTypes[id]
	if !ok {
		return nil
	}

	newBody.Weapons = make(map[int]*body.WeaponSlot)

	for number, slot := range bodyTypes[id].Weapons {
		newBody.Weapons[number] = &body.WeaponSlot{
			Number:  slot.Number,
			XAttach: slot.XAttach,
			YAttach: slot.YAttach,
		}
	}

	newBody.Texture = newBody.Name
	newBody.SetWheelsPositions()
	return &newBody
}

func GetAllBody() map[int]body.Body {
	allBodies := make(map[int]body.Body)

	for id, bodyType := range bodyTypes {
		bodyType.SetWheelsPositions()
		allBodies[id] = bodyType
	}

	return allBodies
}

func GetRandomBody() *body.Body {
	allBodies := make([]int, 0)

	for _, bodyType := range bodyTypes {

		if bodyType.OnlyConfig {
			continue
		}

		allBodies = append(allBodies, bodyType.ID)
	}

	return GetNewBody(allBodies[rand.Intn(len(allBodies))])
}
