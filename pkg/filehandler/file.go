package filehandler

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var Location string

func CreateConfigYaml[T any](cfg *T) error {
	path, ok := findFile()
	if !ok {
		return fmt.Errorf("no config found at location %s", Location)
	}
	c, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file: %s", err)
	}
	if err := yaml.Unmarshal(c, &cfg); err != nil {
		return fmt.Errorf("error unmarshalling config: %s", err)
	}
	return nil
}

func findFile() (string, bool) {
	abs, err := filepath.Abs(Location)
	if err != nil {
		return "", false
	}

	file, err := os.Open(abs)
	if err != nil {
		return "", false
	}

	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	return abs, true
}
