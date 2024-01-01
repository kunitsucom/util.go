package cockroachdb

import (
	"strings"

	"github.com/kunitsucom/util.go/exp/database/sql/ddl/internal"
)

// MEMO: https://www.cockroachlabs.com/docs/stable/create-table //diff:ignore-line-postgres-cockroach

var _ Stmt = (*CreateTableStmt)(nil)

type CreateTableStmt struct {
	Comment     string
	Indent      string
	IfNotExists bool
	Name        *ObjectName
	Columns     []*Column
	Constraints Constraints
	Options     []*Option
}

func (s *CreateTableStmt) GetNameForDiff() string {
	return s.Name.StringForDiff()
}

//nolint:cyclop
func (s *CreateTableStmt) String() string {
	var str string
	if s.Comment != "" {
		for _, s := range strings.Split(s.Comment, "\n") {
			str += "-- " + s + "\n"
		}
	}
	str += "CREATE TABLE "
	if s.IfNotExists {
		str += "IF NOT EXISTS "
	}
	str += s.Name.String() + " (\n"
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
		lastConstraint := len(s.Constraints) - 1
		for i, v := range s.Constraints {
			str += Indent
			str += v.String()
			if i != lastConstraint {
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
