package httpz_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/kunitsuinc/util.go/net/httpz"
	"github.com/kunitsuinc/util.go/testz"
)

func TestRequestBodyBuffer(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		expect := "test_request_body"
		r := httptest.NewRequest(http.MethodPost, "http://util.go/net/httpz", bytes.NewBufferString(expect))
		buf, err := httpz.RequestBodyBuffer(r)
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		actual := buf.String()
		if expect != actual {
			t.Errorf("expect != actual: %v", actual)
		}
	})

	t.Run("failure()", func(t *testing.T) {
		t.Parallel()

		expect := (*bytes.Buffer)(nil)
		r := httptest.NewRequest(http.MethodPost, "http://util.go/net/httpz", testz.NewReadWriter(0, testz.ErrTestError))
		actual, err := httpz.RequestBodyBuffer(r)
		if !errors.Is(err, testz.ErrTestError) {
			t.Errorf("err != nil: %v", err)
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v", actual)
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

		h := httpz.NewRequestBodyBufferHandler(http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				actualBuf, ok = httpz.ContextRequestBodyBuffer(r.Context())
			}),
			func(rw http.ResponseWriter, r *http.Request, err error) {
				rw.WriteHeader(http.StatusInternalServerError)
				_, _ = rw.Write([]byte(err.Error()))
			},
		)

		r := httptest.NewRequest(http.MethodPost, "http://util.go/net/httpz", bytes.NewBufferString(expect))
		h.ServeHTTP(actualResponse, r)

		if !ok {
			t.Errorf("!ok")
		}

		actual := actualBuf.String()
		if expect != actual {
			t.Errorf("expect != actual: %s", actual)
		}
	})

	t.Run("success(WithBufferingSkipLimit)", func(t *testing.T) {
		t.Parallel()

		expect := "<nil>"
		var actualBuf *bytes.Buffer
		var ok bool
		actualResponse := &httptest.ResponseRecorder{}

		h := httpz.NewRequestBodyBufferHandler(http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				actualBuf, ok = httpz.ContextRequestBodyBuffer(r.Context())
			}),
			func(rw http.ResponseWriter, r *http.Request, err error) {
				rw.WriteHeader(http.StatusInternalServerError)
				_, _ = rw.Write([]byte(err.Error()))
			},
			httpz.WithRequestBodyBufferingSkipLimit(1),
		)

		r := httptest.NewRequest(http.MethodPost, "http://util.go/net/httpz", bytes.NewBufferString("over_limit_string"))
		h.ServeHTTP(actualResponse, r)

		if ok {
			t.Errorf("ok")
		}

		actual := actualBuf.String()
		if expect != actual {
			t.Errorf("expect != actual: %s", actual)
		}
	})

	t.Run("failure(errorHandler)", func(t *testing.T) {
		t.Parallel()

		expect := "<nil>"
		var actualBuf *bytes.Buffer
		var ok bool
		actualResponse := &httptest.ResponseRecorder{}

		h := httpz.NewRequestBodyBufferHandler(http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				actualBuf, ok = httpz.ContextRequestBodyBuffer(r.Context())
			}),
			func(rw http.ResponseWriter, r *http.Request, err error) {
				rw.WriteHeader(http.StatusInternalServerError)
				_, _ = rw.Write([]byte(err.Error()))
			},
			httpz.WithRequestBodyBufferingSkipLimit(100),
		)

		r := httptest.NewRequest(http.MethodPost, "http://util.go/net/httpz", testz.NewReadWriter(0, testz.ErrTestError))
		h.ServeHTTP(actualResponse, r)

		if ok {
			t.Errorf("ok")
		}

		actual := actualBuf.String()
		if expect != actual {
			t.Errorf("expect != actual: %s", actual)
		}
	})
}
