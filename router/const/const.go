package _const

const (
	CellSize     = 32
	DiscreteSize = 4 * CellSize
	Gravity      = 980.0 / 2
	AmmoRadius   = 3

	ServerTick        = 32
	ServerBulletTick  = 32
	ServerTickSecPart = float64(1000) / ServerTick

	// SpriteSize размеры спрайтов юнитов и обьектов, что бы правильно расчитывать расположение оружия и снаряжения на теле/корпусе
	SpriteSize  = 128.00
	SpriteSize2 = SpriteSize / 2

	LaserWeapon    = "laser"
	FirearmsWeapon = "firearms"
	MissileWeapon  = "missile"
)

var MasterInit = false
var MapBinItems = map[string]int{
	"unit":            3,
	"dynamic_objects": 6,
	"object":          9,
	"mark":            10,
	"bullet":          11,
}
