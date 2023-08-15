package syncz_test

import (
	"io"
	"testing"

	syncz "github.com/kunitsucom/util.go/sync"
)

func TestOnce_Do(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		var once syncz.Once
		actual := 0
		const expect1 = 1
		for i := 0; i <= 10; i++ {
			if err := once.Do(func() error {
				actual++
				return nil
			}); err != nil {
				t.Errorf("❌: err != nil: %v", err)
			}
		}
		if actual != expect1 {
			t.Errorf("❌: actual != expect: %v != %v", actual, expect1)
		}

		once.Reset()
		const expect2 = 2
		for i := 0; i <= 10; i++ {
			if err := once.Do(func() error {
				actual++
				return nil
			}); err != nil {
				t.Errorf("❌: err != nil: %v", err)
			}
		}
		if actual != expect2 {
			t.Errorf("❌: actual != expect: %v != %v", actual, expect1)
		}
	})

	t.Run("failure()", func(t *testing.T) {
		t.Parallel()
		var once syncz.Once
		actual := 0
		const expect = 10
		for i := 1; i <= 10; i++ {
			i := i
			if err := once.Do(func() error {
				actual = i
				return io.EOF // any error
			}); err == nil {
				t.Errorf("❌: err == nil")
			}
		}
		if actual != expect {
			t.Errorf("❌: actual != expect: %v != %v", actual, expect)
		}
	})
}
