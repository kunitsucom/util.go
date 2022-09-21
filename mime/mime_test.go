package mime_test

import (
	"strings"
	"testing"

	"github.com/kunitsuinc/util.go/mime"
	"github.com/kunitsuinc/util.go/test/fixture"
)

func TestDetectContentType(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		const expect = "text/html; charset=utf-8"
		actual, err := mime.DetectContentType(strings.NewReader("<!DOCTYPE html>"))
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		r := &fixture.Reader{Err: fixture.ErrFixtureError}
		if _, err := mime.DetectContentType(r); err == nil {
			t.Errorf("err == nil")
		}
	})
}
