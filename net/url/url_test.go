package urlz_test

import (
	"testing"

	urlz "github.com/kunitsucom/util.go/net/url"
)

func TestNewValues(t *testing.T) {
	t.Parallel()

	t.Run("normal", func(t *testing.T) {
		t.Parallel()

		const expect = "AddKey1=AddValue1&AddKey2=AddValue2&AddKey3=AddValue3&SetKey1=SetValue1&SetKey2=SetValue2"

		actual := urlz.NewValues(
			urlz.Add("AddKey1", "AddValue1"),
			urlz.Add("AddKey2", "AddValue2"),
			urlz.Add("SetKey1", "BeforeValue"),
			urlz.Set("SetKey1", "SetValue1"),
			urlz.Add("AddKey3", "AddValue3"),
			urlz.Set("SetKey2", "SetValue2"),
		).Encode()

		if actual != expect {
			t.Errorf("❌: expect != actual: %s != %s", expect, actual)
		}
	})
}
