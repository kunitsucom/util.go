package require

import (
	"testing"

	"github.com/kunitsucom/util.go/testing/internal"
)

// NoError asserts that err is nil.
func NoError(tb testing.TB, err error) bool {
	tb.Helper()

	return internal.NoError(tb, tb.Fatalf, err)
}

// Error asserts that err is not nil.
func Error(tb testing.TB, err error) bool {
	tb.Helper()

	return internal.Error(tb, tb.Fatalf, err)
}

// True asserts that value is true.
func True(tb testing.TB, value bool) bool {
	tb.Helper()

	return internal.True(tb, tb.Fatalf, value)
}

// Equal asserts that expected and actual are deeply equal.
func Equal(tb testing.TB, expected, actual interface{}) bool {
	tb.Helper()

	return internal.Equal(tb, tb.Fatalf, expected, actual)
}

// NotEqual asserts that expected and actual are not deeply equal.
func NotEqual(tb testing.TB, expected, actual interface{}) bool {
	tb.Helper()

	return internal.NotEqual(tb, tb.Fatalf, expected, actual)
}
