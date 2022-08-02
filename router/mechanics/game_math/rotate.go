package game_math

import "math"

func Rotate(unitRotate, needRotate *float64, step float64) float64 {

	PrepareAngle(unitRotate)
	PrepareAngle(needRotate)

	countRotateAngle := 0.0

	for i := 0.0; i < step*10; i++ {

		if math.Round(*unitRotate*10) != math.Round(*needRotate*10) {

			if directionRotate(*unitRotate, *needRotate) {
				*unitRotate += 0.1
				if *unitRotate >= 360 {
					*unitRotate -= 360
				}

				countRotateAngle += 0.1

			} else {
				*unitRotate -= 0.1
				if *unitRotate < 0 {
					*unitRotate += 360
				}

				countRotateAngle += 0.1

			}

		} else {
			return countRotateAngle
		}
	}

	return countRotateAngle
}

func directionRotate(unitAngle, needAngle float64) bool {

	PrepareAngle(&unitAngle)
	PrepareAngle(&needAngle)

	// true ++
	// false --
	count := 0
	direction := false

	if unitAngle < needAngle {
		for unitAngle < needAngle {
			count++
			direction = true
			unitAngle++
		}
	} else {
		for unitAngle > needAngle {
			count++
			direction = false
			needAngle++
		}
	}

	if direction {
		return count <= 180
	} else {
		return !(count <= 180)
	}
}

func PrepareAngle(angle *float64) {
	if *angle < 0 {
		*angle += 360
	}

	if *angle >= 360 {
		*angle -= 360
	}

	if *angle < 0 {
		*angle += 360
	}

	if *angle >= 360 {
		*angle -= 360
	}
}
