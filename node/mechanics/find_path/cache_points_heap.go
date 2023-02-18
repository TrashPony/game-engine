package find_path

import "sync"

type pointsHeapType struct {
	points sync.Pool
}

var pointHead = pointsHeapType{
	points: sync.Pool{
		New: func() interface{} { return &Point{} },
	},
}

func (h *pointsHeapType) Push(x *Point) {
	x.parent = nil
	h.points.Put(x)
}

func (h *pointsHeapType) Pop() *Point {
	return h.points.Get().(*Point)
}

type minFArrayHeapType struct {
	points sync.Pool
}

var minFArrayHeap = minFArrayHeapType{
	points: sync.Pool{
		New: func() interface{} { return []*Point{} },
	},
}

func (h *minFArrayHeapType) Push(x []*Point) {
	x = x[:0]
	h.points.Put(x)
}

func (h *minFArrayHeapType) Pop() []*Point {
	return h.points.Get().([]*Point)
}
