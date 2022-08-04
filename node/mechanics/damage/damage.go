package damage

import (
	"github.com/TrashPony/game-engine/node/mechanics/factories/units"
	battle2 "github.com/TrashPony/game-engine/router/mechanics/game_objects/battle"
	_map "github.com/TrashPony/game-engine/router/mechanics/game_objects/map"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/player"
)

type Object struct {
	TypeTarget string      `json:"type_target"`
	IdTarget   int         `json:"id_target"`
	Damage     int         `json:"damage"`
	TypeDamage string      `json:"type_damage"`
	Dead       bool        `json:"dead"`
	Obj        interface{} `json:"obj"`
	X          int         `json:"x"`
	Y          int         `json:"y"`
	PushPower  int         `json:"push_power"`
}

func CollisionDamage(typeTarget string, idTarget, damage, areaCovers int, mp *_map.Map, x, y int, z float64, startX, startY, pushPower int) []*Object {
	if areaCovers > 0 {
		return Explosion(mp, typeTarget, idTarget, damage, areaCovers, x, y, z, startX, startY, pushPower)
	} else {
		return []*Object{{TypeTarget: typeTarget, IdTarget: idTarget, Damage: damage, PushPower: pushPower}}
	}
}

func Damage(objs []*Object, b *battle2.Battle, bulletOwnerType string, bulletOwnerID, x, y, bulletAmmoID, bulletWeaponID, bulletEquipID int) []*Object {

	if objs == nil || len(objs) == 0 {
		return nil
	}

	var damageDealer *player.Player
	if bulletOwnerType == "unit" {
		u := units.Units.GetUnitByIDAndMapID(bulletOwnerID, b.Map.Id)
		if u != nil {
			damageDealer = b.GetPlayerByID(u.OwnerID)
		}
	}

	if bulletOwnerType == "object" {
		o := b.Map.GetDynamicObjectsByID(bulletOwnerID)
		if o != nil {
			damageDealer = b.GetPlayerByID(o.OwnerID)
		}
	}

	for _, obj := range objs {
		obj.Dead, obj.Obj, obj.X, obj.Y, obj.Damage = damage(b, obj.TypeTarget, obj.IdTarget, obj.Damage, obj.PushPower,
			damageDealer, bulletOwnerType, bulletOwnerID, bulletAmmoID, x, y, bulletWeaponID, bulletEquipID)
	}

	return objs
}

func damage(b *battle2.Battle, typeTarget string, idTarget, damage, pushPower int, damageDealer *player.Player, bulletOwnerType string, bulletOwnerID, bulletAmmoID, bulletX, bulletY, bulletWeaponID, bulletEquipID int) (bool, interface{}, int, int, int) {

	if typeTarget == "object" || typeTarget == "static_object" {
		obj := b.Map.GetDynamicObjectsByID(idTarget)

		if obj != nil {

			pushUnit(obj.GetPhysicalModel(), bulletX, bulletY, pushPower)

			if damageDealer != nil && damageDealer.TeamID == obj.TeamID {
				return false, obj, obj.GetX(), obj.GetY(), 0
			}

			if obj.Immortal {
				// это неуязвимые обьекты но с действиями
				return false, nil, 0, 0, damage
			}

			obj.SetHP(obj.GetHP() - damage)

			if obj.GetHP() < 0 {
				b.Map.RemoveDynamicObject(obj)
				return true, obj, obj.GetX(), obj.GetY(), damage
			}

			return false, obj, obj.GetX(), obj.GetY(), damage
		}
	}

	if typeTarget == "unit" {

		damageUnit := units.Units.GetUnitByIDAndMapID(idTarget, b.Map.Id)
		if damageUnit == nil {
			return false, damageUnit, 0, 0, damage
		}

		pushUnit(damageUnit.GetPhysicalModel(), bulletX, bulletY, pushPower)

		ownerUser := b.GetPlayerByID(damageUnit.OwnerID)
		if damageDealer != nil && ownerUser != nil && damageDealer.TeamID == ownerUser.TeamID && damageDealer.GetID() != ownerUser.GetID() {
			return false, damageUnit, damageUnit.GetX(), damageUnit.GetY(), 0
		}

		if damageUnit.HP > 0 {
			damage = damageUnit.SetDamage(damage)

			if damageUnit.HP <= 0 {
				damageUnit.HP = 0
				units.Units.RemoveUnitByID(damageUnit.GetID(), damageUnit.MapID)
				return true, damageUnit, damageUnit.GetX(), damageUnit.GetY(), damage
			}

			return false, damageUnit, damageUnit.GetX(), damageUnit.GetY(), damage
		}
	}

	return false, nil, 0, 0, damage
}
