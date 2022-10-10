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

type responseWriterBufferHandler = func(statusCode int, header http.Header, responseWriterBuffer *bytes.Buffer)

type ResponseWriterBufferHandler struct {
	next                        http.Handler
	responseWriterBufferHandler responseWriterBufferHandler
}

type ResponseWriterBufferHandlerOption func(h *ResponseWriterBufferHandler)

func NewResponseWriterBufferHandler(next http.Handler, responseWriterHandler responseWriterBufferHandler, opts ...ResponseWriterBufferHandlerOption) *ResponseWriterBufferHandler {
	h := &ResponseWriterBufferHandler{
		next:                        next,
		responseWriterBufferHandler: responseWriterHandler,
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *ResponseWriterBufferHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rwb := newResponseWriterBuffer(rw)

	h.next.ServeHTTP(rwb, r)

	h.responseWriterBufferHandler(rwb.statusCode, rwb.Header(), rwb.responseWriterBuffer)
}
