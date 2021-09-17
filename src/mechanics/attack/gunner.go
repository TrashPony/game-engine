package attack

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/body"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/burst_of_shots"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/target"
)

type Gunner interface {
	GetX() int
	GetY() int
	GetRotate() float64
	GetGunRotate(int) float64
	GetWeaponSlot(int) *body.WeaponSlot
	RangeWeaponSlots() <-chan *body.WeaponSlot
	GetFirePos(int) *game_math.Positions
	GetMapHeight() float64
	GetWeaponFirePos(int) []*game_math.Positions
	GetWeaponTarget() *target.Target
	SetWeaponTarget(*target.Target)
	GetDamage(int) int
	GetWeaponMaxRange(float64, int) (int, float64)
	GetWeaponMinRange(float64, int) (int, float64)
	GetWeaponAccuracy(int) int
	GetWeaponPosInMap(int) (int, int)
	GetBurstOfShots() *burst_of_shots.BurstOfShots
	GetGunRotateSpeed(int) int
	SetGunRotate(float64, int)
}
