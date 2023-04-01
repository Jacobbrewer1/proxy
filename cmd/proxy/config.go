package main

import "github.com/jacobbrewer1/reverse-proxy/pkg/logging"

const (
	appName     = `proxy_server`
	logFilePath = `./`
)

type Config struct {
	ListeningPort string
	Redirect      string
	loggingConfig *logging.Config
}

func NewConfig(loggingConfig *logging.Config) *Config {
	return &Config{
		ListeningPort: "8443",
		Redirect:      `http://localhost:3000`,
		loggingConfig: loggingConfig,
	}
}
