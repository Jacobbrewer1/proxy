package config

import (
	"github.com/jacobbrewer1/reverse-proxy/pkg/dataacess"
	"github.com/jacobbrewer1/reverse-proxy/pkg/dataacess/connection"
	"github.com/jacobbrewer1/reverse-proxy/pkg/filehandler"
	"github.com/jacobbrewer1/reverse-proxy/pkg/logging"
)

const (
	AppName     = `udp_proxy`
	LogFilePath = `./`
)

type Config struct {
	ListeningPort  string `yaml:"listening_port"`
	MonitoringPort string `yaml:"monitoring_port"`

	RedisDb       *connection.RedisDb `yaml:"redis_db"`
	LoggingConfig *logging.Config
}

func (c *Config) setConnections() {
	c.RedisDb.Client(0).Ping()

	dataacess.Connections.SetRedisDb(c.RedisDb)
}

func NewConfig(loggingConfig *logging.Config) (*Config, error) {
	var cfg Config
	if err := filehandler.CreateConfigYaml[Config](&cfg); err != nil {
		return nil, err
	}
	cfg.LoggingConfig = loggingConfig
	cfg.setConnections()
	return &cfg, nil
}
