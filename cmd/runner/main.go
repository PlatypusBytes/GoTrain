package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	configDir := "configs"
	err := filepath.WalkDir(configDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".yaml") {
			fmt.Printf("Running config: %s\n", path)

			// Ensure the binary is built before running
			if _, err := os.Stat("./bin/critical_speed"); os.IsNotExist(err) {
				log.Fatal("Binary 'critical_speed' not found in './bin/'. Please build the project first.")
			}
			cmd := exec.Command("./bin/critical_speed", "-config", path)

			if err := cmd.Run(); err != nil {
				log.Printf("Failed on config %s: %v\n", path, err)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking through config directory: %v", err)
	}
}
