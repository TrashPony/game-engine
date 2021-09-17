package game_loop_move

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/dynamic_map_object"
	_map "github.com/TrashPony/game_engine/src/mechanics/game_objects/map"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/unit"
	"github.com/TrashPony/game_engine/src/web_socket"
	"sync"
	"time"
)

var minInputTime = int64(1000) * 1000000

func Unit(mp *_map.Map, units []*unit.Unit, mapObjects map[int]*dynamic_map_object.Object) *web_socket.GameLoopMessages {

	msgs := &web_socket.GameLoopMessages{}

	wg := sync.WaitGroup{}
	for _, mUnit := range units {
		moveUnit(mUnit, mp, msgs, units, mapObjects)
	}

	for _, mUnit := range units {
		// паралельные вычесления для поиска пути, т.к. поиск пути на 1 юнита занмается 1-2 мс если их будет много то они не уложатся в тик
		// это безопасно т.к. при вычеслениях мы делаем копию физ. модели и не происходит дата рейс при проврки колизий с другими ботами которые ищуют путь
		// TODO а лучше это вообще вынести поиск пути отдкльно и считать его раз в секунду например
		wg.Add(1)
		go moveGlobalUnit(mUnit, units, &wg)
	}

	wg.Wait()

	return msgs
}

func moveUnit(mUnit *unit.Unit, mp *_map.Map, msgs *web_socket.GameLoopMessages, units []*unit.Unit, mapObjects map[int]*dynamic_map_object.Object) {

	if mUnit == nil {
		return
	}

	// клиент отвалился останавливаем движение
	if (time.Now().UnixNano() - mUnit.GetPhysicalModel().WASD.GetUpdate()) > minInputTime {
		mUnit.GetPhysicalModel().WASD.SetAllFalse()
	}

	xVelocity, yVelocity := mUnit.GetPhysicalModel().GetVelocity()

	if mUnit.GetPhysicalModel().GetPowerMove() == 0 && mUnit.GetPhysicalModel().GetReverse() == 0 &&
		(mUnit.GetPhysicalModel().GetAngularVelocity() < 0.001 && mUnit.GetPhysicalModel().GetAngularVelocity() > -0.001) &&
		(xVelocity < 0.1 && xVelocity > -0.1) && (yVelocity < 0.1 && yVelocity > -0.1) {

		mUnit.GetPhysicalModel().SetAngularVelocity(0)
		mUnit.GetPhysicalModel().SetVelocity(0, 0)
	}

	WasdMove(mUnit.GetPhysicalModel(), mUnit.GetGunner(), "unit", mUnit.ID, mp, msgs, units, mapObjects)
}

func SetUnitsPos(units []*unit.Unit) {
	for _, gameUnit := range units {

		posFunc := gameUnit.GetPhysicalModel().GetPosFunc()
		if posFunc != nil {
			posFunc()
			gameUnit.GetPhysicalModel().SetPosFunc(nil)
		}
	}
}
