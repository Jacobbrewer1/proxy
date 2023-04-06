package main

import (
	"fmt"
	"github.com/jacobbrewer1/reverse-proxy/cmd/proxy/config"
	"github.com/jacobbrewer1/reverse-proxy/cmd/proxy/monitoring"
	"github.com/jacobbrewer1/reverse-proxy/pkg/dataacess"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type proxyServer struct {
	logger *slog.Logger
	proxy  *httputil.ReverseProxy
	cfg    *config.Config
}

func (p *proxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	redis := dataacess.NewRedisDal(r.Context(), 0)
	if redis == nil {
		p.logger.Error("redis client came back nil")
	}
	dbUrl, err := redis.GetValue("test")
	if err != nil {
		p.logger.Error("error fetching url from redis", slog.String("err", err.Error()))
		return
	}
	redirect, err := url.Parse(dbUrl)
	if err != nil {
		p.logger.Error("error parsing url", slog.String("err", err.Error()))
		return
	}
	prox := httputil.NewSingleHostReverseProxy(redirect)
	prox.ServeHTTP(w, r)
}

func newProxyServer(logger *slog.Logger, cfg *config.Config) *proxyServer {
	return &proxyServer{
		logger: logger,
		cfg:    cfg,
	}
}

type servers struct {
	secureServer     *http.Server
	httpServer       *http.Server
	monitoringServer *http.Server
}

func newServers(logger *slog.Logger, cfg *config.Config) *servers {
	return &servers{
		httpServer:       newHttpServer(logger, cfg),
		secureServer:     newHttpSecureServer(logger, cfg),
		monitoringServer: newMonitoringServer(cfg),
	}
}

func newHttpSecureServer(logger *slog.Logger, cfg *config.Config) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ListeningPortHttps),
		Handler:      middleware(http.DefaultServeMux, logger, cfg),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorLog:     log.Default(),
	}
}

func middleware(_ http.Handler, logger *slog.Logger, cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := prometheus.NewTimer(monitoring.RequestDuration.WithLabelValues(r.Host))
		defer t.ObserveDuration()

		monitoring.TotalRequests.WithLabelValues(r.Host).Inc()

		writer := &clientWriter{ResponseWriter: w}
		p := newProxyServer(logger, cfg)
		p.ServeHTTP(writer, r)

		monitoring.ResponseStatus.WithLabelValues(r.Host, fmt.Sprintf("%d", writer.statusCode)).Inc()
	})
}

func newHttpServer(logger *slog.Logger, cfg *config.Config) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ListeningPortHttp),
		Handler:      httpRedirect(logger, cfg),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorLog:     log.Default(),
	}
}

func httpRedirect(logger *slog.Logger, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			logger.Error("error splitting port from host", slog.String("err", err.Error()))
		}
		u := r.URL
		u.Host = net.JoinHostPort(host, cfg.ListeningPortHttps)
		u.Scheme = "https"
		log.Println(u.String())
		http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
	}
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
