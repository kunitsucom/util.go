package postgres

import stringz "github.com/kunitsucom/util.go/strings"

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
	String() string
}

type DDL struct {
	Stmts []Stmt
}

func (d *DDL) String() string { return stringz.JoinStringers("", d.Stmts...) }

type IdentBuilder struct {
	QuotationMarks []string
}

type Ident struct {
	Name          string
	QuotationMark string
	Raw           string
}

func (i Ident) IsQuoted() bool {
	return i.QuotationMark != ""
}

func (i Ident) String() string {
	return i.Raw
}

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
