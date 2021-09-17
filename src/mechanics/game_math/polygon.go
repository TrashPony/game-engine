package game_math

import "sync"

type Polygon struct {
	RotateSides      []*SideRec `json:"-"`
	startSides       []*SideRec
	CenterX, CenterY float64
	Angle            float64
	mx               sync.RWMutex
}

func (r *Polygon) GetCenter() (float64, float64) {
	return r.CenterX, r.CenterY
}

func (r *Polygon) GetAngle() float64 {
	return r.Angle
}

func (r *Polygon) GetRotateSides() []*SideRec {
	return r.RotateSides
}

type SideRec struct {
	XY1 *Point `json:"xy_1"`
	XY2 *Point `json:"xy_2"`
}

func (r *Polygon) Rotate(rotate float64) {

	r.Angle = rotate

	rotateSide := func(rotateSide, side *SideRec, x0, y0 float64, rotate float64) {
		rotateSide.XY1.X, rotateSide.XY1.Y = RotatePoint(side.XY1.X, side.XY1.Y, x0, y0, rotate)
	}

	for i, side := range r.startSides {
		rotateSide(r.RotateSides[i], side, r.CenterX, r.CenterY, rotate)
	}
}

func (r *Polygon) DetectCollisionRectToCircle(centerCircle *Point, radius int) bool {

	// A - [0]1 B - [1]1 C = [2]1 D = [3]1

	if r.DetectPointInRectangle(centerCircle.X, centerCircle.Y) {
		// цент находится внутри прямоуголника, пересекается
		return true
	}

	/*
			intersectCircle(S, (A, B)) or
		    intersectCircle(S, (B, C)) or
		    intersectCircle(S, (C, D)) or
		    intersectCircle(S, (D, A))
	*/

	intersect1, _, _, _, _ := IntersectVectorToCircle(r.RotateSides[0].XY1, r.RotateSides[1].XY1, centerCircle, radius)
	if intersect1 {
		return true
	}

	intersect2, _, _, _, _ := IntersectVectorToCircle(r.RotateSides[1].XY1, r.RotateSides[2].XY1, centerCircle, radius)
	if intersect2 {
		return true
	}

	intersect3, _, _, _, _ := IntersectVectorToCircle(r.RotateSides[2].XY1, r.RotateSides[3].XY1, centerCircle, radius)
	if intersect3 {
		return true
	}

	intersect4, _, _, _, _ := IntersectVectorToCircle(r.RotateSides[3].XY1, r.RotateSides[0].XY1, centerCircle, radius)
	if intersect4 {
		return true
	}

	// пересекается 1 из сторон
	return false
}

func (r *Polygon) DetectCollisionRectToRect(r2 *Polygon) bool {

	centerX1, centerY1 := r.GetCenter()
	if r2.DetectPointInRectangle(centerX1, centerY1) {
		// цент находится внутри прямоуголника, пересекается
		return true
	}

	centerX2, centerY2 := r2.GetCenter()
	if r.DetectPointInRectangle(centerX2, centerY2) {
		// цент находится внутри прямоуголника, пересекается
		return true
	}

	if r.CenterX == r2.CenterX && r.CenterY == r2.CenterY {
		// при одинаковом прямоугольнике и одинаковым центром, не будет пересечений и колизия будет не найдена
		// поэтому это тут
		return true
	}

	intersection := func(ax1, ay1, ax2, ay2, bx1, by1, bx2, by2 float64) bool {
		v1 := (bx2-bx1)*(ay1-by1) - (by2-by1)*(ax1-bx1)
		v2 := (bx2-bx1)*(ay2-by1) - (by2-by1)*(ax2-bx1)
		v3 := (ax2-ax1)*(by1-ay1) - (ay2-ay1)*(bx1-ax1)
		v4 := (ax2-ax1)*(by2-ay1) - (ay2-ay1)*(bx2-ax1)

		return (v1*v2 < 0) && (v3*v4 < 0)
	}

	for _, side1 := range r.RotateSides {
		for _, side2 := range r2.RotateSides {

			if intersection(
				side1.XY1.X, side1.XY1.Y, side1.XY2.X, side1.XY2.Y,
				side2.XY1.X, side2.XY1.Y, side2.XY2.X, side2.XY2.Y,
			) {
				return true
			}
		}
	}

	return false
}

func (r *Polygon) DetectPointInRectangle(x, y float64) bool {

	dot := func(u, v *Point) float64 {
		return u.X*v.X + u.Y*v.Y
	}

	// A - [0]1 B - [1]1 C = [2]1 D = [3]1
	//0 ≤ AP·AB ≤ AB·AB and 0 ≤ AP·AD ≤ AD·AD
	// TODO  invalid memory address or nil pointer dereference, у квадрата по какой то причин все 4 стороны были nil r.RotateSides
	AB := Vector(r.RotateSides[0].XY1, r.RotateSides[1].XY1)
	AM := Vector(r.RotateSides[0].XY1, &Point{X: x, Y: y})
	BC := Vector(r.RotateSides[1].XY1, r.RotateSides[2].XY1)
	BM := Vector(r.RotateSides[1].XY1, &Point{X: x, Y: y})

	return 0 <= dot(AB, AM) && dot(AB, AM) <= dot(AB, AB) && 0 <= dot(BC, BM) && dot(BC, BM) <= dot(BC, BC)
}

func GetCenterRect(x, y, height, width float64) *Polygon {

	// делем на 2 что бы центр квадрата был в х у
	height = height / 2
	width = width / 2

	return GetRect(x, y, height, width)
}

func (r *Polygon) UpdateCenterRect(x, y, height, width float64) *Polygon {
	height = height / 2
	width = width / 2

	// не самый оптимальный код, но это экономит выделение памяти

	// A
	r.startSides[0].XY1.X = x - width
	r.startSides[0].XY1.Y = y - height

	// B
	r.startSides[1].XY1.X = x - width
	r.startSides[1].XY1.Y = y + height

	// C
	r.startSides[2].XY1.X = x + width
	r.startSides[2].XY1.Y = y + height

	// D
	r.startSides[3].XY1.X = x + width
	r.startSides[3].XY1.Y = y - height

	r.CenterX = x
	r.CenterY = y

	return r
}

func GetRect(x, y, height, width float64) *Polygon {

	a := &Point{X: x - width, Y: y - height}
	b := &Point{X: x - width, Y: y + height}
	c := &Point{X: x + width, Y: y + height}
	d := &Point{X: x + width, Y: y - height}

	ra := &Point{X: x - width, Y: y - height}
	rb := &Point{X: x - width, Y: y + height}
	rc := &Point{X: x + width, Y: y + height}
	rd := &Point{X: x + width, Y: y - height}

	return &Polygon{
		startSides: []*SideRec{
			{XY1: a, XY2: b},
			{XY1: b, XY2: c},
			{XY1: c, XY2: d},
			{XY1: d, XY2: a},
		},
		RotateSides: []*SideRec{
			{XY1: ra, XY2: rb},
			{XY1: rb, XY2: rc},
			{XY1: rc, XY2: rd},
			{XY1: rd, XY2: ra},
		},
		CenterX: x,
		CenterY: y,
	}
}
