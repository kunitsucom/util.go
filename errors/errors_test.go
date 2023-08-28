package errorz_test

import (
	"testing"

	errorz "github.com/kunitsucom/util.go/errors"
	testingz "github.com/kunitsucom/util.go/testing"
)

func TestContains(t *testing.T) {
	t.Parallel()
	t.Run("success,nil", func(t *testing.T) {
		t.Parallel()
		err := (error)(nil)
		if errorz.Contains(err, "testingz: test error") {
			t.Errorf("❌: err not contain %s: %v", "testingz: test error", err)
		}
	})

	t.Run("success,testingz.ErrTestError", func(t *testing.T) {
		t.Parallel()
		err := testingz.ErrTestError
		if !errorz.Contains(err, "testingz: test error") {
			t.Errorf("❌: err not contain `%s`: %v", "testingz: test error", err)
		}
	})
}

func TestHasPrefix(t *testing.T) {
	t.Parallel()
	t.Run("success,nil", func(t *testing.T) {
		t.Parallel()
		err := (error)(nil)
		if errorz.HasPrefix(err, "testingz: test error") {
			t.Errorf("❌: err not has prefix %s: %v", "testingz: test error", err)
		}
	})

	t.Run("success,testingz.ErrTestError", func(t *testing.T) {
		t.Parallel()
		err := testingz.ErrTestError
		if !errorz.HasPrefix(err, "testingz: test error") {
			t.Errorf("❌: err not has prefix `%s`: %v", "testingz: test error", err)
		}
	})
}

func TestHasSuffix(t *testing.T) {
	t.Parallel()
	t.Run("success,nil", func(t *testing.T) {
		t.Parallel()
		err := (error)(nil)
		if errorz.HasSuffix(err, "testingz: test error") {
			t.Errorf("❌: err not has suffix %s: %v", "testingz: test error", err)
		}
	})

	t.Run("success,testingz.ErrTestError", func(t *testing.T) {
		t.Parallel()
		err := testingz.ErrTestError
		if !errorz.HasSuffix(err, "testingz: test error") {
			t.Errorf("❌: err not has suffix `%s`: %v", "testingz: test error", err)
		}
	})
}

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
	})
}
