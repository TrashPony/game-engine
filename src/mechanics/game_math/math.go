package game_math

import (
	"math"
	"math/rand"
)

func GetBetweenDist(fromX, fromY, toX, toY int) float64 {
	var dx = float64(toX - fromX)
	var dy = float64(toY - fromY)
	return math.Sqrt(dx*dx + dy*dy)
}

func GetBetweenDistFloat(fromX, fromY, toX, toY float64) float64 {
	var dx = toX - fromX
	var dy = toY - fromY
	return math.Sqrt(dx*dx + dy*dy)
}

func DegToRadian(angle float64) float64 {
	return angle * math.Pi / 180
}

func RadianToDeg(radian float64) float64 {
	return radian * (180 / math.Pi)
}

func GetBetweenAngle(x, y, targetX, targetY float64) float64 {
	//math.Atan2 куда у - текущие у, куда х - текущие х, получаем угол
	needRad := math.Atan2(y-targetY, x-targetX)
	return RadianToDeg(needRad)
}

func RotatePoint(x, y, x0, y0, rotate float64) (newX, newY float64) {
	// поворачиваем квадрат по формуле (x0:y0 - центр)
	//X = (x — x0) * cos(alpha) — (y — y0) * sin(alpha) + x0;
	//Y = (x — x0) * sin(alpha) + (y — y0) * cos(alpha) + y0;

	alpha := rotate * math.Pi / 180
	newX = (x-x0)*Cos(alpha) - (y-y0)*Sin(alpha) + x0
	newY = (x-x0)*Sin(alpha) + (y-y0)*Cos(alpha) + y0
	return
}

func GetBetweenDistLinePoint(xPoint, yPoint, x1Line, y1Line, x2Line, y2Line int) int {

	A := y1Line - y2Line
	B := x1Line - x2Line
	C := y1Line*x2Line - y2Line*x1Line

	dist := int(float64(A*xPoint+B*yPoint+C) / math.Sqrt(float64(A*A+B*B)))
	if dist < 0 {
		dist *= -1
	}
	return dist
}

func VectorToPointBySpeed(x1, y1, x2, y2, speed float64) (int, int) {
	// метод находит точку которая находится между точкой x1y1 и x2y2 на дистанции speed(в пикселях) от точки x1y1
	// возвращает инты потому что пикселы не могут быть дробным числом
	angle := GetBetweenAngle(x2, y2, x1, y1)
	radRotate := angle * math.Pi / 180

	return int(x1 + speed*Cos(radRotate)), int(y1 + speed*Sin(radRotate))
}

func VectorToAngleBySpeed(x1, y1, speed, angle float64) (int, int) {
	// метод находит точку которая находится между точкой отклененной от точки x1y1 в сторону angle на speed пикселей
	radRotate := angle * math.Pi / 180
	return int(x1 + speed*Cos(radRotate)), int(y1 + speed*Sin(radRotate))
}

func GetRangeRand(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func GetRotateVelocity(vx, vy float64) float64 {
	return math.Atan2(vy, vx)
}

func ShortestBetweenAngle(angle1, angle2 float64) int {
	// метод вовзращает разницу углов в формате -180:180
	difference := int(angle2 - angle1)

	if difference == 0 {
		return 0
	}

	var times = math.Floor(float64(difference-(-180)) / 360)

	return difference - int(times*360)
}

func SpeedAndAngleToVelocity(speed float64, radian float64) (float64, float64) {
	return speed * Cos(radian), speed * Sin(radian)
}
