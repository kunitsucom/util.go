package postgres

import "github.com/kunitsucom/util.go/exp/database/sql/ddl/internal"

// MEMO: https://www.cockroachlabs.com/docs/stable/drop-index

var _ Stmt = (*DropIndexStmt)(nil)

type DropIndexStmt struct {
	IfExists bool
	Name     *Ident
}

func (s *DropIndexStmt) GetPlainName() string {
	return s.Name.PlainString()
}

func (s *DropIndexStmt) String() string {
	str := "DROP INDEX "
	if s.IfExists {
		str += "IF EXISTS " //nolint:goconst
	}
	str += s.Name.String() + ";\n"
	return str
}

func (s *DropIndexStmt) PlainString() string {
	str := "DROP INDEX "
	if s.IfExists {
		str += "IF EXISTS "
	}
	str += s.Name.PlainString() + ";\n"
	return str
}

func (*DropIndexStmt) isStmt()            {}
func (s *DropIndexStmt) GoString() string { return internal.GoString(*s) }
