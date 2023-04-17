package main

import "net/http"

type clientWriter struct {
	http.ResponseWriter
	statusCode      int
	bytesWritten    uint64
	isHeaderWritten bool
}

func (c *clientWriter) Write(p []byte) (bytes int, err error) {
	if !o.isHeaderWritten {
		o.WriteHeader(http.StatusOK)
	}
	bytes, err = o.ResponseWriter.Write(p)
	o.bytesWritten += uint64(bytes)
	return
}

func (c *clientWriter) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.isHeaderWritten {
		return
	}
	o.isHeaderWritten = true
	o.statusCode = code
}
