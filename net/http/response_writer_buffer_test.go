package httpz_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	httpz "github.com/kunitsucom/util.go/net/http"
	testz "github.com/kunitsucom/util.go/test"
)

func ResponseWriterBufferHandlerTestOption(_ *httpz.ResponseWriterBufferHandler) {
	// noop
}

func TestResponseWriterBufferHandler(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		expect := "200 map[Test-Header:[TestString]] test_request_body http://util.go/net/httpz"
		var actual string
		actualResponse := &httptest.ResponseRecorder{}

		middleware := httpz.NewResponseWriterBufferHandler(
			func(rwb *httpz.ResponseWriterBuffer, r *http.Request) {
				actual = fmt.Sprintf("%d %v %s %s", rwb.StatusCode(), rwb.Header(), rwb.Buffer, r.RequestURI)
			},
		).Middleware

		middleware(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Test-Header", "TestString")
			_, _ = io.Copy(rw, r.Body)
		})).ServeHTTP(
			actualResponse,
			httptest.NewRequest(http.MethodPost, "http://util.go/net/httpz", bytes.NewBufferString("test_request_body")),
		)

		if expect != actual {
			t.Errorf("❌: expect != actual: %s", actual)
		}
	})

	t.Run("failure()", func(t *testing.T) {
		t.Parallel()

		expect := "200 map[Test-Header:[TestString]] test_request_body http://util.go/net/httpz"
		var actual string

		middleware := httpz.NewResponseWriterBufferHandler(
			func(rwb *httpz.ResponseWriterBuffer, r *http.Request) {
				actual = fmt.Sprintf("%d %v %s %s", rwb.StatusCode(), rwb.Header(), rwb.Buffer, r.RequestURI)
			},
			ResponseWriterBufferHandlerTestOption,
		).Middleware

		header := http.Header{}
		middleware(http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				rw.WriteHeader(http.StatusOK)
				rw.Header().Set("Test-Header", "TestString")
				_, _ = io.Copy(rw, r.Body)
			},
		)).ServeHTTP(
			&testz.ResponseWriter{
				WriteFunc: func(p []byte) (n int, err error) {
					return 0, testz.ErrTestError
				},
				HeaderFunc: func() http.Header {
					return header
				},
				WriteHeaderFunc: func(statusCode int) {
					// do nothing
				},
			},
			httptest.NewRequest(http.MethodPost, "http://util.go/net/httpz", bytes.NewBufferString("test_request_body")),
		)

		if expect != actual {
			t.Errorf("❌: expect(%s) != actual(%s)", expect, actual)
		}
	})
}
