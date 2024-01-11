package mysql

import (
	"sort"
	"strings"

	"github.com/kunitsucom/util.go/exp/database/sql/ddl/internal"
	stringz "github.com/kunitsucom/util.go/strings"
)

type Constraint interface {
	isConstraint()
	GetName() *Ident
	GoString() string
	String() string
	StringForDiff() string
}

type Constraints []Constraint

func (constraints Constraints) Append(constraint Constraint) Constraints {
	for i := range constraints {
		if constraints[i].GetName().Name == constraint.GetName().Name {
			constraints[i] = constraint
			return constraints
		}
	}
	constraints = append(constraints, constraint)

	sort.Slice(constraints, func(left, right int) bool {
		_, leftIsPrimaryKeyConstraint := constraints[left].(*PrimaryKeyConstraint)
		switch {
		case leftIsPrimaryKeyConstraint:
			return true
		default:
			return false
		}
	})

	return constraints
}

// PrimaryKeyConstraint represents a PRIMARY KEY constraint.
type PrimaryKeyConstraint struct {
	Name    *Ident
	Columns []*ColumnIdent
}

var _ Constraint = (*PrimaryKeyConstraint)(nil)

func (*PrimaryKeyConstraint) isConstraint()      {}
func (c *PrimaryKeyConstraint) GetName() *Ident  { return c.Name }
func (c *PrimaryKeyConstraint) GoString() string { return internal.GoString(*c) }
func (c *PrimaryKeyConstraint) String() string {
	var str string
	// MEMO: MySQL does not support naming PRIMARY KEY constraints.
	// if c.Name != nil {
	// 	str += "CONSTRAINT " + c.Name.String() + " "
	// }
	str += "PRIMARY KEY"
	str += " (" + stringz.JoinStringers(", ", c.Columns...) + ")"
	return str
}

func (c *PrimaryKeyConstraint) StringForDiff() string {
	var str string
	// MEMO: MySQL does not support naming PRIMARY KEY constraints.
	// if c.Name != nil {
	// 	str += "CONSTRAINT " + c.Name.StringForDiff() + " "
	// }
	str += "PRIMARY KEY"
	str += " ("
	for i, v := range c.Columns {
		if i != 0 {
			str += ", "
		}
		str += v.StringForDiff()
	}
	str += ")"
	return str
}

// ForeignKeyConstraint represents a FOREIGN KEY constraint.
type ForeignKeyConstraint struct {
	Name       *Ident
	Columns    []*ColumnIdent
	Ref        *Ident
	RefColumns []*ColumnIdent
}

var _ Constraint = (*ForeignKeyConstraint)(nil)

func (*ForeignKeyConstraint) isConstraint()      {}
func (c *ForeignKeyConstraint) GetName() *Ident  { return c.Name }
func (c *ForeignKeyConstraint) GoString() string { return internal.GoString(*c) }
func (c *ForeignKeyConstraint) String() string {
	var str string
	if c.Name != nil {
		str += "CONSTRAINT " + c.Name.String() + " "
	}
	str += "FOREIGN KEY"
	str += " (" + stringz.JoinStringers(", ", c.Columns...) + ")"
	str += " REFERENCES " + c.Ref.String()
	str += " (" + stringz.JoinStringers(", ", c.RefColumns...) + ")"
	return str
}

func (c *ForeignKeyConstraint) StringForDiff() string {
	var str string
	if c.Name != nil {
		str += "CONSTRAINT " + c.Name.StringForDiff() + " "
	}
	str += "FOREIGN KEY"
	str += " ("
	for i, v := range c.Columns {
		if i != 0 {
			str += ", "
		}
		str += v.StringForDiff()
	}
	str += ")"
	str += " REFERENCES " + c.Ref.Name
	str += " ("
	for i, v := range c.RefColumns {
		if i != 0 {
			str += ", "
		}
		str += v.StringForDiff()
	}
	str += ")"
	return str
}

// IndexConstraint represents a UNIQUE constraint..
type IndexConstraint struct {
	Name    *Ident
	Unique  bool
	Columns []*ColumnIdent
}

var _ Constraint = (*IndexConstraint)(nil)

func (*IndexConstraint) isConstraint()      {}
func (c *IndexConstraint) GetName() *Ident  { return c.Name }
func (c *IndexConstraint) GoString() string { return internal.GoString(*c) }
func (c *IndexConstraint) String() string {
	var str string
	if c.Unique {
		str += "UNIQUE "
	}
	if c.Name != nil {
		str += "KEY " + c.Name.String() + " "
	}
	str += "(" + stringz.JoinStringers(", ", c.Columns...) + ")"
	return str
}

