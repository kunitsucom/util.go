package require

import (
	"testing"

	"github.com/kunitsucom/util.go/testing/internal"
)

// NoError asserts that err is nil.
func NoError(tb testing.TB, err error) {
	tb.Helper()

	internal.NoError(tb, tb.Fatalf, err)
}

// Error asserts that err is not nil.
func Error(tb testing.TB, err error) {
	tb.Helper()

	internal.Error(tb, tb.Fatalf, err)
}

// True asserts that value is true.
func True(tb testing.TB, value bool) {
	tb.Helper()

	internal.True(tb, tb.Fatalf, value)
}

// Equal asserts that expected and actual are deeply equal.
func Equal(tb testing.TB, expected, actual interface{}) {
	tb.Helper()

	internal.Equal(tb, tb.Fatalf, expected, actual)
}

// NotEqual asserts that expected and actual are not deeply equal.
func NotEqual(tb testing.TB, expected, actual interface{}) {
	tb.Helper()

	internal.NotEqual(tb, tb.Fatalf, expected, actual)
}
