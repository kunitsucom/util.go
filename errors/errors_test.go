package errorz_test

import (
	"io"
	"net"
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

type testTimeoutError struct {
	err     error
	timeout bool
}

func (e testTimeoutError) Error() string {
	return e.err.Error()
}

func (e testTimeoutError) Timeout() bool {
	return e.timeout
}

func (e testTimeoutError) Temporary() bool {
	return false
}

func TestIsNetTimeout(t *testing.T) {
	t.Parallel()

	t.Run("success,true", func(t *testing.T) {
		t.Parallel()
		err := testTimeoutError{err: net.ErrClosed, timeout: true}
		if !errorz.IsNetTimeout(err) {
			t.Errorf("❌: err is net timeout: %v", err)
		}
	})

	t.Run("success,false,net.Error", func(t *testing.T) {
		t.Parallel()
		err := testTimeoutError{err: net.ErrClosed, timeout: false}
		if errorz.IsNetTimeout(err) {
			t.Errorf("❌: err is net timeout: %v", err)
		}
	})

	t.Run("success,false,error", func(t *testing.T) {
		t.Parallel()
		err := io.EOF
		if errorz.IsNetTimeout(err) {
			t.Errorf("❌: err is not net timeout: %v", err)
		}
	})
}
