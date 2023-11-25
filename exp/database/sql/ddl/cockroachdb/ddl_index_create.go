package cockroachdb

import (
	"github.com/kunitsucom/util.go/exp/database/sql/ddl/internal"
	stringz "github.com/kunitsucom/util.go/strings"
)

// MEMO: https://www.cockroachlabs.com/docs/stable/create-index

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
		str += "UNIQUE " //nolint:goconst
	}
	str += "INDEX " //nolint:goconst
	if s.IfNotExists {
		str += "IF NOT EXISTS "
	}
	str += s.Name.String() + " ON " + s.TableName.String() + " (" + stringz.JoinStringers(", ", s.Columns...) + ");\n"
	return str
}

func (*CreateIndexStmt) isStmt()            {}
func (s *CreateIndexStmt) GoString() string { return internal.GoString(*s) }
