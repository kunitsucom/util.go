package httpz_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	httpz "github.com/kunitsucom/util.go/net/http"
	testingz "github.com/kunitsucom/util.go/testing"
)

func TestRequestBodyBuffer(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		expect := "test_request_body"
		r := httptest.NewRequest(http.MethodPost, "http://util.go/net/httpz", bytes.NewBufferString(expect))
		buf, err := httpz.RequestBodyBuffer(r)
		if err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
		actual := buf.String()
		if expect != actual {
			t.Errorf("❌: expect != actual: %v", actual)
		}
	})

	t.Run("failure()", func(t *testing.T) {
		t.Parallel()

		expect := (*bytes.Buffer)(nil)
		r := httptest.NewRequest(http.MethodPost, "http://util.go/net/httpz", &testingz.ReadWriter{
			ReadFunc: func(p []byte) (n int, err error) {
				return 0, testingz.ErrTestError
			},
			WriteFunc: func(p []byte) (n int, err error) {
				return 0, testingz.ErrTestError
			},
		})
		actual, err := httpz.RequestBodyBuffer(r)
		if !errors.Is(err, testingz.ErrTestError) {
			t.Errorf("❌: err != nil: %v", err)
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v", actual)
		}
	})
}

func TestNewRequestBodyBufferHandler(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		expect := "test_request_body"
		var actualBuf *bytes.Buffer
		var ok bool
		actualResponse := &httptest.ResponseRecorder{}

		middleware := httpz.NewRequestBodyBufferHandler(
			func(rw http.ResponseWriter, r *http.Request, err error) {
				rw.WriteHeader(http.StatusInternalServerError)
				_, _ = rw.Write([]byte(err.Error()))
			},
		).Middleware

		r := httptest.NewRequest(http.MethodPost, "http://util.go/net/httpz", bytes.NewBufferString(expect))
		middleware(http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				actualBuf, ok = httpz.ContextRequestBodyBuffer(r.Context())
			})).
			ServeHTTP(actualResponse, r)
		if !ok {
			t.Errorf("❌: !ok")
		}

		actual := actualBuf.String()
		if expect != actual {
			t.Errorf("❌: expect != actual: %s", actual)
		}
	})

	t.Run("success(WithBufferingSkipLimit)", func(t *testing.T) {
		t.Parallel()

		expect := "<nil>"
		var actualBuf *bytes.Buffer
		var ok bool
		actualResponse := &httptest.ResponseRecorder{}

		middleware := httpz.NewRequestBodyBufferHandler(
			func(rw http.ResponseWriter, r *http.Request, err error) {
				rw.WriteHeader(http.StatusInternalServerError)
				_, _ = rw.Write([]byte(err.Error()))
			},
			httpz.WithRequestBodyBufferingSkipLimit(1),
		).Middleware

		r := httptest.NewRequest(http.MethodPost, "http://util.go/net/httpz", bytes.NewBufferString("over_limit_string"))
		middleware((http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				actualBuf, ok = httpz.ContextRequestBodyBuffer(r.Context())
			}))).
			ServeHTTP(actualResponse, r)

		if ok {
			t.Errorf("❌: ok")
		}

		actual := actualBuf.String()
		if expect != actual {
			t.Errorf("❌: expect != actual: %s", actual)
		}
	})

	t.Run("failure(errorHandler)", func(t *testing.T) {
		t.Parallel()

		expect := "<nil>"
		var actualBuf *bytes.Buffer
		var ok bool
		actualResponse := &httptest.ResponseRecorder{}

		middleware := httpz.NewRequestBodyBufferHandler(
			func(rw http.ResponseWriter, r *http.Request, err error) {
				rw.WriteHeader(http.StatusInternalServerError)
				_, _ = rw.Write([]byte(err.Error()))
			},
			httpz.WithRequestBodyBufferingSkipLimit(100),
		).Middleware

		r := httptest.NewRequest(http.MethodPost, "http://util.go/net/httpz", &testingz.ReadWriter{
			ReadFunc: func(p []byte) (n int, err error) {
				return 0, testingz.ErrTestError
			},
			WriteFunc: func(p []byte) (n int, err error) {
				return 0, testingz.ErrTestError
			},
		})
		middleware(http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				actualBuf, ok = httpz.ContextRequestBodyBuffer(r.Context())
			})).
			ServeHTTP(actualResponse, r)

		if ok {
			t.Errorf("❌: ok")
		}

		actual := actualBuf.String()
		if expect != actual {
			t.Errorf("❌: expect != actual: %s", actual)
		}
	})
}
