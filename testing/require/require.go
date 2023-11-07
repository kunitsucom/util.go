package require

import (
	"testing"

	"github.com/kunitsucom/util.go/testing/internal"
)

// NoError asserts that err is nil.
func NoError(tb testing.TB, err error) (success bool) {
	tb.Helper()

	return internal.NoError(tb, tb.Fatalf, err)
}

// Error asserts that err is not nil.
func Error(tb testing.TB, err error) (success bool) {
	tb.Helper()

	return internal.Error(tb, tb.Fatalf, err)
}

// ErrorsIs asserts that err is target.
func ErrorsIs(tb testing.TB, err, target error) (success bool) {
	tb.Helper()

	return internal.ErrorsIs(tb, tb.Fatalf, err, target)
}

// ErrorsContains asserts that err contains substr.
func ErrorsContains(tb testing.TB, err error, substr string) (success bool) {
	tb.Helper()

	return internal.ErrorsContains(tb, tb.Fatalf, err, substr)
}

// True asserts that value is true.
func True(tb testing.TB, value bool) (success bool) {
	tb.Helper()

	return internal.True(tb, tb.Fatalf, value)
}

// False asserts that value is false.
func False(tb testing.TB, value bool) (success bool) {
	tb.Helper()

	return internal.False(tb, tb.Fatalf, value)
}

// Equal asserts that expected and actual are deeply equal.
func Equal(tb testing.TB, expected, actual interface{}) (success bool) {
	tb.Helper()

	return internal.Equal(tb, tb.Fatalf, expected, actual)
}

// NotEqual asserts that expected and actual are not deeply equal.
func NotEqual(tb testing.TB, expected, actual interface{}) (success bool) {
	tb.Helper()

	return internal.NotEqual(tb, tb.Fatalf, expected, actual)
}

// Nil asserts that value is nil.
func Nil(tb testing.TB, value interface{}) (success bool) {
	tb.Helper()

	return internal.Nil(tb, tb.Fatalf, value)
}

// NotNil asserts that value is not nil.
func NotNil(tb testing.TB, value interface{}) (success bool) {
	tb.Helper()

	return internal.NotNil(tb, tb.Fatalf, value)
}
