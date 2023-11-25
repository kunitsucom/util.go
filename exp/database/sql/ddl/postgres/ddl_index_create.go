package postgres

import (
	"github.com/kunitsucom/util.go/exp/database/sql/ddl/internal"
	stringz "github.com/kunitsucom/util.go/strings"
)

// MEMO: https://www.cockroachlabs.com/docs/stable/create-index

var _ Stmt = (*CreateIndexStmt)(nil)

type CreateIndexStmt struct {
	Unique      bool
	IfNotExists bool
	Schema      *Ident
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
	if s.Schema != nil {
		str += s.Schema.String() + "."
	}
	str += s.Name.String() + " ON " + s.TableName.String() + " (" + stringz.JoinStringers(", ", s.Columns...) + ");\n"
	return str
}

func (s *CreateIndexStmt) PlainString() string {
	str := "CREATE "
	if s.Unique {
		str += "UNIQUE " //nolint:goconst
	}
	str += "INDEX " //nolint:goconst
	if s.IfNotExists {
		str += "IF NOT EXISTS "
	}
	if s.Schema != nil {
		str += s.Schema.PlainString() + "."
	}
	str += s.Name.PlainString() + " ON " + s.TableName.PlainString() + " ("
	for i, c := range s.Columns {
		if i > 0 {
			str += ", "
		}
		str += c.PlainString()
	}
	str += ");\n"
	return str
}

func (*CreateIndexStmt) isStmt()            {}
func (s *CreateIndexStmt) GoString() string { return internal.GoString(*s) }
