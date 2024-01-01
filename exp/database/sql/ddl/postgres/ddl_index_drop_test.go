package postgres

import (
	"testing"

	"github.com/kunitsucom/util.go/testing/require"
)

func TestDropIndexStmt_GetPlainName(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		stmt := &DropIndexStmt{Name: &ObjectName{Name: &Ident{Name: "test", QuotationMark: `"`, Raw: `"test"`}}}
		expected := "test"
		actual := stmt.GetPlainName()

		require.Equal(t, expected, actual)
	})
}

func TestDropIndexStmt_String(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		stmt := &DropIndexStmt{
			IfExists: true,
			Name:     &ObjectName{Name: &Ident{Name: "test", QuotationMark: `"`, Raw: `"test"`}},
		}
		expected := `DROP INDEX IF EXISTS "test";` + "\n"
		actual := stmt.String()

		require.Equal(t, expected, actual)

		t.Logf("✅: %s: stmt: %#v", t.Name(), stmt)
	})
}
