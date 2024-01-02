package main

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	appName = "proxier"
)

type configuration struct {
	Resources map[string]*resource `yaml:"resources"`
}

type resource struct {
	Endpoint string `yaml:"endpoint"`
	Method   string `yaml:"method"`
	Redirect string `yaml:"redirect"`

	Auth *auth `yaml:"auth"`
}

func (r *resource) isValid() error {
	if r.Endpoint == "" {
		return fmt.Errorf("endpoint is required")
	}
	if r.Redirect == "" {
		return fmt.Errorf("destination url is required")
	}
	if r.Method == "" {
		return fmt.Errorf("method is required")
	}
	return nil
}

type auth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func newConfiguration() (*configuration, error) {
	viper.AddConfigPath("data")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error loading config file: %s", err)
	}

	type yamlConfig struct {
		Resources []*resource `yaml:"resources"`
	}

	yamlCfg := new(yamlConfig)
	if err := viper.Unmarshal(&yamlCfg); err != nil {
		return nil, fmt.Errorf("error reading config file: %s", err)
	}
	if len(yamlCfg.Resources) == 0 {
		return nil, fmt.Errorf("no resources defined")
	}

	resources := make(map[string]*resource)
	for _, r := range yamlCfg.Resources {
		if err := r.isValid(); err != nil {
			return nil, fmt.Errorf("invalid resource: %w", err)
		}
		resources[r.Endpoint] = r
	}

	cfg := &configuration{
		Resources: resources,
	}

	return cfg, nil
}
