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

// Test computation of the dispersion curve for a layered soil
func TestDispersionSoil(t *testing.T) {

	layers := []Layer{
		{
			Density:       1900,
			YoungsModulus: 50.6666666667e6,
			PoissonRatio:  0.3333333333333,
			Thickness:     5,
		},
		{
			Density:       1900,
			YoungsModulus: 202.666666667e6,
			PoissonRatio:  0.3333333333333,
			Thickness:     10,
		},
		{
			Density:       1900,
			YoungsModulus: 456e6,
			PoissonRatio:  0.3,
			Thickness:     15,
		},
		{
			Density:       1900,
			YoungsModulus: 819.666666667e6,
			PoissonRatio:  0.3333333333333,
			Thickness:     math.Inf(1),
		},
	}

	omega := math_utils.Linspace(0.1, 250, 100)
	phase_velocity := SoilDispersion(layers, omega)

	// Read the expected results from a JSON file
	expectedResults, _ := os.ReadFile("../../testdata/soil_dispersion.json")

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
			t.Errorf("Expected phase_velocity[%d] = %f, got %f", i, *expected.PhaseVelocity[i], *phase_velocity[i])
		}
	}
}

// DispersionResults defines the structure for storing calculation results
type DispersionResults struct {
	Omega         []float64  `json:"omega"`
	PhaseVelocity []*float64 `json:"phase_velocity"`
}
