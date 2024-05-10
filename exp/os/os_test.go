package osz_test

import (
	"os"
	"runtime"
	"testing"

	osz "github.com/kunitsucom/util.go/exp/os"
)

var _, callerFile, _, _ = runtime.Caller(0)

//nolint:paralleltest
func TestReadlinkAndReadFile(t *testing.T) {
	t.Run("success()", func(t *testing.T) {
		expect := callerFile
		linkTestFile := expect + ".symlink"
		_ = os.Symlink(expect, linkTestFile)
		t.Cleanup(func() {
			_ = os.Remove(linkTestFile)
		})
		actual, _, err := osz.ReadlinkAndReadFile(linkTestFile)
		if err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure()", func(t *testing.T) {
		for _, path := range []string{"path_not_exist"} {
			if _, _, err := osz.ReadlinkAndReadFile(path); err == nil {
				t.Errorf("❌: path `%s` should not exist", path)
			}
		}
	})
}
