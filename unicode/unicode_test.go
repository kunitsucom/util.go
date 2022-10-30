package unicodez_test

import (
	"testing"

	unicodez "github.com/kunitsuinc/util.go/unicode"
)

func TestTrimNonGraphic(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := "expect"
		actual := unicodez.TrimNonGraphic(string('\u0000') + "e" + string('\u0001') + "x" + string('\u0002') + "p" + string('\u0003') + "e" + string('\u0004') + "c" + string('\u0005') + "t" + string('\u0006'))
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		expect := "expect"
		actual := unicodez.TrimNonGraphic(expect)
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}
