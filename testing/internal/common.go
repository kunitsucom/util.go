package internal

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/kunitsucom/util.go/exp/diff/simplediff"
)

func NoError(tb testing.TB, printf func(format string, args ...any), err error) {
	tb.Helper()

	if err != nil {
		printf("❌: err != nil: %+v", err)
	}
}

func Error(tb testing.TB, printf func(format string, args ...any), err error) {
	tb.Helper()

	if err == nil {
		printf("❌: err == nil: %+v", err)
	}
}

func True(tb testing.TB, printf func(format string, args ...any), value bool) {
	tb.Helper()

	if !value {
		printf("❌: value == false")
	}
}

func False(tb testing.TB, printf func(format string, args ...any), value bool) {
	tb.Helper()

	if value {
		printf("❌: value == true")
	}
}

func Equal(tb testing.TB, printf func(format string, args ...any), expected, actual interface{}) {
	tb.Helper()

	if reflect.DeepEqual(expected, actual) {
		printf("❌: expected != actual:\n---EXPECTED\n+++ACTUAL\n%s", simplediff.Diff(fmt.Sprintf("%+v", expected), fmt.Sprintf("%+v", actual)))
	}
}

func NotEqual(tb testing.TB, printf func(format string, args ...any), expected, actual interface{}) {
	tb.Helper()

	if !reflect.DeepEqual(expected, actual) {
		printf("❌: expected == actual:\n---EXPECTED\n+++ACTUAL\n%s", simplediff.Diff(fmt.Sprintf("%+v", expected), fmt.Sprintf("%+v", actual)))
	}
}
