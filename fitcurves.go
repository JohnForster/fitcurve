package fitcurves

func FitCurves(points []Point, tolerance float64) []Bezier {

	// Filter duplicate points
	points = dedupeSlice(points)

	if len(points) < 2 {
		return make([]Bezier, 0)
	}

	len := len(points)
	left_tangent := createTangent(points[1], points[0])
	right_tangent := createTangent(points[len-2], points[len-1])

	return fitCubic(points, left_tangent, right_tangent, tolerance)
}

func dedupeSlice[T comparable](sliceList []T) []T {
	dedupeMap := make(map[T]struct{})
	list := []T{}

	for _, slice := range sliceList {
		if _, exists := dedupeMap[slice]; !exists {
			dedupeMap[slice] = struct{}{}
			list = append(list, slice)
		}
	}

	return list
}

// Find the unit vector describing the tangent between two points, eg the
// two points at the beginning or at the end of the curve.
func createTangent(p0 Point, p1 Point) Vec2 {
	v := p0.Subtract(p1)
	return v.Normalize()
}

func fitCubic(points []Point, leftTangent Vec2, rightTangent Vec2, targetError float64) []Bezier {
	// Use heuristic if region only has two points
	if len(points) == 2 {
		p0, p1 := points[0], points[1]
		const MAGIC_NUMBER float64 = 3.0

		dist := p0.Subtract(p1).Length() / MAGIC_NUMBER
		c1 := p0.Translate(leftTangent.Mult(dist))
		c2 := p1.Translate(rightTangent.Mult(dist))
		bezier := Bezier{p0, c1, c2, p1}
		return []Bezier{bezier}
	}

	params := chordLengthParameterize(points)

	bezier, maxError, splitPoint := generate(points, params, params, leftTangent, rightTangent)

	if maxError == 0.0 || maxError < targetError {
		return []Bezier{bezier}
	}

	// If maxError is relatively close, try reparameterization and iteration
	if maxError < targetError*targetError {
		const MAX_ITERATIONS int = 20
		newParams, prevSplit, prevErr := params, splitPoint, maxError
		for i := 0; i < MAX_ITERATIONS; i++ {
			newParams = reparameterize(bezier, points, newParams)
			bezier, maxError, splitPoint = generate(points, params, newParams, leftTangent, rightTangent)

			if maxError < targetError {
				return []Bezier{bezier}
			} else
			// If the development of the fitted curve grinds to a halt,
			// we abort this attempt (and try a shorter curve):
			if splitPoint == prevSplit {
				errChange := maxError / prevErr
				if errChange > 0.9999 && errChange < 1.0001 {
					break
				}
			}

			prevErr, prevSplit = maxError, splitPoint
		}
	}

	//Fitting failed -- split at max error point and fit recursively
	beziers := []Bezier{}

	//To create a smooth transition from one curve segment to the next, we
	//calculate the line between the points directly before and after the
	//center, and use that as the tangent both to and from the center point.
	centerVector := points[splitPoint-1].Subtract(points[splitPoint+1])

	//However, this won't work if they're the same point, because the line we
	//want to use as a tangent would be 0. Instead, we calculate the line from
	//that "double-point" to the center point, and use its tangent.
	if centerVector.X == 0 && centerVector.Y == 0 {
		//[x,y] -> [-y,x]: http://stackoverflow.com/a/4780141/1869660
		centerVector = points[splitPoint-1].Subtract(points[splitPoint])
		centerVector = Vec2{X: -centerVector.Y, Y: centerVector.X}
	}

	// To and From point in opposite directions
	toCenterTangent := centerVector.Normalize()
	fromCenterTangent := toCenterTangent.Mult(-1)

	beziers = append(beziers, fitCubic(points[0:splitPoint+1], leftTangent, toCenterTangent, targetError)...)
	beziers = append(beziers, fitCubic(points[splitPoint:], fromCenterTangent, rightTangent, targetError)...)

	return beziers
}

// Calculate parameter values (u) for each point on the parametric curve Q
// using chord-length parameterization, where Q(u) ≈ Q(t).
func chordLengthParameterize(points []Point) []float64 {
	var params []float64
	for i, p := range points {
		var u float64
		if i == 0 {
			u = 0
		} else {
			prevU := params[i-1]
			prevP := points[i-1]
			diff := p.Subtract(prevP)
			dist := diff.Length()
			u = prevU + dist
		}

		params = append(params, u)
	}

	maxU := params[len(params)-1]
	for i, u := range params {
		params[i] = u / maxU
	}

	return params
}

