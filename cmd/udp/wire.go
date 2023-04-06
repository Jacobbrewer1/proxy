//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jacobbrewer1/reverse-proxy/pkg/logging"
)

func InitializeApp() (*App, error) {
	wire.Build(
		wire.Value(logging.LoggerName("test")),
		wire.Value(logging.LoggerPath("./")),
		logging.NewConfig,
		logging.CommonLogger,
		NewApp,
	)
	return &App{}, nil
}
