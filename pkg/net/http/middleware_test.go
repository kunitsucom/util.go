package httpz_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	httpz "github.com/kunitsuinc/util.go/pkg/net/http"
	testz "github.com/kunitsuinc/util.go/pkg/test"
)

func testMiddleware(num int) func(http.Handler) http.Handler {
	return func(original http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if _, err := rw.Write([]byte(fmt.Sprintf("middleware %d preProcess\n", num))); err != nil {
				panic(fmt.Errorf("rw.Write: %w", err))
			}

			original.ServeHTTP(rw, r)

			if _, err := rw.Write([]byte(fmt.Sprintf("middleware %d postProcess\n", num))); err != nil {
				panic(fmt.Errorf("rw.Write: %w", err))
			}
		})
	}
}

func TestMiddlewares(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mids := httpz.Middlewares(testMiddleware(1), testMiddleware(2), testMiddleware(3)).
			Middlewares(testMiddleware(4))
		mids = httpz.Middlewares(mids, testMiddleware(5))

		expect := "middleware 5 preProcess\nmiddleware 4 preProcess\nmiddleware 3 preProcess\nmiddleware 2 preProcess\nmiddleware 1 preProcess\ntest_request_body\nmiddleware 1 postProcess\nmiddleware 2 postProcess\nmiddleware 3 postProcess\nmiddleware 4 postProcess\nmiddleware 5 postProcess\n"
		request := bytes.NewBufferString("test_request_body\n")
		r := httptest.NewRequest(http.MethodPost, "http://util.go/net/httpz", request)
		hander := mids(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = io.Copy(w, r.Body)
		}))

		response := bytes.NewBuffer(nil)
		hander.ServeHTTP(testz.NewResponseWriter(response, nil, 0, nil), r)

		actual := response.String()
		if expect != actual {
			t.Errorf("❌: expect != actual:\n-%s\n+%s", expect, actual)
		}
	})
}
