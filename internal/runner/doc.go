// Package runner provides functionality for executing critical speed calculations
// in parallel across multiple YAML configuration files.
//
// The runner package orchestrates batch processing of railway critical speed analyses
// by spawning multiple worker goroutines to process configuration files concurrently.
// It invokes the critical_speed functionality for each configuration file and provides
// progress tracking and concurrency control.
//
// # Features
//
//   - Recursive directory traversal to discover all YAML configuration files
//   - Configurable worker pool for parallel processing
//   - Real-time progress tracking with visual progress bar
//   - Atomic counting for thread-safe progress reporting
//
// # Usage
//
// The runner can be used as a library by calling the Run function:
//
//	import "github.com/PlatypusBytes/GoTrain/internal/runner"
//
//	err := runner.Run("/path/to/configs", 4)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Or via the command-line interface:
//
//	go run cmd/runner/main.go -dir path/to/configs -workers 4
//
// Using the compiled binary:
//
//	./bin/runner -dir path/to/configs -workers 4
//
// # Command-line Flags
//
//	-dir string
//		Required. Directory containing YAML configuration files.
//		The runner will recursively search for all .yaml files.
//
//	-workers int
//		Optional. Number of parallel workers (default: number of logical CPUs).
//		Controls the level of concurrency for processing configuration files.
//
// # Requirements
//
//   - Configuration files must have the `.yaml` extension
//   - Files must follow the GoTrain configuration format (see configs/sample_config.yaml)
//
// # Example
//
// To process all configuration files in the testdata/batch directory using 8 workers:
//
//	err := runner.Run("testdata/batch", 8)
//	if err != nil {
//		log.Fatalf("Batch processing failed: %v", err)
//	}
//
// The runner will display a progress bar showing the processing status:
//
//	[==========================                        ] 50.00% (5/10)
package runner
