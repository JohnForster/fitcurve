package fitcurves

type Bezier struct {
	p0 Point
	p1 Point
	p2 Point
	p3 Point
}

// Evaluates cubic bezier at parameter t
func (b Bezier) Q(t float64) Point {

	tx := 1.0 - t

	x := ((b.p0.x * tx * tx * tx) +
		(b.p1.x * 3 * tx * tx * t)) +
		((b.p2.x * 3 * tx * t * t) +
			(b.p3.x * t * t * t))

	y := ((b.p0.y * tx * tx * tx) +
		(b.p1.y * 3 * tx * tx * t)) +
		((b.p2.y * 3 * tx * t * t) +
			(b.p3.y * t * t * t))

	return Point{x, y}
}

// Evaluates cubic bezier first derivative at t
func (b Bezier) QPrime(t float64) Point {
	tx := 1 - t
	d1 := b.p1.Subtract(b.p0)
	d2 := b.p2.Subtract(b.p1)
	d3 := b.p3.Subtract(b.p2)
	x := (d1.x * 3 * tx * tx) +
		(d2.x * 6 * tx * t) +
		(d3.x * 3 * t * t)
	y := (d1.y * 3 * tx * tx) +
		(d2.y * 6 * tx * t) +
		(d3.y * 3 * t * t)
	return Point{x, y}
}

func (b Bezier) QPrimePrime(t float64) Point {
	tx := 1 - t

	x := ((b.p0.x + (b.p2.x - (b.p1.x * 2))) * (6 * tx)) +
		((b.p1.x + (b.p3.x - (b.p2.x * 2))) * (6 * t))

	y := ((b.p0.y + (b.p2.y - (b.p1.y * 2))) * (6 * tx)) +
		((b.p1.y + (b.p3.y - (b.p2.y * 2))) * (6 * t))

	return Point{x, y}
}
