package main

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func proxyHandler(target *url.URL, endpoint string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Proxying request", slog.String("url", r.URL.String()))

		// Update the headers to allow for SSL redirection
		r.URL.Host = target.Host
		r.URL.Scheme = target.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = target.Host

		path := r.URL.Path
		r.URL.Path = strings.TrimLeft(path, endpoint) // Trim the endpoint from the path

		proxy := httputil.NewSingleHostReverseProxy(target)
		// Note that ServeHttp is non-blocking and uses a go routine under the hood
		proxy.ServeHTTP(w, r)
	}
}
