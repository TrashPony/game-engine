package physical_model

import (
	"time"
)

type WASD struct {
	w      bool
	a      bool
	s      bool
	d      bool
	z      bool
	sp     bool
	st     bool
	update int64
}

func (wasd *WASD) SetAllFalse() {
	wasd.Set(false, false, false, false, false, false, false)
}

func (wasd *WASD) Set(w, a, s, d, sp, st, z bool) {

	wasd.w = w
	wasd.a = a
	wasd.s = s
	wasd.d = d
	wasd.sp = sp
	wasd.st = st
	wasd.z = z

	wasd.update = time.Now().UnixNano()
}

func (wasd *WASD) GetW() bool {
	return wasd.w
}

func (wasd *WASD) GetA() bool {
	return wasd.a
}

func (wasd *WASD) GetS() bool {
	return wasd.s
}

func (wasd *WASD) GetD() bool {
	return wasd.d
}

func (wasd *WASD) GetZ() bool {
	return wasd.z
}

func (wasd *WASD) GetSpace() bool {
	return wasd.sp
}

func (wasd *WASD) GetShift() bool {
	return wasd.st
}

func (wasd *WASD) GetUpdate() int64 {
	return wasd.update
}
