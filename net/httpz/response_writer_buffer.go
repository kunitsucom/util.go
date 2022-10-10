package httpz

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type responseWriterBuffer struct {
	http.ResponseWriter
	statusCode           int
	responseWriterBuffer *bytes.Buffer
}

func newResponseWriterBuffer(rw http.ResponseWriter) *responseWriterBuffer {
	return &responseWriterBuffer{
		ResponseWriter:       rw,
		responseWriterBuffer: bytes.NewBuffer(nil),
	}
}

func (rwb *responseWriterBuffer) WriteHeader(status int) {
	rwb.statusCode = status
	rwb.ResponseWriter.WriteHeader(status)
}

func (rwb *responseWriterBuffer) Write(p []byte) (int, error) {
	n, err := io.MultiWriter(rwb.responseWriterBuffer, rwb.ResponseWriter).Write(p)
	if err != nil {
		return n, fmt.Errorf("io.MultiWriter().Write: %w", err)
	}

	return n, nil
}

type ResponseWriterBufferHandler struct {
	responseWriterBufferHandler func(statusCode int, header http.Header, responseWriterBuffer *bytes.Buffer)
}

type ResponseWriterBufferHandlerOption func(h *ResponseWriterBufferHandler)

func NewResponseWriterBufferHandler(responseWriterBufferHandler func(statusCode int, header http.Header, responseWriterBuffer *bytes.Buffer), opts ...ResponseWriterBufferHandlerOption) *ResponseWriterBufferHandler {
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

		h.responseWriterBufferHandler(rwb.statusCode, rwb.Header(), rwb.responseWriterBuffer)
	})
}
