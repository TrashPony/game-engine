package bullet

import (
	"github.com/TrashPony/game-engine/node/binary_msg"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/ammo"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/target"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/weapon"
)

type Bullet struct {
	ID          int            `json:"id"`
	Weapon      *weapon.Weapon `json:"-"`
	Ammo        *ammo.Ammo     `json:"ammo"`
	EquipID     int            `json:"-"`
	Rotate      float64        `json:"rotate"`
	X           int            `json:"x"`
	Y           int            `json:"y"`
	Z           float64        `json:"z"` // определяет "высоту" пули (сильнее отдалять тени)
	Speed       int            `json:"speed"`
	Target      *target.Target `json:"target"`
	ChaseTarget *target.Target `json:"-"`
	OwnerID     int            `json:"owner_id"`   // какой игрок стрелял
	OwnerType   string         `json:"owner_type"` // unit, structure
	OwnerTeamID int            `json:"-"`
	IgnoreOwner bool           `json:"ignore_owner"`
	MaxRange    int            `json:"max_range"`
	FirePos     int            `json:"-"`
	MapID       int            `json:"map_id"`
	HP          int            `json:"destroy"`
	StartX      int            `json:"-"`
	StartY      int            `json:"-"`
	StartZ      float64        `json:"-"`
	StartRadian float64        `json:"start_radian"`
	Damage      int            `json:"-"`
	MaxAngle    float64        `json:"-"`
	MinAngle    float64        `json:"-"`

	RealX            float64        `json:"-"`
	RealY            float64        `json:"-"`
	DistanceTraveled float64        `json:"-"`
	RealSpeed        float64        `json:"-"`
	RadRotate        float64        `json:"-"`
	Attributes       map[string]int `json:"-"`

	XVelocity float64 `json:"x_velocity"`
	YVelocity float64 `json:"y_velocity"`

	CacheJson      []byte `json:"-"`
	CreateJsonTime int64  `json:"-"`
	end            bool
}

func (b *Bullet) GetUpdateData(mapTime int64) []byte {
	return []byte{}
}

func (b *Bullet) GetX() int {
	return b.X
}

func (b *Bullet) SetX(x int) {
	b.X = x
}

func (b *Bullet) GetY() int {
	return b.Y
}

func (b *Bullet) SetY(y int) {
	b.Y = y
}

func (b *Bullet) GetZ() float64 {
	return b.Z
}

func (b *Bullet) SetZ(z float64) {
	b.Z = z
}

func (b *Bullet) GetRotate() float64 {
	return b.Rotate
}

func (b *Bullet) SetRotate(rotate float64) {
	b.Rotate = rotate
}

func (b *Bullet) SetID(id int) {
	b.ID = id
}

func (b *Bullet) GetEnd() bool {
	return b.end
}

func (b *Bullet) SetEnd(end bool) {
	b.end = end
}

func (b *Bullet) GetBytes(mapTime int64) []byte {

	if b.CreateJsonTime == mapTime && len(b.CacheJson) > 0 {
		return b.CacheJson
	}

	command := []byte{}

	command = append(command, byte(b.Ammo.ID))
	command = append(command, binary_msg.GetIntBytes(b.ID)...)
	command = append(command, binary_msg.GetIntBytes(b.GetX())...)
	command = append(command, binary_msg.GetIntBytes(b.GetY())...)
	command = append(command, binary_msg.GetIntBytes(int(b.GetZ()))...)
	command = append(command, binary_msg.GetIntBytes(int(b.GetRotate()))...)

	b.CacheJson = command
	b.CreateJsonTime = mapTime

	return command
}
