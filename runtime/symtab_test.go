package runtimez_test

import (
	"testing"

	runtimez "github.com/kunitsucom/util.go/runtime"
)

func TestFuncName(t *testing.T) {
	t.Parallel()
	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		if expect, actual := "runtime_test.testFuncName", testFuncName(0); expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestFullFuncName(t *testing.T) {
	t.Parallel()
	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		if expect, actual := "github.com/kunitsucom/util.go/runtime_test.testFullFuncName", testFullFuncName(0); expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func testFuncName(skip int) string {
	return runtimez.FuncName(skip)
}

func testFullFuncName(skip int) string {
	return runtimez.FullFuncName(skip)
}
