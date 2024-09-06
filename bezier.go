package fitcurve

type Bezier struct {
	p0 Point
	c1 Point
	c2 Point
	p1 Point
}

// Evaluates cubic bezier at parameter t
func (b Bezier) Q(t float64) Point {
	// fmt.Printf("** b: %v\n", b)
	// fmt.Printf("** t: %v\n", t)

	tx := 1.0 - t
	// pA := (b.p0.x * tx * tx * tx)
	// pB := (b.c1.x * (3 * tx * tx * t))
	// pC := (b.c2.x * 3 * tx * t * t)
	// pD := (b.p1.x * t * t * t)

	// fmt.Printf("** pA: %v\n", pA)
	// fmt.Printf("** pB: %v\n", pB)
	// fmt.Printf("** pC: %v\n", pC)
	// fmt.Printf("** pD: %v\n", pD)

	// pBb := 3 * tx * tx * t
	// pBa := b.c1.x
	// fmt.Printf("pBb, pBa: %v, %v\n", pBb, pBa)

	x := ((b.p0.x * tx * tx * tx) +
		(b.c1.x * 3 * tx * tx * t)) +
		((b.c2.x * 3 * tx * t * t) +
			(b.p1.x * t * t * t))

	y := ((b.p0.y * tx * tx * tx) +
		(b.c1.y * 3 * tx * tx * t)) +
		((b.c2.y * 3 * tx * t * t) +
			(b.p1.y * t * t * t))

	return Point{x, y}
}

// Evaluates cubic bezier first derivative at t
func (b Bezier) QPrime(t float64) Point {
	tx := 1 - t
	d1 := b.c1.Subtract(b.p0)
	d2 := b.c2.Subtract(b.c1)
	d3 := b.p1.Subtract(b.c2)
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

	x := ((b.p0.x + (b.c2.x - (b.c1.x * 2))) * (6 * tx)) +
		((b.c1.x + (b.p1.x - (b.c2.x * 2))) * (6 * t))

	y := ((b.p0.y + (b.c2.y - (b.c1.y * 2))) * (6 * tx)) +
		((b.c1.y + (b.p1.y - (b.c2.y * 2))) * (6 * t))

	return Point{x, y}
}
