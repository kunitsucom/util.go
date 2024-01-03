package postgres

import (
	"strings"

	"github.com/kunitsucom/util.go/exp/database/sql/ddl/internal"
)

// MEMO: https://www.cockroachlabs.com/docs/stable/drop-index

var _ Stmt = (*DropIndexStmt)(nil)

type DropIndexStmt struct {
	Comment  string
	IfExists bool
	Name     *ObjectName
}

func (s *DropIndexStmt) GetNameForDiff() string {
	return s.Name.StringForDiff()
}

func (s *DropIndexStmt) String() string {
	var str string
	if s.Comment != "" {
		for _, v := range strings.Split(s.Comment, "\n") {
			str += CommentPrefix + v + "\n"
		}
	}
	str += "DROP INDEX "
	if s.IfExists {
		str += "IF EXISTS "
	}
	str += s.Name.String() + ";\n"
	return str
}

func (*DropIndexStmt) isStmt()            {}
func (s *DropIndexStmt) GoString() string { return internal.GoString(*s) }
