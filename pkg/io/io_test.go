package ioz_test

import (
	"bytes"
	"strings"
	"testing"

	ioz "github.com/kunitsuinc/util.go/pkg/io"
	testz "github.com/kunitsuinc/util.go/pkg/test"
)

func TestReadAllString(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		expect := "test string"
		actual, err := ioz.ReadAllString(bytes.NewBufferString(expect))
		if err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure()", func(t *testing.T) {
		t.Parallel()

		_, err := ioz.ReadAllString(testz.NewReadWriter(nil, 0, testz.ErrTestError))
		if err == nil {
			t.Errorf("❌: err == nil: %v", err)
		}
		expect := testz.ErrTestError
		if !strings.Contains(err.Error(), expect.Error()) {
			t.Errorf("❌: err: %v != %v", err, expect)
		}
	})
}
