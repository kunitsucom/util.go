package errorz_test

import (
	"errors"
	"testing"

	errorz "github.com/kunitsucom/util.go/errors"
	testingz "github.com/kunitsucom/util.go/testing"
)

func TestIsRetryable(t *testing.T) {
	t.Parallel()
	t.Run("success,nil", func(t *testing.T) {
		t.Parallel()
		err := (error)(nil)
		if errorz.IsRetryable(err) {
			t.Errorf("❌: err is retryable: %v", err)
		}
	})

	t.Run("success,testingz.ErrTestError", func(t *testing.T) {
		t.Parallel()
		err := testingz.ErrTestError
		if errorz.IsRetryable(err) {
			t.Errorf("❌: err is retryable: %v", err)
		}
	})

	t.Run("success,IsRetryable", func(t *testing.T) {
		t.Parallel()
		err := errorz.WithRetryable(testingz.ErrTestError, true)
		if !errorz.IsRetryable(err) {
			t.Errorf("❌: err is not retryable: %v", err)
		}

		var e interface{ Unwrap() error }
		if !errors.As(err, &e) {
			t.Errorf("❌: err is not Unwrap: %v", err)
		}
		if !errors.Is(err, testingz.ErrTestError) {
			t.Errorf("❌: err is retryable: %v", err)
		}
	})
}
