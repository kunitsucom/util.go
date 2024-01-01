package errorz_test

import (
	"io"
	"testing"

	errorz "github.com/kunitsucom/util.go/errors"
)

func TestPanicOrIgnore(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("❌: %s: expected not to panic", t.Name())
			}
		}()

		errorz.PanicOrIgnore(nil, io.EOF)
		errorz.PanicOrIgnore(io.EOF, io.EOF)
	})

	t.Run("failure,", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("❌: %s: expected to panic", t.Name())
			}
		}()

		errorz.PanicOrIgnore(io.EOF)
	})
}
