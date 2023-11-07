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

// ErrorsIs asserts that err is target.
func ErrorsIs(tb testing.TB, err, target error) (success bool) {
	tb.Helper()

	return internal.ErrorsIs(tb, tb.Errorf, err, target)
}

// ErrorsContains asserts that err contains substr.
func ErrorsContains(tb testing.TB, err error, substr string) (success bool) {
	tb.Helper()

	return internal.ErrorsContains(tb, tb.Errorf, err, substr)
}

// True asserts that value is true.
func True(tb testing.TB, value bool) (success bool) {
	tb.Helper()

	return internal.True(tb, tb.Errorf, value)
}

// False asserts that value is false.
func False(tb testing.TB, value bool) (success bool) {
	tb.Helper()

	return internal.False(tb, tb.Errorf, value)
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

// Nil asserts that value is nil.
func Nil(tb testing.TB, value interface{}) (success bool) {
	tb.Helper()

	return internal.Nil(tb, tb.Errorf, value)
}

// NotNil asserts that value is not nil.
func NotNil(tb testing.TB, value interface{}) (success bool) {
	tb.Helper()

	return internal.NotNil(tb, tb.Errorf, value)
}
