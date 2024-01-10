package mysql

import (
	"github.com/kunitsucom/util.go/exp/database/sql/ddl/internal"
	stringz "github.com/kunitsucom/util.go/strings"
)

const (
	Indent        = "    "
	CommentPrefix = "-- "
)

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
	GetNameForDiff() string
	String() string
}

type DDL struct {
	Stmts []Stmt
}

func (d *DDL) String() string {
	if d == nil {
		return ""
	}
	return stringz.JoinStringers("", d.Stmts...)
}

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

func (i *Ident) StringForDiff() string {
	if i == nil {
		return ""
	}
	return i.Name
}

type ColumnIdent struct {
	Ident *Ident
	Order *Order
}

type Order struct{ Desc bool }

func (i *ColumnIdent) GoString() string { return internal.GoString(*i) }

func (i *ColumnIdent) String() string {
	str := i.Ident.String()
	if i.Order != nil {
		if i.Order.Desc {
			str += " DESC"
		}
		// MEMO: If not DESC, it is ASC by default.
		// else {
		// str += " ASC"
		// }
	}
	return str
}

func (i *ColumnIdent) StringForDiff() string {
	str := i.Ident.StringForDiff()
	if i.Order != nil && i.Order.Desc {
		str += " DESC"
	}
	// MEMO: If not DESC, it is ASC by default.
	// else {
	// str += " ASC"
	// }
	return str
}

type DataType struct {
	Name   string
	Type   TokenType
	Idents []*Ident
}

func (s *DataType) String() string {
	if s == nil {
		return ""
	}
	str := s.Name
	if len(s.Idents) > 0 {
		str += "(" + stringz.JoinStringers(", ", s.Idents...) + ")"
	}
	return str
}

func (s *DataType) StringForDiff() string {
	if s == nil {
		return ""
	}
	var str string
	if s.Type != "" {
		str += string(s.Type)
	} else {
		str += string(TOKEN_ILLEGAL)
	}

	if len(s.Idents) > 0 {
		str += "("
		for i, ident := range s.Idents {
			if i > 0 {
				str += ", "
			}
			str += ident.StringForDiff()
		}
		str += ")"
	}

	return str
}
