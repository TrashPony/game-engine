package physical_model

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_math"
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/obstacle_point"
	"math"
)

// PhysicalModel структура которая встраивается в обьекты которые должны двигатся
type PhysicalModel struct {
	RealX           float64                         `json:"real_x"`
	RealY           float64                         `json:"real_y"`
	NextX           float64                         `json:"-"`
	NextY           float64                         `json:"-"`
	ChassisType     string                          `json:"-"`
	X               int                             `json:"x"`
	Y               int                             `json:"y"`
	Z               float64                         `json:"z"`                // высота над землей (для наземных целей всегда 0)
	Rotate          float64                         `json:"rotate"`           // текущий угол поворота
	PowerMove       float64                         `json:"power_move"`       // силя которая тянет вперед, текущая
	Reverse         float64                         `json:"reverse"`          // сила которая тянет назад, текущая
	AngularVelocity float64                         `json:"angular_velocity"` // скорость поворота, текущая
	XVelocity       float64                         `json:"x_velocity"`       // вектор х
	YVelocity       float64                         `json:"y_velocity"`       // вектор у
	NeedZ           float64                         `json:"need_z"`           // высота которую должен набрать транспорт
	Speed           float64                         `json:"speed"`            // -- макс скорость вперед
	ReverseSpeed    float64                         `json:"reverse_speed"`    // -- макс скорость назад
	PowerFactor     float64                         `json:"power_factor"`     // -- сила ускорения вперед
	ReverseFactor   float64                         `json:"reverse_factor"`   // -- сила ускорения назад
	TurnSpeed       float64                         `json:"turn_speed"`       // -- скорость поворота в радианах
	WASD            WASD                            `json:"-"`                // обьект который говорти когда нажата какая клавиша
	MoveDrag        float64                         `json:"-"`                // сопротивление земли при движение (XVelocity * MoveDrag), (YVelocity * MoveDrag)
	AngularDrag     float64                         `json:"-"`                // сопротивление земли при повороте (AngularVelocity * AngularDrag)
	Weight          float64                         `json:"-"`                // вес
	Height          float64                         `json:"height"`           // высота обьекта
	Length          float64                         `json:"length"`           // длинна обьекта
	Width           float64                         `json:"width"`            // ширина обькта
	Radius          int                             `json:"radius"`           // радиус окружности обьекта
	GeoData         []*obstacle_point.ObstaclePoint `json:"-"`
	PosFunc         func()                          `json:"-"` // функция для принятия положения в конце сервертика
	Type            string                          `json:"type"`
	ID              int                             `json:"id"`
	Fly             bool                            `json:"-"`
	polygon         *game_math.Polygon
	newxPoligon     *game_math.Polygon // todo полигон для проверки следующей позиции, что бы не создавать каждый раз заного
}

func (m *PhysicalModel) GetChassisType() string {
	return m.ChassisType
}

func (m *PhysicalModel) IsFly() bool {
	return m.Fly
}

func (m *PhysicalModel) GetID() int {
	return m.ID
}

func (m *PhysicalModel) GetType() string {
	return m.Type
}

func (m *PhysicalModel) SetNextPos(x, y float64) {
	m.NextX = x
	m.NextY = y
}

func (m *PhysicalModel) GetNextPos() (float64, float64) {
	if m.NextX == 0 && m.NextY == 0 {
		m.NextX = float64(m.GetX())
		m.NextY = float64(m.GetY())
	}
	return m.NextX, m.NextY
}

func (m *PhysicalModel) GetPolygon() *game_math.Polygon {

	if m.polygon == nil {
		m.polygon = game_math.GetCenterRect(m.RealX, m.RealY, m.GetLength()*2, m.GetWidth()*2)
	}

	return m.polygon
}

func (m *PhysicalModel) SetPolygon(p *game_math.Polygon) {
	m.polygon = p
}

func (m *PhysicalModel) GetNextPolygon() *game_math.Polygon {

	if m.newxPoligon == nil {
		m.newxPoligon = game_math.GetCenterRect(m.RealX, m.RealY, m.GetLength()*2, m.GetWidth()*2)
	}

	return m.newxPoligon
}

func (m *PhysicalModel) SetNextPolygon(p *game_math.Polygon) {
	m.newxPoligon = p
}

func (m *PhysicalModel) GetX() int {
	return m.X
}

func (m *PhysicalModel) GetY() int {
	return m.Y
}

func (m *PhysicalModel) MultiplyVelocity(x float64, y float64) {
	m.XVelocity *= x
	m.YVelocity *= y
}

func (m *PhysicalModel) AddVelocity(x float64, y float64) {
	m.XVelocity += x
	m.YVelocity += y
}

