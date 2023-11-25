package cockroachdb

import (
	"testing"

	"github.com/kunitsucom/util.go/testing/require"
)

func Test_isStmt(t *testing.T) {
	t.Parallel()

	(&CreateTableStmt{}).isStmt()
	(&DropTableStmt{}).isStmt()
	(&AlterTableStmt{}).isStmt()
}

func TestIdent_String(t *testing.T) {
	t.Parallel()

	ident := &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}
	expected := ident.Raw
	actual := ident.String()

	require.Equal(t, expected, actual)

	t.Logf("✅: ident: %#v", ident)
}
