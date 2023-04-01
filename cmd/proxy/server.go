package main

import (
	"fmt"
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
	cfg    *Config
}

func (p *proxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend, err := url.Parse(p.cfg.Redirect)
	if err != nil {
		p.logger.Error(err.Error())
	}
	prox := httputil.NewSingleHostReverseProxy(backend)
	prox.ServeHTTP(w, r)
}

func newProxyServer(logger *slog.Logger, cfg *Config) *proxyServer {
	return &proxyServer{
		logger: logger,
		cfg:    cfg,
	}
}

func newHttpServer(logger *slog.Logger, cfg *Config) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ListeningPort),
		Handler:      middleware(http.DefaultServeMux, logger, cfg),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorLog:     log.Default(),
	}
}

func middleware(_ http.Handler, logger *slog.Logger, cfg *Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ob := &resp{
			ResponseWriter: w,
		}
		s := newProxyServer(logger, cfg)
		s.ServeHTTP(ob, r)
		logger.Debug("proxy result",
			slog.Group("details",
				slog.String("requestedFrom", r.RemoteAddr),
				slog.String("host", r.Host),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("protocol", r.Proto),
				slog.Int("responseStatus", ob.status),
				slog.Uint64("bytesWritten", ob.written),
				slog.String("referer", r.Referer()),
				slog.String("agent", r.UserAgent()),
			),
		)
	})
}
