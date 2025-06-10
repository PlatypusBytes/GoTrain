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

			cmd := exec.Command("./bin/critical_speed", "-config", path)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

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