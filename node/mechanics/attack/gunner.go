package attack

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/body"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/burst_of_shots"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/position"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/target"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/visible_objects"
	"sync"
)

type Gunner interface {
	GetX() int
	GetY() int
	GetRotate() float64
	GetGunRotate(int) float64
	GetWeaponSlot(int) *body.WeaponSlot
	RangeWeaponSlots() map[int]*body.WeaponSlot
	GetFirePos(int) *position.Positions
	GetMapHeight() float64
	GetWeaponFirePos(int) []*position.Positions
	GetWeaponFirePosOne(int, int) *position.Positions
	GetWeaponTarget() *target.Target
	SetWeaponTarget(*target.Target)
	GetDamage(int) int
	GetWeaponMaxRange(float64, int, bool) (int, float64)
	GetWeaponMinRange(float64, int) (int, float64)
	GetWeaponAccuracy(int) int
	GetWeaponPosInMap(int) (int, int)
	GetBurstOfShots() *burst_of_shots.BurstOfShots
	GetGunRotateSpeed(int) int
	SetGunRotate(float64, int)
	GetWeaponReloadAmmoTime(int) int
	GetWeaponReloadTime(int) int
	GetMinDamage(int) int
	GetMaxDamage(int) int
	GetTeamID() int
	UnsafeRangeVisibleObjects() ([]*visible_objects.VisibleObject, *sync.RWMutex)
}
