package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/PlatypusBytes/GoTrain/internal/runner"
)

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
