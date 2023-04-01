//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jacobbrewer1/reverse-proxy/pkg/logging"
)

func InitializeApp() (*App, error) {
	wire.Build(
		wire.Value(logging.LoggerName(appName)),
		wire.Value(logging.LoggerPath(logFilePath)),
		logging.NewConfig,
		logging.CommonLogger,
		NewConfig,
		newProxyServer,
		newHttpServer,
		newApp,
	)
	return &App{}, nil
}
