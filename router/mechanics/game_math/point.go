package game_math

type Point struct {
	X float64
	Y float64
}

func PointInVector(a, b *Point, x, y float64) bool {

	dx1 := b.X - a.X
	dy1 := b.Y - a.Y

	dx := x - a.X
	dy := y - a.Y

	//вычеслям площадь треуголника, если точка принадлежит вектору то площать будет 0
	s := dx1*dy - dx*dy1

	// однако если точка находится в векторе но оч далеко то будет ошибка
	// поэтому надо проверить что бы точка была точно внутри вектора с помощью измерения растоиня
	vectorDist := GetBetweenDist(int(a.X), int(a.Y), int(b.X), int(b.Y))
	aDist := GetBetweenDist(int(a.X), int(a.Y), int(x), int(y))
	bDist := GetBetweenDist(int(b.X), int(b.Y), int(x), int(y))

	return s <= 0.1 && s >= -0.1 && (vectorDist >= aDist && vectorDist >= bDist)
}
