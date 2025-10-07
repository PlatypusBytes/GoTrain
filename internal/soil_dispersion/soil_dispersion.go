package soil_dispersion

import (
	"math"
	"math/cmplx"

	math_utils "github.com/PlatypusBytes/GoTrain/pkg/utils"
)

// Layer represents a layer in a soil profile with its physical properties.
// It includes density, Young's modulus, Poisson's ratio, thickness,
// compressional wave speed, and shear wave speed.
type Layer struct {
	Density                float64 // Density of the layer [kg/m^3]
	YoungsModulus          float64 // Young's modulus of the layer [Pa]
	PoissonRatio           float64 // Poisson's ratio of the layer
	Thickness              float64 // Thickness of the layer [m]
	CompressionalWaveSpeed float64 // Compressional wave speed [m/s]
	ShearWaveSpeed         float64 // Shear wave speed [m/s]
}

// WaveSpeed calculates the compressional and shear wave speeds for the Layer
// based on its Young's modulus, Poisson's ratio, and density.
// The calculated values are stored in the Layer's CompressionalWaveSpeed and
// ShearWaveSpeed fields.
func (l *Layer) WaveSpeed() {
	shear_modulus := l.YoungsModulus / (2 * (1 + l.PoissonRatio))
	p_modulus := l.YoungsModulus * (1 - l.PoissonRatio) / ((1 + l.PoissonRatio) * (1 - 2*l.PoissonRatio))
	l.CompressionalWaveSpeed = math.Sqrt(p_modulus / l.Density)
	l.ShearWaveSpeed = math.Sqrt(shear_modulus / l.Density)
}

// SoilDispersion calculates the phase velocity dispersion curve for a soil profile
// using a numerical root-finding approach. It finds the phase speed for each frequency
// in the provided omega array by iterating over a range of compressional wave speeds.
// It returns a slice of pointers to float64, allowing for null values in the output.
// The function uses a fast method to compute the dispersion relation for each frequency.
//
// Parameters:
//   - layers: A slice of Layer structs representing the soil profile.
//   - omega: A slice of angular frequencies [rad/s] at which to compute phase velocities.
//
// Returns:
//   - A slice of pointers to float64, where each pointer corresponds to the phase speed
//     for the respective frequency in omega. If no solution is found, the pointer will be nil.
//
// Note: The function assumes that the layers have been initialized with their physical properties
// (density, Young's modulus, Poisson's ratio, thickness) and that the WaveSpeed method has been
// called to compute the wave speeds for each layer.
func SoilDispersion(layers []Layer, omega []float64) []float64 {

	// find the minimum & maximum compressional wave speed in layers
	min_shear_wave_speed := math.Inf(1)
	max_shear_wave_speed := math.Inf(-1)
	for _, layer := range layers {
		if layer.ShearWaveSpeed < min_shear_wave_speed {
			min_shear_wave_speed = layer.ShearWaveSpeed
		}
		if layer.ShearWaveSpeed > max_shear_wave_speed {
			max_shear_wave_speed = layer.ShearWaveSpeed
		}
	}

	c_min := 0.5 * min_shear_wave_speed
	c_max := max_shear_wave_speed
	c_list := math_utils.Linspace(c_min, c_max, int((c_max-c_min)/0.01))

	phase_speed := make([]float64, len(omega))

	for i := range omega {
		// Initialize with nan
		phase_speed[i] = math.NaN()

		d_1 := dispersionFastDelta(layers, omega[i], c_list[0])
		for j := range len(c_list) - 1 {
			d_2 := dispersionFastDelta(layers, omega[i], c_list[j+1])
			if d_1*d_2 < 0 {
				// When solution is found, create a value and set it
				value := (c_list[j-1] + c_list[j]) / 2
				phase_speed[i] = value
				break
			}
			d_1 = d_2
		}
	}
	return phase_speed
}

