package math_utils

import (
	"fmt"
	"math"
)

// Brent finds a root of a function f in the interval [a, b] using Brent's method.
// It returns the root and an error if the method fails to converge.
//
// f must be continuous in [a, b], and f(a) and f(b) must have opposite signs,
// indicating a root exists in the interval (by the intermediate value theorem).
//
// The tolerance for convergence is  based on machine epsilon.
//
// Parameters:
//
//	a, b  - interval bounds (must bracket a root, i.e., f(a)*f(b) < 0)
//	f     - function for which the root is to be found
//
// Returns:
//
//	root  - the estimated root
//	error - an error if convergence fails or inputs are invalid
func Brent(a, b, tol float64, f func(float64) float64) (float64, error) {
	// Maximum number of iterations
	max_nb_iterations := 1000

	eps := math.Nextafter(1.0, 2.0) - 1.0
	if tol < eps {
		tol = eps
	}

	// Function evaluations at interval endpoints
	fa := f(a)
	fb := f(b)

	// Check if the interval brackets a root
	if fa*fb >= 0 {
		return 0, fmt.Errorf("root not bracketed: f(a) and f(b) must have opposite signs")
	}

	// If one of the endpoints is the root, return it immediately
	if fa == 0 {
		return a, nil
	}
	if fb == 0 {
		return b, nil
	}

	// Make sure that b is the point with the smaller function value
	if math.Abs(fa) < math.Abs(fb) {
		a, b = b, a
		fa, fb = fb, fa
	}

	// Initialize variables for the algorithm
	c := a     // c is the point with the next-to-smallest function value
	fc := fa   // fc is f(c)
	d := b - a // d is used for the step size
	e := d     // e is the distance moved in the step before last

	// Main iteration loop
	for iter := 0; iter < max_nb_iterations; iter++ {
		// Convergence check
		delta := 2*eps*math.Abs(b) + tol
		m := 0.5 * (c - b)

		// Check if we've converged
		if math.Abs(m) <= delta || fb == 0 {
			return b, nil // Converged to the root
		}

		// Decide which method to use
		useSecant := true

		// Check if we need to use bisection or an interpolation
		if math.Abs(e) >= delta && math.Abs(fa) > math.Abs(fb) {
			// Try inverse quadratic interpolation
			s := fb / fa
			var p, q float64

			if a == c {
				// Use linear interpolation (secant method) instead
				p = 2 * m * s
				q = 1 - s
			} else {
				// Use inverse quadratic interpolation
				q = fa / fc
				r := fb / fc
				p = s * (2*m*q*(q-r) - (b-a)*(r-1))
				q = (q - 1) * (r - 1) * (s - 1)
			}

			// Check if p/q is in bounds
			if p > 0 {
				q = -q
			} else {
				p = -p
			}

			// Accept the interpolated value if it's within bounds and
			// represents a sufficiently small step
			if 2*p < 3*m*q-math.Abs(delta*q) && p < math.Abs(0.5*e*q) {
				e = d
				d = p / q
				useSecant = false
			}
		}

		// If interpolation was rejected, use bisection
		if useSecant {
			e = m
			d = e
		}

		// Update a to be the previous best approximation
		a = b
		fa = fb

		// Update b using the chosen step
		if math.Abs(d) > delta {
			b += d
		} else if m > 0 {
			b += delta
		} else {
			b -= delta
		}

		// Evaluate function at new point
		fb = f(b)

		// Update c, fc for the next iteration based on the signs of f(a) and f(b)
		if fa*fb < 0 {
			c = a
			fc = fa
		}
	}

	return 0, fmt.Errorf("max iterations %d reached", max_nb_iterations)
}

// Linspace returns an array of n evenly spaced values over the interval [start, end].
// This function mimics the behavior of numpy's linspace function.
//
// Parameters:
//
//	start - the starting value of the sequence
//	end   - the end value of the sequence
//	n     - number of samples to generate
//
// Returns:
//
//	[]float64 - array of evenly spaced values
func Linspace(start, end float64, n int) []float64 {
	if n <= 0 {
		return []float64{}
	}

	if n == 1 {
		return []float64{start}
	}

	result := make([]float64, n)
	step := (end - start) / float64(n-1)

	for i := 0; i < n; i++ {
		result[i] = start + float64(i)*step
	}

	// Ensure the last element is exactly end
	if n > 1 {
		result[n-1] = end
	}

	return result
}
