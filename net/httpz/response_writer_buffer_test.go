package httpz_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kunitsuinc/util.go/net/httpz"
	"github.com/kunitsuinc/util.go/testz"
)

func ResponseWriterBufferHandlerTestOption(h *httpz.ResponseWriterBufferHandler) {
}

func TestResponseWriterBufferHandler(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		expect := "200 map[Test-Header:[TestString]] test_request_body"
		var actual string
		actualResponse := &httptest.ResponseRecorder{}

		middleware := httpz.NewResponseWriterBufferHandler(
			func(statusCode int, header http.Header, responseWriterBuffer *bytes.Buffer) {
				actual = fmt.Sprintf("%d %v %s", statusCode, header, responseWriterBuffer)
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
			t.Errorf("expect != actual: %s", actual)
		}
	})

	t.Run("failure()", func(t *testing.T) {
		t.Parallel()

		expect := "200 map[Test-Header:[TestString]] test_request_body"
		var actual string

		middleware := httpz.NewResponseWriterBufferHandler(
			func(statusCode int, header http.Header, responseWriterBuffer *bytes.Buffer) {
				actual = fmt.Sprintf("%d %v %s", statusCode, header, responseWriterBuffer)
			},
			ResponseWriterBufferHandlerTestOption,
		).Middleware

		middleware(http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				rw.WriteHeader(http.StatusOK)
				rw.Header().Set("Test-Header", "TestString")
				_, _ = io.Copy(rw, r.Body)
			},
		)).ServeHTTP(
			testz.NewResponseWriter(bytes.NewBuffer(nil), nil, 0, testz.ErrTestError),
			httptest.NewRequest(http.MethodPost, "http://util.go/net/httpz", bytes.NewBufferString("test_request_body")),
		)

		if expect != actual {
			t.Errorf("expect != actual: %s", actual)
		}
	})
}