func generate(points []Point, params0 []float64, params1 []float64, lt Vec2, rt Vec2) (Bezier, float64, int) {
	bezier := generateBezier(points, params1, lt, rt)
	// Find max deviation of points to fitted curve.
	// Here we always use original params because we need to
	// compare the current curve to the source polyline.
	maxError, splitPoint := computeMaxError(points, bezier, params0)

	return bezier, maxError, splitPoint
}

func generateBezier(points []Point, params []float64, lt Vec2, rt Vec2) Bezier {
	var (
		bezier                        Bezier
		a                             []Vec2      // Precomputed rhs for equation
		A                             [][]Vec2    // Precomputed rhs for equation
		C                             [][]float64 // Matrix representing
		X                             []float64   // Matrix representing
		det_C0_C1, det_C0_X, det_X_C1 float64     // Matrix determinants
		alpha_l, alpha_r              float64     // Alpha values

		firstPoint = points[0]
		lastPoint  = points[len(points)-1]
	)

	bezier.P0 = firstPoint
	bezier.P3 = lastPoint

	// Compute the As
	A = make([][]Vec2, len(params))
	for i := range A {
		A[i] = make([]Vec2, len(params))
	}
	for i, u := range params {
		ux := 1 - u
		a = A[i]
		a[0] = lt.Mult(3 * u * (ux * ux))
		a[1] = rt.Mult(3 * ux * (u * u))
	}

	// Create the C and X matrices
	C = [][]float64{{0, 0}, {0, 0}}
	X = []float64{0, 0}
	for i, u := range params {
		a = A[i]

		C[0][0] += a[0].Dot(a[0])
		C[0][1] += a[0].Dot(a[1])
		C[1][0] += a[0].Dot(a[1])
		C[1][1] += a[1].Dot(a[1])

		straightLine := Bezier{P0: firstPoint, P1: firstPoint, P2: lastPoint, P3: lastPoint}

		// Difference between actual point location and a straight line at point u
		tmp := points[i].Subtract(straightLine.Q(u))

		X[0] += a[0].Dot(tmp)
		X[1] += a[1].Dot(tmp)

	}

	// Compute determinants
	det_C0_C1 = (C[0][0] * C[1][1]) - (C[1][0] * C[0][1])
	det_C0_X = (C[0][0] * X[1]) - (C[1][0] * X[0])
	det_X_C1 = (X[0] * C[1][1]) - (X[1] * C[0][1])

	// Derive alpha values
	if det_C0_C1 == 0 {
		alpha_l, alpha_r = 0, 0
	} else {
		alpha_l = det_X_C1 / det_C0_C1
		alpha_r = det_C0_X / det_C0_C1
	}

	// If alpha negative, use the Wu/Barsky heuristic (see text).
	// If alpha is 0, you get coincident control points that lead to
	// divide by zero in any subsequent NewtonRaphsonRootFind() call.
	segLength := firstPoint.Subtract(lastPoint).Length()
	epsilon := 1.0e-6 * segLength
	if alpha_l < epsilon || alpha_r < epsilon {
		// Fall back to rough estimation:
		//   c1 is 1/3 along the segment in the direction of the left tangent
		//   c2 is 1/3 along the segment in the direction of the right tangent
		bezier.P1 = firstPoint.Translate(lt.Mult(segLength / 3.0))
		bezier.P2 = lastPoint.Translate(rt.Mult(segLength / 3.0))
	} else {
		bezier.P1 = firstPoint.Translate(lt.Mult(alpha_l))
		bezier.P2 = lastPoint.Translate(rt.Mult(alpha_r))
	}

	return bezier
}

// Find the maximum squared distance of digitized points to fitted curve.
func computeMaxError(points []Point, bezier Bezier, params []float64) (float64, int) {
	maxDist := 0.0
	splitPoint := len(points) / 2
	const GRANULARITY int = 10
	t_distMap := mapTtoRelativeDistances(bezier, GRANULARITY)

	for i, point := range points {
		t := find_t(params[i], t_distMap, GRANULARITY)

		v := bezier.Q(t).Subtract(point)

		// Just finding max, so no need to sqrt
		dist := v.X*v.X + v.Y*v.Y

		if dist > maxDist {
			maxDist = dist
			splitPoint = i
		}

	}

	return maxDist, splitPoint
}

