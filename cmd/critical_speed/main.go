// Package main provides the command-line interface for the critical speed analysis tool.
//
// The critical speed calculator analyzes railway vehicle dynamics to determine
// critical speeds where the vehicle-track interaction becomes unstable. This tool
// reads configuration parameters from a YAML file and performs the necessary
// calculations.
//
// Usage:
//
//	critical_speed -config <path/to/config.yaml>
//
// The configuration file must be provided via the -config flag and should contain
// all necessary parameters for the critical speed analysis.
package main

import (
	"flag"
	"log"

	"github.com/PlatypusBytes/GoTrain/internal/critical_speed"
)

// main is the entry point for the critical speed analysis application.
// It parses command-line flags, validates the configuration file path,
// and executes the critical speed calculation.
//
// The program requires a single flag:
//   - config: Path to the YAML configuration file (required)
//
// If the configuration file is not provided or if an error occurs during
// execution, the program will terminate with a fatal error message.
func main() {
	configPath := flag.String("config", "", "Path to configuration YAML file (required)")
	flag.Parse()

	if *configPath == "" {
		log.Fatal("Error: You must provide a configuration file path using -config")
	}

	if err := critical_speed.Run(*configPath); err != nil {
		log.Fatal(err)
	}
}
