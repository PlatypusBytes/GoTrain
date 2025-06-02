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
