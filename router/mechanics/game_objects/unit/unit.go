package unit

import (
	"github.com/TrashPony/game-engine/node/binary_msg"
	_const "github.com/TrashPony/game-engine/router/const"
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/behavior_rule"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/body"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/burst_of_shots"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/coordinate"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/gunner"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/physical_model"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/target"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/visible_objects"
	"sync"
	"time"
)

type Unit struct {
	ID              int        `json:"id"`
	OwnerID         int        `json:"owner_id"`
	MapID           int        `json:"map_id"`
	HP              int        `json:"hp"`
	Body            *body.Body `json:"body"`
	TeamID          int        `json:"team_id"`
	RespawnTime     int        `json:"-"`
	Respawn         bool       `json:"-"`
	movePath        *MovePath  // специальный обьект для пути
	moveMx          sync.Mutex
	findPathTimeOut int
	BehaviorRules   *behavior_rule.BehaviorRules `json:"-"`
	Role            string                       `json:"-"`
	physicalModel   *physical_model.PhysicalModel
	gunner          *gunner.Gunner
	BurstOfShots    *burst_of_shots.BurstOfShots `json:"-"`
	weaponTarget    *target.Target
	OldStateMsg     map[string][]byte `json:"-"` // TODO небольшой костыль что бы не обновлять постоянно стейт на фронте
	LastDamage      LastDamage        `json:"-"` // последний урон конкретно от другово игрока
	LastDamageTime  int64             `json:"-"` // время последнего урона неважно от кого
	LastFireTime    int64             `json:"-"` // время последнего выстрела, включая активные модули
	CacheJson       []byte            `json:"-"`
	CreateJsonTime  int64             `json:"-"`
	visibleObjects  *visible_objects.VisibleObjectsStore
	mx              sync.RWMutex
}

func (u *Unit) GetUpdateData(mapTime int64) []byte {
	// TODO оптимизация и совмещение с src/game_loop/game_loop_view/update_messages.go
	command := []byte{}
	command = append(command, binary_msg.GetIntBytes(u.HP)...)
	command = append(command, binary_msg.GetIntBytes(u.TeamID)...)
	command = append(command, binary_msg.GetIntBytes(u.GetRangeView())...)
	command = append(command, binary_msg.GetIntBytes(u.GetRangeRadar())...)
	return command
}

func (u *Unit) GetTeamID() int {
	return u.TeamID
}

type LastDamage struct {
	UserID     int   `json:"user_id"`
	TimeDamage int64 `json:"time_damage"`
}

func (u *Unit) SetLastDamage(userId int) {
	if userId == u.OwnerID {
		return
	}

	u.LastDamage.UserID = userId
	u.LastDamage.TimeDamage = time.Now().UnixNano()
}

func (u *Unit) GetLastDamage() int {
	// если урон наносился больше чем 30 сек назад то он не учитывается
	if time.Now().UnixNano()-u.LastDamage.TimeDamage > int64(time.Second*30) {
		return 0
	}

	return u.LastDamage.UserID
}

func (u *Unit) GetOwnerID() int {
	return u.OwnerID
}

func (u *Unit) GetWeight() float64 {
	return u.GetPhysicalModel().GetWeight()
}

func (u *Unit) AddVelocity(x float64, y float64) {
	u.GetPhysicalModel().AddVelocity(x, y)
}

func (u *Unit) UpdatePhysicalModel() {

	if u.physicalModel == nil {
		u.initPhysicalModel() // инициируем по умолчания, и ток в методе UpdatePhysicalModel уже докидываем скилы и тд
	}

	u.physicalModel.ID = u.ID
	u.physicalModel.Speed = u.GetMoveMaxPower() / _const.ServerTickSecPart
	u.physicalModel.ReverseSpeed = u.GetMaxReverse() / _const.ServerTickSecPart
	u.physicalModel.PowerFactor = u.GetPowerFactor() / _const.ServerTickSecPart
	u.physicalModel.ReverseFactor = u.GetReverseFactor() / _const.ServerTickSecPart
	u.physicalModel.TurnSpeed = u.GetTurnSpeed() / _const.ServerTickSecPart
}

