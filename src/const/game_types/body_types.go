package game_types

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/body"
)

// TODO тестовые данные
var BodyTypes = map[int]body.Body{
	1: {
		Name:    "apd_light",
		Texture: "apd_light_skin_1",
		MaxHP:   100,
		Scale:   40,
		Length:  64,
		Height:  64,
		Width:   71,
		Radius:  77,
		Weapons: map[int]*body.WeaponSlot{
			1: {
				Number:  1,
				XAttach: 43,
				YAttach: 64,
			},
			2: {
				Number:  2,
				XAttach: 104,
				YAttach: 64,
			},
		},
	},
}

func GetNewBody(id int) *body.Body {
	newBody := BodyTypes[id]
	newBody.Weapons = make(map[int]*body.WeaponSlot)

	for number, slot := range BodyTypes[id].Weapons {
		newBody.Weapons[number] = &body.WeaponSlot{
			Number:  slot.Number,
			XAttach: slot.XAttach,
			YAttach: slot.YAttach,
		}
	}

	return &newBody
}
