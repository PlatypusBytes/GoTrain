// Package math_utils provides mathematical utility functions for numerical computations
// used throughout the GoTrain project.
//
// The package implements various numerical methods including:
//   - Root finding algorithms (Brent's method)
//   - Linear space generator (similar to numpy's linspace)
//   - Line intersection calculations
//
// These utilities support the core computational needs of the dispersion analysis
// and other mathematical operations required for railway modeling.
//
// # Brent's Method
//
// The Brent function implements Brent's method for root finding, which combines
// the bisection method, secant method, and inverse quadratic interpolation to
// efficiently find roots of continuous functions.
//
// # Linear Space Generation
//
// The Linspace function generates evenly spaced values over a specified interval,
// similar to NumPy's linspace function, useful for frequency and wavenumber arrays.
//
// # Line Intersection
//
// The LineIntersection function computes the intersection point of two line segments,
// used in critical speed calculations where dispersion curves intersect.
package math_utils
