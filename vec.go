package fitcurve

import "math"

type Vec2 struct {
	x float64
	y float64
}

func (v Vec2) Normalize() Vec2 {
	l := v.Length()
	if l == 0 {
		return Vec2{}
	}
	return Vec2{
		x: v.x / l,
		y: v.y / l,
	}
}

func (v Vec2) Length() float64 {
	return math.Hypot(v.x, v.y)
}

func (v Vec2) Mult(n float64) Vec2 {
	return Vec2{
		x: v.x * n,
		y: v.y * n,
	}
}

func (v Vec2) Dot(v1 Vec2) float64 {
	return v.x*v1.x + v.y*v1.y
}

func (v Vec2) Subtract(v1 Vec2) Vec2 {
	return Vec2{
		x: v.x - v1.x,
		y: v.y - v1.y,
	}
}
