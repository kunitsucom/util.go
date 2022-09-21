package util_test

import (
	"errors"
	"testing"

	"github.com/kunitsuinc/util.go"
)

var errSomeError = errors.New("some error")

func newInt(i int, err error) (int, error) {
	return i, err
}

func newInt2(i1 int, i2 int, err error) (int, int, error) {
	return i1, i2, err
}

func newInt3(i1 int, i2 int, i3 int, err error) (int, int, int, error) {
	return i1, i2, i3, err
}

func newInt4(i1 int, i2 int, i3 int, i4 int, err error) (int, int, int, int, error) {
	return i1, i2, i3, i4, err
}

func TestMust(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		i := util.Must(newInt(1, nil))

		if i <= 0 {
			t.Errorf("i <= 0")
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("recover: err == nil")
			}
		}()

		_ = util.Must(newInt(1, errSomeError))
	})
}

func TestMust2(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		i1, i2 := util.Must2(newInt2(1, 1, nil))

		if i1 <= 0 || i2 <= 0 {
			t.Errorf("i1 <= 0 || i2 <= 0")
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("recover: err == nil")
			}
		}()

		_, _ = util.Must2(newInt2(1, 1, errSomeError))
	})
}

func TestMust3(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		i1, i2, i3 := util.Must3(newInt3(1, 1, 1, nil))

		if i1 <= 0 || i2 <= 0 || i3 <= 0 {
			t.Errorf("i1 <= 0 || i2 <= 0 || i3 <= 0")
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("recover: err == nil")
			}
		}()

		_, _, _ = util.Must3(newInt3(1, 1, 1, errSomeError))
	})
}

func TestMust4(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		i1, i2, i3, i4 := util.Must4(newInt4(1, 1, 1, 1, nil))

		if i1 <= 0 || i2 <= 0 || i3 <= 0 || i4 <= 0 {
			t.Errorf("i1 <= 0 || i2 <= 0 || i3 <= 0 || i4 <= 0")
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("recover: err == nil")
			}
		}()

		_, _, _, _ = util.Must4(newInt4(1, 1, 1, 1, errSomeError))
	})
}
