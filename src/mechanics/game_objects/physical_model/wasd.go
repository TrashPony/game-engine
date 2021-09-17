package physical_model

import "time"

type WASD struct {
	w      bool
	a      bool
	s      bool
	d      bool
	update int64
}

func (wasd *WASD) SetAllFalse() {
	wasd.w = false
	wasd.a = false
	wasd.s = false
	wasd.d = false
}

func (wasd *WASD) Set(w, a, s, d bool) {

	wasd.w = w
	wasd.a = a
	wasd.s = s
	wasd.d = d

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

func (wasd *WASD) GetUpdate() int64 {
	return wasd.update
}
