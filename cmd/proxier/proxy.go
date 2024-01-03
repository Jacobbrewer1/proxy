package main

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/Jacobbrewer1/proxy/pkg/logging"
)

func proxyHandler(dest string) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Proxying request", slog.String("url", r.URL.String()))

		// Parse the target URL.
		target, err := url.Parse(dest)
		if err != nil {
			slog.Error("Error parsing destination url", slog.String(logging.KeyError, err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Host = target.Host
		r.URL.Host = target.Host
		r.URL.Scheme = target.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))

		// Prevent the proxy from adding a trailing slash.
		r.URL.Path = target.Path
		target.Path = ""

		proxy := httputil.NewSingleHostReverseProxy(target)
		// Note that ServeHttp is non-blocking and uses a go routine under the hood
		proxy.ServeHTTP(w, r)
	}
}
