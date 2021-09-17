package game_engine

import (
	"fmt"
	"github.com/TrashPony/game_engine/src/mechanics/create_battle"
	"github.com/TrashPony/game_engine/src/mechanics/debug"
	"github.com/TrashPony/game_engine/src/mechanics/factories/quick_battles"
	"github.com/TrashPony/game_engine/src/mechanics/factories/units"
	"github.com/TrashPony/game_engine/src/mechanics/find_path"
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/coordinate"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/target"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
	"math"
	"math/rand"
	"testing"
)

func BenchmarkFindPathTime(t *testing.B) {

	debug.Store.Move = true

	battle := create_battle.CreateBattle(nil, 1)
	quick_battles.Battles.AddNewGame(battle)

	gameUnit := &unit.Unit{ID: units.Units.GetBotID(), OwnerID: 0, MapID: battle.Map.Id, HP: 100}
	gameUnit.GetPhysicalModel().SetPos(100, 100, float64(rand.Intn(360)))
	gameUnit.MovePath = &unit.MovePath{
		NeedFindPath: true,
		Path:         &[]*coordinate.Coordinate{{X: battle.Map.XSize - 100, Y: battle.Map.YSize - 100}},
		FollowTarget: &target.Target{X: battle.Map.XSize - 100, Y: battle.Map.YSize - 100, Radius: 50},
	}
	units.Units.AddUnit(gameUnit)

	unitsArray := []*unit.Unit{gameUnit}

	for i := 0; i < 20; i++ {
		gameUnit := &unit.Unit{ID: units.Units.GetBotID(), OwnerID: 0, MapID: battle.Map.Id, HP: 100}
		gameUnit.GetPhysicalModel().SetPos(float64(rand.Intn(battle.Map.XSize)), float64(rand.Intn(battle.Map.YSize)), float64(rand.Intn(360)))
		units.Units.AddUnit(gameUnit)
		unitsArray = append(unitsArray, gameUnit)
	}

	t.Run("find_path", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := find_path.LeftHandAlgorithm(
				gameUnit,
				float64(gameUnit.GetPhysicalModel().GetX()), float64(gameUnit.GetPhysicalModel().GetY()),
				float64(gameUnit.MovePath.FollowTarget.GetX()), float64(gameUnit.MovePath.FollowTarget.GetY()),
				unitsArray)

			if err != nil {
				b.Error(err)
			}
		}
	})
}

func BenchmarkFastSinAndCos(t *testing.B) {
	t.Run("sin", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			math.Cos(game_math.DegToRadian(float64(i)))
		}
	})

	t.Run("fas_sin", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			game_math.Sin(game_math.DegToRadian(float64(i)))
		}
	})

	t.Run("cos", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			math.Cos(game_math.DegToRadian(float64(i)))
		}
	})

	t.Run("fas_cos", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			game_math.Cos(game_math.DegToRadian(float64(i)))
		}
	})

	t.Run("math", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fmt.Println(i)
			fmt.Println(game_math.Cos(game_math.DegToRadian(float64(i))))
			fmt.Println(math.Cos(game_math.DegToRadian(float64(i))))
			fmt.Println("-----------------")
		}
	})
}
