//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jacobbrewer1/reverse-proxy/cmd/udp/config"
	"github.com/jacobbrewer1/reverse-proxy/pkg/logging"
)

func InitializeApp() (*App, error) {
	wire.Build(
		wire.Value(logging.LoggerName(config.AppName)),
		wire.Value(logging.LoggerPath(config.LogFilePath)),
		logging.NewConfig,
		logging.CommonLogger,
		config.NewConfig,
		NewApp,
	)
	return &App{}, nil
}