func (u *Unit) initPhysicalModel() {
	// todo тестовые параметры
	u.physicalModel = &physical_model.PhysicalModel{
		Speed:         u.GetMoveMaxPower() / _const.ServerTickSecPart,
		ReverseSpeed:  u.GetMaxReverse() / _const.ServerTickSecPart,
		PowerFactor:   u.GetPowerFactor() / _const.ServerTickSecPart,
		ReverseFactor: u.GetReverseFactor() / _const.ServerTickSecPart,
		TurnSpeed:     u.GetTurnSpeed() / _const.ServerTickSecPart,
		WASD:          physical_model.WASD{},
		MoveDrag:      u.Body.MoveDrag,
		AngularDrag:   u.Body.AngularDrag,
		Weight:        u.Body.Weight,
		ChassisType:   u.Body.ChassisType,
	}

	// применяем настройки размера к обьектам колизий
	sizeOffset := float64(u.Body.Scale) / 100
	u.physicalModel.Height = float64(u.Body.Height) * sizeOffset
	u.physicalModel.Length = float64(u.Body.Length) * sizeOffset
	u.physicalModel.Width = float64(u.Body.Width) * sizeOffset
	u.physicalModel.Radius = int(float64(u.Body.Radius) * sizeOffset)
	u.HP = u.Body.MaxHP
	u.Body.CurrentPower = u.Body.MaxPower
}

// ячейки которые отображены на панеле быстрого доступа игрока
type EquipSell struct {
	Number int    `json:"-"`
	Source string `json:"source"` // squadInventory, Constructor, "empty" - убрать все выделения

	// если body то эти параметры важны
	TypeSlot int `json:"type_slot"` // 0,1,2,3
	Slot     int `json:"slot"`

	// эти параметры берутся каждый тик по состоянию снаряжения дл фронтенда
	StartReload  int64 `json:"sr"`
	EndReload    int64 `json:"er"`
	AmmoQuantity int   `json:"aq"`
	On           bool  `json:"on"`
}

func (u *Unit) GetType() string {
	return "unit"
}

func (u *Unit) GetPower() int {
	return u.Body.CurrentPower
}

func (u *Unit) SetPower(power int) {
	u.Body.CurrentPower = power
}

func (u *Unit) GetMapID() int {
	return u.MapID
}

func (u *Unit) SetMapID(mapID int) {
	u.MapID = mapID
}

func (u *Unit) GetX() int {
	return u.GetPhysicalModel().GetX()
}

func (u *Unit) GetY() int {
	return u.GetPhysicalModel().GetY()
}

type MovePath struct {
	path         *[]*coordinate.Coordinate
	followTarget *target.Target
	currentPoint int
	needFindPath bool
}

func (u *Unit) GetID() int {
	return u.ID
}

func (u *Unit) GetBurstOfShots() *burst_of_shots.BurstOfShots {
	if u.BurstOfShots == nil {
		u.BurstOfShots = &burst_of_shots.BurstOfShots{}
	}

	return u.BurstOfShots
}

func (u *Unit) GetGunner() *gunner.Gunner {
	if u.gunner == nil {
		u.gunner = &gunner.Gunner{GunUser: u}
	}

	return u.gunner
}

func (u *Unit) GetPhysicalModel() *physical_model.PhysicalModel {
	if u.physicalModel == nil {
		u.initPhysicalModel()
	}

	return u.physicalModel
}

func (u *Unit) GetCopyPhysicalModel() *physical_model.PhysicalModel {
	pm := *u.physicalModel
	return &pm
}

