package errorz_test

import (
	"testing"

	errorz "github.com/kunitsucom/util.go/errors"
	testingz "github.com/kunitsucom/util.go/testing"
)

func TestContains(t *testing.T) {
	t.Parallel()
	t.Run("success(nil)", func(t *testing.T) {
		t.Parallel()
		err := (error)(nil)
		if errorz.Contains(err, "testingz: test error") {
			t.Errorf("❌: err not contain %s: %v", "testingz: test error", err)
		}
	})

	t.Run("success(testingz.ErrTestError)", func(t *testing.T) {
		t.Parallel()
		err := testingz.ErrTestError
		if !errorz.Contains(err, "testingz: test error") {
			t.Errorf("❌: err not contain `%s`: %v", "testingz: test error", err)
		}
	})
}
