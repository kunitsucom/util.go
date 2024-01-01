package cockroachdb

import "github.com/kunitsucom/util.go/exp/database/sql/ddl/internal"

// MEMO: https://www.cockroachlabs.com/docs/stable/alter-table

type AlterTableAction interface {
	isAlterTableAction()
	GoString() string
}

// RenameTable represents ALTER TABLE table_name RENAME TO new_table_name.
type RenameTable struct {
	NewName *ObjectName
}

func (*RenameTable) isAlterTableAction() {}

func (s *RenameTable) GoString() string { return internal.GoString(*s) }

// RenameConstraint represents ALTER TABLE table_name RENAME COLUMN.
type RenameConstraint struct {
	Name    *Ident
	NewName *Ident
}

func (*RenameConstraint) isAlterTableAction() {}

func (s *RenameConstraint) GoString() string { return internal.GoString(*s) }

// RenameColumn represents ALTER TABLE table_name RENAME COLUMN.
type RenameColumn struct {
	Name    *Ident
	NewName *Ident
}

func (*RenameColumn) isAlterTableAction() {}

func (s *RenameColumn) GoString() string { return internal.GoString(*s) }

// AddColumn represents ALTER TABLE table_name ADD COLUMN.
type AddColumn struct {
	Column *Column
}

func (*AddColumn) isAlterTableAction() {}

func (s *AddColumn) GoString() string { return internal.GoString(*s) }

// DropColumn represents ALTER TABLE table_name DROP COLUMN.
type DropColumn struct {
	Name *Ident
}

func (*DropColumn) isAlterTableAction() {}

func (s *DropColumn) GoString() string { return internal.GoString(*s) }

// AlterColumn represents ALTER TABLE table_name ALTER COLUMN.
type AlterColumn struct {
	Name   *Ident
	Action AlterColumnAction
}

func (*AlterColumn) isAlterTableAction() {}

func (s *AlterColumn) GoString() string { return internal.GoString(*s) }

type AlterColumnAction interface {
	isAlterColumnAction()
	GoString() string
}

// AlterColumnSetDataType represents ALTER TABLE table_name ALTER COLUMN column_name SET DATA TYPE.
type AlterColumnSetDataType struct {
	DataType *DataType
}

func (*AlterColumnSetDataType) isAlterColumnAction() {}

func (s *AlterColumnSetDataType) GoString() string { return internal.GoString(*s) }

// AlterColumnSetDefault represents ALTER TABLE table_name ALTER COLUMN column_name SET DEFAULT.
type AlterColumnSetDefault struct {
	Default *Default
}

func (*AlterColumnSetDefault) isAlterColumnAction() {}

func (s *AlterColumnSetDefault) GoString() string { return internal.GoString(*s) }

// AlterColumnDropDefault represents ALTER TABLE table_name ALTER COLUMN column_name DROP DEFAULT.
type AlterColumnDropDefault struct{}

func (*AlterColumnDropDefault) isAlterColumnAction() {}

func (s *AlterColumnDropDefault) GoString() string { return internal.GoString(*s) }

// AlterColumnSetNotNull represents ALTER TABLE table_name ALTER COLUMN column_name SET NOT NULL.
type AlterColumnSetNotNull struct{}

func (*AlterColumnSetNotNull) isAlterColumnAction() {}

func (s *AlterColumnSetNotNull) GoString() string { return internal.GoString(*s) }

// AlterColumnDropNotNull represents ALTER TABLE table_name ALTER COLUMN column_name DROP NOT NULL.
type AlterColumnDropNotNull struct{}

func (*AlterColumnDropNotNull) isAlterColumnAction() {}

func (s *AlterColumnDropNotNull) GoString() string { return internal.GoString(*s) }

// AddConstraint represents ALTER TABLE table_name ADD CONSTRAINT.
type AddConstraint struct {
	Constraint Constraint
	NotValid   bool
}

func (*AddConstraint) isAlterTableAction() {}

func (s *AddConstraint) GoString() string { return internal.GoString(*s) }

// DropConstraint represents ALTER TABLE table_name DROP CONSTRAINT.
type DropConstraint struct {
	Name *Ident
}

func (*DropConstraint) isAlterTableAction() {}

func (s *DropConstraint) GoString() string { return internal.GoString(*s) }

// AlterConstraint represents ALTER TABLE table_name ALTER CONSTRAINT.
type AlterConstraint struct {
	Name              *Ident
	Deferrable        bool
	InitiallyDeferred bool
}

func (*AlterConstraint) isAlterTableAction() {}

func (s *AlterConstraint) GoString() string { return internal.GoString(*s) }

var _ Stmt = (*AlterTableStmt)(nil)

type AlterTableStmt struct {
	Name   *ObjectName
	Action AlterTableAction
}

func (*AlterTableStmt) isStmt() {}

func (s *AlterTableStmt) GetPlainName() string {
	return s.Name.StringForDiff()
}

//nolint:cyclop,funlen
func (s *AlterTableStmt) String() string {
	str := "ALTER TABLE "
	str += s.Name.String() + " "
	switch a := s.Action.(type) {
	case *RenameTable:
		str += "RENAME TO "
		str += a.NewName.String()
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
		if a.NotValid {
			str += " NOT VALID"
		}
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

func (s *AlterTableStmt) GoString() string { return internal.GoString(*s) }
