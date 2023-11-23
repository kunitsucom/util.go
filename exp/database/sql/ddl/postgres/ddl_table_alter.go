package postgres

// MEMO: https://www.postgresql.jp/docs/11/sql-altertable.html

type AlterTableAction interface {
	isAlterTableAction()
}

// RenameTable represents ALTER TABLE table_name RENAME TO new_table_name.
type RenameTable struct {
	NewName *Ident
}

func (*RenameTable) isAlterTableAction() {}

// RenameConstraint represents ALTER TABLE table_name RENAME COLUMN.
type RenameConstraint struct {
	Name    *Ident
	NewName *Ident
}

func (*RenameConstraint) isAlterTableAction() {}

// RenameColumn represents ALTER TABLE table_name RENAME COLUMN.
type RenameColumn struct {
	Name    *Ident
	NewName *Ident
}

func (*RenameColumn) isAlterTableAction() {}

// AddColumn represents ALTER TABLE table_name ADD COLUMN.
type AddColumn struct {
	Column *Column
}

func (*AddColumn) isAlterTableAction() {}

// DropColumn represents ALTER TABLE table_name DROP COLUMN.
type DropColumn struct {
	Name *Ident
}

func (*DropColumn) isAlterTableAction() {}

// AlterColumn represents ALTER TABLE table_name ALTER COLUMN.
type AlterColumn struct {
	Name   *Ident
	Action AlterColumnAction
}

func (*AlterColumn) isAlterTableAction() {}

type AlterColumnAction interface {
	isAlterColumnAction()
}

// AlterColumnSetDataType represents ALTER TABLE table_name ALTER COLUMN column_name SET DATA TYPE.
type AlterColumnSetDataType struct {
	DataType *DataType
}

func (*AlterColumnSetDataType) isAlterColumnAction() {}

// AlterColumnSetDefault represents ALTER TABLE table_name ALTER COLUMN column_name SET DEFAULT.
type AlterColumnSetDefault struct {
	Default *Default
}

func (*AlterColumnSetDefault) isAlterColumnAction() {}

// AlterColumnDropDefault represents ALTER TABLE table_name ALTER COLUMN column_name DROP DEFAULT.
type AlterColumnDropDefault struct{}

func (*AlterColumnDropDefault) isAlterColumnAction() {}

// AlterColumnSetNotNull represents ALTER TABLE table_name ALTER COLUMN column_name SET NOT NULL.
type AlterColumnSetNotNull struct{}

func (*AlterColumnSetNotNull) isAlterColumnAction() {}

// AlterColumnDropNotNull represents ALTER TABLE table_name ALTER COLUMN column_name DROP NOT NULL.
type AlterColumnDropNotNull struct{}

func (*AlterColumnDropNotNull) isAlterColumnAction() {}

// AddConstraint represents ALTER TABLE table_name ADD CONSTRAINT.
type AddConstraint struct {
	Constraint Constraint
	NotValid   bool
}

func (*AddConstraint) isAlterTableAction() {}

// DropConstraint represents ALTER TABLE table_name DROP CONSTRAINT.
type DropConstraint struct {
	Name *Ident
}

func (*DropConstraint) isAlterTableAction() {}

// AlterConstraint represents ALTER TABLE table_name ALTER CONSTRAINT.
type AlterConstraint struct {
	Name              *Ident
	Deferrable        bool
	InitiallyDeferred bool
}

func (*AlterConstraint) isAlterTableAction() {}

var _ Stmt = (*AlterTableStmt)(nil)

type AlterTableStmt struct {
	Indent    string
	TableName *Ident
	Action    AlterTableAction
}

func (*AlterTableStmt) isStmt() {}

//nolint:cyclop
func (s *AlterTableStmt) String() string {
	str := "ALTER TABLE " +
		s.TableName.String() + " "

	switch a := s.Action.(type) {
	case *RenameTable:
		str += "RENAME TO " + a.NewName.String()
	case *RenameColumn:
		str += "RENAME COLUMN " + a.Name.String() + " TO " + a.NewName.String()
	case *RenameConstraint:
		str += "RENAME CONSTRAINT " + a.Name.String() + " TO " + a.NewName.String()
	case *AddColumn:
		str += "ADD COLUMN " + a.Column.String()
	case *DropColumn:
		str += "DROP COLUMN " + a.Name.String()
	case *AlterColumn:
		str += "ALTER COLUMN " + a.Name.String() + " "
		switch ca := a.Action.(type) {
		case *AlterColumnSetDataType:
			str += "SET DATA TYPE " + ca.DataType.String()
		case *AlterColumnSetDefault:
			str += "SET " + ca.Default.String()
		case *AlterColumnDropDefault:
			str += "DROP DEFAULT"
		case *AlterColumnSetNotNull:
			str += "SET NOT NULL"
		case *AlterColumnDropNotNull:
			str += "DROP NOT NULL"
		}
	case *AddConstraint:
		str += "ADD " + a.Constraint.String()
	case *DropConstraint:
		str += "DROP CONSTRAINT " + a.Name.String()
	case *AlterConstraint:
		str += "ALTER CONSTRAINT " + a.Name.String() + " "
		if a.Deferrable {
			str += "DEFERRABLE"
		} else {
			str += "NOT DEFERRABLE"
		}
		if a.InitiallyDeferred {
			str += " INITIALLY DEFERRED"
		} else {
			str += " INITIALLY IMMEDIATE"
		}
	}

	return str + ";\n"
}
