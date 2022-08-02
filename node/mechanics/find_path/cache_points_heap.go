package find_path

import "sync"

type pointsHeapType struct {
	points []*Point
	mx     sync.Mutex
}

var pointHead = pointsHeapType{}

func (h *pointsHeapType) Push(x *Point) {
	h.mx.Lock()
	defer h.mx.Unlock()
	x.parent = nil
	h.points = append(h.points, x)
}

func (h *pointsHeapType) Pop() *Point {
	h.mx.Lock()
	defer h.mx.Unlock()

	old := h.points

	if len(old) == 0 {
		return &Point{}
	}

	n := len(old)
	x := old[n-1]
	h.points = old[0 : n-1]
	return x
}

type minFArrayHeapType struct {
	points [][]*Point
	mx     sync.Mutex
}

var minFArrayHeap = minFArrayHeapType{}

func (h *minFArrayHeapType) Push(x []*Point) {
	h.mx.Lock()
	defer h.mx.Unlock()

	x = x[:0]
	h.points = append(h.points, x)
}

func (h *minFArrayHeapType) Pop() []*Point {
	h.mx.Lock()
	defer h.mx.Unlock()

	old := h.points

	if len(old) == 0 {
		return make([]*Point, 0, 128)
	}

	n := len(old)
	x := old[n-1]
	h.points = old[0 : n-1]

	return x
}
