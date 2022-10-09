package httpz

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

func RequestBodyBuffer(r *http.Request) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)

	if _, err := buf.ReadFrom(r.Body); err != nil {
		return nil, fmt.Errorf("(*bytes.Buffer).ReadFrom: %w", err)
	}

	r.Body = io.NopCloser(bytes.NewBuffer(buf.Bytes()))

	return buf, nil
}

func ContextRequestBodyBuffer(ctx context.Context) (body *bytes.Buffer, ok bool) {
	v, ok := ctx.Value(keyRequestBodyBuffer).(*bytes.Buffer)
	return v, ok
}

func ContextWithRequestBodyBuffer(parent context.Context, body *bytes.Buffer) context.Context {
	return context.WithValue(parent, keyRequestBodyBuffer, body)
}

type RequestBodyBufferHandler struct {
	next               http.Handler
	bufferingSkipLimit int64
	errorHandler       func(rw http.ResponseWriter, r *http.Request, err error)
}

type RequestBodyBufferHandlerOption func(h *RequestBodyBufferHandler)

func WithBufferingSkipLimit(bufferingSkipLimit int64) RequestBodyBufferHandlerOption {
	return func(h *RequestBodyBufferHandler) {
		h.bufferingSkipLimit = bufferingSkipLimit
	}
}

const DefaultBufferingSkipLimit = 1 << 20 // 1 MiB

func NewRequestBodyBufferHandler(next http.Handler, errorHandler func(rw http.ResponseWriter, r *http.Request, err error), opts ...RequestBodyBufferHandlerOption) *RequestBodyBufferHandler {
	h := &RequestBodyBufferHandler{
		next:               next,
		bufferingSkipLimit: DefaultBufferingSkipLimit,
		errorHandler:       errorHandler,
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *RequestBodyBufferHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.ContentLength < h.bufferingSkipLimit {
		buf, err := RequestBodyBuffer(r)
		if err != nil {
			h.errorHandler(rw, r, err)
			return
		}

		h.next.ServeHTTP(rw, r.WithContext(ContextWithRequestBodyBuffer(r.Context(), buf)))
		return
	}

	h.next.ServeHTTP(rw, r)
}
