package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/Jacobbrewer1/proxy/pkg/request"
	"github.com/gorilla/mux"
)

// The app is the main application.
type app struct {
	// r is the router.
	r *mux.Router

	// srv is the server.
	srv *http.Server

	// cfg is the configuration.
	cfg *configuration
}

func newApp(_ *slog.Logger, r *mux.Router, cfg *configuration) *app {
	return &app{
		r:   r,
		cfg: cfg,
	}
}

func (a *app) run() error {
	if err := a.init(); err != nil {
		return fmt.Errorf("error initializing app: %w", err)
	}

	if err := a.srv.ListenAndServe(); err != nil {
		return fmt.Errorf("error running server: %w", err)
	}

	return nil
}

func (a *app) init() error {
	for _, r := range a.cfg.Resources {
		dest, err := url.Parse(r.Redirect)
		if err != nil {
			return fmt.Errorf("error parsing destination url: %w", err)
		}

		if err := r.isValid(); err != nil {
			return fmt.Errorf("invalid resource: %w", err)
		}

		authOpt := AuthOptionNone

		if r.Auth != nil {
			authOpt = AuthOptionRequired
		}

		ph := proxyHandler(dest, r.Endpoint)
		a.r.HandleFunc(r.Endpoint, a.middlewareHttp(ph, authOpt)).Methods(r.Method)
	}

	a.r.NotFoundHandler = request.NotFoundHandler()
	a.r.MethodNotAllowedHandler = request.MethodNotAllowedHandler()

	a.srv = &http.Server{
		Addr:    ":8080",
		Handler: a.r,
	}

	return nil
}
