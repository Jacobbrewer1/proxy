package main

import (
	"fmt"
	"github.com/jacobbrewer1/reverse-proxy/cmd/proxy/config"
	"golang.org/x/exp/slog"
	"net/http"
)

type App struct {
	logger *slog.Logger
	cfg    *config.Config
	server *http.Server
	proxy  *proxyServer
}

func (a *App) start() error {
	a.logger.Info(fmt.Sprintf("listening at %s", a.server.Addr))
	http.Handle("/", a.proxy)
	return a.server.ListenAndServe()
}

func newApp(logger *slog.Logger, server *http.Server, proxyServer *proxyServer, cfg *config.Config) *App {
	return &App{
		logger: logger,
		cfg:    cfg,
		server: server,
		proxy:  proxyServer,
	}
}
