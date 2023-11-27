package postgres

import "github.com/kunitsucom/util.go/exp/database/sql/ddl/internal"

// MEMO: https://www.postgresql.jp/docs/11/sql-createtable.html

var _ Stmt = (*DropTableStmt)(nil)

type DropTableStmt struct {
	IfExists bool
	Name     *ObjectName
}

func (s *DropTableStmt) GetPlainName() string {
	return s.Name.StringForDiff()
}

func (s *DropTableStmt) String() string {
	str := "DROP TABLE "
	if s.IfExists {
		str += "IF EXISTS "
	}
	str += s.Name.String() + ";\n"
	return str
}

func (*DropTableStmt) isStmt()            {}
func (s *DropTableStmt) GoString() string { return internal.GoString(*s) }
