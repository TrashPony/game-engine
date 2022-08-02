package game_math

import (
	"math"
)

type collider interface {
	GetX() int
	GetY() int
	GetVelocity() (float64, float64)
	GetVelocityRotate() float64
	GetCurrentSpeed() float64
	GetDirection() bool
	SetPowerMove(float64)
	SetVelocity(float64, float64)
	AddVelocity(float64, float64)
	GetPowerMove() float64
	SetNextPos(x, y float64)
	GetNextPos() (float64, float64)
	GetRealPos() (float64, float64)
	GetRotate() float64
	GetAngularVelocity() float64
	SetAngularVelocity(float64)
	GetType() string
}

func CollisionReactionBallBall(collider1, collider2 collider, weight1, weight2, pf1, pf2, x2, y2 float64) {
	// https://www-plasmaphysics-org-uk.translate.goog/programs/coll2d_cpp.htm?_x_tr_sl=en&_x_tr_tl=ru&_x_tr_hl=ru&_x_tr_pto=wapp
	// https://gamedev.stackexchange.com/questions/5906/collision-resolution
	cX1, cY1 := collider1.GetNextPos()
	cX2, cY2 := collider2.GetNextPos()

	x1, y1 := collider1.GetRealPos()
	//x2, y2 := collider2.GetRealPos()

	m1, m2 := weight1, weight2

	startXV1, startYV1 := collider1.GetVelocity()
	startXV2, startYV2 := collider2.GetVelocity()

	m21 := m2 / m1
	x21 := cX2 - cX1
	y21 := cY2 - cY1
	vx21 := startXV2 - startXV1
	vy21 := startYV2 - startYV1

	vx_cm := (m1*startXV1 + m2*startXV2) / (m1 + m2)
	vy_cm := (m1*startYV1 + m2*startYV2) / (m1 + m2)

	direct := (vx21*x21 + vy21*y21)
	if direct >= 0 || (collider1.GetType() == "object" && collider2.GetType() == "object") {
		if weight2 >= 20000 || (collider1.GetType() == "object" && collider2.GetType() == "object") {
			fixAdhesion(collider1, collider2, startXV1, startXV2, startYV1, startYV2, m1, m2, x1, y1, x2, y2)
			return
		}

		return
	}

	var sign float64

	fy21 := 1.0e-12 * math.Abs(y21)
	if math.Abs(x21) < fy21 {
		if x21 < 0 {
			sign = -1
		} else {
			sign = 1
		}
		x21 = fy21 * sign
	}

	a := y21 / x21
	dvx2 := -2 * (vx21 + a*vy21) / ((1 + a*a) * (1 + m21))
	vx2 := startXV2 + dvx2
	vy2 := startYV2 + a*dvx2
	vx1 := startXV1 - m21*dvx2
	vy1 := startYV1 - a*m21*dvx2

	R := 1.0

	vx1 = (vx1-vx_cm)*R + vx_cm
	vy1 = (vy1-vy_cm)*R + vy_cm
	vx2 = (vx2-vx_cm)*R + vx_cm
	vy2 = (vy2-vy_cm)*R + vy_cm

	collider1.SetVelocity(vx1, vy1)
	collider2.SetVelocity(vx2, vy2)

	collider1.SetNextPos(x1+vx1, y1+vy1)
	collider2.SetNextPos(x2+vx2, y2+vy2)

	collider1.SetPowerMove(collider1.GetPowerMove() / 1.5)

	collisionRotate(collider1, collider2, startXV1, startYV1, startXV2, startYV2, x1, x2, y1, y2, vx1, vy1)

	return
}

func collisionRotate(collider1, collider2 collider, startXV1, startYV1, startXV2, startYV2, x1, x2, y1, y2, vx1, vy1 float64) {

	if collider1.GetType() == "unit" && collider2.GetType() == "object" {
		return
	}

	v1 := Vector{X: startXV2, Y: startYV2}
	v2 := Vector{X: x2 - x1, Y: y2 - y1}
	v3 := v1.Norm().Sub(v2.Norm())

	angle := GetBetweenAngle(x1+v3.X, v3.Y+y1, x1, y1) - GetBetweenAngle(x2, y2, x1, y1)
	PrepareAngle(&angle)

	v4 := Vector{X: startXV1, Y: startYV1}
	v5 := Vector{X: vx1, Y: vy1}
	l := v4.Sub(&v5)

	len1 := 2 - v3.Len()
	v := ((l.Len() / 100) + (v3.Len() / 100)) / 2
	if math.IsNaN(v) || len1 < 0.02 {
		return
	}

	if angle > 0 && angle < 180 {
		collider1.SetAngularVelocity(collider1.GetAngularVelocity() + v)
	} else {
		collider1.SetAngularVelocity(collider1.GetAngularVelocity() + v*-1)
	}
}

// иногда модели слипают и не могут разлипнуть, тогда мы их просто выталкиваем друг из друга
func fixAdhesion(collider1, collider2 collider, XVelocity1, XVelocity2, YVelocity1, YVelocity2, m1, m2, cX1, cY1, cX2, cY2 float64) (float64, float64) {

	speed1 := math.Sqrt((XVelocity1 * XVelocity1) + (YVelocity1 * YVelocity1))
	speed2 := math.Sqrt((XVelocity2 * XVelocity2) + (YVelocity2 * YVelocity2))

	minSpeed := 3.0
	if speed1 < minSpeed {
		speed1 = minSpeed
	}

	if speed2 < minSpeed {
		speed2 = minSpeed
	}

	if m1 > m2 {
		speed1 = (speed1) * m2 / m1
	} else {
		speed2 = (speed2) * m1 / m2
	}

	angle1 := DegToRadian(GetBetweenAngle(cX1, cY1, cX2, cY2))
	angle2 := DegToRadian(GetBetweenAngle(cX2, cY2, cX1, cY1))

	vx1, vy1 := speed1*Cos(angle1), speed1*Sin(angle1)
	vx2, vy2 := speed2*Cos(angle2), speed2*Sin(angle2)

	collider1.SetPowerMove(collider1.GetPowerMove() / 1.5)
	collider2.SetPowerMove(collider1.GetPowerMove() / 1.5)

	collider1.SetVelocity(vx1, vy1)
	collider2.SetVelocity(vx2, vy2)

	return vx1, vy1
}
