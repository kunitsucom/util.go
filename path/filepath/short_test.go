package filepathz_test

import (
	"testing"

	filepathz "github.com/kunitsucom/util.go/path/filepath"
)

func TestShort(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		if expect, actual := "to/file", filepathz.Short("/path/to/file"); expect != actual {
			t.Errorf("expect(%v) != actual(%v)", expect, actual)
		}
	})
}
