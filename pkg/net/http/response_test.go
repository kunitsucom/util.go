package httpz_test

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	httpz "github.com/kunitsuinc/util.go/pkg/net/http"
	testz "github.com/kunitsuinc/util.go/pkg/test"
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
			Body: io.NopCloser(testz.NewReadWriter(bytes.NewBuffer(nil), 0, testz.ErrTestError)),
		})
		if err == nil {
			t.Errorf("❌: httpz.NewBufferFromResponseBody: err == nil")
		}
	})
}
