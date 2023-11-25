package cockroachdb

import (
	"github.com/kunitsucom/util.go/exp/database/sql/ddl/internal"
	stringz "github.com/kunitsucom/util.go/strings"
)

// MEMO: https://www.postgresql.jp/docs/11/sql-createtable.html

var _ Stmt = (*CreateIndexStmt)(nil)

type CreateIndexStmt struct {
	Unique      bool
	IfNotExists bool
	Name        *Ident
	TableName   *Ident
	Columns     []*ColumnIdent
}

func (s *CreateIndexStmt) GetPlainName() string {
	return s.Name.PlainString()
}

func (s *CreateIndexStmt) String() string {
	str := "CREATE "
	if s.Unique {
		str += "UNIQUE "
	}
	str += "INDEX "
	if s.IfNotExists {
		str += "IF NOT EXISTS "
	}
	str += s.Name.String() + " ON " + s.TableName.String() + " (" + stringz.JoinStringers(", ", s.Columns...) + ");\n"
	return str
}

func (*CreateIndexStmt) isStmt()            {}
func (s *CreateIndexStmt) GoString() string { return internal.GoString(*s) }
