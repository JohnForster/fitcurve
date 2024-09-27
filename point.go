package fitcurves

type Point struct {
	X float64
	Y float64
}

func NewPoint(x, y float64) Point {
	return Point{x, y}
}

func (p Point) Subtract(p1 Point) Vec2 {
	return Vec2{
		X: p.X - p1.X,
		Y: p.Y - p1.Y,
	}
}

func (p Point) Translate(v Vec2) Point {
	return Point{
		X: p.X + v.X,
		Y: p.Y + v.Y,
	}
}

func (p Point) Mult(p1 Point) float64 {
	return p.X*p1.X + p.Y*p1.Y
}

func (p Point) Vec() Vec2 {
	return Vec2{
		X: p.X,
		Y: p.Y,
	}
}
