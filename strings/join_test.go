package stringz_test

import (
	"testing"

	stringz "github.com/kunitsucom/util.go/strings"
)

func TestJoin(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		const expect = "a_ _b_ _c"
		actual := stringz.Join("_ _", "a", "b", "c")

		if expect != actual {
			t.Errorf("❌: expect(%s) != actual(%s)", expect, actual)
		}
	})
}
