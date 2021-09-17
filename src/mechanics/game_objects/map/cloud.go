package _map

import (
	"encoding/json"
)

type Cloud struct {
	TypeID   int     `json:"type_id"`
	Speed    float64 `json:"-"`
	Alpha    float64 `json:"alpha"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Angle    int     `json:"angle"`
	ID       int     `json:"id"`
	SizeMapX int     `json:"-"`
	SizeMapY int     `json:"-"`
	IDMap    int     `json:"m"`
}

func (c *Cloud) GetX() float64 {
	return c.X
}

func (c *Cloud) SetX(x float64) {
	c.X = x
}

func (c *Cloud) GetY() float64 {
	return c.Y
}

func (c *Cloud) SetY(y float64) {
	c.Y = y
}

func (c *Cloud) GetAngle() int {
	return c.Angle
}

func (c *Cloud) SetAngle(angle int) {
	c.Angle = angle
}

func (c *Cloud) GetJSON() string {
	jsonCloud, err := json.Marshal(c)
	if err != nil {
		println("clud json", err.Error())
	}

	return string(jsonCloud)
}
