package config

import "github.com/jacobbrewer1/reverse-proxy/pkg/logging"

const (
	AppName     = `proxy_server`
	LogFilePath = `./`
)

type Config struct {
	ListeningPort string
	Redirect      string
	LoggingConfig *logging.Config
}

func NewConfig(loggingConfig *logging.Config) *Config {
	return &Config{
		ListeningPort: "8443",
		Redirect:      `http://localhost:3000`,
		LoggingConfig: loggingConfig,
	}
}
