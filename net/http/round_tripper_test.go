package httpz_test

import (
	"net/http"
	"testing"

	httpz "github.com/kunitsucom/util.go/net/http"
	testingz "github.com/kunitsucom/util.go/testing"
	"github.com/kunitsucom/util.go/testing/assert"
)

func TestRoundTripFunc_RoundTrip(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		a := httpz.RoundTripFunc(func(_ *http.Request) (*http.Response, error) {
			return nil, testingz.ErrTestError
		})

		_, err := a.RoundTrip(nil) //nolint:bodyclose
		assert.ErrorIs(t, err, testingz.ErrTestError)
	})
}
