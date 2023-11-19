package postgres

import (
	"fmt"
	"log"
	"strconv"
)

type Stmt interface{}

const indent = "    "

// CreateTableStmt はCREATE TABLE文を表す構造体です。
type CreateTableStmt struct {
	TableName   Literal
	IfNotExists bool
	Columns     []*TableColumn
	Constraints []*TableConstraint
}

func (s *CreateTableStmt) String() string {
	str := "CREATE TABLE "
	if s.IfNotExists {
		str += "IF NOT EXISTS "
	}
	str += s.TableName.String() + " (\n"

	hasConstraints := len(s.Constraints) > 0

	columnNameFormat := "%-" + func() string {
		max := 0
		for _, column := range s.Columns {
			if len(column.Name.String()) > max {
				max = len(column.Name.String())
			}
		}
		return strconv.Itoa(max)
	}() + "s"
	dataTypeFormat := "%-" + func() string {
		max := 0
		for _, column := range s.Columns {
			if len(column.DataType) > max {
				max = len(column.DataType)
			}
		}
		return strconv.Itoa(max)
	}() + "s"
	for i, column := range s.Columns {
		log.Printf("🚧: column: %+v\n", column)
		str += indent + fmt.Sprintf(columnNameFormat, column.Name) + " " + fmt.Sprintf(dataTypeFormat, column.DataType) + " " + column.ColumnConstraint
		if i < len(s.Columns)-1 || hasConstraints {
			str += ",\n"
		}
	}

	for _, constraint := range s.Constraints {
		log.Printf("🚧: constraint: %+v\n", constraint)
		str += indent + constraint.String()
	}

	str += "\n)\n"

	return str
}

// TableColumn はテーブルのカラムを表す構造体です。
type TableColumn struct {
	Name             Literal
	NameWidth        int
	DataType         string
	ColumnConstraint string
}

// TableConstraint はテーブルの制約を表す構造体です。
type TableConstraint struct {
	ConstraintType string
	Columns        []Literal
}

func (s *TableConstraint) String() string {
	str := s.ConstraintType
	if len(s.Columns) > 0 {
		str += " ("
		for i, column := range s.Columns {
			str += column.String()
			if i < len(s.Columns)-1 {
				str += ", "
			}
		}
		str += ")"
	}
	return str
}
