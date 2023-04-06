package main

import (
	"fmt"
	"github.com/jacobbrewer1/reverse-proxy/cmd/udp/config"
	"github.com/jacobbrewer1/reverse-proxy/cmd/udp/monitoring"
	"log"
	"net/http"
	"time"
)

func (a *App) startMonitoring() error {
	if err := monitoring.Register(); err != nil {
		return err
	}
	a.logger.Info(fmt.Sprintf("Listening for prometheus at %s", a.monitoringServer.Addr))
	return a.monitoringServer.ListenAndServe()
}

func newMonitoringServer(cfg *config.Config) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.MonitoringPort),
		Handler:      monitoring.Handler(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorLog:     log.Default(),
	}
}
