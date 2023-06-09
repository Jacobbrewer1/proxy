//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jacobbrewer1/reverse-proxy/cmd/reverseproxy/config"
	"github.com/jacobbrewer1/reverse-proxy/pkg/logging"
)

func InitializeApp() (*App, error) {
	wire.Build(
		wire.Value(logging.LoggerName(config.AppName)),
		wire.Value(logging.LoggerPath(config.LogFilePath)),
		logging.NewConfig,
		logging.CommonLogger,
		config.NewConfig,
		newProxyServer,
		newServers,
		newApp,
	)
	return &App{}, nil
}
