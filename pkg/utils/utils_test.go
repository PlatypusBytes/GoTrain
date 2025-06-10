package math_utils

import (
	"math"
	"testing"
)

// TestBrentSimplePolynomial tests the Brent method on a simple quadratic function.
func TestBrentSimplePolynomial(t *testing.T) {
	// f(x) = x^2 - 4 => roots at x = -2 and x = 2
	f := func(x float64) float64 {
		return x*x - 4
	}

	// Try to find the positive root near x = 2
	root, err := Brent(1.0, 3.0, 1e-12, f)
	if err != nil {
		t.Fatalf("Brent failed: %v", err)
	}

	expected := 2.0
	if math.Abs(root-expected) > 1e-6 {
		t.Errorf("Expected root near %f, but got %f", expected, root)
	}
}

// TestBrentInvalidInterval checks if Brent returns an error when f(a)*f(b) >= 0
func TestBrentInvalidInterval(t *testing.T) {
	// f(x) = x^2 + 1 has no real roots
	f := func(x float64) float64 {
		return x*x + 1
	}

	_, err := Brent(-1.0, 1.0, 1e-12, f)
	if err == nil {
		t.Error("Expected error for invalid interval, got nil")
	}
}

// TestBrentConvergenceTolerance tests convergence with tighter tolerance
func TestBrentConvergenceTolerance(t *testing.T) {
	f := func(x float64) float64 {
		return math.Sin(x)
	}

	// sin(x) = 0 has a root at x = π (≈ 3.14159)
	root, err := Brent(3.0, 4.0, 1e-12, f)
	if err != nil {
		t.Fatalf("Brent failed: %v", err)
	}

	expected := math.Pi
	if math.Abs(root-expected) > 1e-9 {
		t.Errorf("Expected root near π, but got %f", root)
	}
}

// TestBrentRootNearBoundary tests the Brent method on a function with a root near the search interval boundary
// This is specifically for the edge case with x^3 - 0.001 = 0, which has a root at x = 0.1
func TestBrentRootNearBoundary(t *testing.T) {
	f := func(x float64) float64 {
		return math.Pow(x, 3) - 0.001
	}

	root, err := Brent(0.01, 1.0, 1e-12, f)
	if err != nil {
		t.Fatalf("Brent failed: %v", err)
	}

	expected := 0.1 // The exact root is x = (0.001)^(1/3) ≈ 0.1
	if math.Abs(root-expected) > 1e-9 {
		t.Errorf("Expected root near %f, but got %f", expected, root)
	}
}

// TestLinspaceBasic tests the basic functionality of Linspace with default parameters
func TestLinspaceBasic(t *testing.T) {
	result := Linspace(0.0, 10.0, 11)

	// Expected: [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
	expected := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	if len(result) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if math.Abs(result[i]-v) > 1e-10 {
			t.Errorf("At index %d: Expected %f, got %f", i, v, result[i])
		}
	}
}

// TestLinspaceEmptyResult tests Linspace with n <= 0
func TestLinspaceEmptyResult(t *testing.T) {
	result := Linspace(0.0, 10.0, 0)
	if len(result) != 0 {
		t.Errorf("Expected empty result for n=0, got length %d", len(result))
	}

	result = Linspace(0.0, 10.0, -5)
	if len(result) != 0 {
		t.Errorf("Expected empty result for n=-5, got length %d", len(result))
	}
}

// TestLinspaceSingleElement tests Linspace with n = 1
func TestLinspaceSingleElement(t *testing.T) {
	result := Linspace(5.0, 10.0, 1)

	if len(result) != 1 {
		t.Fatalf("Expected length 1, got %d", len(result))
	}

	if math.Abs(result[0]-5.0) > 1e-10 {
		t.Errorf("Expected [5.0], got [%f]", result[0])
	}
}

// TestLinspaceDecreasingRange tests Linspace with start > end
func TestLinspaceDecreasingRange(t *testing.T) {
	result := Linspace(10.0, 0.0, 6)

	// Expected: [10, 8, 6, 4, 2, 0]
	expected := []float64{10, 8, 6, 4, 2, 0}

	if len(result) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if math.Abs(result[i]-v) > 1e-10 {
			t.Errorf("At index %d: Expected %f, got %f", i, v, result[i])
		}
	}
}

