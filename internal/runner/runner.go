// Runner is a package for executing critical speed calculations in parallel across multiple YAML
// configuration files. It invokes the critical_speed  binary for each configuration file and
// provides progress tracking and concurrency control.
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
package runner

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	critical_speed "github.com/PlatypusBytes/GoTrain/internal/critical_speed"
)

// Job represents a single YAML file to process by the critical_speed binary.
type Job struct {
	path string
}

// worker processes jobs from the jobs channel concurrently.
func worker(id int, jobs <-chan Job, wg *sync.WaitGroup, processedCount *atomic.Int64) {
	defer wg.Done()

	for job := range jobs {

		// Execute the critical_speed with the YAML file
		if err := critical_speed.Run(job.path); err != nil {
			log.Printf("Worker %d: Failed on config %s: %v\n", id, job.path, err)
		}

		processedCount.Add(1)
	}
}

// reportProgress prints the current processing progress with a visual progress bar.
func reportProgress(processed *atomic.Int64, total int64, done <-chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			count := processed.Load()
			percent := float64(count) / float64(total) * 100
			width := 50
			bar := strings.Repeat("=", int(float64(width)*float64(count)/float64(total)))
			padding := strings.Repeat(" ", width-len(bar))
			fmt.Printf("\r[%s%s] %.2f%% (%d/%d)", bar, padding, percent, count, total)
		case <-done:
			return
		}
	}
}

// Run sets up the runner for parallel processing of YAML configuration files.
func Run(configDir string, numWorkers int) error {

	// Create job channel
	jobs := make(chan Job, 100)

	var wg sync.WaitGroup
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
		return fmt.Errorf("error walking through config directory: %v", err)
	}
	if len(yamlFiles) == 0 {
		return fmt.Errorf("no YAML configuration files found in directory: %s", configDir)
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

	wg.Wait()
	close(done)

	fmt.Printf("\nCompleted processing %d YAML files\n", processedCount.Load())
	return nil
}
