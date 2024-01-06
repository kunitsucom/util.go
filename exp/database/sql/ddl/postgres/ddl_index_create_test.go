package postgres

import (
	"testing"

	"github.com/kunitsucom/util.go/testing/require"
)

func TestCreateIndexStmt_GetNameForDiff(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		stmt := &CreateIndexStmt{Name: &ObjectName{Name: &Ident{Name: "test", QuotationMark: `"`, Raw: `"test"`}}}
		expected := "test"
		actual := stmt.GetNameForDiff()

		require.Equal(t, expected, actual)
	})
}

func TestCreateIndexStmt_String(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		stmt := &CreateIndexStmt{
			Comment:     "test comment content",
			IfNotExists: true,
			Name:        &ObjectName{Name: &Ident{Name: "test", QuotationMark: `"`, Raw: `"test"`}},
			TableName:   &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
			Using:       []*Ident{{Name: "btree", QuotationMark: ``, Raw: `btree`}},
			Columns: []*ColumnIdent{
				{
					Ident: &Ident{Name: "id", QuotationMark: `"`, Raw: `"id"`},
				},
			},
		}
		expected := `-- test comment content
CREATE INDEX IF NOT EXISTS "test" ON "users" USING btree ("id");
`

		actual := stmt.String()

		require.Equal(t, expected, actual)

		t.Logf("✅: %s: stmt: %#v", t.Name(), stmt)
	})
}
