// Command critical_speed is a command-line tool for computing the critical train speed
// on railway tracks using dispersion analysis. It supports both ballast and slab track
// models and reads physical parameters from a YAML configuration file.
//
// The tool performs the following steps:
//   - Parses a YAML configuration file describing track and soil parameters.
//   - Computes the dispersion curves of the railway track.
//   - Computes the dispersion curves of the soil layers.
//   - Identifies the critical speed where the track and soil phase velocities intersect.
//   - Outputs the results (omega, phase velocities, critical values) to a structured JSON file.
//
// Usage:
//
//	go run cmd/critical_speed/main.go -config path/to/config.yaml
//
// Or using the compiled binary:
//
//	./bin/critical_speed -config path/to/config.yaml
//
// Required flags:
//
//	-config string
//	 	Path to the YAML configuration file defining model parameters.
//
// Configuration:
// The YAML file must specify the track type ("ballast" or "slabtrack"), the frequency range,
// track structure parameters, soil layer properties, and output file destination.
//
// For a complete example configuration file, see:
//
//	./configs/sample_config.yaml
//	...
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	soil_dispersion "github.com/PlatypusBytes/GoTrain/internal/soil_dispersion"
	track_dispersion "github.com/PlatypusBytes/GoTrain/internal/track_dispersion"
	math_utils "github.com/PlatypusBytes/GoTrain/pkg/utils"
	"gopkg.in/yaml.v3"
)

// Config contains the configuration for the dispersion calculation.
// It contains all necessary parameters to define track type, frequency range,
// and physical properties of either ballast or slab tracks.
type Config struct {
	TrackType string `yaml:"track_type"` // Type of track: "ballast" or "slabtrack"
	Frequency struct {
		Min    float64 `yaml:"min"`    // Minimum angular frequency for calculation [rad/s]
		Max    float64 `yaml:"max"`    // Maximum angular frequency for calculation [rad/s]
		Points int     `yaml:"points"` // Number of angular frequency points to calculate
	} `yaml:"frequency"`
	BallastTrack struct {
		EIRail        float64 `yaml:"EI_rail"`        // Rail bending stiffness [N·m²]
		MRail         float64 `yaml:"m_rail"`         // Rail mass per unit length [kg/m]
		KRailPad      float64 `yaml:"k_rail_pad"`     // Railpad stiffness [N/m]
		CRailPad      float64 `yaml:"c_rail_pad"`     // Railpad damping [N·s/m]
		MSleeper      float64 `yaml:"m_sleeper"`      // Sleeper (distributed) mass [kg/m]
		EBallast      float64 `yaml:"E_ballast"`      // Young's modulus of ballast [Pa]
		HBallast      float64 `yaml:"h_ballast"`      // Ballast layer thickness [m]
		WidthSleeper  float64 `yaml:"width_sleeper"`  // Half-track width [m]
		RhoBallast    float64 `yaml:"rho_ballast"`    // Ballast density [kg/m³]
		SoilStiffness float64 `yaml:"soil_stiffness"` // Soil spring stiffness [N/m]
	} `yaml:"ballast_track"`
	SlabTrack struct {
		EIRail        float64 `yaml:"EI_rail"`        // Rail bending stiffness [N·m²]
		MRail         float64 `yaml:"m_rail"`         // Rail mass per unit length [kg/m]
		EISlab        float64 `yaml:"EI_slab"`        // Slab bending stiffness [N·m²]
		MSlab         float64 `yaml:"m_slab"`         // Slab mass per unit length [kg/m]
		KRailPad      float64 `yaml:"k_rail_pad"`     // Railpad stiffness [N/m]
		CRailPad      float64 `yaml:"c_rail_pad"`     // Railpad damping [N·s/m]
		SoilStiffness float64 `yaml:"soil_stiffness"` // Soil spring stiffness [N/m]
	} `yaml:"slab_track"`
	SoilLayers []SoilLayer `yaml:"soil_layers"` // Array of soil layers
	Output     struct {
		FileName string `yaml:"file_name"` // Name of the output JSON file
	} `yaml:"output"`
}

// DispersionResults defines the structure for storing calculation results
type DispersionResults struct {
	Omega              []float64 `json:"omega"`
	TrackPhaseVelocity []float64 `json:"track_phase_velocity"`
	SoilPhaseVelocity  []float64 `json:"soil_phase_velocity"`
	CriticalOmega      float64   `json:"critical_omega"`
	CriticalVelocity   float64   `json:"critical_velocity"`
}

// SoilLayer defines the structure for a soil layer
type SoilLayer struct {
	Thickness    float64 `yaml:"thickness"`     // Thickness of the soil layer [m]
	Density      float64 `yaml:"density"`       // Density of the soil layer [kg/m³]
	YoungModulus float64 `yaml:"young_modulus"` // Young's modulus of the soil layer [Pa]
	PoissonRatio float64 `yaml:"poisson_ratio"` // Poisson's ratio of the soil layer
}

// createBallastTrackParams creates ballast track parameters from config.
//
// Parameters:
//   - config: The configuration structure containing ballast track parameters
//
// Returns:
//   - track_dispersion.BallastTrackParameters: A struct with parameters for ballast track dispersion calculations
func createBallastTrackParams(config Config) track_dispersion.BallastTrackParameters {
	return track_dispersion.BallastTrackParameters{
		EIRail:        config.BallastTrack.EIRail,
		MRail:         config.BallastTrack.MRail,
		KRailPad:      config.BallastTrack.KRailPad,
		CRailPad:      config.BallastTrack.CRailPad,
		MSleeper:      config.BallastTrack.MSleeper,
		EBallast:      config.BallastTrack.EBallast,
		HBallast:      config.BallastTrack.HBallast,
		WidthSleeper:  config.BallastTrack.WidthSleeper,
		RhoBallast:    config.BallastTrack.RhoBallast,
		SoilStiffness: config.BallastTrack.SoilStiffness,
	}
}

