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
	Name        *ObjectName
	TableName   *ObjectName
	Columns     []*ColumnIdent
}

func (s *CreateIndexStmt) GetNameForDiff() string {
	return s.Name.StringForDiff()
}

func (s *CreateIndexStmt) String() string {
	str := "CREATE "
	if s.Unique {
		str += "UNIQUE " //nolint:goconst
	}
	str += "INDEX " //nolint:goconst
	if s.IfNotExists {
		str += "IF NOT EXISTS " //nolint:goconst
	}
	str += s.Name.String() + " ON " + s.TableName.String() + " (" + stringz.JoinStringers(", ", s.Columns...) + ");\n"
	return str
}

func (s *CreateIndexStmt) StringForDiff() string {
	str := "CREATE "
	if s.Unique {
		str += "UNIQUE " //nolint:goconst
	}
	str += "INDEX " //nolint:goconst
	if s.IfNotExists {
		str += "IF NOT EXISTS " //nolint:goconst
	}
	str += s.Name.StringForDiff() + " ON " + s.TableName.StringForDiff() + " ("
	for i, c := range s.Columns {
		if i > 0 {
			str += ", "
		}
		str += c.StringForDiff()
	}
	str += ");\n"
	return str
}

func (*CreateIndexStmt) isStmt()            {}
func (s *CreateIndexStmt) GoString() string { return internal.GoString(*s) }
