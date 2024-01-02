//go:build wireinject
// +build wireinject

package main

import (
	"github.com/Jacobbrewer1/proxy/pkg/logging"
	"github.com/google/wire"
	"github.com/gorilla/mux"
)

func initializeApp() (*app, error) {
	wire.Build(
		wire.Value(logging.AppName(appName)),
		logging.NewConfig,
		logging.CommonLogger,
		mux.NewRouter,
		newConfiguration,
		newApp,
	)
	return new(app), nil
}
