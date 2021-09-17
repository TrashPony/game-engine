package gunner

import (
	"github.com/TrashPony/game_engine/src/mechanics/game_math"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/body"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/burst_of_shots"
	"github.com/TrashPony/game_engine/src/mechanics/game_objects/target"
)

type GunUser interface {
	GetWeaponSlot(slotNumber int) *body.WeaponSlot
	RangeWeaponSlots() <-chan *body.WeaponSlot
	GetMapHeight() float64
	GetRotate() float64
	GetX() int
	GetY() int
	GetScale() int
	GetWeaponTarget() *target.Target
	SetWeaponTarget(target *target.Target)
	GetBurstOfShots() *burst_of_shots.BurstOfShots
}

type Gunner struct {
	GunUser GunUser
}

func (g *Gunner) GetX() int {
	return g.GunUser.GetX()
}

func (g *Gunner) GetY() int {
	return g.GunUser.GetY()
}

func (g *Gunner) GetRotate() float64 {
	return g.GunUser.GetRotate()
}

func (g *Gunner) GetWeaponSlot(slotNumber int) *body.WeaponSlot {
	return g.GunUser.GetWeaponSlot(slotNumber)
}

func (g *Gunner) RangeWeaponSlots() <-chan *body.WeaponSlot {
	return g.GunUser.RangeWeaponSlots()
}

func (g *Gunner) GetMapHeight() float64 {
	return g.GunUser.GetMapHeight()
}

func (g *Gunner) GetBurstOfShots() *burst_of_shots.BurstOfShots {
	return g.GunUser.GetBurstOfShots()
}

func (g *Gunner) GetGunRotate(slotNumber int) float64 {
	weaponSlot := g.GunUser.GetWeaponSlot(slotNumber)
	if weaponSlot.GetAmmo() == nil {
		return 0
	}

	return weaponSlot.GetGunRotate()
}

func (g *Gunner) GetFirePos(slotNumber int) *game_math.Positions {
	// отдае координату откуда стрелять
	weaponSlot := g.GunUser.GetWeaponSlot(slotNumber)
	if weaponSlot.GetAmmo() == nil {
		return &game_math.Positions{}
	}

	return g.GetWeaponFirePos(slotNumber)[weaponSlot.GetLastFirePosition()]
}

func (g *Gunner) GetWeaponFirePos(slotNumber int) []*game_math.Positions {

	weaponSlot := g.GunUser.GetWeaponSlot(slotNumber)
	if weaponSlot.GetAmmo() == nil {
		return nil
	}

	return game_math.GetWeaponFirePositions(
		g.GunUser.GetX(), g.GunUser.GetY(), g.GunUser.GetScale(), g.GunUser.GetRotate(), weaponSlot.GetGunRotate(),
		weaponSlot.Weapon.XAttach, weaponSlot.Weapon.YAttach,
		weaponSlot.Weapon.FirePositions,
		float64(weaponSlot.XAttach),
		float64(weaponSlot.YAttach),
	)
}

func (g *Gunner) GetWeaponTarget() *target.Target {
	return g.GunUser.GetWeaponTarget()
}

func (g *Gunner) SetWeaponTarget(target *target.Target) {
	g.GunUser.SetWeaponTarget(target)
}

// возарщает урон простой и урон по снаряжению
func (g *Gunner) GetDamage(slotNumber int) int {

	weaponSlot := g.GunUser.GetWeaponSlot(slotNumber)
	if weaponSlot.GetAmmo() == nil {
		return 0
	}

	return game_math.GetRangeRand(weaponSlot.GetAmmo().MinDamage, weaponSlot.GetAmmo().MaxDamage)
}

func (g *Gunner) GetWeaponMaxRange(lvlMap float64, slotNumber int) (int, float64) {

	weaponSlot := g.GunUser.GetWeaponSlot(slotNumber)
	if weaponSlot == nil || weaponSlot.GetAmmo() == nil {
		return 0, 0
	}

	if weaponSlot.Weapon.Type == "laser" || weaponSlot.Weapon.Type == "missile" {
		return weaponSlot.Weapon.MaxRange, float64(weaponSlot.Weapon.MaxAngle)
	}

	bulletSpeed := weaponSlot.Weapon.BulletSpeed + weaponSlot.GetAmmo().BulletSpeed

	angle := 0.0

	if weaponSlot.Weapon.Artillery {
		angle = float64(weaponSlot.Weapon.MinAngle)
	} else {
		angle = float64(weaponSlot.Weapon.MaxAngle)
	}

	maxRange := game_math.GetMaxDistBallisticWeapon(
		float64(bulletSpeed),
		lvlMap,
		lvlMap+g.GunUser.GetMapHeight(),
		game_math.DegToRadian(angle),
		weaponSlot.Weapon.Type,
	)

	return int(maxRange), angle
}

func (g *Gunner) GetWeaponMinRange(lvlMap float64, slotNumber int) (int, float64) {

	weaponSlot := g.GunUser.GetWeaponSlot(slotNumber)
	if weaponSlot == nil || weaponSlot.GetAmmo() == nil {
		return 0, 0
	}

	if weaponSlot.Weapon.Type == "laser" || weaponSlot.Weapon.Type == "missile" {
		return weaponSlot.Weapon.MinRange, float64(weaponSlot.Weapon.MinAngle)
	}

	bulletSpeed := weaponSlot.Weapon.BulletSpeed + weaponSlot.GetAmmo().BulletSpeed

	angle := 0.0
	if !weaponSlot.Weapon.Artillery {
		angle = float64(weaponSlot.Weapon.MinAngle)
	} else {
		angle = float64(weaponSlot.Weapon.MaxAngle)
	}

	minRange := game_math.GetMaxDistBallisticWeapon(
		float64(bulletSpeed),
		lvlMap,
		lvlMap+g.GunUser.GetMapHeight(),
		game_math.DegToRadian(angle),
		weaponSlot.Weapon.Type,
	)

	return int(minRange), angle
}

func (g *Gunner) GetWeaponAccuracy(slotNumber int) int {
	weaponSlot := g.GunUser.GetWeaponSlot(slotNumber)
	if weaponSlot != nil && weaponSlot.Weapon != nil {
		return weaponSlot.Weapon.Accuracy
	}

	return 999
}

func (g *Gunner) GetWeaponPosInMap(slotNumber int) (int, int) {

	weaponSlot := g.GunUser.GetWeaponSlot(slotNumber)
	if weaponSlot == nil {
		return 0, 0
	}

	return game_math.GetWeaponPosInMap(
		g.GunUser.GetX(), g.GunUser.GetY(), g.GunUser.GetScale(),
		float64(weaponSlot.XAttach),
		float64(weaponSlot.YAttach),
		g.GunUser.GetRotate())
}

func (g *Gunner) GetGunRotateSpeed(slotNumber int) int {

	weaponSlot := g.GunUser.GetWeaponSlot(slotNumber)
	if weaponSlot == nil {
		return 0
	}

	return weaponSlot.Weapon.RotateSpeed
}

func (g *Gunner) SetGunRotate(angle float64, slotNumber int) {
	weaponSlot := g.GunUser.GetWeaponSlot(slotNumber)
	if weaponSlot == nil {
		return
	}

	weaponSlot.SetGunRotate(angle)
}
