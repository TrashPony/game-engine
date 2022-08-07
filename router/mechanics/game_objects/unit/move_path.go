package unit

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/coordinate"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/target"
)

type MovePath struct {
	path         *[]*coordinate.Coordinate
	followTarget *target.Target
	currentPoint int
	needFindPath bool
}

func (u *Unit) GetMovePathState() (*target.Target, *[]*coordinate.Coordinate, int, bool) {
	u.moveMx.Lock()
	defer u.moveMx.Unlock()

	if u.movePath == nil {
		return nil, nil, 0, false
	}

	return u.movePath.followTarget, u.movePath.path, u.movePath.currentPoint, u.movePath.needFindPath
}

func (u *Unit) NextMovePoint() {
	u.moveMx.Lock()
	defer u.moveMx.Unlock()

	if u.movePath == nil {
		return
	}

	u.movePath.currentPoint++
}

func (u *Unit) SetFindMovePath() {
	u.moveMx.Lock()
	defer u.moveMx.Unlock()

	if u.movePath != nil {
		u.movePath.needFindPath = true
	}
}

func (u *Unit) RemoveMovePath() {
	u.moveMx.Lock()
	defer u.moveMx.Unlock()

	u.movePath = nil
}

func (u *Unit) SetMovePath(path *[]*coordinate.Coordinate) {
	u.moveMx.Lock()
	defer u.moveMx.Unlock()

	if u.movePath == nil {
		return
	}

	u.movePath.needFindPath = false
	u.movePath.path = path
	u.movePath.currentPoint = 0
}

func (u *Unit) SetMovePathTarget(t *target.Target) {
	u.moveMx.Lock()
	defer u.moveMx.Unlock()

	u.movePath = &MovePath{
		needFindPath: true,
		path:         &[]*coordinate.Coordinate{{X: t.X, Y: t.Y}},
		followTarget: t,
	}
}

func (u *Unit) GetFollowTarget() *target.Target {
	if u.movePath != nil {
		return u.movePath.followTarget
	}

	return nil
}
