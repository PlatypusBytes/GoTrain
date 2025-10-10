// Package GoTrain is a high-performance Go library for analyzing critical speeds in
// railway systems, focusing on soil and track dispersion analysis.
//
// # Overview
//
// GoTrain computes the speed at which critical train speed occurs.
// The critical train speed is the speed at which the train speed matches the phase velocity of
// waves propagating through the track-soil system.
//
// # Key Features
//
//   - Critical speed calculation for railway track-soil systems
//   - Support for both ballast and slab track configurations
//   - Multi-layered soil profile modelling with elastic properties
//   - High-performance parallel batch processing capabilities
//   - JSON output format for integration with other tools
//
// # Background
//
// This project is based on TrainCritSpeed (https://github.com/PlatypusBytes/TrainCritSpeed),
// originally implemented in Python. GoTrain reimplements the core functionality in Go,
// providing improved performance and native concurrency support.
//
// Main differences from TrainCritSpeed:
//   - Computes only the fundamental mode for subsurface layers
//   - Does not generate dispersion field plots
//   - Significantly faster execution with Go's performance characteristics
//   - Built-in parallel processing for batch operations
//
// For advanced features like higher-order modes and dispersion field visualization,
// please use the original TrainCritSpeed (https://github.com/PlatypusBytes/TrainCritSpeed) Python implementation.
//
// # Methodology
//
// The critical speed computation is based on established scientific methods:
//
// Critical Speed Analysis:
// Mezher, S. B., Connolly, D. P., Woodward, P. K., Laghrouche, O., Pombo, J., & Costa, P. A. (2016).
// "Railway critical velocity - Analytical prediction and analysis".
// Transportation Geotechnics, 6, 84–96.
// https://doi.org/10.1016/j.trgeo.2015.09.002
//
// The critical speed is identified at the intersection point of the track and soil
// dispersion curves, where the phase velocities match.
//
// Soil Dispersion Computation:
// Buchen, P. W., & Ben-Hador, R. (1996).
// "Free-mode surface-wave computations".
// Geophysical Journal International, 124(3), 869–887.
// https://doi.org/10.1111/j.1365-246X.1996.tb05642.x
//
// The soil dispersion curves are computed using the Fast Delta Matrix method, which
// efficiently handles multi-layered soil profiles with varying elastic properties.
//
// # Architecture
//
// The package is organized into several key components:
//
//   - internal/critical_speed: Core critical speed analysis engine
//   - internal/runner: Parallel batch processor for multiple configurations
//   - internal/soil_dispersion: Soil dispersion curve computation (Fast Delta Matrix)
//   - internal/track_dispersion: Track dispersion curve computation (ballast & slab tracks)
//   - pkg/utils: Mathematical utilities (Brent's method, linear interpolation, etc.)
//
// # Commands
//
// GoTrain provides two main command-line tools:
//
// Critical Speed Calculator (cmd/critical_speed):
//
// Analyzes a single railway configuration and computes dispersion curves and critical speed.
//
//	# Single configuration analysis
//	./critical_speed -config configs/sample_config.yaml
//
// The output is a JSON file containing omega values, track phase velocities, soil phase
// velocities, critical omega, and critical velocity.
//
// Batch Runner (cmd/runner):
//
// Processes multiple YAML configuration files in parallel with configurable worker pools.
// Automatically discovers all .yaml files in a directory tree and processes them concurrently.
//
//	# Process multiple configurations with 4 workers
//	./runner -dir testdata/batch -workers 4
//
// The runner displays a real-time progress bar and processes files concurrently for
// maximum throughput.
//
// # Library Usage
//
// GoTrain can be used as a library in your Go applications:
//
// Single Configuration Analysis:
//
//	import "github.com/PlatypusBytes/GoTrain/internal/critical_speed"
//
//	func main() {
//		err := critical_speed.Run("configs/my_config.yaml", true)
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//
// Batch Processing:
//
//	import "github.com/PlatypusBytes/GoTrain/internal/runner"
//
//	func main() {
//		// Process all YAML files in directory with 8 workers
//		err := runner.Run("configs_directory", 8)
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//
// # Configuration
//
// Configuration files use YAML format and must specify:
//   - Track type: "ballast" or "slabtrack"
//   - Frequency range: min, max, and number of points
//   - Track parameters: rail, sleeper/slab, railpad properties
//   - Soil layers: multi-layer profile with elastic properties
//   - Output: JSON filename for results
//
// See configs/sample_config.yaml for a complete annotated example.
//
// # Output Format
//
// Results are saved as JSON files with the following structure:
//
//	{
//		"omega": [1.0, 4.14, 7.28, ...],
//		"track_phase_velocity": [245.3, 251.7, ...],
//		"soil_phase_velocity": [183.5, 185.2, ...],
//		"critical_omega": 125.66,
//		"critical_velocity": 198.45
//	}
//
// Where:
//   - omega: Angular frequencies [rad/s]
//   - track_phase_velocity: Phase velocities in track system [m/s]
//   - soil_phase_velocity: Phase velocities in soil layers [m/s]
//   - critical_omega: Critical angular frequency [rad/s]
//   - critical_velocity: Critical train speed [m/s]
//
// # Installation
//
// Download Pre-built Binaries (Recommended):
//
// Download the latest release for your platform from:
// https://github.com/PlatypusBytes/GoTrain/releases
//
// Available for Linux, and Windows (amd64 architecture).
//
// Build from Source:
//
//	git clone https://github.com/PlatypusBytes/GoTrain.git
//	cd GoTrain
//	make build
//
// # References
//
// For more detailed information:
//   - GitHub: https://github.com/PlatypusBytes/GoTrain
//   - Original Python version: https://github.com/PlatypusBytes/TrainCritSpeed
//   - Package documentation: https://pkg.go.dev/github.com/PlatypusBytes/GoTrain
package gotrain