func (c *IndexConstraint) StringForDiff() string {
	var str string
	if c.Unique {
		str += "UNIQUE "
	}
	if c.Name != nil {
		str += "KEY " + c.Name.StringForDiff() + " "
	}
	str += "("
	for i, v := range c.Columns {
		if i != 0 {
			str += ", "
		}
		str += v.StringForDiff()
	}
	str += ")"
	return str
}

// CheckConstraint represents a CHECK constraint.
type CheckConstraint struct {
	Name *Ident
	Expr *Expr
}

var _ Constraint = (*CheckConstraint)(nil)

func (*CheckConstraint) isConstraint()      {}
func (c *CheckConstraint) GetName() *Ident  { return c.Name }
func (c *CheckConstraint) GoString() string { return internal.GoString(*c) }
func (c *CheckConstraint) String() string {
	var str string
	if c.Name != nil {
		str += "CONSTRAINT " + c.Name.String() + " "
	}
	str += "CHECK "
	str += c.Expr.String()
	return str
}

func (c *CheckConstraint) StringForDiff() string {
	var str string
	if c.Name != nil {
		str += "CONSTRAINT " + c.Name.StringForDiff() + " "
	}
	str += "CHECK "
	for i, v := range c.Expr.Idents {
		if i != 0 {
			str += " "
		}
		str += v.StringForDiff()
	}
	return str
}

func NewObjectName(name string) *ObjectName {
	objName := &ObjectName{}

	tableName := NewRawIdent(name)
	switch name := strings.Split(tableName.Name, "."); len(name) { //nolint:exhaustive
	case 2:
		// CREATE TABLE "schema.table"
		objName.Schema = NewRawIdent(tableName.QuotationMark + name[0] + tableName.QuotationMark)
		objName.Name = NewRawIdent(tableName.QuotationMark + name[1] + tableName.QuotationMark)
	default:
		// CREATE TABLE "table"
		objName.Name = tableName
	}

	return objName
}

type ObjectName struct {
	Schema *Ident
	Name   *Ident
}

func (t *ObjectName) String() string {
	if t == nil {
		return ""
	}
	if t.Schema != nil {
		return t.Name.QuotationMark + t.Schema.StringForDiff() + "." + t.Name.StringForDiff() + t.Name.QuotationMark
	}
	return t.Name.String()
}

func (t *ObjectName) StringForDiff() string {
	if t == nil {
		return ""
	}
	if t.Schema != nil {
		return t.Schema.StringForDiff() + "." + t.Name.StringForDiff()
	}
	return t.Name.StringForDiff()
}

type Column struct {
	Name     *Ident
	DataType *DataType
	Default  *Default
	NotNull  bool
}

type Default struct {
	Value *Expr
}

func (d *Expr) Append(idents ...*Ident) *Expr {
	if d == nil {
		d = &Expr{Idents: idents}
		return d
	}
	d.Idents = append(d.Idents, idents...)
	return d
}

type Expr struct {
	Idents []*Ident
}

//nolint:cyclop
func (d *Expr) String() string {
	if d == nil || len(d.Idents) == 0 {
		return ""
	}

	var str string
	for i := range d.Idents {
		switch {
		case i == 0 ||
			d.Idents[i-1].String() == "(" || d.Idents[i].String() == "(" ||
			d.Idents[i].String() == ")" ||
			d.Idents[i-1].String() == "::" || d.Idents[i].String() == "::" ||
			d.Idents[i-1].String() == ":::" || d.Idents[i].String() == ":::" ||
			d.Idents[i].String() == ",":
			// noop
		default:
			str += " "
		}
		str += d.Idents[i].String()
	}

	return str
}

func (d *Default) GoString() string { return internal.GoString(*d) }

func (d *Default) String() string {
	if d == nil {
		return ""
	}
	if d.Value != nil {
		return "DEFAULT " + d.Value.String()
	}
	return ""
}

func (d *Default) StringForDiff() string {
	if d == nil {
		return ""
	}
	if e := d.Value; e != nil {
		str := "DEFAULT "
		for i, v := range d.Value.Idents {
			if i != 0 {
				str += " "
			}
			str += v.StringForDiff()
		}
		return str
	}
	return ""
}

func (c *Column) String() string {
	str := c.Name.String() + " " +
		c.DataType.String()
	if c.NotNull {
		str += " NOT NULL"
	}
	if c.Default != nil {
		str += " " + c.Default.String()
	}
	return str
}

func (c *Column) GoString() string { return internal.GoString(*c) }

type Option struct {
	Name  string
	Value *Ident
}

func (o *Option) String() string {
	if o.Value == nil {
		return ""
	}
	return o.Name + "=" + o.Value.String()
}

func (o *Option) GoString() string { return internal.GoString(*o) }
