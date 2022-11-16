package osz_test

import (
	"os"
	"runtime"
	"testing"

	osz "github.com/kunitsuinc/util.go/pkg/os"
)

var _, callerFile, _, _ = runtime.Caller(1)

func TestExists(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		for _, path := range []string{".", ".."} {
			if !osz.Exists(path) {
				t.Errorf("path `%s` should exist", path)
			}
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		for _, path := range []string{"path_not_exist"} {
			if osz.Exists(path) {
				t.Errorf("path `%s` should not exist", path)
			}
		}
	})
}

func TestIsDir(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		for _, path := range []string{".", ".."} {
			if !osz.IsDir(path) {
				t.Errorf("path `%s` should be dir", path)
			}
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		for _, path := range []string{"path_not_exist", callerFile} {
			if osz.IsDir(path) {
				t.Errorf("path `%s` should not exist", path)
			}
		}
	})
}

func TestCheckDir(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		for _, path := range []string{".", ".."} {
			if err := osz.CheckDir(path); err != nil {
				t.Errorf("err != nil: %v", err)
			}
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		for _, path := range []string{"path_not_exist", callerFile} {
			if err := osz.CheckDir(path); err == nil {
				t.Errorf("err == nil")
			}
		}
	})
}

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
			t.Errorf("err != nil: %v", err)
		}
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure()", func(t *testing.T) {
		for _, path := range []string{"path_not_exist"} {
			if _, _, err := osz.ReadlinkAndReadFile(path); err == nil {
				t.Errorf("path `%s` should not exist", path)
			}
		}
	})
}