// dispersionFastDelta computes the dispersion relation for a given frequency
// and compressional wave speed using a fast method. It calculates the determinant
// of a matrix representing the track-soil system and returns the real part of the result.
// This function is optimized for performance and uses complex arithmetic to handle
// the wave propagation characteristics in the soil layers.
//
// Parameters:
//   - layers: A slice of Layer structs representing the soil profile.
//   - omega: Angular frequency [rad/s] at which to compute the dispersion relation.
//   - c: Compressional wave speed [m/s] to evaluate the dispersion relation.
//
// Returns:
//   - The real part of the determinant, representing the dispersion relation for the given frequency and compressional wave speed.
func dispersionFastDelta(layers []Layer, omega float64, c float64) float64 {

	// Calculate the wavenumber for each compressional wave speed
	wavenumber := omega / c

	// re-compute values for the first layer
	beta0 := layers[0].ShearWaveSpeed
	t_value := 2 - math.Pow(c/beta0, 2)
	mu0 := layers[0].Density * math.Pow(beta0, 2)

	// Initialize X1 with complex values
	X1 := []complex128{
		complex(mu0*mu0*2*t_value, 0),
		complex(mu0*mu0*-math.Pow(t_value, 2), 0),
		complex(0, 0),
		complex(0, 0),
		complex(mu0*mu0*-4, 0),
	}

	// Compute the terms for the halfspace (last layer)
	_, _, _, _, r_h, s_h := computeTerms(c, wavenumber, layers[len(layers)-1].Thickness, layers[len(layers)-1].CompressionalWaveSpeed, layers[len(layers)-1].ShearWaveSpeed)

	// Process each layer except the last one
	for i := 0; i < len(layers)-1; i++ {
		current_layer := layers[i]
		next_layer := layers[i+1]

		// Calculate layer properties directly when needed
		gamma := math.Pow(current_layer.ShearWaveSpeed/c, 2)
		gamma_next := math.Pow(next_layer.ShearWaveSpeed/c, 2)
		C_alpha, S_alpha, C_beta, S_beta, r, s := computeTerms(c, wavenumber, layers[i].Thickness, layers[i].CompressionalWaveSpeed, layers[i].ShearWaveSpeed)

		epsilon := next_layer.Density / current_layer.Density
		eta := 2 * (gamma - epsilon*gamma_next)

		a := epsilon + eta
		a_prime := a - 1
		b := 1 - eta
		b_prime := b - 1

		// Extract X1 components
		x1 := X1[0]
		x2 := X1[1]
		x3 := X1[2]
		x4 := X1[3]
		x5 := X1[4]

		// Calculate intermediate values using complex math
		p1 := C_beta*x2 + s*S_beta*x3
		p2 := C_beta*x4 + s*S_beta*x5
		p3 := 1/s*S_beta*x2 + C_beta*x3
		p4 := 1/s*S_beta*x4 + C_beta*x5

		q1 := C_alpha*p1 - r*S_alpha*p2
		q2 := -1/r*S_alpha*p3 + C_alpha*p4
		q3 := C_alpha*p3 - r*S_alpha*p4
		q4 := -1/r*S_alpha*p1 + C_alpha*p2

		y1 := complex(a_prime, 0)*x1 + complex(a, 0)*q1
		y2 := complex(a, 0)*x1 + complex(a_prime, 0)*q2
		z1 := complex(b, 0)*x1 + complex(b_prime, 0)*q1
		z2 := complex(b_prime, 0)*x1 + complex(b, 0)*q2

		// Update X1 for next iteration
		X1 = []complex128{
			complex(b_prime, 0)*y1 + complex(b, 0)*y2,
			complex(a, 0)*y1 + complex(a_prime, 0)*y2,
			complex(epsilon, 0) * q3,
			complex(epsilon, 0) * q4,
			complex(b_prime, 0)*z1 + complex(b, 0)*z2,
		}
	}

	// Calculate determinant using complex values
	r_h_cmplx := r_h //complex(r_h, 0)
	s_h_cmplx := s_h //complex(s_h, 0)
	D := X1[1] + s_h_cmplx*X1[2] - r_h_cmplx*(X1[3]+s_h_cmplx*X1[4])

	// Return the real part as the result
	return real(D)
}

// computeTerms calculates the terms needed for the dispersion relation
// based on the compressional wave speed, wavenumber, layer thickness,
// and the compressional and shear wave speeds of the layer.
// It returns the terms C_alpha, S_alpha, C_beta, S_beta, r, and s.
//
// Parameters:
//   - c: Compressional wave speed [m/s]
//   - wavenumber: Wavenumber [1/m]
//   - thickness: Thickness of the layer [m]
//   - compressionalWave: Compressional wave speed of the layer [m/s]
//   - shearWaveSpeed: Shear wave speed of the layer [m/s]
//
// Returns:
//   - C_alpha: Complex term for P-wave
//   - S_alpha: Complex term for P-wave
//   - C_beta: Complex term for S-wave
//   - S_beta: Complex term for S-wave
//   - r: Real term for P-wave
//   - s: Real term for S-wave
func computeTerms(c float64, wavenumber float64, thickness float64, compressionalWave float64, shearWaveSpeed float64) (complex128, complex128, complex128, complex128, complex128, complex128) {

	var r, s complex128
	var C_alpha, C_beta complex128
	var S_alpha, S_beta complex128

	r = cmplx.Sqrt(complex((1 - math.Pow(c/compressionalWave, 2)), 0))
	s = cmplx.Sqrt(complex((1 - math.Pow(c/shearWaveSpeed, 2)), 0))

	complex_wavenb := complex(wavenumber, 0)
	complex_thickness := complex(thickness, 0)

	C_alpha = cmplx.Cosh(complex_wavenb * r * complex_thickness)
	S_alpha = cmplx.Sinh(complex_wavenb * r * complex_thickness)

	C_beta = cmplx.Cosh(complex_wavenb * s * complex_thickness)
	S_beta = cmplx.Sinh(complex_wavenb * s * complex_thickness)

	return C_alpha, S_alpha, C_beta, S_beta, r, s
}
