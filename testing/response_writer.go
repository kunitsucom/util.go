package testingz

import (
	"net/http"
)

var _ http.ResponseWriter = (*ResponseWriter)(nil)

type ResponseWriter struct {
	WriteFunc       func(p []byte) (n int, err error)
	HeaderFunc      func() http.Header
	WriteHeaderFunc func(statusCode int)
}

func (rw *ResponseWriter) Header() http.Header {
	return rw.HeaderFunc()
}

func (rw *ResponseWriter) Write(p []byte) (n int, err error) {
	return rw.WriteFunc(p)
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.WriteHeaderFunc(statusCode)
}
