package stringz_test

import (
	"testing"

	stringz "github.com/kunitsucom/util.go/strings"
)

func TestInsertPrefix(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		const before = `a
b
c

`
		const after = `prefix a
prefix b
prefix c
prefix 
`

		if expected, actual := after, stringz.AddPrefix("prefix ", before, "\n"); expected != actual {
			t.Errorf("❌: expected != actual:\n---EXPECTED\n%s\n+++ACTUAL\n%s", expected, actual)
		}
	})

	t.Run("success,EmptyPrefix", func(t *testing.T) {
		t.Parallel()

		const before = `a`
		const after = `a`
		if expected, actual := after, stringz.AddPrefix("", before, "\n"); expected != actual {
			t.Errorf("❌: expected != actual:\n---EXPECTED\n%s\n+++ACTUAL\n%s", expected, actual)
		}
	})

	t.Run("success,EmptyS", func(t *testing.T) {
		t.Parallel()

		const before = ``
		const after = `prefix `
		if expected, actual := after, stringz.AddPrefix("prefix ", before, "\n"); expected != actual {
			t.Errorf("❌: expected != actual:\n---EXPECTED\n%s\n+++ACTUAL\n%s", expected, actual)
		}
	})

	t.Run("success,EmptySep", func(t *testing.T) {
		t.Parallel()

		const before = `a
b
c`
		const after = `prefix a
b
c`

		if expected, actual := after, stringz.AddPrefix("prefix ", before, ""); expected != actual {
			t.Errorf("❌: expected != actual:\n---EXPECTED\n%s\n+++ACTUAL\n%s", expected, actual)
		}
	})
}
