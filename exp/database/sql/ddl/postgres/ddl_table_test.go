package postgres

import (
	"testing"

	"github.com/kunitsucom/util.go/testing/require"
)

func Test_isConstraint(t *testing.T) {
	t.Parallel()

	(&PrimaryKeyConstraint{}).isConstraint()
	(&ForeignKeyConstraint{}).isConstraint()
	(&UniqueConstraint{}).isConstraint()
	(&CheckConstraint{}).isConstraint()
}

func TestDefault_String(t *testing.T) {
	t.Parallel()

	t.Run("success,nil", func(t *testing.T) {
		t.Parallel()

		d := (*Default)(nil)
		expected := ""
		actual := d.String()
		require.Equal(t, expected, actual)

		t.Logf("✅: d: %#v", d)
	})
	t.Run("success,nilnil", func(t *testing.T) {
		t.Parallel()

		d := &Default{}
		expected := ""
		actual := d.String()
		require.Equal(t, expected, actual)

		t.Logf("✅: d: %#v", d)
	})
	t.Run("success,DEFAULT_VALUE", func(t *testing.T) {
		t.Parallel()

		d := &Default{Value: &DefaultValue{[]*Ident{{Name: "now()", Raw: "now()"}}}}
		expected := "DEFAULT now()"
		actual := d.String()
		require.Equal(t, expected, actual)

		t.Logf("✅: d: %#v", d)
	})
	t.Run("success,DEFAULT_EXPR", func(t *testing.T) {
		t.Parallel()

		d := &Default{Value: &DefaultValue{[]*Ident{{Name: "(", Raw: "("}, {Name: "age", Raw: "age"}, {Name: ">=", Raw: ">="}, {Name: "0", Raw: "0"}, {Name: ")", Raw: ")"}}}}
		expected := "DEFAULT (age >= 0)"
		actual := d.String()
		require.Equal(t, expected, actual)

		t.Logf("✅: d: %#v", d)
	})
}
