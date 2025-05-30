// Package main provides a command-line tool for calculating
// the critical speeds of railway tracks using dispersion analysis.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	track_dispersion "github.com/PlatypusBytes/GoTrain/internal/track_dispersion"
	"github.com/PlatypusBytes/GoTrain/pkg/math_utils"
	"gopkg.in/yaml.v3"
)

// Config holds the configuration for the dispersion calculation
type Config struct {
    TrackType string `yaml:"track_type"`
    Frequency struct {
        Min    float64 `yaml:"min"`
        Max    float64 `yaml:"max"`
        Points int     `yaml:"points"`
    } `yaml:"frequency"`
    BallastTrack struct {
        EIRail        float64 `yaml:"ei_rail"`
        MRail         float64 `yaml:"m_rail"`
        KRailPad      float64 `yaml:"k_rail_pad"`
        CRailPad      float64 `yaml:"c_rail_pad"`
        MSleeper      float64 `yaml:"m_sleeper"`
        EBallast      float64 `yaml:"e_ballast"`
        HBallast      float64 `yaml:"h_ballast"`
        WidthSleeper  float64 `yaml:"width_sleeper"`
        RhoBallast    float64 `yaml:"rho_ballast"`
        SoilStiffness float64 `yaml:"soil_stiffness"`
    } `yaml:"ballast_track"`
    SlabTrack struct {
        EIRail        float64 `yaml:"ei_rail"`
        MRail         float64 `yaml:"m_rail"`
        EISlab        float64 `yaml:"ei_slab"`
        MSlab         float64 `yaml:"m_slab"`
        KRailPad      float64 `yaml:"k_rail_pad"`
        CRailPad      float64 `yaml:"c_rail_pad"`
        SoilStiffness float64 `yaml:"soil_stiffness"`
    } `yaml:"slab_track"`
    Output struct {
        FileName string `yaml:"file_name"`
    } `yaml:"output"`
}

// DispersionResults defines the structure for storing calculation results
type DispersionResults struct {
    Omega         []float64 `json:"omega"`
    PhaseVelocity []float64 `json:"phase_velocity"`
}

// main executes the critical speed analysis for a railway track based on configuration.
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

    // Create omega values based on configuration
    omega := math_utils.Linspace(
        config.Frequency.Min,
        config.Frequency.Max,
        config.Frequency.Points,
    )

    // Calculate dispersion curve based on track type
    var phaseVelocity []float64
    var params track_dispersion.TrackParameters

    switch config.TrackType {
    case "ballast":
        params = createBallastTrackParams(config)
    case "slabtrack":
        params = createSlabTrackParams(config)
    default:
        log.Fatalf("Invalid track type: %s. Supported types are 'ballast' or 'slabtrack'", config.TrackType)
    }

    // Calculate the dispersion curve
    phaseVelocity = calculateDispersionCurve(params, omega)

    // Save results to file
    saveResults(omega, phaseVelocity, config.Output.FileName)
    fmt.Printf("Results written successfully to %s\n", config.Output.FileName)
}

// loadConfig loads the configuration from a YAML file
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

// createBallastTrackParams creates ballast track parameters from config
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

// createSlabTrackParams creates slab track parameters from config
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

// calculateDispersionCurve calculates the dispersion curve for the given parameters and frequencies
func calculateDispersionCurve(params track_dispersion.TrackParameters, omega []float64) []float64 {
    return track_dispersion.RailTrackDispersion(params, omega)
}

// saveResults saves the calculation results to a JSON file
func saveResults(omega []float64, phaseVelocity []float64, fileName string) {
    results := DispersionResults{
        Omega:         omega,
        PhaseVelocity: phaseVelocity,
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
