package cockroachdb

import (
	"github.com/kunitsucom/util.go/exp/database/sql/ddl/internal"
	stringz "github.com/kunitsucom/util.go/strings"
)

const Indent = "    "

type Verb string

const (
	VerbCreate   Verb = "CREATE"
	VerbAlter    Verb = "ALTER"
	VerbDrop     Verb = "DROP"
	VerbRename   Verb = "RENAME"
	VerbTruncate Verb = "TRUNCATE"
)

type Object string

const (
	ObjectTable Object = "TABLE"
	ObjectIndex Object = "INDEX"
	ObjectView  Object = "VIEW"
)

type Action string

const (
	ActionAdd    Action = "ADD"
	ActionDrop   Action = "DROP"
	ActionAlter  Action = "ALTER"
	ActionRename Action = "RENAME"
)

type Stmt interface {
	isStmt()
	GetPlainName() string
	String() string
}

type DDL struct {
	Stmts []Stmt
}

func (d *DDL) String() string { return stringz.JoinStringers("", d.Stmts...) }

type Ident struct {
	Name          string
	QuotationMark string
	Raw           string
}

func (i *Ident) GoString() string { return internal.GoString(*i) }

func (i *Ident) String() string {
	if i == nil {
		return ""
	}
	return i.Raw
}

func (i *Ident) PlainString() string {
	if i == nil {
		return ""
	}
	return i.Name
}

type ColumnIdent struct {
	Ident *Ident
	Order *Order
}

type Order struct {
	Desc bool
}

func (i *ColumnIdent) GoString() string { return internal.GoString(*i) }

func (i *ColumnIdent) String() string {
	str := i.Ident.String()
	if i.Order != nil {
		if i.Order.Desc {
			str += " DESC"
		} else {
			str += " ASC"
		}
	}
	return str
}

func (i *ColumnIdent) PlainString() string {
	str := i.Ident.PlainString()
	if i.Order != nil {
		if i.Order.Desc {
			str += " DESC"
		} else {
			str += " ASC"
		}
	}
	return str
}

type DataType struct {
	Name string
	Size string
}

func (s *DataType) String() string {
	if s == nil {
		return ""
	}
	str := s.Name
	if s.Size != "" {
		str += "(" + s.Size + ")"
	}
	return str
}
