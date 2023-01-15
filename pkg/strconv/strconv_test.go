package strconvz_test

import (
	"strconv"
	"testing"

	strconvz "github.com/kunitsuinc/util.go/pkg/strconv"
)

func TestAtoi64(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		const expect = 1
		actual, err := strconvz.Atoi64(strconv.Itoa(expect))
		if err != nil {
			t.Errorf("❌: util.Atoi64: %v", err)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual")
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		const expect = 0
		actual, err := strconvz.Atoi64("failure")
		if err == nil {
			t.Errorf("❌: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual")
		}
	})
}

func TestItoa64(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		const expect = "100000000000"
		actual := strconvz.Itoa64(100000000000)
		if expect != actual {
			t.Errorf("❌: expect != actual: %s != %s", expect, actual)
		}
	})
}
