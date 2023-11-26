package postgres

import (
	"testing"

	"github.com/kunitsucom/util.go/testing/require"
)

func TestDropTableStmt_GetPlainName(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		stmt := &DropTableStmt{Name: &ObjectName{Name: &Ident{Name: "test", QuotationMark: `"`, Raw: `"test"`}}}
		expected := "test"
		actual := stmt.GetPlainName()

		require.Equal(t, expected, actual)

		t.Logf("✅: %s: stmt: %#v", t.Name(), stmt)
	})
}

func TestDropTableStmt_String(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		stmt := &DropTableStmt{
			IfExists: true,
			Name:     &ObjectName{Name: &Ident{Name: "test", QuotationMark: `"`, Raw: `"test"`}},
		}

		expected := `DROP TABLE IF EXISTS "test";` + "\n"
		actual := stmt.String()
		require.Equal(t, expected, actual)

		t.Logf("✅: %s: stmt: %#v", t.Name(), stmt)
	})
}
