package fitcurves

type Bezier struct {
	P0 Point
	P1 Point
	P2 Point
	P3 Point
}

// Evaluates cubic bezier at parameter t
func (b Bezier) Q(t float64) Point {

	tx := 1.0 - t

	x := ((b.P0.X * tx * tx * tx) +
		(b.P1.X * 3 * tx * tx * t)) +
		((b.P2.X * 3 * tx * t * t) +
			(b.P3.X * t * t * t))

	y := ((b.P0.Y * tx * tx * tx) +
		(b.P1.Y * 3 * tx * tx * t)) +
		((b.P2.Y * 3 * tx * t * t) +
			(b.P3.Y * t * t * t))

	return Point{x, y}
}

// Evaluates cubic bezier first derivative at t
func (b Bezier) QPrime(t float64) Point {
	tx := 1 - t
	d1 := b.P1.Subtract(b.P0)
	d2 := b.P2.Subtract(b.P1)
	d3 := b.P3.Subtract(b.P2)
	x := (d1.X * 3 * tx * tx) +
		(d2.X * 6 * tx * t) +
		(d3.X * 3 * t * t)
	y := (d1.Y * 3 * tx * tx) +
		(d2.Y * 6 * tx * t) +
		(d3.Y * 3 * t * t)
	return Point{x, y}
}

func (b Bezier) QPrimePrime(t float64) Point {
	tx := 1 - t

	x := ((b.P0.X + (b.P2.X - (b.P1.X * 2))) * (6 * tx)) +
		((b.P1.X + (b.P3.X - (b.P2.X * 2))) * (6 * t))

	y := ((b.P0.Y + (b.P2.Y - (b.P1.Y * 2))) * (6 * tx)) +
		((b.P1.Y + (b.P3.Y - (b.P2.Y * 2))) * (6 * t))

	return Point{x, y}
}
