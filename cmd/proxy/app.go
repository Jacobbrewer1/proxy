package main

import (
	"fmt"
	"github.com/jacobbrewer1/reverse-proxy/cmd/proxy/config"
	"github.com/jacobbrewer1/reverse-proxy/cmd/proxy/monitoring"
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
	if err := monitoring.Register(); err != nil {
		return err
	}

	http.Handle("/", a.proxy)
	go func() {
		if err := a.listenHttp(); err != nil {
			a.logger.Error("Error listening on http", slog.String("err", err.Error()))
		}
	}()
	go func() {
		if err := a.listenMonitor(); err != nil {
			a.logger.Error("Error listening for monitoring", slog.String("err", err.Error()))
		}
	}()
	return a.listenHttps()
}

func (a *App) listenHttps() error {
	a.logger.Info(fmt.Sprintf("Listening https at %s", a.servers.secureServer.Addr))
	return a.servers.secureServer.ListenAndServeTLS(a.cfg.CertificatePath, a.cfg.PrivateKeyPath)
}

func (a *App) listenHttp() error {
	a.logger.Info(fmt.Sprintf("Listening http at %s", a.servers.httpServer.Addr))
	return a.servers.httpServer.ListenAndServe()
}

func (a *App) listenMonitor() error {
	a.logger.Info(fmt.Sprintf("Listening for prometheus at %s", a.servers.monitoringServer.Addr))
	return a.servers.monitoringServer.ListenAndServe()
}

func newApp(logger *slog.Logger, servers *servers, proxyServer *proxyServer, cfg *config.Config) *App {
	return &App{
		logger:  logger,
		cfg:     cfg,
		servers: servers,
		proxy:   proxyServer,
	}
}
