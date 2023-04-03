package main

import (
	"fmt"
	"github.com/jacobbrewer1/reverse-proxy/cmd/proxy/config"
	"golang.org/x/exp/slog"
	"log"
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
	backend, err := url.Parse(p.cfg.Redirect)
	if err != nil {
		p.logger.Error(err.Error())
	}
	prox := httputil.NewSingleHostReverseProxy(backend)
	prox.ServeHTTP(w, r)
}

func newProxyServer(logger *slog.Logger, cfg *config.Config) *proxyServer {
	return &proxyServer{
		logger: logger,
		cfg:    cfg,
	}
}

type servers struct {
	secureServer *http.Server
	httpServer   *http.Server
}

func newServers(logger *slog.Logger, cfg *config.Config) *servers {
	return &servers{
		httpServer:   newHttpServer(logger, cfg),
		secureServer: newHttpSecureServer(logger, cfg),
	}
}

func newHttpServer(logger *slog.Logger, cfg *config.Config) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ListeningPortHttp),
		Handler:      middleware(http.DefaultServeMux, logger, cfg),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorLog:     log.Default(),
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
		writer := &clientWriter{ResponseWriter: w}
		p := newProxyServer(logger, cfg)
		p.proxy.ServeHTTP(writer, r)
	})
}
