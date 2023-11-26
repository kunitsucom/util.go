package cockroachdb

import (
	"testing"

	"github.com/kunitsucom/util.go/testing/require"
)

func TestCreateIndexStmt_GetPlainName(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		stmt := &CreateIndexStmt{Name: &ObjectName{Name: &Ident{Name: "test", QuotationMark: `"`, Raw: `"test"`}}}
		expected := "test"
		actual := stmt.GetPlainName()

		require.Equal(t, expected, actual)
	})
}

func TestCreateIndexStmt_String(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		stmt := &CreateIndexStmt{
			IfNotExists: true,
			Name:        &ObjectName{Name: &Ident{Name: "test", QuotationMark: `"`, Raw: `"test"`}},
			TableName:   &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
			Columns: []*ColumnIdent{
				{
					Ident: &Ident{Name: "id", QuotationMark: `"`, Raw: `"id"`},
					Order: &Order{Desc: false},
				},
			},
		}
		expected := `CREATE INDEX IF NOT EXISTS "test" ON "users" ("id" ASC);` + "\n"
		actual := stmt.String()

		require.Equal(t, expected, actual)

		t.Logf("✅: %s: stmt: %#v", t.Name(), stmt)
	})
}
