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
	same = same && close(a.P0.X, b.P0.X)
	same = same && close(a.P0.Y, b.P0.Y)
	same = same && close(a.P3.X, b.P3.X)
	same = same && close(a.P3.Y, b.P3.Y)
	same = same && close(a.P1.X, b.P1.X)
	same = same && close(a.P1.Y, b.P1.Y)
	same = same && close(a.P2.X, b.P2.X)
	same = same && close(a.P2.Y, b.P2.Y)
	return same
}

func TestSingleBezier(t *testing.T) {
	expected := Bezier{
		P0: Point{X: 0, Y: 0},
		P1: Point{X: 20.27317402, Y: 20.27317402},
		P2: Point{X: -1.24665147, Y: 0},
		P3: Point{X: 20, Y: 0},
	}

	points := []Point{{X: 0, Y: 0}, {X: 10, Y: 10}, {X: 10, Y: 0}, {X: 20, Y: 0}}
	actual := FitCurves(points, 50)

	verifyMatch([]Bezier{expected}, actual, t)
}

func TestWithDuplicatePoints(t *testing.T) {
	expected := Bezier{
		P0: Point{X: 0, Y: 0},
		P1: Point{X: 20.27317402, Y: 20.27317402},
		P2: Point{X: -1.24665147, Y: 0},
		P3: Point{X: 20, Y: 0},
	}
	points := []Point{{X: 0, Y: 0}, {X: 10, Y: 10}, {X: 10, Y: 0}, {X: 20, Y: 0}, {X: 20, Y: 0}}

	actual := FitCurves(points, 50)

	verifyMatch([]Bezier{expected}, actual, t)
}

func TestMoreComplexPoints(t *testing.T) {
	expected := Bezier{
		P0: Point{X: 244, Y: 92},
		P1: Point{X: 284.2727272958473, Y: 105.42424243194908},
		P2: Point{X: 287.98676736182495, Y: 85},
		P3: Point{X: 297, Y: 85},
	}
	points := []Point{{X: 244, Y: 92}, {X: 247, Y: 93}, {X: 251, Y: 95}, {X: 254, Y: 96}, {X: 258, Y: 97}, {X: 261, Y: 97}, {X: 265, Y: 97}, {X: 267, Y: 97}, {X: 270, Y: 97}, {X: 273, Y: 97}, {X: 281, Y: 97}, {X: 284, Y: 95}, {X: 286, Y: 94}, {X: 289, Y: 92}, {X: 291, Y: 90}, {X: 292, Y: 88}, {X: 294, Y: 86}, {X: 295, Y: 85}, {X: 296, Y: 85}, {X: 297, Y: 85}}

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
		P0: Point{X: -106, Y: 85},
		P1: Point{X: -85.27347011446706, Y: 68.22138056885429},
		P2: Point{X: -167.14381916835873, Y: 103.85618083164127},
		P3: Point{X: -186, Y: 85},
	}

	point := Point{X: -185.0, Y: 86.0}
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
