package track_dispersion

import (
	"encoding/json"
	"github.com/PlatypusBytes/GoTrain/pkg/utils"
	"os"
	"testing"
)

// Test dispersion curve of the Ballasted track
func TestBallastTrack(t *testing.T) {
	// Define ballast track parameters
	ballastParams := BallastTrackParameters{
		EIRail:        1.29e7,
		MRail:         120,
		KRailPad:      5e8,
		CRailPad:      2.5e5,
		MSleeper:      490,
		EBallast:      1.2e8,
		HBallast:      0.35,
		WidthSleeper:  1.25,
		RhoBallast:    1800.0,
		SoilStiffness: 0,
	}

	omega := math_utils.Linspace(0.1, 250, 100)

	// Calculate dispersion curve
	phaseVelocity := RailTrackDispersion(ballastParams, omega)

	// Read the expected results from a JSON file
	expectedResults, _ := os.ReadFile("../../testdata/ballast_track_dispersion.json")

	var expected DispersionResults
	if err := json.Unmarshal(expectedResults, &expected); err != nil {
		t.Fatalf("Failed to unmarshal expected results: %v", err)
	}

	// Compare calculated results with expected results
	for i, v := range expected.Omega {
		if v != omega[i] {
			t.Errorf("Expected omega[%d] = %f, got %f", i, v, omega[i])
		}
		if expected.PhaseVelocity[i] != phaseVelocity[i] {
			t.Errorf("Expected phase_velocity[%d] = %f, got %f", i, expected.PhaseVelocity[i], phaseVelocity[i])
		}
	}

}

func TestSlabTrack(t *testing.T) {
	// Define ballast track parameters
	ballastParams := SlabTrackParameters{
		EIRail:        1.29e7,
		MRail:         120,
		KRailPad:      5e8,
		CRailPad:      2.5e5,
		EISlab:        1.2e8,
		MSlab:         490,
		SoilStiffness: 0,
	}

	omega := math_utils.Linspace(0.1, 250, 100)

	// Calculate dispersion curve
	phaseVelocity := RailTrackDispersion(ballastParams, omega)

	// Read the expected results from a JSON file
	expectedResults, _ := os.ReadFile("../../testdata/slab_track_dispersion.json")

	var expected DispersionResults
	if err := json.Unmarshal(expectedResults, &expected); err != nil {
		t.Fatalf("Failed to unmarshal expected results: %v", err)
	}

	// Compare calculated results with expected results
	for i, v := range expected.Omega {
		if v != omega[i] {
			t.Errorf("Expected omega[%d] = %f, got %f", i, v, omega[i])
		}
		if expected.PhaseVelocity[i] != phaseVelocity[i] {
			t.Errorf("Expected phase_velocity[%d] = %f, got %f", i, expected.PhaseVelocity[i], phaseVelocity[i])
		}
	}

}

// DispersionResults defines the structure for storing calculation results
type DispersionResults struct {
	Omega         []float64 `json:"omega"`
	PhaseVelocity []float64 `json:"phase_velocity"`
}
