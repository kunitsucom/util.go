package stringz_test

import (
	"testing"

	stringz "github.com/kunitsucom/util.go/strings"
	testingz "github.com/kunitsucom/util.go/testing"
)

func TestJoin(t *testing.T) {
	t.Parallel()

	t.Run("success,normal", func(t *testing.T) {
		t.Parallel()

		const expect = "a_ _b_ _c"
		actual := stringz.Join("_ _", "a", "b", "c")

		if expect != actual {
			t.Errorf("❌: expect(%s) != actual(%s)", expect, actual)
		}
	})
}

func TestJoinStringers(t *testing.T) {
	t.Parallel()

	t.Run("success,normal", func(t *testing.T) {
		t.Parallel()

		const expect = "a_ _b_ _c"
		actual := stringz.JoinStringers("_ _", &testingz.Stringer{func() string { return "a" }}, &testingz.Stringer{func() string { return "b" }}, &testingz.Stringer{func() string { return "c" }})

		if expect != actual {
			t.Errorf("❌: expect(%s) != actual(%s)", expect, actual)
		}
	})
}
