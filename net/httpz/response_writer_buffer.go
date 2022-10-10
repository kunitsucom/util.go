package httpz

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type ResponseWriterBuffer struct {
	http.ResponseWriter
	StatusCode int
	Buffer     *bytes.Buffer
}

func newResponseWriterBuffer(rw http.ResponseWriter) *ResponseWriterBuffer {
	return &ResponseWriterBuffer{
		ResponseWriter: rw,
		Buffer:         bytes.NewBuffer(nil),
	}
}

func (rwb *ResponseWriterBuffer) WriteHeader(status int) {
	rwb.StatusCode = status
	rwb.ResponseWriter.WriteHeader(status)
}

func (rwb *ResponseWriterBuffer) Write(p []byte) (int, error) {
	n, err := io.MultiWriter(rwb.Buffer, rwb.ResponseWriter).Write(p)
	if err != nil {
		return n, fmt.Errorf("io.MultiWriter().Write: %w", err)
	}

	return n, nil
}

type ResponseWriterBufferHandler struct {
	responseWriterBufferHandler func(rwb *ResponseWriterBuffer, r *http.Request)
}

type ResponseWriterBufferHandlerOption func(h *ResponseWriterBufferHandler)

func NewResponseWriterBufferHandler(responseWriterBufferHandler func(rwb *ResponseWriterBuffer, r *http.Request), opts ...ResponseWriterBufferHandlerOption) *ResponseWriterBufferHandler {
	h := &ResponseWriterBufferHandler{
		responseWriterBufferHandler: responseWriterBufferHandler,
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *ResponseWriterBufferHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rwb := newResponseWriterBuffer(rw)

		next.ServeHTTP(rwb, r)

		h.responseWriterBufferHandler(rwb, r)
	})
}
