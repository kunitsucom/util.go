package postgres

import "github.com/kunitsucom/util.go/exp/database/sql/ddl/internal"

// MEMO: https://www.postgresql.jp/docs/11/sql-createtable.html

var _ Stmt = (*CreateTableStmt)(nil)

type CreateTableStmt struct {
	Indent      string
	Name        *Ident
	Columns     []*Column
	Constraints Constraints
	Options     []*Option
}

func (s *CreateTableStmt) String() string {
	str := "CREATE TABLE " + s.Name.String() + " (\n"
	lastIndex := len(s.Columns) - 1
	hasConstraint := len(s.Constraints) > 0
	for i, v := range s.Columns {
		str += Indent
		str += v.String()
		if i != lastIndex || hasConstraint {
			str += ",\n"
		} else {
			str += "\n"
		}
	}
	if len(s.Constraints) > 0 {
		lastIndex := len(s.Constraints) - 1
		for i, v := range s.Constraints {
			str += Indent
			str += v.String()
			if i != lastIndex {
				str += ",\n"
			} else {
				str += "\n"
			}
		}
	}
	str += ")"
	if len(s.Options) > 0 {
		str += "\n"
		lastIndex := len(s.Options) - 1
		for i, v := range s.Options {
			str += v.String()
			if i != lastIndex {
				str += ",\n"
			}
		}
	}

	str += ";\n"
	return str
}

func (*CreateTableStmt) isStmt()            {}
func (s *CreateTableStmt) GoString() string { return internal.GoString(*s) }
