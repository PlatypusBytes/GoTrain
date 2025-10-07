// Package main provides the command-line interface for the batch runner tool.
//
// The batch runner processes multiple YAML configuration files in parallel,
// executing critical speed analysis on each file concurrently. This tool is
// designed for high-throughput processing of multiple configurations, utilizing
// worker goroutines to maximize efficiency.
//
// Usage:
//
//	runner -dir <path/to/config/directory> [-workers <n>]
//
// The configuration directory must be provided via the -dir flag and should contain
// one or more YAML configuration files. The tool will recursively search for all
// .yaml files in the specified directory and process them concurrently.
//
// Flags:
//   - dir: Directory containing YAML configuration files (required)
//   - workers: Number of worker goroutines (optional, defaults to number of CPU cores)
//
// The program displays a real-time progress bar showing the percentage of completed
// files and provides summary statistics upon completion.
package main

import (
	"flag"
	"log"
	"runtime"

	runner "github.com/PlatypusBytes/GoTrain/internal/runner"
)

// main is the entry point for the batch runner application.
// It parses command-line flags, validates the configuration directory path,
// and orchestrates parallel processing of YAML configuration files.
//
// The program accepts two flags:
//   - dir: Path to directory containing YAML configuration files (required)
//   - workers: Number of concurrent worker goroutines (optional, defaults to runtime.NumCPU())
//
// If the configuration directory is not provided or if an error occurs during
// execution, the program will terminate with a fatal error message.
func main() {
	configDir := flag.String("dir", "", "Directory containing YAML files (required)")
	workers := flag.Int("workers", runtime.NumCPU(), "Number of worker goroutines")
	flag.Parse()

	if *configDir == "" {
		log.Fatal("You must provide -dir path/to/configs")
	}

	if err := runner.Run(*configDir, *workers); err != nil {
		log.Fatal(err)
	}
}
