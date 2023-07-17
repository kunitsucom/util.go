package httpz_test

import (
	"bytes"
	"strings"
	"testing"

	httpz "github.com/kunitsuinc/util.go/pkg/net/http"
)

func TestNewHeader(t *testing.T) {
	t.Parallel()

	t.Run("normal", func(t *testing.T) {
		t.Parallel()

		const expect = "Add-Key-1: AddValue1\\r\\nAdd-Key-2: AddValue2\\r\\nAdd-Key-3: AddValue3\\r\\nSet-Key-1: SetValue1\\r\\nSet-Key-2: SetValue2\\r\\n"

		buf := bytes.NewBuffer(nil)

		if err := httpz.NewHeader(
			httpz.Add("Add-Key-1", "AddValue1"),
			httpz.Add("Add-Key-2", "AddValue2"),
			httpz.Add("Set-Key-1", "BeforeValue"),
			httpz.Set("Set-Key-1", "SetValue1"),
			httpz.Add("Add-Key-3", "AddValue3"),
			httpz.Set("Set-Key-2", "SetValue2"),
		).Write(buf); err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
		actual := strings.ReplaceAll(buf.String(), "\r\n", "\\r\\n")

		if actual != expect {
			t.Errorf("❌: expect != actual: %s != %s", expect, actual)
		}
	})
}