func (m *PhysicalModel) GetRotate() float64 {
	return m.Rotate
}

func (m *PhysicalModel) GetRealPos() (float64, float64) {
	return m.RealX, m.RealY
}

func (m *PhysicalModel) GetDirection() bool {
	return m.GetPowerMove()-m.GetReverse() > 0
}

func (m *PhysicalModel) GetCurrentSpeed() float64 {
	xVelocity, yVelocity := m.GetVelocity()
	return math.Sqrt(xVelocity*xVelocity + yVelocity*yVelocity)
}

func (m *PhysicalModel) SetPowerMove(powerMove float64) {
	m.PowerMove = powerMove
}

func (m *PhysicalModel) GetHeight() float64 {
	return m.Height
}

func (m *PhysicalModel) SetHeight(height float64) {
	m.Height = height
}

func (m *PhysicalModel) GetLength() float64 {
	return m.Length
}

func (m *PhysicalModel) GetWidth() float64 {
	return m.Width
}

func (m *PhysicalModel) GetRadius() int {
	return m.Radius
}

func (m *PhysicalModel) CheckGrowthPower() bool {
	return m.WASD.GetW()
}

func (m *PhysicalModel) CheckGrowthRevers() bool {
	return m.WASD.GetS()
}

func (m *PhysicalModel) CheckLeftRotate() bool {
	return m.WASD.GetA()
}

func (m *PhysicalModel) CheckRightRotate() bool {
	return m.WASD.GetD()
}

func (m *PhysicalModel) CheckHandBrake() bool {
	return m.WASD.GetSpace()
}

func (m *PhysicalModel) SetReverse(reverse float64) {
	m.Reverse = reverse
}

func (m *PhysicalModel) GetMoveMaxPower() float64 {
	return m.Speed
}

func (m *PhysicalModel) GetMaxReverse() float64 {
	return m.ReverseSpeed
}

func (m *PhysicalModel) GetPowerFactor() float64 {
	return m.PowerFactor
}

func (m *PhysicalModel) GetReverseFactor() float64 {
	return m.ReverseFactor
}

func (m *PhysicalModel) GetTurnSpeed() float64 {
	return m.TurnSpeed
}

func (m *PhysicalModel) GetZ() float64 {
	return m.Z
}

func (m *PhysicalModel) SetZ(z float64) {
	m.Z = z
}

func (m *PhysicalModel) GetNeedZ() float64 {
	return m.NeedZ
}

func (m *PhysicalModel) SetNeedZ(needZ float64) {
	m.NeedZ = needZ
}

func (m *PhysicalModel) GetVelocityRotate() float64 {
	return math.Atan2(m.YVelocity, m.XVelocity)
}

func (m *PhysicalModel) GetVelocity() (float64, float64) {
	return m.XVelocity, m.YVelocity
}

func (m *PhysicalModel) GetPowerMove() float64 {
	return m.PowerMove
}

func (m *PhysicalModel) GetReverse() float64 {
	return m.Reverse
}

func (m *PhysicalModel) GetAngularVelocity() float64 {
	return m.AngularVelocity
}

func (m *PhysicalModel) SetAngularVelocity(angularVelocity float64) {
	m.AngularVelocity = angularVelocity
}

func (m *PhysicalModel) SetVelocity(x float64, y float64) {
	m.XVelocity, m.YVelocity = x, y
}

func (m *PhysicalModel) GetPosFunc() func() {
	return m.PosFunc
}

func (m *PhysicalModel) SetPosFunc(fun func()) {
	m.PosFunc = fun
}

func (m *PhysicalModel) GetWeight() float64 {
	return m.Weight
}

func (m *PhysicalModel) SetWASD(w, a, s, d, sp, st, z bool) {
	m.WASD.Set(w, a, s, d, sp, st, z)
}

func (m *PhysicalModel) GetMoveDrag() float64 {
	return m.MoveDrag
}

func (m *PhysicalModel) GetAngularDrag() float64 {
	return m.AngularDrag
}

func (m *PhysicalModel) GetGeoData() []*obstacle_point.ObstaclePoint {
	return m.GeoData
}

func (m *PhysicalModel) SetPos(realX, realY, angle float64) {

	m.RealX = realX
	m.RealY = realY
	m.X = int(m.RealX)
	m.Y = int(m.RealY)

	if angle > 360 {
		angle -= 360
	}

	if angle < 0 {
		angle += 360
	}

	m.Rotate = angle
	m.PosFunc = nil
}

func (m *PhysicalModel) GetPhysicalModel() *PhysicalModel {
	return m
}
