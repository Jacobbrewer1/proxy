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

	proxy   *proxyServer
	servers *servers
}

func (a *App) start() error {
	http.Handle("/", a.proxy)
	return a.listenHttps()
}

func (a *App) listenHttps() error {
	a.logger.Info(fmt.Sprintf("listening https at %s", a.servers.secureServer.Addr))
	return a.servers.secureServer.ListenAndServeTLS(a.cfg.CertificatePath, a.cfg.PrivateKeyPath)
}

func newApp(logger *slog.Logger, servers *servers, proxyServer *proxyServer, cfg *config.Config) *App {
	return &App{
		logger:  logger,
		cfg:     cfg,
		servers: servers,
		proxy:   proxyServer,
	}
}
