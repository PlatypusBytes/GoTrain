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

// Job represents a single YAML configuration file to be processed.
// It contains the file path that will be passed to the critical_speed analyzer.
type Job struct {
	path string // Path to the YAML configuration file
}

// worker processes jobs from the jobs channel concurrently.
// It continuously reads Job items from the jobs channel, executes the critical_speed
// analyzer on each configuration file, and increments the processed count.
// If an error occurs during processing, it logs the error but continues with the next job.
// The worker signals completion to the WaitGroup when the jobs channel is closed.
//
// Parameters:
//   - id: Unique identifier for the worker goroutine (used in error logging)
//   - jobs: Receive-only channel from which Job items are read for processing
//   - wg: WaitGroup used to signal when the worker has completed all jobs
//   - processedCount: Atomic counter incremented for each successfully processed job
func worker(id int, jobs <-chan Job, wg *sync.WaitGroup, processedCount *atomic.Int64) {
	defer wg.Done()

	for job := range jobs {

		// Execute the critical_speed with the YAML file
		if err := critical_speed.Run(job.path, false); err != nil {
			log.Printf("Worker %d: Failed on config %s: %v\n", id, job.path, err)
		}

		processedCount.Add(1)
	}
}

// reportProgress prints the current processing progress with a visual progress bar.
// It runs in a separate goroutine and updates the console every second with a progress bar
// showing the percentage of completed jobs. The progress bar has a fixed width of 50 characters
// and displays both percentage completion and absolute counts (processed/total).
// The function terminates when a signal is received on the done channel.
//
// Parameters:
//   - processed: Atomic counter tracking the number of processed jobs (read concurrently)
//   - total: Total number of jobs to be processed
//   - done: Receive-only channel that signals when progress reporting should stop
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

// Run orchestrates parallel processing of YAML configuration files in the specified directory.
// It spawns numWorkers goroutines to process files concurrently and displays a progress bar.
//
// Parameters:
//   - configDir: Directory path to search for YAML configuration files (searched recursively)
//   - numWorkers: Number of concurrent workers to spawn for parallel processing
//
// Returns:
//   - error: An error if directory traversal fails or no YAML files are found
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
