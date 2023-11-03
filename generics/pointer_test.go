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
