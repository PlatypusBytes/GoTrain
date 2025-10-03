package soil_dispersion

import (
	"encoding/json"
	"github.com/PlatypusBytes/GoTrain/pkg/utils"
	"math"
	"os"
	"testing"
)

// Test computation of the wave speed
func TestWaveSpeed(t *testing.T) {

	layers := []Layer{
		{
			Density:       1900,
			YoungsModulus: 20e6,
			PoissonRatio:  0.2,
			Thickness:     1,
		},
		{
			Density:       2000,
			YoungsModulus: 30e6,
			PoissonRatio:  0.25,
			Thickness:     2,
		},
		{
			Density:       2200,
			YoungsModulus: 40e6,
			PoissonRatio:  0.3,
			Thickness:     3,
		},
	}

	for i := range layers {
		layers[i].WaveSpeed()
		p_modulus := layers[i].YoungsModulus * (1 - layers[i].PoissonRatio) / ((1 + layers[i].PoissonRatio) * (1 - 2*layers[i].PoissonRatio))
		shear_modulus := layers[i].YoungsModulus / (2 * (1 + layers[i].PoissonRatio))
		compression_wave_speed := math.Sqrt(p_modulus / layers[i].Density)
		shear_wave_speed := math.Sqrt(shear_modulus / layers[i].Density)
		if layers[i].CompressionalWaveSpeed != compression_wave_speed {
			t.Errorf("Layer %d: Compressional wave speed should be %f and not %f", i, compression_wave_speed,
				layers[i].CompressionalWaveSpeed)
		}
		if layers[i].ShearWaveSpeed != shear_wave_speed {
			t.Errorf("Layer %d: Shear wave speed should be %f and not %f", i, shear_wave_speed,
				layers[i].ShearWaveSpeed)
		}
	}
}

// TestDispersionSoil validates the computation of dispersion curves for a layered soil system.
// This test is based on the example from Foti et al. (2014), "Surface Wave Methods for Near-Surface
// Site Characterization" (CRC Press, pp 94, Fig 2.30).
func TestDispersionSoil_1(t *testing.T) {

	E0, nu0 := ComputeElasticProperties(1900, 100, 200)
	E1, nu1 := ComputeElasticProperties(1900, 200, 400)
	E2, nu2 := ComputeElasticProperties(1900, 300, 600)
	E3, nu3 := ComputeElasticProperties(1900, 400, 800)

	layers := []Layer{
		{
			Density:       1900,
			YoungsModulus: E0,
			PoissonRatio:  nu0,
			Thickness:     5,
		},
		{
			Density:       1900,
			YoungsModulus: E1,
			PoissonRatio:  nu1,
			Thickness:     10,
		},
		{
			Density:       1900.,
			YoungsModulus: E2,
			PoissonRatio:  nu2,
			Thickness:     15,
		},
		{
			Density:       1900.,
			YoungsModulus: E3,
			PoissonRatio:  nu3,
			Thickness:     math.Inf(1),
		},
	}

	for i := range layers {
		layers[i].WaveSpeed() // Calculate wave speeds
	}

	omega := math_utils.Linspace(1, 50*2*math.Pi, 100)
	phase_velocity := SoilDispersion(layers, omega)

	// Read the expected results from a JSON file
	expectedResults, _ := os.ReadFile("../../testdata/soil_dispersion_1.json")

	var expected DispersionResults
	if err := json.Unmarshal(expectedResults, &expected); err != nil {
		t.Fatalf("Failed to unmarshal expected results: %v", err)
	}

	// Compare calculated results with expected results
	for i, v := range expected.Omega {
		if v != omega[i] {
			t.Errorf("Expected omega[%d] = %f, got %f", i, v, omega[i])
		}
		if expected.PhaseVelocity[i] != phase_velocity[i] {
			t.Errorf("Expected phase_velocity[%d] = %f, got %f", i, expected.PhaseVelocity[i], phase_velocity[i])
		}
	}
}


// TestDispersionSoil validates the computation of dispersion curves for a layered soil system.
// This test is based on the example from Mezher et al. (2016), Figure 15a.
func TestDispersionSoil_2(t *testing.T) {

	layers := []Layer{
		{
			Density:       2000,
			YoungsModulus: 30e6,
			PoissonRatio:  0.35,
			Thickness:     2,
		},
		{
			Density:       2000,
			YoungsModulus: 40e6,
			PoissonRatio:  0.35,
			Thickness:     10,
		},
		{
			Density:       2000,
			YoungsModulus: 75e6,
			PoissonRatio:  0.4,
			Thickness:     math.Inf(1),
		},
	}

	for i := range layers {
		layers[i].WaveSpeed() // Calculate wave speeds
	}

	omega := math_utils.Linspace(1, 50*2*math.Pi, 100)
	phase_velocity := SoilDispersion(layers, omega)

	// Read the expected results from a JSON file
	expectedResults, _ := os.ReadFile("../../testdata/soil_dispersion_2.json")

	var expected DispersionResults
	if err := json.Unmarshal(expectedResults, &expected); err != nil {
		t.Fatalf("Failed to unmarshal expected results: %v", err)
	}

	// Compare calculated results with expected results
	for i, v := range expected.Omega {
		if v != omega[i] {
			t.Errorf("Expected omega[%d] = %f, got %f", i, v, omega[i])
		}
		if expected.PhaseVelocity[i] != phase_velocity[i] {
			t.Errorf("Expected phase_velocity[%d] = %f, got %f", i, expected.PhaseVelocity[i], phase_velocity[i])
		}
	}
}


// DispersionResults defines the structure for storing dispersion curve calculation results.
// It contains two arrays:
//   - Omega: Angular frequencies (rad/s) at which phase velocities are computed
//   - PhaseVelocity: Corresponding phase velocities (m/s) for each frequency
//
// This structure is used for JSON serialization and deserialization of test results.
type DispersionResults struct {
	Omega         []float64 `json:"omega"`
	PhaseVelocity []float64 `json:"phase_velocity"`
}

// Helper function to compute the Young's modulus and Poisson ratio from the shear and compression wave velocities.
// Given the material density and both shear and compressional wave speeds, this function computes
// the corresponding elastic properties using the following relationships:
//   - Shear modulus G = ρ * Vs²
//   - Poisson's ratio ν = (Vp² - 2Vs²) / (2(Vp² - Vs²))
//   - Young's modulus E = 2G(1 + ν)
//
// Parameters:
//   - density: Material density (kg/m³)
//   - shear_wave_speed: Shear wave velocity Vs (m/s)
//   - compressional_wave_speed: Compressional wave velocity Vp (m/s)
//
// Returns:
//   - youngs_modulus: Young's modulus E (Pa)
//   - poisson_ratio: Poisson's ratio ν (dimensionless)
func ComputeElasticProperties(density, shear_wave_speed, compressional_wave_speed float64) (float64, float64) {
	shear_modulus := density * math.Pow(shear_wave_speed, 2)
	poisson_ratio := (math.Pow(compressional_wave_speed, 2) - 2*math.Pow(shear_wave_speed, 2)) / (2 * (math.Pow(compressional_wave_speed, 2) - math.Pow(shear_wave_speed, 2)))
	youngs_modulus := 2 * shear_modulus * (1 + poisson_ratio)
	return youngs_modulus, poisson_ratio
}
