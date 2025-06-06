package soil_dispersion

import (
	"math"

	math_utils "github.com/PlatypusBytes/GoTrain/pkg/utils"
)

type Layer struct {
	Density                float64 // Density of the layer [kg/m^3]
	YoungsModulus          float64 // Young's modulus of the layer [Pa]
	PoissonRatio           float64 // Poisson's ratio of the layer
	Thickness              float64 // Thickness of the layer [m]
	CompressionalWaveSpeed float64 // Compressional wave speed [m/s]
	ShearWaveSpeed         float64 // Shear wave speed [m/s]
}

func (l *Layer) WaveSpeed() {
	shear_modulus := l.YoungsModulus / (2 * (1 + l.PoissonRatio))
	p_modulus := l.YoungsModulus * (1 - l.PoissonRatio) / ((1 + l.PoissonRatio) * (1 - 2*l.PoissonRatio))
	l.CompressionalWaveSpeed = math.Sqrt(p_modulus / l.Density)
	l.ShearWaveSpeed = math.Sqrt(shear_modulus / l.Density)
}

func SoilDispersion(layers []Layer, omega []float64) []*float64 {

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

	c_min := 0.6 * min_shear_wave_speed
	c_max := 1.6 * max_shear_wave_speed
	c_list := math_utils.Linspace(c_min, c_max, 1000)

	// Use pointers to allow null values in JSON
	phase_speed := make([]*float64, len(omega))

	for i := range omega {
		// Initialize with nil (will be represented as null in JSON)
		phase_speed[i] = nil

		d_1 := dispersionFastDelta(layers, omega[i], c_list[0])

		for j := 1; j < len(c_list); j++ {
			d_2 := dispersionFastDelta(layers, omega[i], c_list[j])
			if d_1*d_2 < 0 {
				// When solution is found, create a value and set it
				value := (c_list[j-1] + c_list[j]) / 2
				phase_speed[i] = &value
				break
			}
			d_1 = d_2
		}
	}

	return phase_speed
}

func dispersionFastDelta(layers []Layer, omega float64, c float64) float64 {

	// Calculate the wavenumber for each compressional wave speed
	wavenumber := omega / c

	// re-compute values for the first layer
	beta0 := layers[0].ShearWaveSpeed
	t_value := 2 - math.Pow(c/beta0, 2)
	mu0 := layers[0].Density * math.Pow(beta0, 2)

	// Initialize X1 with complex values
	X1 := []complex128{
		complex(mu0*2*t_value, 0),
		complex(mu0*-math.Pow(t_value, 2), 0),
		complex(0, 0),
		complex(0, 0),
		complex(mu0*-4, 0),
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
		p1 := complex(C_beta, 0)*x2 + complex(s, 0)*S_beta*x3
		p2 := complex(C_beta, 0)*x4 + complex(s, 0)*S_beta*x5
		p3 := complex(1/s, 0)*S_beta*x2 + complex(C_beta, 0)*x3
		p4 := complex(1/s, 0)*S_beta*x4 + complex(C_beta, 0)*x5

		q1 := complex(C_alpha, 0)*p1 - complex(r, 0)*S_alpha*p2
		q2 := complex(-1/r, 0)*S_alpha*p3 + complex(C_alpha, 0)*p4
		q3 := complex(C_alpha, 0)*p3 - complex(r, 0)*S_alpha*p4
		q4 := complex(-1/r, 0)*S_alpha*p1 + complex(C_alpha, 0)*p2

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

	// Calculate final determinant using complex values
	r_h_cmplx := complex(r_h, 0)
	s_h_cmplx := complex(s_h, 0)
	D := X1[1] + s_h_cmplx*X1[2] - r_h_cmplx*(X1[3]+s_h_cmplx*X1[4])

	// Return the real part as the result (should be real for a physical solution)
	return real(D)
}

func computeTerms(c float64, wavenumber float64, thickness float64, compressionalWave float64, shearWaveSpeed float64) (float64, complex128, float64, complex128, float64, float64) {
	// Declare variables at function level
	var r, s float64
	var C_alpha, C_beta float64
	var S_alpha, S_beta complex128

	epsilon := 1e-200 // small value to avoid division by zero

	// P-wave terms - equivalent to Python:
	// r = np.where(c < c_p, np.sqrt(1 - (c / c_p)**2), np.where(c == c_p, epsilon, np.sqrt((c / c_p)**2 - 1)))
	if c < compressionalWave {
		r = math.Sqrt(1 - math.Pow(c/compressionalWave, 2))
		C_alpha = math.Cosh(wavenumber * r * thickness)
		S_alpha = complex(math.Sinh(wavenumber*r*thickness), 0)
	} else if c == compressionalWave {
		r = epsilon
		C_alpha = math.Cosh(wavenumber * r * thickness)
		S_alpha = complex(math.Sinh(wavenumber*r*thickness), 0)
	} else {
		r = math.Sqrt(math.Pow(c/compressionalWave, 2) - 1)
		C_alpha = math.Cos(wavenumber * r * thickness)
		// Using complex number with imaginary part for S_alpha when c > c_p
		// Equivalent to Python: 1j * np.sin(k * r * d)
		S_alpha = complex(0, math.Sin(wavenumber*r*thickness))
	}

	// S-wave terms - equivalent to Python:
	// s = np.where(c < c_s, np.sqrt(1 - (c / c_s)**2), np.where(c == c_s, epsilon, np.sqrt((c / c_s)**2 - 1)))
	if c < shearWaveSpeed {
		s = math.Sqrt(1 - math.Pow(c/shearWaveSpeed, 2))
		C_beta = math.Cosh(wavenumber * s * thickness)
		S_beta = complex(math.Sinh(wavenumber*s*thickness), 0)
	} else if c == shearWaveSpeed {
		s = epsilon
		C_beta = math.Cosh(wavenumber * s * thickness)
		S_beta = complex(math.Sinh(wavenumber*s*thickness), 0)
	} else {
		s = math.Sqrt(math.Pow(c/shearWaveSpeed, 2) - 1)
		C_beta = math.Cos(wavenumber * s * thickness)
		// Using complex number with imaginary part for S_beta when c > c_s
		// Equivalent to Python: 1j * np.sin(k * s * d)
		S_beta = complex(0, math.Sin(wavenumber*s*thickness))
	}

	return C_alpha, S_alpha, C_beta, S_beta, r, s
}
