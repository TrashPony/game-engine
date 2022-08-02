package obstacle_point

import (
	"encoding/json"
	"sync"
	"sync/atomic"
)

type ObstaclePoint struct {
	X          int32   `json:"x"`
	Y          int32   `json:"y"`
	Radius     int32   `json:"radius"`
	Move       bool    `json:"move"`     // если тру то это только для пуль
	Resource   bool    `json:"resource"` // не влияет на колизии с миром, но влияет на то что низя строить, или куда навести мыш
	ParentID   int32   `json:"-"`
	ParentType string  `json:"-"`
	Key        string  `json:"-"`
	Height     float64 `json:"height"`
	mx         sync.RWMutex
}

func (o *ObstaclePoint) GetX() int {
	return int(atomic.LoadInt32(&o.X))
}

func (o *ObstaclePoint) SetX(x int) {
	atomic.StoreInt32(&o.X, int32(x))
}

func (o *ObstaclePoint) GetY() int {
	return int(atomic.LoadInt32(&o.Y))
}

func (o *ObstaclePoint) SetY(y int) {
	atomic.StoreInt32(&o.Y, int32(y))
}

func (o *ObstaclePoint) GetRadius() int {
	return int(atomic.LoadInt32(&o.Radius))
}

func (o *ObstaclePoint) SetRadius(radius int) {
	atomic.StoreInt32(&o.Radius, int32(radius))
}

func (o *ObstaclePoint) GetMove() bool {
	return o.Move
}

func (o *ObstaclePoint) SetMove(move bool) {
	o.Move = move
}

func (o *ObstaclePoint) GetResource() bool {
	return o.Resource
}

func (o *ObstaclePoint) SetResource(resource bool) {
	o.Resource = resource
}

func (o *ObstaclePoint) GetParentID() int {
	return int(atomic.LoadInt32(&o.ParentID))
}

func (o *ObstaclePoint) SetParentID(parentID int) {
	atomic.StoreInt32(&o.ParentID, int32(parentID))

}

func (o *ObstaclePoint) GetParentType() string {
	return o.ParentType
}

func (o *ObstaclePoint) SetParentType(parentType string) {
	o.ParentType = parentType
}

func (o *ObstaclePoint) GetKey() string {
	return o.Key
}

func (o *ObstaclePoint) SetKey(key string) {
	o.Key = key
}

func (o *ObstaclePoint) GetHeight() float64 {
	return o.Height
}

func (o *ObstaclePoint) SetHeight(height float64) {
	o.Height = height
}

func (o *ObstaclePoint) GetJSON() string {
	jsonPoint, err := json.Marshal(struct {
		X        int32   `json:"x"`
		Y        int32   `json:"y"`
		Radius   int32   `json:"radius"`
		Move     bool    `json:"move"`
		Resource bool    `json:"resource"`
		Height   float64 `json:"height"`
	}{
		X:        o.X,
		Y:        o.Y,
		Radius:   o.Radius,
		Move:     o.Move,
		Resource: o.Resource,
		Height:   o.Height,
	})
	if err != nil {
		println("geo point to json: ", err.Error())
	}

	return string(jsonPoint)
}
