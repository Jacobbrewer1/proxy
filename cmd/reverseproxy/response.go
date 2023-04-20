package main

import "net/http"

type clientWriter struct {
	http.ResponseWriter
	statusCode      int
	bytesWritten    uint64
	isHeaderWritten bool
}

func (c *clientWriter) Write(p []byte) (bytes int, err error) {
	if !c.isHeaderWritten {
		c.WriteHeader(http.StatusOK)
	}
	bytes, err = c.ResponseWriter.Write(p)
	c.bytesWritten += uint64(bytes)
	return
}

func (c *clientWriter) WriteHeader(code int) {
	c.ResponseWriter.WriteHeader(code)
	if c.isHeaderWritten {
		return
	}
	c.isHeaderWritten = true
	c.statusCode = code
}
