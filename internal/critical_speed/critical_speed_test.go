package critical_speed

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

const TOL = 1e-3

// Test the computation of critical speed using a sample configuration file.
// This is an integration test that compares the output against expected values.
func TestRunWithSampleConfig(t *testing.T) {
	tmpFile := filepath.Join("dispersion_results.json")
	configPath := "../../testdata/sample_config.yaml"

	err := Run(configPath)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Check if the output file was created
	if _, err := os.Stat(tmpFile); err != nil {
		t.Errorf("expected output file %s not created", tmpFile)
	}

	// Load the output file and check contents
	data, _ := os.ReadFile(tmpFile)

	var results map[string]interface{}
	if err := json.Unmarshal(data, &results); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}

	// Check for expected keys and values
	expectedKeys := []string{"omega", "track_phase_velocity", "soil_phase_velocity", "critical_omega", "critical_velocity"}
	for _, key := range expectedKeys {
		if _, exists := results[key]; !exists {
			t.Errorf("expected key %s not found in results", key)
		}
	}
	expected_speed := 78.223
	if speed, ok := results["critical_velocity"].(float64); !ok {
		t.Errorf("critical_velocity is not a float64")
	} else if diff := speed - expected_speed; diff < -TOL || diff > TOL {
		t.Errorf("unexpected critical_velocity: got %v, want %v (tolerance %v)", speed, expected_speed, TOL)
	}

	expectedOmega := 63.017
	if omega, ok := results["critical_omega"].(float64); !ok {
		t.Errorf("critical_omega is not a float64")
	} else if diff := omega - expectedOmega; diff < -TOL || diff > TOL {
		t.Errorf("unexpected critical_omega: got %v, want %v (tolerance %v)", omega, expectedOmega, TOL)
	}

	os.Remove(tmpFile)

}
