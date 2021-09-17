package game_math

// это очень грубая в точности оптимизация
var (
	fact2 = 1 / Factorial(2)
	fact3 = 1 / Factorial(3)
	fact4 = 1 / Factorial(4)
	fact5 = 1 / Factorial(5)
	fact6 = 1 / Factorial(6)
	fact7 = 1 / Factorial(7)
	fact8 = 1 / Factorial(8)
	fact9 = 1 / Factorial(9)
)

func Factorial(n float64) float64 {
	if n == 0 {
		return 1
	}
	return n * Factorial(n-1)
}

// Cos(x) = 1 - x^2 / 2! + x^4 / 4! - x^6 / 6! To infinity...
func Cos(radian float64) float64 { // +
	//return math.Cos(radian)
	a := 1.0
	if radian > 3.14 {
		a = -1
		radian -= 3.14
	}

	if radian < -3.14 {
		a = -1
		radian += 3.14
	}

	return a * (1 - (radian*radian)*fact2 + (radian*radian*radian*radian)*fact4 - (radian*radian*radian*radian*radian*radian)*fact6 + (radian*radian*radian*radian*radian*radian*radian*radian)*fact8)
}

// Sin(x) = x - x^3 / 3! + x^5 / 5! - x^7/7! To infinity...
func Sin(radian float64) float64 {
	//return math.Sin(radian)
	a := 1.0
	if radian > 3.14 {
		a = -1
		radian -= 3.14
	}

	if radian < -3.14 {
		a = -1
		radian += 3.14
	}

	return a * (radian - (radian*radian*radian)*fact3 + (radian*radian*radian*radian*radian)*fact5 - (radian*radian*radian*radian*radian*radian*radian)*fact7 + (radian*radian*radian*radian*radian*radian*radian*radian*radian)*fact9)
}
