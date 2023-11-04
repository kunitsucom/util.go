package internal

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/kunitsucom/util.go/exp/diff/simplediff"
	stringz "github.com/kunitsucom/util.go/strings"
)

func NoError(tb testing.TB, printf func(format string, args ...any), err error) (success bool) {
	tb.Helper()

	if err != nil {
		printf("❌: err != nil: %+v", err)
		return false
	}
	return true
}

func Error(tb testing.TB, printf func(format string, args ...any), err error) (success bool) {
	tb.Helper()

	if err == nil {
		printf("❌: err == nil")
		return false
	}
	return true
}

func ErrorsIs(tb testing.TB, printf func(format string, args ...any), err, target error) (success bool) {
	tb.Helper()

	if !errors.Is(err, target) {
		printf("❌: err != target:\n--- TARGET\n+++ ERROR\n%s\n%s",
			stringz.AddPrefix("-", fmt.Sprintf("%v", target), "\n"), //nolint:perfsprint
			stringz.AddPrefix("+", fmt.Sprintf("%+v", err), "\n"),
		)
		return false
	}
	return true
}

func True(tb testing.TB, printf func(format string, args ...any), value bool) (success bool) {
	tb.Helper()

	if !value {
		printf("❌: value == false")
		return false
	}
	return true
}

func False(tb testing.TB, printf func(format string, args ...any), value bool) (success bool) {
	tb.Helper()

	if value {
		printf("❌: value == true")
		return false
	}
	return true
}

func Equal(tb testing.TB, printf func(format string, args ...any), expected, actual interface{}) (success bool) {
	tb.Helper()

	if !reflect.DeepEqual(expected, actual) {
		printf("❌: expected != actual:\n--- EXPECTED\n+++ ACTUAL\n%s", simplediff.Diff(fmt.Sprintf("%+v", expected), fmt.Sprintf("%+v", actual)))
		return false
	}
	return true
}

func NotEqual(tb testing.TB, printf func(format string, args ...any), expected, actual interface{}) (success bool) {
	tb.Helper()

	if reflect.DeepEqual(expected, actual) {
		printf("❌: expected == actual:\n--- EXPECTED\n+++ ACTUAL\n%s", simplediff.Diff(fmt.Sprintf("%+v", expected), fmt.Sprintf("%+v", actual)))
		return false
	}
	return true
}
