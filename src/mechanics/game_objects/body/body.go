package body

type Body struct {
	Name    string              `json:"name"`
	Texture string              `json:"texture"`
	MaxHP   int                 `json:"max_hp"`
	Scale   int                 `json:"scale"`
	Length  int                 `json:"length"`
	Width   int                 `json:"width"`
	Height  int                 `json:"height"`
	Radius  int                 `json:"radius"`
	Weapons map[int]*WeaponSlot `json:"-"`
}
