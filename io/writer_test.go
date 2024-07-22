package ioz_test

import (
	"bytes"
	"testing"

	ioz "github.com/kunitsucom/util.go/io"
	"github.com/kunitsucom/util.go/testing/assert"
)

func TestWriteFunc_Write(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		buf := new(bytes.Buffer)
		w := ioz.WriteFunc(func(p []byte) (n int, err error) {
			return buf.Write(append([]byte("prefix "), p...))
		})
		_, _ = w.Write([]byte("test"))
		assert.Equal(t, "prefix test", buf.String())
	})
}
