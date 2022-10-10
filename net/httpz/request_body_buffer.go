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
	bufferingSkipLimit int64
	errorHandler       func(rw http.ResponseWriter, r *http.Request, err error)
}

type RequestBodyBufferHandlerOption func(h *RequestBodyBufferHandler)

func WithRequestBodyBufferingSkipLimit(bufferingSkipLimit int64) RequestBodyBufferHandlerOption {
	return func(h *RequestBodyBufferHandler) {
		h.bufferingSkipLimit = bufferingSkipLimit
	}
}

const DefaultRequestBodyBufferingSkipLimit = 1 << 20 // 1 MiB

func NewRequestBodyBufferHandler(errorHandler func(rw http.ResponseWriter, r *http.Request, err error), opts ...RequestBodyBufferHandlerOption) func(http.Handler) http.Handler {
	h := &RequestBodyBufferHandler{
		bufferingSkipLimit: DefaultRequestBodyBufferingSkipLimit,
		errorHandler:       errorHandler,
	}

	for _, opt := range opts {
		opt(h)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if r.ContentLength < h.bufferingSkipLimit {
				buf, err := RequestBodyBuffer(r)
				if err != nil {
					h.errorHandler(rw, r, err)
					return
				}

				next.ServeHTTP(rw, r.WithContext(ContextWithRequestBodyBuffer(r.Context(), buf)))
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}
