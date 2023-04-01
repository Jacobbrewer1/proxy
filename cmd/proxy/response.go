package main

import "net/http"

type resp struct {
	http.ResponseWriter
	status      int
	written     uint64
	wroteHeader bool
}

func (o *resp) Write(p []byte) (bytes int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	bytes, err = o.ResponseWriter.Write(p)
	o.written += uint64(bytes)
	return
}

func (o *resp) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}
