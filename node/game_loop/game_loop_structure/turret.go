package game_loop_structure

import (
	"github.com/TrashPony/game-engine/node/ai/check_target"
	"github.com/TrashPony/game-engine/node/mechanics/actual_target"
	"github.com/TrashPony/game-engine/node/mechanics/attack"
	"github.com/TrashPony/game-engine/node/mechanics/factories/units"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/dynamic_map_object"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/target"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/unit"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/visible_objects"
	"math"
	"math/rand"
	"sync"
)

func TurretTarget(turret *dynamic_map_object.Object, b *battle.Battle) {

	// таймаут в секунду
	if !turret.GetUpdateWeaponTarget() {
		return
	}

	if !turret.CheckVisibleObjectStore() {
		turret.SetVisibleObjectStore(b.Teams[turret.TeamID].GetVisibleObjectStore())
	}

	if (turret.GetHP() <= 0 && !turret.Immortal) || b.Map.Exit {
		return
	}

	if turret.GetWeaponSlot(1).GetAmmo() != nil {
		getVisibleTarget(turret, b)
	}
}

func getVisibleTarget(turret *dynamic_map_object.Object, b *battle.Battle) {

	if turret.ForceTarget != nil {
		if turret.GetWeaponTarget() == nil {
			turret.SetWeaponTarget(&target.Target{Type: "map", X: turret.ForceTarget.X + game_math.GetRangeRand(turret.ForceTarget.Radius*-1, turret.ForceTarget.Radius), Y: turret.ForceTarget.Y + game_math.GetRangeRand(turret.ForceTarget.Radius*-1, turret.ForceTarget.Radius), Attack: turret.ForceTarget.Attack})
		}
		return
	}

	team := b.Teams[turret.TeamID]
	// проверять что цель в пределах дальности
	if !attack.CheckFire(turret.GetGunner(), "object", turret.GetID(), b, turret.GetWeaponTarget(), turret.GetWeaponSlot(1), true, nil) {
		turret.SetWeaponTarget(nil)
	}

	currentTarget := turret.GetWeaponTarget()

	if currentTarget != nil {
		xWeapon, yWeapon := turret.GetGunner().GetWeaponPosInMap(1)
		if !attack.CheckFireToTarget(currentTarget, turret.GetWeaponSlot(1), xWeapon, yWeapon, turret.GetWeaponSlot(1).GetGunRotate(), true) {
			turret.SetWeaponTarget(nil)
			currentTarget = nil
		}
	}

	if currentTarget != nil && currentTarget.Type != "map" {
		// турель потеряла цель
		if team == nil || team.GetVisibleObjectByTypeAndID(currentTarget.Type, currentTarget.ID) == nil {
			turret.SetWeaponTarget(nil)
			currentTarget = nil
		}
	}

	if currentTarget != nil && currentTarget.Type == "unit" && currentTarget.Attack {
		// цель вражеских юнитов считается самой приоритетной
		// поэтому если она уже установлена то не меняем цель
		hostileUnit := units.Units.GetUnitByIDAndMapID(currentTarget.ID, b.Map.Id)

		if hostileUnit != nil {
			hostile, _, _ := check_target.CheckTarget(team, currentTarget, b.Map, turret.TeamID)
			if hostile {
				return
			} else {
				turret.SetWeaponTarget(nil)
				currentTarget = nil
			}
		} else {
			turret.SetWeaponTarget(nil)
			currentTarget = nil
		}
	}

	typeTarget, idTarget := FindHostile(turret, b)
	if typeTarget != "" {
		if typeTarget == "unit" {
			turret.SetWeaponTarget(&target.Target{Type: typeTarget, ID: idTarget, Attack: true})
			actual_target.GetXYZTarget(turret.GetGunner(), turret.GetWeaponTarget(), b.Map, turret.GetWeaponSlot(1))
			return
		} else {
			if currentTarget != nil && currentTarget.Type == "object" && !currentTarget.Attack {
				return
			} else {
				turret.SetWeaponTarget(&target.Target{Type: typeTarget, ID: idTarget, Attack: true})
				return
			}
		}
	}

	if currentTarget == nil {
		// берем цель просто что бы смотреть на нее)
		getFollowTarget(turret, b.Map)
	}
}

