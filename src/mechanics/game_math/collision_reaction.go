package game_math

import (
	"math"
	"math/rand"
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
}

func CollisionReactionBallBall(collider1, collider2 collider, noSetVelocityCollider1 bool, weight1, weight2, pf1, pf2 float64) (int, int) {

	theta1 := collider1.GetVelocityRotate()
	theta2 := collider2.GetVelocityRotate()

	cX1 := float64(collider1.GetX())
	cY1 := float64(collider1.GetY())

	cX2 := float64(collider2.GetX())
	cY2 := float64(collider2.GetY())

	phi := math.Atan2(cY2-cY1, cX2-cX1)

	v1, v2 := collider1.GetCurrentSpeed(), collider2.GetCurrentSpeed()
	m1, m2 := weight1/1000, weight2/1000

	minSpeed := 1.0
	if v1 < minSpeed && v2 < minSpeed {
		v1 = minSpeed
		v2 = minSpeed

		theta1 = phi
		theta2 = phi + DegToRadian(180)
	}

	startXV1, startYV1 := collider1.GetVelocity()
	startXV2, startYV2 := collider2.GetVelocity()

	XVelocity1 := (v1*Cos(theta1-phi)*(m1-m2)+2*m2*v2*Cos(theta2-phi))/(m1+m2)*Cos(phi) + v1*Sin(theta1-phi)*Cos(phi+math.Pi/2)
	YVelocity1 := (v1*Cos(theta1-phi)*(m1-m2)+2*m2*v2*Cos(theta2-phi))/(m1+m2)*Sin(phi) + v1*Sin(theta1-phi)*Sin(phi+math.Pi/2)
	XVelocity2 := (v2*Cos(theta2-phi)*(m2-m1)+2*m1*v1*Cos(theta1-phi))/(m1+m2)*Cos(phi) + v2*Sin(theta2-phi)*Cos(phi+math.Pi/2)
	YVelocity2 := (v2*Cos(theta2-phi)*(m2-m1)+2*m1*v1*Cos(theta1-phi))/(m1+m2)*Sin(phi) + v2*Sin(theta2-phi)*Sin(phi+math.Pi/2)

	if !noSetVelocityCollider1 {
		collider1.SetPowerMove(collider1.GetPowerMove() - (pf1 * 1.5))
	}

	collider2.SetPowerMove(collider2.GetPowerMove() - (pf2 * 1.5))

	newX1 := cX1 + XVelocity1
	newY1 := cY1 + YVelocity1

	newX2 := cX2 + XVelocity2
	newY2 := cY2 + YVelocity2

	if GetBetweenDistFloat(newX1, newY1, newX2, newY2) < GetBetweenDistFloat(cX1, cY1, cX2, cY2) {
		fixAdhesion(collider1, collider2, XVelocity1, XVelocity2, YVelocity1, YVelocity2, m1, m2, cX1, cY1, cX2, cY2, noSetVelocityCollider1)
	} else {
		if m2 == 9999999 {
			if !noSetVelocityCollider1 {
				collider1.SetPowerMove(collider1.GetPowerMove() * 0.25)
				collider1.SetVelocity(XVelocity1/2, YVelocity1/2)
			}
		} else {
			if !noSetVelocityCollider1 {
				collider1.SetVelocity(XVelocity1, YVelocity1)
			}
			collider2.SetVelocity(XVelocity2, YVelocity2)
		}
	}

	if cX1 == cX2 && cY1 == cY2 {
		fixAdhesion(collider1, collider2, XVelocity1, XVelocity2, YVelocity1, YVelocity2, m1, m2, cX1, cY1, cX2, cY2, noSetVelocityCollider1)
	}

	return getDamage(startXV1, startYV1, XVelocity1, YVelocity1, startXV2, startYV2, XVelocity2, YVelocity2, m1, m2)
}

// иногда модели слипают и не могут разлипнуть
func fixAdhesion(collider1, collider2 collider, XVelocity1, XVelocity2, YVelocity1, YVelocity2, m1, m2, cX1, cY1, cX2, cY2 float64, noSetVelocityCollider1 bool) {
	var x1, y1, x2, y2 float64

	cX1 -= float64(GetRangeRand(-10, 10))
	cY1 -= float64(GetRangeRand(-10, 10))

	speed1 := math.Sqrt((XVelocity1 * XVelocity1) + (YVelocity1 * YVelocity1))
	speed2 := math.Sqrt((XVelocity2 * XVelocity2) + (YVelocity2 * YVelocity2))

	minSpeed := 2.0
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

	x1, y1 = speed1*Cos(angle1), speed1*Sin(angle1)
	x2, y2 = speed2*Cos(angle2), speed2*Sin(angle2)

	if m2 == 9999999 {
		x1 = x1 / 3
		y1 = y1 / 3
		if !noSetVelocityCollider1 {
			collider1.SetPowerMove(0)
		}
	} else {
		if !noSetVelocityCollider1 {
			collider1.SetPowerMove(collider1.GetPowerMove() / 2)
		}
	}

	// костыль из за того что однвременно 2 могут просчитываться поочереди и быть collider1
	if rand.Intn(2) == 0 {
		collider2.SetVelocity(x1, y1)
		if !noSetVelocityCollider1 {
			collider1.SetVelocity(x2, y2)
		}
	} else {
		if !noSetVelocityCollider1 {
			collider1.SetVelocity(x1, y1)
		}
		collider2.SetVelocity(x2, y2)
	}
}

func getDamage(startXV1, startYV1, XVelocity1, YVelocity1, startXV2, startYV2, XVelocity2, YVelocity2, m1, m2 float64) (int, int) {
	getD := func(startXV, startYV, XVelocity, YVelocity float64) int {
		diffX := startXV - XVelocity
		if XVelocity > startXV {
			diffX = XVelocity - startXV
		}

		diffY := startYV - YVelocity
		if YVelocity > startYV {
			diffY = YVelocity - startYV
		}

		return int(diffX) + int(diffY)
	}

	if m1 == 9999999 {
		return getD(startXV2, startYV2, XVelocity2, YVelocity2), getD(startXV2, startYV2, XVelocity2, YVelocity2)
	}

	if m2 == 9999999 {
		return getD(startXV1, startYV1, XVelocity1, YVelocity1), getD(startXV1, startYV1, XVelocity1, YVelocity1)
	}

	return getD(startXV1, startYV1, XVelocity1, YVelocity1), getD(startXV2, startYV2, XVelocity2, YVelocity2)
}
