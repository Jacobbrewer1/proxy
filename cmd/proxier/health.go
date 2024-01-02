package main

import (
	"net/http"
)

func healthHandler() Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
