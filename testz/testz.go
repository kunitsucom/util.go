package testz

import (
	"errors"
	"io"
	"net/http"
)

var ErrTestError = errors.New("testz: test error")

var _ io.ReadWriter = (*ReadWriter)(nil)

type ReadWriter struct {
	N   int
	Err error
}

func NewReadWriter(n int, err error) *ReadWriter {
	return &ReadWriter{
		N:   n,
		Err: err,
	}
}

func (t *ReadWriter) Read(p []byte) (n int, err error) {
	return t.N, t.Err
}

func (t *ReadWriter) Write(p []byte) (n int, err error) {
	return t.N, t.Err
}

var _ http.ResponseWriter = (*ResponseWriter)(nil)

type ResponseWriter struct {
	header http.Header

	N   int
	Err error

	StatusCode int
}

func NewResponseWriter(initial http.Header, n int, err error) *ResponseWriter {
	header := make(map[string][]string)

	for key, value := range initial {
		header[key] = value
	}

	return &ResponseWriter{
		header: header,
		N:      n,
		Err:    err,
	}
}

func (t *ResponseWriter) Header() http.Header {
	return t.header
}

func (t *ResponseWriter) Write(p []byte) (n int, err error) {
	return t.N, t.Err
}

func (t *ResponseWriter) WriteHeader(statusCode int) {
	t.StatusCode = statusCode
}
