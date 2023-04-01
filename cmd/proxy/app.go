package main

import (
	"fmt"
	"golang.org/x/exp/slog"
	"net/http"
)

type App struct {
	logger *slog.Logger
	cfg    *Config
	server *http.Server
	proxy  *proxyServer
}

func (a *App) start() error {
	a.logger.Info(fmt.Sprintf("listening at %s", a.server.Addr))
	http.Handle("/", a.proxy)
	return a.server.ListenAndServe()
}

func newApp(logger *slog.Logger, server *http.Server, proxyServer *proxyServer, cfg *Config) *App {
	return &App{
		logger: logger,
		cfg:    cfg,
		server: server,
		proxy:  proxyServer,
	}
}
