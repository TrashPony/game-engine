package game_math

import "math"

type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (v *Vector) Norm() *Vector {
	/*normalized vector with length of 1 */
	vLen := v.Len()
	return &Vector{X: v.X / vLen, Y: v.Y / vLen}
}

func (v *Vector) Resize(len float64) *Vector {
	/*resized vector to given length */
	return v.Norm().Scale(len)
}

func (v *Vector) Scale(fac float64) *Vector {
	/*scaled vector */
	return &Vector{X: v.X * fac, Y: v.Y * fac}
}

func (v *Vector) Len() float64 {
	/*length of the vector */
	return math.Sqrt(v.LenSqrt())
}

func (v *Vector) LenSqrt() float64 {
	/*non squareroot length of the vector*/
	return v.X*v.X + v.Y*v.Y
}

func (v *Vector) Copy() *Vector {
	/* returns hard copied vector */
	return &Vector{X: v.X, Y: v.Y}
}
func (v *Vector) Sub(vec *Vector) *Vector {
	/* substracts this-vec */
	return &Vector{X: v.X - vec.X, Y: v.Y - vec.Y}
}

func (v *Vector) Add(vec *Vector) *Vector {
	/* sums this + vec */
	return &Vector{X: v.X + vec.X, Y: v.Y + vec.Y}
}

func (v *Vector) VecTo(vec *Vector) *Vector {
	/* gets new vector from this to vec */
	return vec.Sub(v)
}
