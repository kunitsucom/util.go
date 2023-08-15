package testz

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

var ErrTestError = errors.New("testz: test error")

var _ io.ReadWriter = (*ReadWriter)(nil)

type ReadWriter struct {
	Buffer *bytes.Buffer

	N   int
	Err error
}

func NewReadWriter(buf *bytes.Buffer, n int, err error) *ReadWriter {
	return &ReadWriter{
		Buffer: buf,
		N:      n,
		Err:    err,
	}
}

func (t *ReadWriter) Read(p []byte) (n int, err error) {
	if t.Err != nil {
		return t.N, t.Err
	}

	//nolint:wrapcheck
	return t.Buffer.Read(p)
}

func (t *ReadWriter) Write(p []byte) (n int, err error) {
	if t.Err != nil {
		return t.N, t.Err
	}

	//nolint:wrapcheck
	return t.Buffer.Write(p)
}

var _ http.ResponseWriter = (*ResponseWriter)(nil)

type ResponseWriter struct {
	Buffer *bytes.Buffer

	header http.Header

	N   int
	Err error

	StatusCode int
}

func NewResponseWriter(buf *bytes.Buffer, initial http.Header, n int, err error) *ResponseWriter {
	header := make(map[string][]string)

	for key, value := range initial {
		header[key] = value
	}

	return &ResponseWriter{
		Buffer: buf,
		header: header,
		N:      n,
		Err:    err,
	}
}

func (t *ResponseWriter) Header() http.Header {
	return t.header
}

func (t *ResponseWriter) Write(p []byte) (n int, err error) {
	if t.Err != nil {
		return t.N, t.Err
	}

	//nolint:wrapcheck
	return t.Buffer.Write(p)
}

func (t *ResponseWriter) WriteHeader(statusCode int) {
	t.StatusCode = statusCode
}
