package config

import (
	"github.com/jacobbrewer1/reverse-proxy/pkg/dataacess"
	"github.com/jacobbrewer1/reverse-proxy/pkg/dataacess/connection"
	"github.com/jacobbrewer1/reverse-proxy/pkg/logging"
)

const (
	AppName     = `reverse_proxy`
	LogFilePath = `./`
)

type Config struct {
	ListeningPortHttp  string `yaml:"listening_port_http"`
	ListeningPortHttps string `yaml:"listening_port_https"`
	MonitoringPort     string `yaml:"monitoring_port"`
	CertificatePath    string `yaml:"certificate_path"`
	PrivateKeyPath     string `yaml:"private_key_path"`

	RedisDb       *connection.RedisDb `yaml:"redis_db"`
	LoggingConfig *logging.Config
}

func (c *Config) setConnections() error {
	if _, err := c.RedisDb.Conn(0); err != nil {
		return err
	}
	dataacess.Connections.SetRedisDb(c.RedisDb)
	return nil
}

func NewConfig(loggingConfig *logging.Config) (*Config, error) {
	cfg := &Config{
		ListeningPortHttp:  "80",
		ListeningPortHttps: "443",
		MonitoringPort:     "45454",
		CertificatePath:    "./certs/fullchain1.pem",
		PrivateKeyPath:     "./certs/privkey1.pem",
		RedisDb: &connection.RedisDb{
			Host: "localhost",
			Port: "6379",
		},
		LoggingConfig: nil,
	}
	cfg.LoggingConfig = loggingConfig
	if err := cfg.setConnections(); err != nil {
		return nil, err
	}
	return cfg, nil
}