// createSlabTrackParams creates slab track parameters from config.
//
// Parameters:
//   - config: The configuration structure containing slab track parameters
//
// Returns:
//   - track_dispersion.SlabTrackParameters: A struct with parameters for slab track dispersion calculations
func createSlabTrackParams(config Config) track_dispersion.SlabTrackParameters {
	return track_dispersion.SlabTrackParameters{
		EIRail:        config.SlabTrack.EIRail,
		MRail:         config.SlabTrack.MRail,
		EISlab:        config.SlabTrack.EISlab,
		MSlab:         config.SlabTrack.MSlab,
		KRailPad:      config.SlabTrack.KRailPad,
		CRailPad:      config.SlabTrack.CRailPad,
		SoilStiffness: config.SlabTrack.SoilStiffness,
	}
}

// createSoilLayers converts the soil layers from the config to soil_dispersion.Layer format
//
// Parameters:
//   - config: The configuration structure containing soil layer parameters
//
// Returns:
//   - []soil_dispersion.Layer: A slice of soil_dispersion.Layer objects
func createSoilLayers(config Config) []soil_dispersion.Layer {
	layers := make([]soil_dispersion.Layer, len(config.SoilLayers))

	for i, soilLayer := range config.SoilLayers {
		layer := soil_dispersion.Layer{
			Thickness:     soilLayer.Thickness,
			Density:       soilLayer.Density,
			YoungsModulus: soilLayer.YoungModulus,
			PoissonRatio:  soilLayer.PoissonRatio,
		}
		layer.WaveSpeed() // Calculate wave speeds
		layers[i] = layer
	}

	return layers
}

// saveResults saves the calculation results to a JSON file.
//
// Parameters:
//   - omega: Array of angular frequencies [rad/s]
//   - trackphaseVelocity: Array of phase velocities for the track [m/s]
//   - soilPhaseVelocity: Array of phase velocities for the soil layers [m/s], can contain nil values
//   - criticalOmega: Critical angular frequency [rad/s]
//   - criticalSpeed: Critical train speed [m/s]
//   - fileName: Path and name of the output JSON file
//
// The function creates directories as needed and writes the results
// in a structured JSON format.
func saveResults(omega []float64, trackPhaseVelocity []float64, soilPhaseVelocity []float64, criticalOmega float64, criticalSpeed float64, fileName string) {
	results := DispersionResults{
		Omega:              omega,
		TrackPhaseVelocity: trackPhaseVelocity,
		SoilPhaseVelocity:  soilPhaseVelocity,
		CriticalOmega:      criticalOmega,
		CriticalVelocity:   criticalSpeed,
	}

	jsonData, err := json.MarshalIndent(results, "", "\t")
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(fileName)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Error creating directory: %v", err)
		}
	}

	// Write to file
	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing JSON to file: %v", err)
	}
}

// loadConfig loads the configuration from a YAML file.
//
// Parameters:
//   - configPath: Path to the YAML configuration file
//
// Returns:
//   - Config: The loaded configuration structure
//   - error: An error if the file cannot be read or parsed
func loadConfig(configPath string) (Config, error) {

	var config Config

	// Read the configuration file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse YAML data
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("failed to parse YAML: %v", err)
	}

	return config, nil
}

// main executes the critical speed analysis for a railway track based on configuration.
// The function:
//   - Processes command-line flags for configuration path
//   - Loads the track parameters from a YAML configuration file
//   - Computes dispersion curves for either ballast or slab track
//   - Computes dispersion curve for the soil layered system (ToDo)
//   - Compute the critical train speed (ToDo)
//   - Saves the results to a JSON file
//
// Command-line flags:
//
//	-config: Path to YAML configuration file (required)
func main() {
	// Parse command line arguments
	configPath := flag.String("config", "", "Path to configuration YAML file (required)")
	flag.Parse()

	// Check if config file path is provided
	if *configPath == "" {
		log.Fatal("Error: You must provide a configuration file path using the -config flag")
	}

	// Load configuration
	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Create omega values based on configuration file
	omega := math_utils.Linspace(
		config.Frequency.Min,
		config.Frequency.Max,
		config.Frequency.Points,
	)

	var params track_dispersion.TrackParameters

	switch config.TrackType {
	case "ballast":
		params = createBallastTrackParams(config)
	case "slabtrack":
		params = createSlabTrackParams(config)
	default:
		log.Fatalf("Invalid track type: %s. Supported types are 'ballast' or 'slabtrack'", config.TrackType)
	}

	// Calculate the dispersion curve for the track
	phaseVelocity := track_dispersion.RailTrackDispersion(params, omega)

	// Process soil layers if provided
	soilLayers := createSoilLayers(config)

	// Calculate the dispersion curve for the soil layers
	soilPhaseVelocity := soil_dispersion.SoilDispersion(soilLayers, omega)

	// Compute the critical train speed
	omegaCrit, phaseVelocityCrit, err := math_utils.InterceptLines(omega, phaseVelocity, soilPhaseVelocity)
	if err != nil {
		log.Fatalf("Error calculating critical speed: %v", err)
	}

	// Save results to file
	saveResults(omega, phaseVelocity, soilPhaseVelocity, omegaCrit, phaseVelocityCrit, config.Output.FileName)
	fmt.Printf("Results written successfully to %s\n", config.Output.FileName)
}
