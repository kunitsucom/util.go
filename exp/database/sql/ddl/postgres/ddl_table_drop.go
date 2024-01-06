package postgres

import (
	"strings"

	"github.com/kunitsucom/util.go/exp/database/sql/ddl/internal"
)

// MEMO: https://www.postgresql.jp/docs/11/sql-createtable.html

var _ Stmt = (*DropTableStmt)(nil)

type DropTableStmt struct {
	Comment  string
	IfExists bool
	Name     *ObjectName
}

func (s *DropTableStmt) GetNameForDiff() string {
	return s.Name.StringForDiff()
}

func (s *DropTableStmt) String() string {
	var str string
	if s.Comment != "" {
		comments := strings.Split(s.Comment, "\n")
		for i := range comments {
			if comments[i] != "" {
				str += CommentPrefix + comments[i] + "\n"
			}
		}
	}
	str += "DROP TABLE "
	if s.IfExists {
		str += "IF EXISTS "
	}
	str += s.Name.String() + ";\n"
	return str
}

func (*DropTableStmt) isStmt()            {}
func (s *DropTableStmt) GoString() string { return internal.GoString(*s) }
