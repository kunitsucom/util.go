package genericz_test

import (
	"testing"

	genericz "github.com/kunitsucom/util.go/generics"
)

func TestPointer(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		actual := genericz.Pointer("test")
		if *actual != "test" {
			t.Errorf("❌: *actual != test: %v", *actual)
		}
	})
}

func TestIsZero(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		if actual := !genericz.IsZero(""); actual {
			t.Errorf("❌: !genericz.IsZero(\"\"): %v", actual)
		}
		if actual := genericz.IsZero("test"); actual {
			t.Errorf("❌: genericz.IsZero(\"test\"): %v", actual)
		}
	})
}

func TestZero(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		if actual := genericz.Zero("1"); actual != "" {
			t.Errorf("❌: genericz.Zero(\"1\"): %v", actual)
		}
		if actual := genericz.Zero(1); actual != 0 {
			t.Errorf("❌: genericz.Zero(1): %v", actual)
		}
	})
}
