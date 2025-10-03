// Command runner is a batch processing tool for executing critical speed calculations
// in parallel across multiple YAML configuration files. It invokes the critical_speed
// binary for each configuration file and provides progress tracking and concurrency control.
//
// The tool performs the following:
//   - Walks a directory recursively to discover all `.yaml` configuration files.
//   - Spawns a configurable number of worker goroutines.
//   - Each worker runs `critical_speed` with a given YAML file.
//
// Usage:
//
//	go run cmd/runner/main.go -dir path/to/configs -workers 4
//
// Or using the compiled binary:
//
//	./bin/runner -dir path/to/configs -workers 4
//
// Flags:
//
//	-dir string
//	 	Required. Directory containing YAML configuration files.
//	-workers int
//	 	Optional. Number of parallel workers (default: number of logical CPUs).
//
// Notes:
//   - Ensure the critical_speed binary is already compiled and present in ./bin/ directory.
//     On Windows, this will be critical_speed.exe; on other platforms, critical_speed.
//   - Files must have the `.yaml` extension and be properly formatted.
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// getBinaryPath returns the appropriate binary path for the current platform.
// On Windows, it appends .exe extension; on other platforms, it returns the path as-is.
func getBinaryPath(basePath string) string {
	if runtime.GOOS == "windows" {
		return basePath + ".exe"
	}
	return basePath
}

// Job represents a single YAML file to process by the critical_speed binary.
// It encapsulates the path to the configuration file.
type Job struct {
	path string // Path to the YAML configuration file
}

// worker processes jobs from the jobs channel concurrently.
// It executes the critical_speed binary for each YAML configuration file
// and tracks processing statistics.
//
// Parameters:
//
//	id             - Worker identifier for logging purposes
//	jobs           - Channel from which jobs are received
//	wg             - WaitGroup for synchronization of worker completion
//	processedCount - Atomic counter tracking the number of processed files
func worker(id int, jobs <-chan Job, wg *sync.WaitGroup, processedCount *atomic.Int64) {
	defer wg.Done()

	for job := range jobs {
		// Run the command for this YAML file
		binaryPath := getBinaryPath("./bin/critical_speed")
		cmd := exec.Command(binaryPath, "-config", job.path)

		if err := cmd.Run(); err != nil {
			log.Printf("Worker %d: Failed on config %s: %v\n", id, job.path, err)
		}

		// Increment processed count
		processedCount.Add(1)
	}
}

// reportProgress prints the current processing progress with a visual progress bar.
// It updates at regular intervals and terminates when processing is complete.
//
// Parameters:
//
//	processed - Atomic counter tracking the number of processed files
//	total     - Total number of files to process
//	done      - Channel signaling when all processing is complete
func reportProgress(processed *atomic.Int64, total int64, done <-chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			count := processed.Load()
			percent := float64(count) / float64(total) * 100

			// Create progress bar
			width := 50
			bar := strings.Repeat("=", int(float64(width)*float64(count)/float64(total)))
			padding := strings.Repeat(" ", width-len(bar))

			fmt.Printf("\r[%s%s] %.2f%% (%d/%d)", bar, padding, percent, count, total)
		case <-done:
			return
		}
	}
}

// main function to set up the runner for parallel processing of YAML configuration files.
// It parses command-line flags, validates input, discovers YAML files, and
// coordinates worker goroutines to process the files concurrently.
func main() {
	// Command line flags
	configDirPtr := flag.String("dir", "", "Directory containing YAML files (required)")
	workersPtr := flag.Int("workers", runtime.NumCPU(), "Number of worker goroutines")
	flag.Parse()

	// Ensure directory is specified
	if *configDirPtr == "" {
		log.Fatal("Required flag -dir not specified. Please provide the directory containing YAML files.")
	}

	configDir := *configDirPtr
	numWorkers := *workersPtr

	fmt.Printf("Starting processing with %d workers\n", numWorkers)

	// Ensure the binary is built before running
	binaryPath := getBinaryPath("./bin/critical_speed")
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		log.Fatalf("Binary '%s' not found. Please build the project first.", filepath.Base(binaryPath))
	}

	// Create job channel
	jobs := make(chan Job, 100)

	// WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup

	// Track progress
	var processedCount atomic.Int64
	var totalFiles atomic.Int64

	// Start workers
	for i := range numWorkers {
		wg.Add(1)
		go worker(i, jobs, &wg, &processedCount)
	}

	// Collect YAML files
	yamlFiles := []string{}
	err := filepath.WalkDir(configDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".yaml") {
			yamlFiles = append(yamlFiles, path)
			totalFiles.Add(1)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking through config directory: %v", err)
	}

	total := totalFiles.Load()
	fmt.Printf("Found %d YAML files to process\n", total)

	// Start progress reporting goroutine
	done := make(chan struct{})
	go reportProgress(&processedCount, total, done)

	// Send jobs to workers
	for _, path := range yamlFiles {
		jobs <- Job{path: path}
	}
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()

	// Signal progress reporting to finish
	close(done)

	fmt.Printf("\nCompleted processing %d YAML files\n", processedCount.Load())
}