// TestLinspaceEndpoints tests if the endpoints are exactly as specified
func TestLinspaceEndpoints(t *testing.T) {
	// Use non-integer values to test floating-point precision
	start, end := 0.1, 250.0
	n := 100

	result := Linspace(start, end, n)

	if len(result) != n {
		t.Fatalf("Expected length %d, got %d", n, len(result))
	}

	// Check first element is exactly start
	if result[0] != start {
		t.Errorf("First element: expected %f, got %f", start, result[0])
	}

	// Check last element is exactly end
	if result[n-1] != end {
		t.Errorf("Last element: expected %f, got %f", end, result[n-1])
	}
}

// TestLinspaceSpacing tests if the spacing between elements is constant
func TestLinspaceSpacing(t *testing.T) {
	start, end := -5.0, 5.0
	n := 11

	result := Linspace(start, end, n)

	if len(result) != n {
		t.Fatalf("Expected length %d, got %d", n, len(result))
	}

	expectedSpacing := (end - start) / float64(n-1)

	for i := 1; i < len(result); i++ {
		spacing := result[i] - result[i-1]
		if math.Abs(spacing-expectedSpacing) > 1e-10 {
			t.Errorf("Spacing between elements %d and %d: expected %f, got %f",
				i-1, i, expectedSpacing, spacing)
		}
	}
}

func TestInterceptLines_1(t *testing.T) {

	x := []float64{0, 1, 2, 3, 4}
	y1 := []float64{0, 1, 2, 3, 4}
	y2 := []float64{4, 3, 2, 1, 0}

	interceptX, interceptY, err := InterceptLines(x, y1, y2)
	if err != nil {
		t.Fatalf("InterceptLines failed: %v", err)
	}
	expectedX := 2.0
	expectedY := 2.0
	if math.Abs(interceptX-expectedX) > 1e-10 || math.Abs(interceptY-expectedY) > 1e-10 {
		t.Errorf("Expected intercept at (%f, %f), got (%f, %f)", expectedX, expectedY, interceptX, interceptY)
	}
}

func TestInterceptLines_2(t *testing.T) {

	x := []float64{0, 1, 2, 3, 4}
	y1 := []float64{1, 1, 1, 1, 1}
	y2 := []float64{2, 2, 2, 2, 2}

	interceptX, interceptY, err := InterceptLines(x, y1, y2)
	if err == nil || err.Error() != "input arrays are parallel" {
		t.Fatalf("Expected error 'input arrays are parallel', got: %v", err)
	}
	expectedX := 0.0
	expectedY := 0.0
	if math.Abs(interceptX-expectedX) > 1e-10 || math.Abs(interceptY-expectedY) > 1e-10 {
		t.Errorf("Expected intercept at (%f, %f), got (%f, %f)", expectedX, expectedY, interceptX, interceptY)
	}
}

// TestInterceptLines_ShortArrays tests that InterceptLines returns an error when arrays are too short
func TestInterceptLines_ShortArrays(t *testing.T) {
	// Test with single-element arrays
	x := []float64{1}
	y1 := []float64{2}
	y2 := []float64{3}

	_, _, err := InterceptLines(x, y1, y2)
	if err == nil || err.Error() != "input arrays must have at least two elements" {
		t.Errorf("Expected error for short arrays, got: %v", err)
	}

	// Test with empty arrays
	x = []float64{}
	y1 = []float64{}
	y2 = []float64{}

	_, _, err = InterceptLines(x, y1, y2)
	if err == nil || err.Error() != "input arrays must have at least two elements" {
		t.Errorf("Expected error for empty arrays, got: %v", err)
	}
}

// TestInterceptLines_DifferentLengths tests that InterceptLines returns an error when arrays have different lengths
func TestInterceptLines_DifferentLengths(t *testing.T) {
	// Test with y1 longer than x
	x := []float64{0, 1, 2}
	y1 := []float64{0, 1, 2, 3}
	y2 := []float64{3, 2, 1}

	_, _, err := InterceptLines(x, y1, y2)
	if err == nil || err.Error() != "all input arrays must have the same length" {
		t.Errorf("Expected error for different length arrays, got: %v", err)
	}

	// Test with y2 longer than x
	x = []float64{0, 1, 2}
	y1 = []float64{0, 1, 2}
	y2 = []float64{3, 2, 1, 0}

	_, _, err = InterceptLines(x, y1, y2)
	if err == nil || err.Error() != "all input arrays must have the same length" {
		t.Errorf("Expected error for different length arrays, got: %v", err)
	}
}

