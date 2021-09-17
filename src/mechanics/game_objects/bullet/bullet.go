package bullet

import (
	"encoding/json"
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/ammo"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/target"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/weapon"
	"math"
	"sync"
)

type Bullet struct {
	UUID        string               `json:"-"`
	ID          int                  `json:"id"`
	Weapon      *weapon.Weapon       `json:"-"`
	Ammo        *ammo.Ammo           `json:"ammo"`
	Rotate      float64              `json:"rotate"`
	X           int                  `json:"x"`
	Y           int                  `json:"y"`
	Z           float64              `json:"z"` // определяет "высоту" пули (сильнее отдалять тени)
	Speed       int                  `json:"speed"`
	Target      *target.Target       `json:"target"`
	OwnerID     int                  `json:"owner_id"`   // какой игрок стрелял
	OwnerType   string               `json:"owner_type"` // unit, structure
	IgnoreOwner bool                 `json:"ignore_owner"`
	MaxRange    int                  `json:"max_range"`
	FirePos     *game_math.Positions `json:"-"`
	MapID       int                  `json:"map_id"`
	HP          int                  `json:"destroy"`
	StartX      int                  `json:"-"`
	StartY      int                  `json:"-"`
	StartZ      float64              `json:"-"`
	StartRadian float64              `json:"start_radian"`
	Damage      int                  `json:"-"`
	MaxAngle    float64              `json:"-"`
	MinAngle    float64              `json:"-"`
	//
	RealX            float64 `json:"-"`
	RealY            float64 `json:"-"`
	DistanceTraveled float64 `json:"-"`
	RealSpeed        float64 `json:"-"`
	RadRotate        float64 `json:"-"`

	AngularVelocity float64 `json:"angular_velocity"`
	XVelocity       float64 `json:"x_velocity"`
	YVelocity       float64 `json:"y_velocity"`

	CacheJson      string `json:"-"`
	CreateJsonTime int64  `json:"-"`
	end            bool
	mx             sync.RWMutex
}

func (b *Bullet) AddVelocity(x float64, y float64) {
	b.XVelocity += x
	b.YVelocity += y
}

func (b *Bullet) GetVelocityRotate() float64 {
	return math.Atan2(b.YVelocity, b.XVelocity)
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

func (b *Bullet) GetID() int {
	return b.ID
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

func (b *Bullet) GetJSON(mapTime int64) string {

	if b.CreateJsonTime == mapTime && b.CacheJson != "" {
		return b.CacheJson
	}

	jsonBullet, err := json.Marshal(struct {
		TypeID int `json:"type_id"`
		ID     int `json:"id"`
		X      int `json:"x"`
		Y      int `json:"y"`
		Z      int `json:"z"`
		Rotate int `json:"r"`
		MapID  int `json:"m"`
	}{
		TypeID: b.Ammo.ID,
		ID:     b.ID,
		X:      b.GetX(),
		Y:      b.GetY(),
		Z:      int(b.GetZ()),
		Rotate: int(b.GetRotate()),
		MapID:  b.MapID,
	})

	if err != nil {
		println("bullet to json: ", err.Error())
	}

	b.CacheJson = string(jsonBullet)
	b.CreateJsonTime = mapTime

	return b.CacheJson
}
