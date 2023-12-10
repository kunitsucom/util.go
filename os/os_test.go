package osz_test

import (
	"path/filepath"
	"runtime"
	"testing"

	osz "github.com/kunitsucom/util.go/os"
)

var callerFile = func() string {
	_, callerFile, _, _ := runtime.Caller(0) //nolint:dogsled
	return filepath.Base(callerFile)
}()

func TestExists(t *testing.T) {
	t.Parallel()
	t.Run("success,", func(t *testing.T) {
		t.Parallel()
		for _, path := range []string{".", ".."} {
			if !osz.Exists(path) {
				t.Errorf("❌: path `%s` should exist", path)
			}
		}
	})

	t.Run("failure,", func(t *testing.T) {
		t.Parallel()
		for _, path := range []string{"path_not_exist"} {
			if osz.Exists(path) {
				t.Errorf("❌: path `%s` should not exist", path)
			}
		}
	})
}

func TestIsFile(t *testing.T) {
	t.Parallel()
	t.Run("success,", func(t *testing.T) {
		t.Parallel()
		for _, path := range []string{callerFile} {
			if !osz.IsFile(path) {
				t.Errorf("❌: path `%s` should be file", path)
			}
		}
	})

	t.Run("failure,", func(t *testing.T) {
		t.Parallel()
		for _, path := range []string{"path_not_exist", "."} {
			if osz.IsFile(path) {
				t.Errorf("❌: path `%s` should not exist", path)
			}
		}
	})
}

func TestIsDir(t *testing.T) {
	t.Parallel()
	t.Run("success,", func(t *testing.T) {
		t.Parallel()
		for _, path := range []string{".", ".."} {
			if !osz.IsDir(path) {
				t.Errorf("❌: path `%s` should be dir", path)
			}
		}
	})

	t.Run("failure,", func(t *testing.T) {
		t.Parallel()
		for _, path := range []string{"path_not_exist", callerFile} {
			if osz.IsDir(path) {
				t.Errorf("❌: path `%s` should not exist", path)
			}
		}
	})
}

func TestCheckDir(t *testing.T) {
	t.Parallel()
	t.Run("success,", func(t *testing.T) {
		t.Parallel()
		for _, path := range []string{".", ".."} {
			if err := osz.CheckDir(path); err != nil {
				t.Errorf("❌: err != nil: %v", err)
			}
		}
	})

	t.Run("failure,", func(t *testing.T) {
		t.Parallel()
		for _, path := range []string{"path_not_exist", callerFile} {
			if err := osz.CheckDir(path); err == nil {
				t.Errorf("❌: err == nil")
			}
		}
	})
}
