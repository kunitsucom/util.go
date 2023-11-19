package postgres

// MEMO: https://www.postgresql.jp/docs/11/sql-altertable.html

type AlterTableStmt struct {
	Indent     string
	Name       Ident
	Action     string
	Object     string
	Column     Column
	Constraint Constraint
}

func (s *AlterTableStmt) String() string {
	str := "ALTER TABLE " +
		s.Name.String() + " " +
		s.Action + " "
	switch s.Action {
	case "ADD":
		str += s.Object + " "
		switch s.Object {
		case "COLUMN":
			str += s.Column.String()
		case "CONSTRAINT":
			str += s.Constraint.String()
		}
	case "DROP":
		str += s.Object + " "
		switch s.Object {
		case "COLUMN":
			str += s.Column.Name.String()
		case "CONSTRAINT":
			str += s.Constraint.GetName().String()
		}
	}
	return str
}

func (*AlterTableStmt) isStmt() {}
