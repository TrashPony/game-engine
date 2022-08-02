package game_math

import (
	"math"
)

func VectorSub(p1, p2 *Point) *Point {
	return &Point{
		X: p2.X - p1.X,
		Y: p2.Y - p1.Y,
	}
}

func IntersectVectorToCircle(a, b, centerCircle *Point, radius int) (intersect bool, x1, y1, x2, y2 float64) {
	// https://stackoverflow.com/questions/1073336/circle-line-segment-collision-detection-algorithm
	// вычисляем расстояние между A и B
	var LAB = math.Sqrt(fastPow(b.X-a.X) + fastPow(b.Y-a.Y))

	// вычислить вектор направления D от A до B
	var Dx = (b.X - a.X) / LAB
	var Dy = (b.Y - a.Y) / LAB

	// compute the value t of the closest point to the circle center (Cx, Cy)
	var t = Dx*(centerCircle.X-a.X) + Dy*(centerCircle.Y-a.Y)

	// This is the projection of C on the line from A to B.

	// вычислить координаты точки E на прямой и ближайшей к C
	xE, yE := (t*Dx)+a.X, (t*Dy)+a.Y

	// высчитывает растояние от E до центра круга
	var LEC = math.Sqrt(fastPow(xE-centerCircle.X) + fastPow(yE-centerCircle.Y))

	// проверяем что бы проекционная точка была ближе радиуса
	if int(LEC) < radius {
		// compute distance from t to circle intersection point
		var dt = math.Sqrt(fastPow(float64(radius)) - fastPow(LEC))
		// ищем первую точку пересечения
		x1, y1 := (t-dt)*Dx+a.X, (t-dt)*Dy+a.Y
		// и вторую
		x2, y2 := (t+dt)*Dx+a.X, (t+dt)*Dy+a.Y

		intersect := PointInVector(a, b, x1, y1) || PointInVector(a, b, x2, y2)
		return intersect, x1, y1, x2, y2
	}

	if int(LEC) == radius { // else test if the line is tangent to circle
		// прямая прилегает к окружности 1 точка пересечени
		return PointInVector(a, b, xE, yE), xE, yE, 0, 0
	}

	return false, 0, 0, 0, 0
}

func fastPow(a float64) float64 {
	return a * a
}