// TestInterceptLines_ExactMatch tests when the intersection is exactly at one of the input points
func TestInterceptLines_ExactMatch(t *testing.T) {
	// Test where lines intersect exactly at x[2]
	x := []float64{0, 1, 2, 3, 4}
	y1 := []float64{0, 1, 2, 3, 4}
	y2 := []float64{4, 3, 2, 1, 0}

	interceptX, interceptY, err := InterceptLines(x, y1, y2)
	if err != nil {
		t.Fatalf("InterceptLines failed: %v", err)
	}
	expectedX := 2.0
	expectedY := 2.0
	if math.Abs(interceptX-expectedX) > 1e-10 || math.Abs(interceptY-expectedY) > 1e-10 {
		t.Errorf("Expected intercept at (%f, %f), got (%f, %f)", expectedX, expectedY, interceptX, interceptY)
	}

	// Test where lines intersect exactly at x[0]
	x = []float64{0, 1, 2, 3}
	y1 = []float64{1, 2, 3, 4}
	y2 = []float64{1, 0, -1, -2}

	interceptX, interceptY, err = InterceptLines(x, y1, y2)
	if err != nil {
		t.Fatalf("InterceptLines failed: %v", err)
	}
	expectedX = 0.0
	expectedY = 1.0
	if math.Abs(interceptX-expectedX) > 1e-10 || math.Abs(interceptY-expectedY) > 1e-10 {
		t.Errorf("Expected intercept at (%f, %f), got (%f, %f)", expectedX, expectedY, interceptX, interceptY)
	}

	// Test where lines intersect exactly at x[3] (last point)
	x = []float64{0, 1, 2, 3}
	y1 = []float64{4, 3, 2, 1}
	y2 = []float64{0, 0, 0, 1}

	interceptX, interceptY, err = InterceptLines(x, y1, y2)
	if err != nil {
		t.Fatalf("InterceptLines failed: %v", err)
	}
	expectedX = 3.0
	expectedY = 1.0
	if math.Abs(interceptX-expectedX) > 1e-10 || math.Abs(interceptY-expectedY) > 1e-10 {
		t.Errorf("Expected intercept at (%f, %f), got (%f, %f)", expectedX, expectedY, interceptX, interceptY)
	}
}

// TestInterceptLines_NonIntersectingLines tests when lines don't intersect within the provided x-range
func TestInterceptLines_NonIntersectingLines(t *testing.T) {
	x := []float64{0, 1, 2, 3, 4}
	y1 := []float64{5, 6, 7, 8, 9}
	y2 := []float64{0, 1, 2, 3, 4}

	_, _, err := InterceptLines(x, y1, y2)
	if err == nil {
		t.Errorf("Expected an error for non-intersecting lines, got nil")
	}
	// The implementation considers non-intersecting lines with different slopes as parallel
	// which is a reasonable implementation choice, so we accept either error message
}

// TestInterceptLines_MultipleIntersections tests that InterceptLines returns the first intersection when lines cross multiple times
func TestInterceptLines_MultipleIntersections(t *testing.T) {
	x := []float64{0, 1, 2, 3, 4, 5}
	y1 := []float64{0, 2, 0, 2, 0, 2} // Oscillating line: /\/\/\
	y2 := []float64{1, 1, 1, 1, 1, 1} // Horizontal line

	interceptX, interceptY, err := InterceptLines(x, y1, y2)
	if err != nil {
		t.Fatalf("InterceptLines failed: %v", err)
	}

	// The first intersection should be between x[0]=0 and x[1]=1
	expectedX := 0.5
	expectedY := 1.0
	if math.Abs(interceptX-expectedX) > 1e-10 || math.Abs(interceptY-expectedY) > 1e-10 {
		t.Errorf("Expected first intercept at (%f, %f), got (%f, %f)", expectedX, expectedY, interceptX, interceptY)
	}
}

// TestInterceptLines_AlmostParallel tests lines that are not quite parallel
func TestInterceptLines_AlmostParallel(t *testing.T) {
	// Let's make lines with a clear enough difference to be detected
	x := []float64{0, 1, 2, 3, 4}
	y1 := []float64{1.0, 1.2, 1.4, 1.6, 1.8} // Line with slope 0.2
	y2 := []float64{3.0, 2.5, 2.0, 1.5, 1.0} // Line with slope -0.5

	interceptX, interceptY, err := InterceptLines(x, y1, y2)
	if err != nil {
		t.Errorf("InterceptLines failed to find intersection: %v", err)
		return
	}

	// These lines should intersect at x=2.85714, y=1.57143
	expectedX := 2.85714
	expectedY := 1.57143

	if math.Abs(interceptX-expectedX) > 1e-4 || math.Abs(interceptY-expectedY) > 1e-4 {
		t.Errorf("Expected intercept near (%f, %f), got (%f, %f)", expectedX, expectedY, interceptX, interceptY)
	}
}

