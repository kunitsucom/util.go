package osz_test

import (
	"testing"

	"github.com/kunitsuinc/util.go/osz"
)

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
		for _, path := range []string{"path_not_exist"} {
			if osz.IsDir(path) {
				t.Errorf("path `%s` should not exist", path)
			}
		}
	})
}
