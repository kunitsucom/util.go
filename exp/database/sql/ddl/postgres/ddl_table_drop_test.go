package postgres

import (
	"testing"

	"github.com/kunitsucom/util.go/testing/require"
)

func TestDropTableStmt_GetNameForDiff(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		stmt := &DropTableStmt{Name: &ObjectName{Name: &Ident{Name: "test", QuotationMark: `"`, Raw: `"test"`}}}
		expected := "test"
		actual := stmt.GetNameForDiff()

		require.Equal(t, expected, actual)

		t.Logf("✅: %s: stmt: %#v", t.Name(), stmt)
	})
}

func TestDropTableStmt_String(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		stmt := &DropTableStmt{
			Comment:  "test comment content",
			IfExists: true,
			Name:     &ObjectName{Name: &Ident{Name: "test", QuotationMark: `"`, Raw: `"test"`}},
		}

		expected := `-- test comment content
DROP TABLE IF EXISTS "test";
`
		actual := stmt.String()
		require.Equal(t, expected, actual)

		t.Logf("✅: %s: stmt: %#v", t.Name(), stmt)
	})
}
