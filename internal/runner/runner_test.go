package runner

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"testing"
)

// Test that if there are no YAML files, Run returns an appropriate error.
func TestRunWithNoConfigs(t *testing.T) {
	dir := t.TempDir() // empty dir
	err := Run(dir, 2)

	expectedMsg := "no YAML configuration files found in directory"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("expected error message to contain %q, got: %v", expectedMsg, err)
	}
}

// Test that Run processes YAML files without error.
func TestRunWithConfigs(t *testing.T) {

	BinaryPath = "../../bin/critical_speed"

	err := Run("../../testdata/batch", 4)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// expected data

	type expectedResult struct {
		criticalSpeed float64
		criticalOmega float64
	}

	var expectedResults = map[int]expectedResult{
		0: {criticalSpeed: 47.209851556707314, criticalOmega: 54.97135460046886},
		1: {criticalSpeed: 47.217618003722116, criticalOmega: 54.975836706978974},
		2: {criticalSpeed: 47.215617196350564, criticalOmega: 54.974682017874876},
		3: {criticalSpeed: 47.21551917331569, criticalOmega: 54.97462544764632},
		4: {criticalSpeed: 47.20975344897503, criticalOmega: 54.971297981360436},
		5: {criticalSpeed: 55.09280906583642, criticalOmega: 59.34618814480292},
		6: {criticalSpeed: 55.09280906583642, criticalOmega: 59.34618814480292},
		7: {criticalSpeed: 55.09255139144127, criticalOmega: 59.34605118652724},
		8: {criticalSpeed: 52.55226570820066, criticalOmega: 57.977738890313006},
		9: {criticalSpeed: 52.55270206278452, criticalOmega: 57.97798020782377},
	}

	// Check for expected output files
	for i := range 10 {
		// load json and compare results
		jsonPath := "tests/dispersion_results_" + strconv.Itoa(i) + ".json"

		// read json
		data, _ := os.ReadFile(jsonPath)

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

			// Verify critical_velocity and critical_omega values
			expected := expectedResults[i]
			if speed, ok := results["critical_velocity"].(float64); ok {
				if speed != expected.criticalSpeed {
					t.Errorf("File %d: expected critical_velocity %.15f, got %.15f", i, expected.criticalSpeed, speed)
				}
			}
			if omega, ok := results["critical_omega"].(float64); ok {
				if omega != expected.criticalOmega {
					t.Errorf("File %d: expected critical_omega %.15f, got %.15f", i, expected.criticalOmega, omega)
				}
			}
			// cleanup
			jsonPath := "tests/dispersion_results_" + strconv.Itoa(i) + ".json"
			os.Remove(jsonPath)

		}

	}
}
