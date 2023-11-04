package assert

import (
	"testing"

	"github.com/kunitsucom/util.go/testing/internal"
)

// NoError asserts that err is nil.
func NoError(tb testing.TB, err error) (success bool) {
	tb.Helper()

	return internal.NoError(tb, tb.Errorf, err)
}

// Error asserts that err is not nil.
func Error(tb testing.TB, err error) (success bool) {
	tb.Helper()

	return internal.Error(tb, tb.Errorf, err)
}

// True asserts that value is true.
func True(tb testing.TB, value bool) (success bool) {
	tb.Helper()

	return internal.True(tb, tb.Errorf, value)
}

// Equal asserts that expected and actual are deeply equal.
func Equal(tb testing.TB, expected, actual interface{}) (success bool) {
	tb.Helper()

	return internal.Equal(tb, tb.Errorf, expected, actual)
}

// NotEqual asserts that expected and actual are not deeply equal.
func NotEqual(tb testing.TB, expected, actual interface{}) (success bool) {
	tb.Helper()

	return internal.NotEqual(tb, tb.Errorf, expected, actual)
}
