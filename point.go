package fitcurve

type Point struct {
	x float64
	y float64
}

func (p Point) Subtract(p1 Point) Vec2 {
	return Vec2{
		x: p.x - p1.x,
		y: p.y - p1.y,
	}
}

func (p Point) Translate(v Vec2) Point {
	return Point{
		x: p.x + v.x,
		y: p.y + v.y,
	}
}

func (p Point) Mult(p1 Point) float64 {
	return p.x*p1.x + p.y*p1.y
}

func (p Point) Vec() Vec2 {
	return Vec2{
		x: p.x,
		y: p.y,
	}
}
