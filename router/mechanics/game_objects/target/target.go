package target

type Target struct {
	Type   string  `json:"type"` // box, unit, map
	UUID   string  `json:"uuid"`
	ID     int     `json:"id"`
	X      int     `json:"x"`
	Y      int     `json:"y"`
	Z      float64 `json:"z"`
	Follow bool    `json:"follow"` // преследовать цель используется для цели атак
	Attack bool    `json:"attack"`
	Force  bool    `json:"force"`
	Radius int     `json:"radius"` // радиус на котором держатся от цели
}

func (t *Target) GetX() int {
	return t.X
}

func (t *Target) SetX(x int) {
	t.X = x
}

func (t *Target) GetY() int {
	return t.Y
}

func (t *Target) SetY(y int) {
	t.Y = y
}

func (t *Target) GetZ() float64 {
	return t.Z
}

func (t *Target) SetZ(z float64) {
	t.Z = z
}

func (t *Target) GetFollow() bool {
	return t.Follow
}

func (t *Target) SetFollow(follow bool) {
	t.Follow = follow
}

func (t *Target) GetCopy() *Target {
	copyTarget := *t
	return &copyTarget
}
