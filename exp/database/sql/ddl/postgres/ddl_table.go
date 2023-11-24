package postgres

import (
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
	return append(constraints, constraint)
}

// PrimaryKeyConstraint represents a PRIMARY KEY constraint.
type PrimaryKeyConstraint struct {
	Name    *Ident
	Columns []*Ident
}

func (*PrimaryKeyConstraint) isConstraint()      {}
func (c *PrimaryKeyConstraint) GetName() *Ident  { return c.Name }
func (c *PrimaryKeyConstraint) GoString() string { return internal.GoString(c) }
func (c *PrimaryKeyConstraint) String() string {
	var str string
	if c.Name != nil {
		str += "CONSTRAINT " + c.Name.String() + " " //nolint:goconst
	}
	str += "PRIMARY KEY"
	str += " (" + stringz.JoinStringers(", ", c.Columns...) + ")"
	return str
}

func (c *PrimaryKeyConstraint) StringForDiff() string {
	var str string
	if c.Name != nil {
		str += "CONSTRAINT " + c.Name.Name + " " //nolint:goconst
	}
	str += "PRIMARY KEY"
	str += " ("
	for i, v := range c.Columns {
		if i != 0 {
			str += ", "
		}
		str += v.Name
	}
	str += ")"
	return str
}

// ForeignKeyConstraint represents a FOREIGN KEY constraint.
type ForeignKeyConstraint struct {
	Name       *Ident
	Columns    []*Ident
	Ref        *Ident
	RefColumns []*Ident
}

func (*ForeignKeyConstraint) isConstraint()      {}
func (c *ForeignKeyConstraint) GetName() *Ident  { return c.Name }
func (c *ForeignKeyConstraint) GoString() string { return internal.GoString(c) }
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
		str += "CONSTRAINT " + c.Name.Name + " "
	}
	str += "FOREIGN KEY"
	str += " ("
	for i, v := range c.Columns {
		if i != 0 {
			str += ", "
		}
		str += v.Name
	}
	str += ")"
	str += " REFERENCES " + c.Ref.Name
	str += " ("
	for i, v := range c.RefColumns {
		if i != 0 {
			str += ", "
		}
		str += v.Name
	}
	str += ")"
	return str
}

// UniqueConstraint represents a UNIQUE constraint.
type UniqueConstraint struct {
	Name    *Ident
	Columns []*Ident
}

func (*UniqueConstraint) isConstraint()      {}
func (c *UniqueConstraint) GetName() *Ident  { return c.Name }
func (c *UniqueConstraint) GoString() string { return internal.GoString(c) }
func (c *UniqueConstraint) String() string {
	var str string
	if c.Name != nil {
		str += "CONSTRAINT " + c.Name.String() + " "
	}
	str += "UNIQUE" //nolint:goconst
	str += " (" + stringz.JoinStringers(", ", c.Columns...) + ")"
	return str
}

func (c *UniqueConstraint) StringForDiff() string {
	var str string
	if c.Name != nil {
		str += "CONSTRAINT " + c.Name.Name + " "
	}
	str += "UNIQUE"
	str += " ("
	for i, v := range c.Columns {
		if i != 0 {
			str += ", "
		}
		str += v.Name
	}
	str += ")"
	return str
}

// CheckConstraint represents a CHECK constraint.
type CheckConstraint struct {
	Name *Ident
	Expr []*Ident
}

func (*CheckConstraint) isConstraint()      {}
func (c *CheckConstraint) GetName() *Ident  { return c.Name }
func (c *CheckConstraint) GoString() string { return internal.GoString(c) }
func (c *CheckConstraint) String() string {
	var str string
	if c.Name != nil {
		str += "CONSTRAINT " + c.Name.String() + " "
	}
	str += "CHECK" //nolint:goconst
	str += " (" + stringz.JoinStringers(" ", c.Expr...) + ")"
	return str
}

func (c *CheckConstraint) StringForDiff() string {
	var str string
	if c.Name != nil {
		str += "CONSTRAINT " + c.Name.Name + " "
	}
	str += "CHECK"
	str += " ("
	for i, v := range c.Expr {
		if i != 0 {
			str += " "
		}
		str += v.Name
	}
	str += ")"
	return str
}

type Column struct {
	Name     *Ident
	DataType *DataType
	Default  *Default
	NotNull  bool
}

type Default struct {
	Value *DefaultValue
}

func (d *DefaultValue) Append(idents ...*Ident) *DefaultValue {
	if d == nil {
		d = &DefaultValue{Idents: idents}
		return d
	}
	d.Idents = append(d.Idents, idents...)
	return d
}

type DefaultValue struct {
	Idents []*Ident
}

//nolint:cyclop
func (d *DefaultValue) String() string {
	if d == nil || len(d.Idents) == 0 {
		return ""
	}

	var str string
	for i := range d.Idents {
		switch {
		case i != 0 && (d.Idents[i-1].String() == "||" || d.Idents[i].String() == "||"):
			str += " "
		case i == 0 ||
			d.Idents[i-1].String() == "(" || d.Idents[i].String() == "(" ||
			d.Idents[i].String() == ")" ||
			d.Idents[i-1].String() == "::" || d.Idents[i].String() == "::" ||
			d.Idents[i].String() == ",":
			// noop
		default:
			str += " "
		}
		str += d.Idents[i].String()
	}

	return str
}

func (d *Default) GoString() string { return internal.GoString(d) }

func (d *Default) String() string {
	if d == nil {
		return ""
	}
	if d.Value != nil {
		return "DEFAULT " + d.Value.String() //nolint:goconst
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
			str += v.Name
		}
		return str
	}
	return ""
}

func (c *Column) String() string {
	str := c.Name.String() + " " +
		c.DataType.String()
	if c.Default != nil {
		str += " " + c.Default.String()
	}
	if c.NotNull {
		str += " NOT NULL"
	}
	return str
}

func (c *Column) GoString() string { return internal.GoString(c) }

type Option struct {
	Str string
}

func (o *Option) GoString() string { return internal.GoString(o) }
