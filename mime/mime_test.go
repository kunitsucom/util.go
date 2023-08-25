package mime_test

import (
	"strings"
	"testing"

	"github.com/kunitsucom/util.go/mime"
	testz "github.com/kunitsucom/util.go/test"
)

func TestDetectContentType(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		const expect = "text/html; charset=utf-8"
		actual, err := mime.DetectContentType(strings.NewReader("<!DOCTYPE html>"))
		if err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		r := &testz.Reader{
			ReadFunc: func(p []byte) (n int, err error) {
				return 0, testz.ErrTestError
			},
		}
		if _, err := mime.DetectContentType(r); err == nil {
			t.Errorf("❌: err == nil")
		}
	})
}
