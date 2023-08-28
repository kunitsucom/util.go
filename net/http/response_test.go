package httpz_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	httpz "github.com/kunitsucom/util.go/net/http"
	testingz "github.com/kunitsucom/util.go/testing"
)

func TestReadResponseBody(t *testing.T) {
	t.Parallel()

	t.Run("normal", func(t *testing.T) {
		t.Parallel()

		const expect = "TestResponseBody"

		body, err := httpz.NewBufferFromResponseBody(&http.Response{
			Body: io.NopCloser(strings.NewReader(expect)),
		})
		if err != nil {
			t.Errorf("❌: httpz.NewBufferFromResponseBody: err != nil: %v", err)
		}

		actual := body.String()
		if actual != expect {
			t.Errorf("❌: expect != actual: %s != %s", expect, actual)
		}
	})

	t.Run("abnormal", func(t *testing.T) {
		t.Parallel()

		_, err := httpz.NewBufferFromResponseBody(&http.Response{
			Body: io.NopCloser(&testingz.ReadWriter{
				ReadFunc: func(p []byte) (n int, err error) {
					return 0, testingz.ErrTestError
				},
			}),
		})
		if !errors.Is(err, testingz.ErrTestError) {
			t.Errorf("❌: httpz.NewBufferFromResponseBody: err == nil")
		}
	})
}
