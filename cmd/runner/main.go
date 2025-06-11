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

// Job represents a single YAML file to process
type Job struct {
	path string
}

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
	if _, err := os.Stat("./bin/critical_speed"); os.IsNotExist(err) {
		log.Fatal("Binary 'critical_speed' not found in './bin/'. Please build the project first.")
	}

	// Create job channel
	jobs := make(chan Job, 100)

	// WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup

	// Track progress
	var processedCount atomic.Int64
	var totalFiles atomic.Int64

	// Start workers
	for i := 0; i < numWorkers; i++ {
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

// worker processes jobs from the jobs channel
func worker(id int, jobs <-chan Job, wg *sync.WaitGroup, processedCount *atomic.Int64) {
	defer wg.Done()

	for job := range jobs {
		// Run the command for this YAML file
		cmd := exec.Command("./bin/critical_speed", "-config", job.path)

		if err := cmd.Run(); err != nil {
			log.Printf("Worker %d: Failed on config %s: %v\n", id, job.path, err)
		}

		// Increment processed count
		processedCount.Add(1)
	}
}

// reportProgress prints the current progress
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