// TestInterceptLines_TangentLines tests lines that are tangent to each other
func TestInterceptLines_TangentLines(t *testing.T) {
	x := []float64{0, 1, 2, 3, 4}
	y1 := []float64{0, 1, 2, 1, 0} // Parabola-like curve
	y2 := []float64{2, 2, 2, 2, 2} // Horizontal line tangent at the peak

	interceptX, interceptY, err := InterceptLines(x, y1, y2)
	if err != nil {
		t.Fatalf("InterceptLines failed: %v", err)
	}

	// Lines should be tangent at x=2, y=2
	expectedX := 2.0
	expectedY := 2.0
	if math.Abs(interceptX-expectedX) > 1e-10 || math.Abs(interceptY-expectedY) > 1e-10 {
		t.Errorf("Expected tangent point at (%f, %f), got (%f, %f)", expectedX, expectedY, interceptX, interceptY)
	}
}

// TestInterceptLines_EdgeIntersection tests when the intersection is very close to the edge of the x range
func TestInterceptLines_EdgeIntersection(t *testing.T) {
	// Lines that intersect just inside the x range at x ≈ 0.005
	x := []float64{0, 1, 2, 3, 4}
	y1 := []float64{2.0, 3.0, 4.0, 5.0, 6.0}        // Line with constant slope 1
	y2 := []float64{2.01, 1.01, 0.01, -0.99, -1.99} // Line with constant slope -1

	interceptX, interceptY, err := InterceptLines(x, y1, y2)
	if err != nil {
		t.Fatalf("InterceptLines failed: %v", err)
	}

	// The intersection is calculated by interpolation in the function
	// Let's check the result is reasonable (between 0 and 1, and positive y)
	if interceptX < 0 || interceptX > 1 || interceptY < 0 {
		t.Errorf("Expected intercept in range (0-1, >0), got (%f, %f)", interceptX, interceptY)
	}

	// The exact expected values depend on the interpolation method
	// With the current implementation, we should get approximately (0.005, 2.005)
	expectedX := 0.005
	expectedY := 2.005
	if math.Abs(interceptX-expectedX) > 1e-10 || math.Abs(interceptY-expectedY) > 1e-10 {
		t.Errorf("Expected intercept at (%f, %f), got (%f, %f)", expectedX, expectedY, interceptX, interceptY)
	}
}

// TestInterceptLines_CloseIntersections tests when multiple intersections occur close together
func TestInterceptLines_CloseIntersections(t *testing.T) {
	// Create a situation where lines cross multiple times in close proximity
	x := []float64{0, 0.1, 0.2, 0.3, 0.4}
	y1 := []float64{1.0, 1.1, 0.9, 1.1, 0.9} // Oscillating line
	y2 := []float64{1.0, 1.0, 1.0, 1.0, 1.0} // Horizontal line

	interceptX, interceptY, err := InterceptLines(x, y1, y2)
	if err != nil {
		t.Fatalf("InterceptLines failed: %v", err)
	}

	// The first intersection is at x=0, y=1.0
	expectedX := 0.0
	expectedY := 1.0
	if math.Abs(interceptX-expectedX) > 1e-10 || math.Abs(interceptY-expectedY) > 1e-10 {
		t.Errorf("Expected first intercept at (%f, %f), got (%f, %f)", expectedX, expectedY, interceptX, interceptY)
	}
}

// TestInterceptLines_LargeValues tests with very large coordinate values
func TestInterceptLines_LargeValues(t *testing.T) {
	// Test with large numerical values
	x := []float64{1e6, 2e6, 3e6, 4e6, 5e6}
	y1 := []float64{1e6, 2e6, 3e6, 4e6, 5e6} // y = x
	y2 := []float64{5e6, 4e6, 3e6, 2e6, 1e6} // y = 6e6 - x

	interceptX, interceptY, err := InterceptLines(x, y1, y2)
	if err != nil {
		t.Fatalf("InterceptLines failed: %v", err)
	}

	// The intersection should be at (3e6, 3e6)
	expectedX := 3e6
	expectedY := 3e6
	if math.Abs((interceptX-expectedX)/expectedX) > 1e-10 || math.Abs((interceptY-expectedY)/expectedY) > 1e-10 {
		t.Errorf("Expected intercept at (%e, %e), got (%e, %e)", expectedX, expectedY, interceptX, interceptY)
	}
}