func getFollowTarget(turret *dynamic_map_object.Object, mp *_map.Map) {
	var targetUnit *unit.Unit

	findTarget := func(objs []*visible_objects.VisibleObject) {
		// если нет юнитов вражеских то пялимя на обычных
		for _, vObj := range objs {

			if vObj.TypeObject != "unit" {
				continue
			}

			gameUnit := units.Units.GetUnitByIDAndMapID(vObj.IDObject, mp.Id)
			if gameUnit == nil {
				continue
			}

			dist := game_math.GetBetweenDist(gameUnit.GetX(), gameUnit.GetY(), turret.GetX(), turret.GetY())
			if dist < 750 {
				fakeTarget := &target.Target{Type: "unit", ID: gameUnit.GetID()}
				actual_target.GetXYZTarget(turret.GetGunner(), fakeTarget, mp, turret.GetWeaponSlot(1))

				xWeapon, yWeapon := turret.GetGunner().GetWeaponPosInMap(1)
				if attack.CheckFireToTarget(fakeTarget, turret.GetWeaponSlot(1), xWeapon, yWeapon, turret.GetWeaponSlot(1).GetGunRotate(), true) {
					if targetUnit != nil {
						dist2 := game_math.GetBetweenDist(targetUnit.GetX(), targetUnit.GetY(), turret.GetX(), turret.GetY())
						if dist2 > dist {
							targetUnit = gameUnit
						}
					} else {
						targetUnit = gameUnit
					}
				}
			}
		}
	}

	objects, mx := turret.UnsafeRangeVisibleObjects()
	mx.RLock()
	findTarget(objects)
	mx.RUnlock()

	if targetUnit != nil && rand.Intn(5) != 0 {
		turret.SetWeaponTarget(&target.Target{Type: "unit", ID: targetUnit.GetID()})
	} else {
		// загадочно смотрим в даль...

		_, _, lvl := mp.GetPosLevel(turret.GetX(), turret.GetY())
		maxRange, _ := turret.GetGunner().GetWeaponMaxRange(lvl, 1, false)

		radRotate := float64(rand.Intn(360)) * math.Pi / 180 // берем рандомный угол
		x := int(float64(maxRange)*game_math.Cos(radRotate)) + turret.GetX()
		y := int(float64(maxRange)*game_math.Sin(radRotate)) + turret.GetY()
		turret.SetWeaponTarget(&target.Target{Type: "map", X: x, Y: y})
	}
}

// TODO повторяющийся код
func FindHostile(turret *dynamic_map_object.Object, b *battle.Battle) (string, int) {

	if turret == nil {
		return "", 0
	}

	findTarget := func(objs []*visible_objects.VisibleObject, mx *sync.RWMutex) (string, int) {

		if mx == nil {
			return "", 0
		}

		mx.RLock()
		defer mx.RUnlock()

		for _, vObj := range objs {

			if vObj.TypeObject == "unit" {

				hostile, _, _ := check_target.CheckUnit(b.Teams[turret.TeamID], vObj.IDObject, b.Map, turret.TeamID)
				if hostile {
					if attack.CheckFire(turret.GetGunner(), "object", turret.GetID(), b, &target.Target{Type: "unit", ID: vObj.IDObject, Attack: true}, turret.GetWeaponSlot(1), true, nil) {
						return "unit", vObj.IDObject
					} else {
						turret.SetWeaponTarget(nil)
					}
				}
			}

			if vObj.TypeObject == "object" {

				hostile, _, _ := check_target.CheckObj(b.Teams[turret.TeamID], vObj.IDObject, b.Map, turret.TeamID)
				if hostile {

					if attack.CheckFire(turret.GetGunner(), "object", turret.GetID(), b, &target.Target{Type: "object", ID: vObj.IDObject, Attack: true}, turret.GetWeaponSlot(1), true, nil) {
						return "object", vObj.IDObject
					} else {
						turret.SetWeaponTarget(nil)
					}
				}
			}
		}

		return "", 0
	}

	team, ok := b.Teams[turret.TeamID]
	if !ok {
		return "", 0
	}

	return findTarget(team.UnsafeRangeVisibleObjects())
}
