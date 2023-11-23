package postgres

import "testing"

func Test_isStmt(t *testing.T) {
	t.Parallel()

	(&CreateTableStmt{}).isStmt()
	(&DropTableStmt{}).isStmt()
	(&AlterTableStmt{}).isStmt()
}
