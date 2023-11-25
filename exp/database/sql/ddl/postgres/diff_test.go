package postgres

import (
	"testing"

	"github.com/kunitsucom/util.go/exp/database/sql/ddl"
	"github.com/kunitsucom/util.go/testing/require"
)

func TestDiff(t *testing.T) {
	t.Parallel()

	t.Run("failure,ddl.ErrNoDifference", func(t *testing.T) {
		t.Parallel()

		before := &DDL{}
		after := &DDL{}
		_, err := Diff(before, after)
		require.ErrorIs(t, err, ddl.ErrNoDifference)
	})

	t.Run("failure,ddl.ErrNotSupported,DropTableStmt", func(t *testing.T) {
		t.Parallel()

		before := &DDL{
			Stmts: []Stmt{
				&DropTableStmt{Name: &Ident{Name: "table_name", Raw: "table_name"}},
			},
		}
		after := (*DDL)(nil)
		_, err := Diff(before, after)
		require.ErrorIs(t, err, ddl.ErrNotSupported)
	})

	t.Run("failure,ddl.ErrNotSupported,DropTableStmt", func(t *testing.T) {
		t.Parallel()

		before := &DDL{
			Stmts: []Stmt{
				&DropTableStmt{Name: &Ident{Name: "table_name", Raw: "table_name"}},
			},
		}
		after := &DDL{}
		_, err := Diff(before, after)
		require.ErrorIs(t, err, ddl.ErrNotSupported)
	})
}
