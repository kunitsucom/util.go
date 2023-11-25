package postgres

import (
	"testing"

	"github.com/kunitsucom/util.go/testing/assert"
)

func TestCreateTableStmt_String(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		stmt := &CreateTableStmt{
			Indent: "  ",
			Name:   &Ident{Name: "test", QuotationMark: `"`, Raw: `"test"`},
			Columns: []*Column{
				{Name: &Ident{Name: "id", Raw: "id"}, DataType: &DataType{Name: "INTEGER"}},
				{Name: &Ident{Name: "name", Raw: "name"}, DataType: &DataType{Name: "VARYING", Size: "255"}},
			},
			Options: []*Option{
				{Name: "TABLESPACE", Value: &Ident{Name: "default_tablespace", Raw: "default_tablespace"}},
				{Name: "LIKE", Value: &Ident{Name: "parent_test", Raw: "parent_test"}},
			},
		}

		expected := `CREATE TABLE "test" (
    id INTEGER,
    name VARYING(255)
)
TABLESPACE default_tablespace,
LIKE parent_test;
`
		actual := stmt.String()
		assert.Equal(t, expected, actual)

		t.Logf("✅: %s: stmt: %#v", t.Name(), stmt)
	})
}

func TestCreateTableStmt_GetPlainName(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		stmt := &CreateTableStmt{Name: &Ident{Name: "test", QuotationMark: `"`, Raw: `"test"`}}
		expected := "test"
		actual := stmt.GetPlainName()

		assert.Equal(t, expected, actual)

		t.Logf("✅: %s: stmt: %#v", t.Name(), stmt)
	})
}
