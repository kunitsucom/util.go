package postgres

import (
	"github.com/kunitsucom/util.go/exp/database/sql/ddl/internal"
	stringz "github.com/kunitsucom/util.go/strings"
)

type DataType struct {
	Name string
	Size string
}

func (s DataType) String() string {
	str := s.Name
	if s.Size != "" {
		str += "(" + s.Size + ")"
	}
	return str
}

type Constraint interface {
	isConstraint()
	GetName() *Ident
	String() string
	GoString() string
}

type PrimaryKeyConstraint struct {
	Name    *Ident
	Columns []*Ident
}

func (*PrimaryKeyConstraint) isConstraint()     {}
func (c *PrimaryKeyConstraint) GetName() *Ident { return c.Name }
func (c PrimaryKeyConstraint) GoString() string { return internal.GoString(c) }
func (c *PrimaryKeyConstraint) String() string {
	var str string
	if c.Name != nil {
		str += "CONSTRAINT " + c.Name.String() + " "
	}
	str += "PRIMARY KEY"
	str += " (" + stringz.JoinStringers(", ", c.Columns...) + ")"
	return str
}

type ForeignKeyConstraint struct {
	Name       *Ident
	Columns    []*Ident
	Ref        *Ident
	RefColumns []*Ident
}

func (*ForeignKeyConstraint) isConstraint()     {}
func (c *ForeignKeyConstraint) GetName() *Ident { return c.Name }
func (c ForeignKeyConstraint) GoString() string { return internal.GoString(c) }
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

type UniqueConstraint struct {
	Name    *Ident
	Columns []*Ident
}

func (*UniqueConstraint) isConstraint()     {}
func (c *UniqueConstraint) GetName() *Ident { return c.Name }
func (c UniqueConstraint) GoString() string { return internal.GoString(c) }
func (c *UniqueConstraint) String() string {
	var str string
	if c.Name != nil {
		str += "CONSTRAINT " + c.Name.String() + " "
	}
	str += "UNIQUE"
	str += " (" + stringz.JoinStringers(", ", c.Columns...) + ")"
	return str
}

type CheckConstraint struct {
	Name *Ident
	Expr []*Ident
}

func (*CheckConstraint) isConstraint()     {}
func (c *CheckConstraint) GetName() *Ident { return c.Name }
func (c CheckConstraint) GoString() string { return internal.GoString(c) }
func (c *CheckConstraint) String() string {
	var str string
	if c.Name != nil {
		str += "CONSTRAINT " + c.Name.String() + " "
	}
	str += "CHECK"
	str += " (" + stringz.JoinStringers(" ", c.Expr...) + ")"
	return str
}

type Column struct {
	Name     *Ident
	DataType DataType
	NotNull  bool
	Default  *Default
}

type Default struct {
	Value *Ident
	Expr  []*Ident
}

func (d *Default) String() string {
	if d.Value != nil {
		return "DEFAULT " + d.Value.String()
	}
	if len(d.Expr) > 0 {
		return "DEFAULT " + "(" + stringz.JoinStringers(", ", d.Expr...) + ")"
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

func (c *Column) GoString() string { return internal.GoString(*c) }

type Option struct {
	Str string
}
