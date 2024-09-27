package fitcurves

import "math"

type Vec2 struct {
	X float64
	Y float64
}

func (v Vec2) Normalize() Vec2 {
	l := v.Length()
	if l == 0 {
		return Vec2{}
	}
	return Vec2{
		X: v.X / l,
		Y: v.Y / l,
	}
}

func (v Vec2) Length() float64 {
	return math.Hypot(v.X, v.Y)
}

func (v Vec2) Mult(n float64) Vec2 {
	return Vec2{
		X: v.X * n,
		Y: v.Y * n,
	}
}

func (v Vec2) Dot(v1 Vec2) float64 {
	return v.X*v1.X + v.Y*v1.Y
}

func (v Vec2) Subtract(v1 Vec2) Vec2 {
	return Vec2{
		X: v.X - v1.X,
		Y: v.Y - v1.Y,
	}
}
