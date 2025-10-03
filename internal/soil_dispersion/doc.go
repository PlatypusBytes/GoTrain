// Package soil_dispersion provides tools to compute phase velocity dispersion curves
// for soil profiles.
//
// The computation of the dispersion curve is based on the Fast Delta Matrix method:
// Buchen, P. W., & Ben-Hador, R. (1996). "Free-mode surface-wave computations".
// Geophysical Journal International, 124(3), 869–887.
// https://doi.org/10.1111/j.1365-246X.1996.tb05642.x
//
// The layers are assumed to be horizontal and infinite, and the last layer is always
// assumed to be a halfspace.
//
// # Layer Representation
//
// The Layer type represents a layer in a soil profile with its physical properties,
// including density, Young's modulus, Poisson's ratio, thickness, compressional wave
// speed, and shear wave speed.
//
// # Dispersion Calculation
//
// The SoilDispersion function calculates the phase velocity dispersion curve for a
// soil profile using a numerical root-finding approach. It finds the phase speed for
// each frequency in the provided omega array by iterating over a range of compressional
// wave speeds and uses the Fast Delta Matrix method to compute the dispersion relation.
//
// # Usage Example
//
//	layers := []soil_dispersion.Layer{
//		{Density: 1900, YoungsModulus: 50e6, PoissonRatio: 0.3, Thickness: 5},
//		{Density: 2000, YoungsModulus: 100e6, PoissonRatio: 0.25, Thickness: 0}, // halfspace
//	}
//	for i := range layers {
//		layers[i].WaveSpeed()
//	}
//	omega := math_utils.Linspace(1, 314, 100)
//	phaseVelocities := soil_dispersion.SoilDispersion(layers, omega)
package soil_dispersion
