package fitcurves

import (
	"math"
	"testing"
)

func verifyMatch(expected []Bezier, actual []Bezier, t *testing.T) {
	if len(expected) != len(actual) {
		t.Fatalf("Different number of curves. Expected %v, received %v", len(expected), len(actual))
	}

	for i, eb := range expected {
		ab := actual[i]
		if !match(ab, eb) {
			t.Fatalf("Beziers with index %v didn't match. Expected %v, received %v", i, eb, ab)
		}
	}
}

func close(a float64, b float64) bool {
	const MAX_ALLOWED_DIFFERENCE float64 = 1.0e-6
	diff := math.Abs(b - a)
	return diff < MAX_ALLOWED_DIFFERENCE
}

func match(a Bezier, b Bezier) bool {
	same := true
	same = same && close(a.p0.x, b.p0.x)
	same = same && close(a.p0.y, b.p0.y)
	same = same && close(a.p1.x, b.p1.x)
	same = same && close(a.p1.y, b.p1.y)
	same = same && close(a.c1.x, b.c1.x)
	same = same && close(a.c1.y, b.c1.y)
	same = same && close(a.c2.x, b.c2.x)
	same = same && close(a.c2.y, b.c2.y)
	return same
}

func TestSingleBezier(t *testing.T) {
	expected := Bezier{
		p0: Point{x: 0, y: 0},
		c1: Point{x: 20.27317402, y: 20.27317402},
		c2: Point{x: -1.24665147, y: 0},
		p1: Point{x: 20, y: 0},
	}

	points := []Point{{x: 0, y: 0}, {x: 10, y: 10}, {x: 10, y: 0}, {x: 20, y: 0}}
	actual := FitCurves(points, 50)

	verifyMatch([]Bezier{expected}, actual, t)
}

func TestWithDuplicatePoints(t *testing.T) {
	expected := Bezier{
		p0: Point{x: 0, y: 0},
		c1: Point{x: 20.27317402, y: 20.27317402},
		c2: Point{x: -1.24665147, y: 0},
		p1: Point{x: 20, y: 0},
	}
	points := []Point{{x: 0, y: 0}, {x: 10, y: 10}, {x: 10, y: 0}, {x: 20, y: 0}, {x: 20, y: 0}}

	actual := FitCurves(points, 50)

	verifyMatch([]Bezier{expected}, actual, t)
}

func TestMoreComplexPoints(t *testing.T) {
	expected := Bezier{
		p0: Point{x: 244, y: 92},
		c1: Point{x: 284.2727272958473, y: 105.42424243194908},
		c2: Point{x: 287.98676736182495, y: 85},
		p1: Point{x: 297, y: 85},
	}
	points := []Point{{x: 244, y: 92}, {x: 247, y: 93}, {x: 251, y: 95}, {x: 254, y: 96}, {x: 258, y: 97}, {x: 261, y: 97}, {x: 265, y: 97}, {x: 267, y: 97}, {x: 270, y: 97}, {x: 273, y: 97}, {x: 281, y: 97}, {x: 284, y: 95}, {x: 286, y: 94}, {x: 289, y: 92}, {x: 291, y: 90}, {x: 292, y: 88}, {x: 294, y: 86}, {x: 295, y: 85}, {x: 296, y: 85}, {x: 297, y: 85}}

	actual := FitCurves(points, 10)

	verifyMatch([]Bezier{expected}, actual, t)
}

func TestUnalignedPointsWithLowTolerance(t *testing.T) {
	expected := []Bezier{
		{
			Point{0, 0},
			Point{3.333333333333333, 3.333333333333333},
			Point{5.285954792089683, 10},
			Point{10, 10},
		},
		{
			Point{10, 10},
			Point{13.333333333333334, 10},
			Point{7.6429773960448415, 2.3570226039551585},
			Point{10, 0},
		},
		{
			Point{10, 0},
			Point{12.3570226, -2.3570226},
			Point{16.66666667, 0},
			Point{20, 0},
		},
	}

	points := []Point{{0, 0}, {10, 10}, {10, 0}, {20, 0}}

	actual := FitCurves(points, 1)

	verifyMatch(expected, actual, t)
}

func TestNewtonRaphsonRootFind(t *testing.T) {
	bezier := Bezier{
		p0: Point{x: -106, y: 85},
		c1: Point{x: -85.27347011446706, y: 68.22138056885429},
		c2: Point{x: -167.14381916835873, y: 103.85618083164127},
		p1: Point{x: -186, y: 85},
	}

	point := Point{x: -185.0, y: 86.0}
	u := 0.9871784373992284

	expected := 0.982463387732839

	actual := findNewtonRaphsonRoot(bezier, point, u)

	if !close(expected, actual) {
		diff := actual - expected
		t.Fatalf("Not close enough. Expected %v, received %v. Diff %v", expected, actual, diff)
	}
}

func TestNewtonRaphsonRootFind2(t *testing.T) {
	bezier := Bezier{Point{244, 92}, Point{268.96666402690425, 100.32222134230143}, Point{279.14260825954716, 85}, Point{297, 85}}

	// i: 1
	point := Point{247, 93}
	u := 0.05389562833843188
	expected := 0.040747450391918696
	actual := findNewtonRaphsonRoot(bezier, point, u)

	if !close(expected, actual) {
		diff := actual - expected
		t.Fatalf("Not close enough. Expected %v, received %v. Diff %v", expected, actual, diff)
	}
}