func mapTtoRelativeDistances(bezier Bezier, n int) []float64 {
	totalLength := 0.0
	prevBt := bezier.P0
	dists := []float64{0.0}
	for i := 1; i <= n; i++ {
		t := float64(i) / float64(n)
		Bt := bezier.Q(t)

		sectionLength := Bt.Subtract(prevBt).Length()
		totalLength += sectionLength

		dists = append(dists, totalLength)
		prevBt = Bt
	}

	normalized := make([]float64, len(dists))
	for i, d := range dists {
		normalized[i] = d / totalLength
	}

	return normalized
}

func find_t(u float64, t_distMap []float64, n int) float64 {
	if u < 0 {
		return 0
	}
	if u > 1 {
		return 1
	}

	/*
	   'u' is a value between 0 and 1 telling us the relative position
	   of a point on the source polyline (linearly from the start (0) to the end (1)).
	   To see if a given curve - 'bez' - is a close approximation of the polyline,
	   we compare such a poly-point to the point on the curve that's the same
	   relative distance along the curve's length.

	   But finding that curve-point takes a little work:
	   There is a function "B(t)" to find points along a curve from the parametric parameter 't'
	   (also relative from 0 to 1: http://stackoverflow.com/a/32841764/1869660
	                               http://pomax.github.io/bezierinfo/#explanation),
	   but 't' isn't linear by length (http://gamedev.stackexchange.com/questions/105230).

	   So, we sample some points along the curve using a handful of values for 't'.
	   Then, we calculate the length between those samples via plain euclidean distance;
	   B(t) concentrates the points around sharp turns, so this should give us a good-enough outline of the curve.
	   Thus, for a given relative distance ('u'), we can now find an upper and lower value
	   for the corresponding 't' by searching through those sampled distances.
	   Finally, we just use linear interpolation to find a better value for the exact 't'.

	   More info:
	       http://gamedev.stackexchange.com/questions/105230/points-evenly-spaced-along-a-bezier-curve
	       http://stackoverflow.com/questions/29438398/cheap-way-of-calculating-cubic-bezier-length
	       http://steve.hollasch.net/cgindex/curves/cbezarclen.html
	       https://github.com/retuxx/tinyspline
	*/

	//Find the two t-s that the current param distance lies between,
	//and then interpolate a somewhat accurate value for the exact t:

	var t float64
	for i := 1; i <= n; i++ {
		if u <= t_distMap[i] {
			tMin := float64(i-1) / float64(n)
			tMax := float64(i) / float64(n)
			lenMin := t_distMap[i-1]
			lenMax := t_distMap[i]

			t = (u-lenMin)/(lenMax-lenMin)*(tMax-tMin) + tMin
			break
		}
	}

	return t
}

func reparameterize(bezier Bezier, points []Point, params []float64) []float64 {
	newParams := []float64{}
	for i, u := range params {
		u1 := findNewtonRaphsonRoot(bezier, points[i], u)
		newParams = append(newParams, u1)
	}

	return newParams
}

// Use Newton-Raphson iteration to find a better root.
func findNewtonRaphsonRoot(bezier Bezier, Point Point, u float64) float64 {
	/*
		Newton's root finding algorithm calculates f(x)=0 by reiterating
		x_n+1 = x_n - f(x_n)/f'(x_n)
		We are trying to find curve parameter u for some point p that minimizes
		the distance from that point to the curve. Distance point to curve is d=q(u)-p.
		At minimum distance the point is perpendicular to the curve.
		We are solving
		f = q(u)-p * q'(u) = 0
		with
		f' = q'(u) * q'(u) + q(u)-p * q''(u)
		gives
				|q(u)-p * q'(u)|
		u_1 = u -  ------------------------------------
				|q'(u)**2 + q(u)-p * q''(u)|
	*/
	Q1u := bezier.QPrime(u)
	Q2u := bezier.QPrimePrime(u)
	d := bezier.Q(u).Subtract(Point)

	numerator := d.Dot(Q1u.Vec())
	denominator := (Q1u.X*Q1u.X + Q1u.Y*Q1u.Y) + 2*(d.X*Q2u.X+d.Y*Q2u.Y)

	if denominator == 0 {
		return u
	} else {
		return u - (numerator / denominator)
	}
}
