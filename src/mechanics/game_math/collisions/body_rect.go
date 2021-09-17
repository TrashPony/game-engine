package collisions

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/physical_model"
)

func GetBodyRect(pm *physical_model.PhysicalModel, x, y, rotate float64, full, min bool) *game_math.Polygon {

	/*
		squad.rectDebag.moveTo(-50, -25);
		squad.rectDebag.lineTo(-50, +25);

		squad.rectDebag.lineTo(-50, +25);
		squad.rectDebag.lineTo(+50, +25);

		squad.rectDebag.lineTo(+50, +25);
		squad.rectDebag.lineTo(+50, -25);

		squad.rectDebag.lineTo(+50, -25);
		squad.rectDebag.lineTo(-50, -25);

		// A - [0] B - [1] C = [2] D = [3]
	*/

	lengthBody, widthBody := pm.GetLength(), pm.GetWidth()

	if full {
		if lengthBody > widthBody {
			widthBody = lengthBody
		} else {
			lengthBody = widthBody
		}
	}

	if min {
		if lengthBody < widthBody {
			widthBody = lengthBody
		} else {
			lengthBody = widthBody
		}
	}

	if pm.Polygon == nil {
		pm.Polygon = game_math.GetCenterRect(x, y, lengthBody*2, widthBody*2)
	} else {
		pm.Polygon.UpdateCenterRect(x, y, lengthBody*2, widthBody*2)
	}
	pm.Polygon.Rotate(rotate)

	return pm.Polygon
}
