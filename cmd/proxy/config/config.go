package config

import (
	"github.com/jacobbrewer1/reverse-proxy/pkg/logging"
)

const (
	AppName     = `proxy`
	LogFilePath = `./`
)

type Config struct {
	ListeningPortHttp  string `yaml:"listening_port_http"`
	ListeningPortHttps string `yaml:"listening_port_https"`
	Redirect           string `yaml:"redirect"`
	runHttps           bool   `yaml:"run_https"`
	CertificatePath    string `yaml:"certificate_path"`
	PrivateKeyPath     string `yaml:"private_key_path"`
	LoggingConfig      *logging.Config
}

func NewConfig(loggingConfig *logging.Config) (*Config, error) {
	cfg, err := createConfig()
	if err != nil {
		return nil, err
	}
	cfg.LoggingConfig = loggingConfig
	return cfg, nil
}
