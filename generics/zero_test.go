package genericz_test

import (
	"testing"

	genericz "github.com/kunitsucom/util.go/generics"
)

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

func TestSliceContentZero(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		{
			expected := "1"
			if actual := genericz.SliceContentZero([]string{expected}); actual != "" {
				t.Errorf("❌: genericz.SliceContentZero([]string{%q}): %v", expected, actual)
			}
		}
		{
			expected := 1
			if actual := genericz.SliceContentZero([]int{expected}); actual != 0 {
				t.Errorf("❌: genericz.SliceContentZero([]int{%d}): %v", expected, actual)
			}
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
