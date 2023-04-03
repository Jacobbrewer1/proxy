package main

import "net/http"

type clientWriter struct {
	http.ResponseWriter
	status      int
	written     uint64
	wroteHeader bool
}

func (o *clientWriter) Write(p []byte) (bytes int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	bytes, err = o.ResponseWriter.Write(p)
	o.written += uint64(bytes)
	return
}

func (o *clientWriter) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}
