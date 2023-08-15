package httpz_test

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	httpz "github.com/kunitsucom/util.go/net/http"
)

func TestNewHeader(t *testing.T) {
	t.Parallel()

	t.Run("normal", func(t *testing.T) {
		t.Parallel()

		const expect = "Add-Key-1: AddValue1\\r\\nAdd-Key-2: AddValue2\\r\\nAdd-Key-3: AddValue3\\r\\nAdd-Key-4: AddValue4\\r\\nSet-Key-1: SetValue1\\r\\nSet-Key-2: SetValue2\\r\\n"

		buf := bytes.NewBuffer(nil)
		h := http.Header{}
		h.Add("Add-Key-1", "AddValue1")

		if err := httpz.NewHeaderBuilder().
			Merge(h).
			Add("Add-Key-2", "AddValue2").
			Add("Add-Key-3", "AddValue3").
			Add("Set-Key-1", "BeforeValue").
			Set("Set-Key-1", "SetValue1").
			Add("Add-Key-4", "AddValue4").
			Set("Set-Key-2", "SetValue2").
			Build().
			Write(buf); err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
		actual := strings.ReplaceAll(buf.String(), "\r\n", "\\r\\n")

		if actual != expect {
			t.Errorf("❌: expect != actual: %s != %s", expect, actual)
		}
	})
}
