package postgres

import (
	"testing"

	"github.com/kunitsucom/util.go/testing/require"
)

func Test_isStmt(t *testing.T) {
	t.Parallel()

	(&CreateTableStmt{}).isStmt()
	(&DropTableStmt{}).isStmt()
	(&AlterTableStmt{}).isStmt()
	(&CreateIndexStmt{}).isStmt()
	(&DropIndexStmt{}).isStmt()
}

func TestIdent_String(t *testing.T) {
	t.Parallel()

	ident := &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}
	expected := ident.Raw
	actual := ident.String()

	require.Equal(t, expected, actual)

	t.Logf("✅: %s: ident: %#v", t.Name(), ident)
}

func TestIdent_PlainString(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		ident := &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}
		expected := ident.Name
		actual := ident.PlainString()

		require.Equal(t, expected, actual)
	})

	t.Run("success,empty", func(t *testing.T) {
		t.Parallel()
		ident := (*Ident)(nil)
		expected := ""
		actual := ident.PlainString()

		require.Equal(t, expected, actual)
	})
}