func (u *Unit) GetBytes(mapTime int64) []byte {

	if u.CreateJsonTime == mapTime && len(u.CacheJson) > 0 {
		return u.CacheJson
	}

	command := []byte{}

	command = append(command, binary_msg.GetIntBytes(u.GetID())...)
	command = append(command, binary_msg.GetIntBytes(u.OwnerID)...)
	command = append(command, binary_msg.GetIntBytes(u.HP)...)
	command = append(command, byte(u.Body.ID))
	command = append(command, byte(u.TeamID))

	// position data
	command = append(command, binary_msg.GetIntBytes(u.GetX())...)
	command = append(command, binary_msg.GetIntBytes(u.GetY())...)
	command = append(command, binary_msg.GetIntBytes(int(u.GetRotate()))...)

	command = append(command, binary_msg.GetIntBytes(u.Body.MaxHP)...)

	command = append(command, byte(len([]byte(u.Body.Texture))))
	command = append(command, []byte(u.Body.Texture)...)

	// weapon data
	for _, unitWeaponSlot := range u.RangeWeaponSlots() {

		command = append(command, byte(unitWeaponSlot.Number))
		command = append(command, binary_msg.GetIntBytes((unitWeaponSlot.GetRealXAttach()))...)
		command = append(command, binary_msg.GetIntBytes((unitWeaponSlot.GetRealYAttach()))...)
		command = append(command, byte(unitWeaponSlot.GetXAnchor()*100))
		command = append(command, byte(unitWeaponSlot.GetYAnchor()*100))
		command = append(command, binary_msg.GetIntBytes(int(unitWeaponSlot.GetGunRotate()))...)

		if unitWeaponSlot.Weapon != nil {
			command = append(command, byte(len([]byte(unitWeaponSlot.Weapon.WeaponTexture))))
			command = append(command, []byte(unitWeaponSlot.Weapon.WeaponTexture)...)
		} else {
			command = append(command, byte(0))
		}
	}

	return command
}

func (u *Unit) GetWeaponSlot(slotNumber int) *body.WeaponSlot {
	return u.Body.Weapons[slotNumber]
}

func (u *Unit) RangeWeaponSlots() map[int]*body.WeaponSlot {
	// мы никогда не пишут в карту слотов оружия поэтому этот метод безопасен (по крайне мере пока)
	return u.Body.Weapons
}

func (u *Unit) GetMapHeight() float64 {
	return u.GetPhysicalModel().GetHeight()
}

func (u *Unit) GetRotate() float64 {
	return u.GetPhysicalModel().GetRotate()
}

func (u *Unit) GetScale() int {
	return u.Body.Scale
}

func (u *Unit) GetWeaponTarget() *target.Target {
	return u.weaponTarget
}

func (u *Unit) SetWeaponTarget(target *target.Target) {
	u.weaponTarget = target
}

func (u *Unit) CheckViewCoordinate(x, y int) (bool, bool) {

	if u.Body == nil || u.HP <= 0 {
		return false, false
	}

	if u.Body.RangeView >= int(game_math.GetBetweenDist(u.GetX(), u.GetY(), x, y)) {
		return true, true
	}

	if u.Body.RangeRadar >= int(game_math.GetBetweenDist(u.GetX(), u.GetY(), x, y)) {
		return false, true
	}

	return false, false
}

func (u *Unit) SetDamage(damage int) int {
	u.HP = u.HP - damage
	return damage
}

func (u *Unit) GetMaxHP() int {
	return u.Body.MaxHP
}

func (u *Unit) GetMaxPower() int {
	return u.Body.MaxPower
}

func (u *Unit) GetMoveMaxPower() float64 {
	return u.Body.Speed
}

func (u *Unit) GetMaxReverse() float64 {
	return u.Body.ReverseSpeed
}

func (u *Unit) GetPowerFactor() float64 {
	return u.Body.PowerFactor
}

func (u *Unit) GetReverseFactor() float64 {
	return u.Body.ReverseFactor
}

func (u *Unit) GetTurnSpeed() float64 {
	return u.Body.TurnSpeed
}

func (u *Unit) GetRangeView() int {
	return u.Body.RangeView
}
func (u *Unit) GetRangeRadar() int {
	return u.Body.RangeRadar
}

func (u *Unit) FindPathTimeOut() bool {
	u.findPathTimeOut -= _const.ServerTick
	if u.findPathTimeOut < 0 {
		u.findPathTimeOut = 1000
		return false
	}

	return true
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

func (u *Unit) SetVisibleObjectStore(v *visible_objects.VisibleObjectsStore) {
	u.visibleObjects = v
}

func (u *Unit) checkVisibleObjectStore() {
	if u.visibleObjects == nil {
		u.visibleObjects = &visible_objects.VisibleObjectsStore{}
	}
}

func (u *Unit) GetVisibleObjects() <-chan *visible_objects.VisibleObject {
	u.checkVisibleObjectStore()
	return u.visibleObjects.GetVisibleObjects()
}

func (u *Unit) UnsafeRangeVisibleObjects() ([]*visible_objects.VisibleObject, *sync.RWMutex) {
	u.checkVisibleObjectStore()
	return u.visibleObjects.UnsafeRangeMapDynamicObjects()
}
