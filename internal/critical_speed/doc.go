// Package critical_speed provides functionality for calculating and analyzing
// the critical speed of trains moving on either ballasted or slab track on top
// of the subsoil.
//
// The package combines soil and track dispersion analysis to determine the critical
// speed at which a moving train causes resonance with the soil-track system. It
// supports both ballast track and slab track configurations with configurable soil
// layer profiles.
//
// # Methodology
//
// The computation is based on the methodology described in:
// Mezher et al. (2016). "Railway critical velocity - Analytical prediction and analysis".
// Transportation Geotechnics, 6, 84–96.
// https://doi.org/10.1016/j.trgeo.2015.09.002
//
// The critical speed occurs at the intersection of the track and soil dispersion
// curves, where the phase velocity of waves in the track matches the phase velocity
// of surface waves in the soil.
//
// # Configuration
//
// The package reads YAML configuration files that specify:
//   - Track type (ballast or slab)
//   - Frequency range for analysis
//   - Track-specific parameters (rail properties, sleeper/slab properties, etc.)
//   - Soil layer profile (thickness, density, elastic properties)
//   - Output file location for results
//
// See configs/sample_config.yaml for a complete configuration example.
//
// # Results
//
// The analysis produces a JSON output file containing:
//   - Angular frequency array (omega)
//   - Track phase velocity dispersion curve
//   - Soil phase velocity dispersion curve
//   - Critical angular frequency (critical_omega)
//   - Critical velocity (critical_velocity)
//
// # Usage
//
// The package can be used as a library by calling the Run function:
//
//	import "github.com/PlatypusBytes/GoTrain/internal/critical_speed"
//
//	err := critical_speed.Run("configs/sample_config.yaml")
//	if err != nil {
//		log.Fatalf("Critical speed calculation failed: %v", err)
//	}
//
// Or via the command-line interface:
//
//	go run cmd/critical_speed/main.go -config configs/sample_config.yaml
//
// Using the compiled binary:
//
//	./bin/critical_speed -config configs/sample_config.yaml
//
// # Example
//
// To analyze a railway system with specific track and soil parameters:
//
//	configPath := "testdata/sample_config.yaml"
//	err := critical_speed.Run(configPath)
//	if err != nil {
//		log.Fatalf("Run failed: %v", err)
//	}
//
// The output JSON file will contain the complete dispersion analysis and the
// computed critical speed for the specified configuration.
package critical_speed
