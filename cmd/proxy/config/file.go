package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var Location string

func createConfig() (*Config, error) {
	path, ok := findFile()
	if !ok {
		return nil, fmt.Errorf("no config found at location %s", Location)
	}
	c, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %s", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(c, &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %s", err)
	}
	return &cfg, nil
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
