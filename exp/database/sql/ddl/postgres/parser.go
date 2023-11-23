package postgres

// MEMO: https://www.postgresql.org/docs/current/ddl-constraints.html
// MEMO: https://www.postgresql.jp/docs/11/ddl-constraints.html

import (
	"fmt"
	"runtime"
	"strings"

	errorz "github.com/kunitsucom/util.go/errors"
	"github.com/kunitsucom/util.go/exp/database/sql/ddl"
	"github.com/kunitsucom/util.go/exp/database/sql/ddl/internal"
	filepathz "github.com/kunitsucom/util.go/path/filepath"
)

//nolint:gochecknoglobals
var quotationMarks = []string{`"`}

func NewIdent(raw string) *Ident {
	for _, q := range quotationMarks {
		if strings.HasPrefix(raw, q) && strings.HasSuffix(raw, q) {
			return &Ident{
				Name:          strings.Trim(raw, q),
				QuotationMark: q,
				Raw:           raw,
			}
		}
	}

	return &Ident{
		Name:          raw,
		QuotationMark: "",
		Raw:           raw,
	}
}

// Parser はSQL文を解析するパーサーです。
type Parser struct {
	l            *Lexer
	currentToken Token
	peekToken    Token
}

// NewParser は新しいParserを生成します。
func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l: l,
	}

	return p
}

// nextToken は次のトークンを読み込みます。
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()

	_, file, line, _ := runtime.Caller(1)
	internal.TraceLog.Printf("🪲: nextToken: caller=%s:%d currentToken: %#v, peekToken: %#v", filepathz.Short(file), line, p.currentToken, p.peekToken)
}

// Parse はSQL文を解析します。
func (p *Parser) Parse() (*DDL, error) { //nolint:ireturn
	p.nextToken() // current = ""
	p.nextToken() // current = CREATE or ALTER or ...

	d := &DDL{}

LabelDDL:
	for {
		switch p.currentToken.Type { //nolint:exhaustive
		case TOKEN_CREATE:
			stmt, err := p.parseCreateStatement()
			if err != nil {
				return nil, errorz.Errorf("parseCreateStatement: %w", err)
			}
			d.Stmts = append(d.Stmts, stmt)
		case TOKEN_CLOSE_PAREN:
			// do nothing
		case TOKEN_SEMICOLON:
			// do nothing
		case TOKEN_EOF:
			break LabelDDL
		default:
			return nil, errorz.Errorf("token=%s: %w", p.currentToken.Type, ddl.ErrUnexpectedToken)
		}

		p.nextToken()
	}
	return d, nil
}

func (p *Parser) parseCreateStatement() (Stmt, error) { //nolint:ireturn
	p.nextToken() // current = TABLE or INDEX or ...

	switch p.currentToken.Type { //nolint:exhaustive
	case TOKEN_TABLE:
		return p.parseCreateTableStmt()
	default:
		return nil, errorz.Errorf("token=%s: %w", p.currentToken.Type, ddl.ErrUnexpectedToken)
	}
}

//nolint:cyclop
func (p *Parser) parseCreateTableStmt() (*CreateTableStmt, error) { //nolint:ireturn
	createTableStmt := &CreateTableStmt{
		Indent: Indent,
	}

	p.nextToken() // current = table_name

	if p.currentToken.Type != TOKEN_IDENT {
		return nil, errorz.Errorf("token=%s: %w", p.currentToken.Type, ddl.ErrUnexpectedToken)
	}

	createTableStmt.Name = NewIdent(p.currentToken.Literal.Str)

	p.nextToken() // current = (

	if p.currentToken.Type != TOKEN_OPEN_PAREN {
		return nil, errorz.Errorf("token=%s: %w", p.currentToken.Type, ddl.ErrUnexpectedToken)
	}

	p.nextToken() // current = column_name

LabelColumns:
	for {
		switch { //nolint:exhaustive
		case p.currentToken.Type == TOKEN_IDENT:
			column, constraints, err := p.parseColumn(createTableStmt.Name)
			if err != nil {
				return nil, errorz.Errorf("parseColumn: %w", err)
			}
			createTableStmt.Columns = append(createTableStmt.Columns, column)
			if len(constraints) > 0 {
				createTableStmt.Constraints = append(createTableStmt.Constraints, constraints...)
			}
		case isConstraint(p.currentToken.Type):
			constraint, err := p.parseTableConstraint(createTableStmt.Name)
			if err != nil {
				return nil, errorz.Errorf("parseConstraint: %w", err)
			}
			createTableStmt.Constraints = append(createTableStmt.Constraints, constraint)
		case p.currentToken.Type == TOKEN_COMMA:
			p.nextToken()
			continue
		case p.currentToken.Type == TOKEN_CLOSE_PAREN:
			switch p.peekToken.Type { //nolint:exhaustive
			case TOKEN_SEMICOLON, TOKEN_EOF:
				break LabelColumns
			default:
				return nil, errorz.Errorf("p.peekToken=%#v: %w", p.peekToken, ddl.ErrUnexpectedToken)
			}
		default:
			return nil, errorz.Errorf("currentToken=%#v: %w", p.currentToken, ddl.ErrUnexpectedToken)
		}
	}

	return createTableStmt, nil
}

//nolint:funlen,cyclop
func (p *Parser) parseColumn(tableName *Ident) (*Column, []Constraint, error) {
	column := &Column{}
	constraints := make([]Constraint, 0)

	if p.currentToken.Type != TOKEN_IDENT {
		return nil, nil, errorz.Errorf("token=%s: %w", p.currentToken.Type, ddl.ErrUnexpectedToken)
	}

	column.Name = NewIdent(p.currentToken.Literal.Str)

	p.nextToken() // current = DATA_TYPE

	switch { //nolint:exhaustive
	case isDataType(p.currentToken.Type):
		column.DataType = &DataType{Name: string(p.currentToken.Type)}
		if p.peekToken.Type == TOKEN_OPEN_PAREN {
			p.nextToken() // current = (
			p.nextToken() // current = 128
			column.DataType.Size = p.currentToken.Literal.Str
			if err := p.checkPeekToken(TOKEN_CLOSE_PAREN); err != nil {
				return nil, nil, errorz.Errorf("checkPeekToken: %w", err)
			}
			p.nextToken() // current = )
		}

		p.nextToken() // current = DEFAULT or NOT or NULL or PRIMARY or UNIQUE or COMMA or ...
	LabelDefaultNotNull:
		for {
			switch p.currentToken.Type { //nolint:exhaustive
			case TOKEN_NOT:
				if err := p.checkPeekToken(TOKEN_NULL); err != nil {
					return nil, nil, errorz.Errorf("checkPeekToken: %w", err)
				}
				p.nextToken() // current = NULL
				column.NotNull = true
			case TOKEN_NULL:
				column.NotNull = false
			case TOKEN_DEFAULT:
				def, err := p.parseColumnDefault()
				if err != nil {
					return nil, nil, errorz.Errorf("parseColumnDefault: %w", err)
				}
				column.Default = def
			default:
				break LabelDefaultNotNull
			}

			p.nextToken()
		}

		cs, err := p.parseColumnConstraints(tableName, column)
		if err != nil {
			return nil, nil, errorz.Errorf("parseConstraint: %w", err)
		}
		if len(cs) > 0 {
			constraints = append(constraints, cs...)
		}
	default:
		return nil, nil, errorz.Errorf("currentToken=%#v: %w", p.currentToken, ddl.ErrUnexpectedToken)
	}

	return column, constraints, nil
}

//nolint:cyclop
func (p *Parser) parseColumnDefault() (*Default, error) {
	def := &Default{}

LabelDefault:
	for {
		switch { //nolint:exhaustive
		case p.peekToken.Type == TOKEN_IDENT:
			def = &Default{Value: NewIdent(p.peekToken.Literal.Str)}
		case p.peekToken.Type == TOKEN_OPEN_PAREN:
			p.nextToken() // current = (
		LabelExpr:
			for {
				switch p.currentToken.Type { //nolint:exhaustive
				case TOKEN_OPEN_PAREN:
					def = &Default{}
					// do nothing
				case TOKEN_IDENT:
					def.Expr = append(def.Expr, NewIdent(p.currentToken.Literal.Str))
				case TOKEN_CLOSE_PAREN:
					break LabelExpr
				default:
					return nil, errorz.Errorf("currentToken=%#v: %w", p.currentToken, ddl.ErrUnexpectedToken)
				}

				p.nextToken()
			}
		case isConstraint(p.peekToken.Type) ||
			p.peekToken.Type == TOKEN_NOT ||
			p.peekToken.Type == TOKEN_NULL ||
			p.peekToken.Type == TOKEN_COMMA ||
			p.peekToken.Type == TOKEN_CLOSE_PAREN:
			break LabelDefault
		default:
			return nil, errorz.Errorf("peekToken=%#v: %w", p.peekToken, ddl.ErrUnexpectedToken)
		}

		p.nextToken()
	}

	return def, nil
}

//nolint:cyclop,funlen,gocognit
func (p *Parser) parseColumnConstraints(tableName *Ident, column *Column) ([]Constraint, error) {
	constraints := make([]Constraint, 0)

LabelConstraints:
	for {
		switch p.currentToken.Type { //nolint:exhaustive
		case TOKEN_PRIMARY:
			if err := p.checkPeekToken(TOKEN_KEY); err != nil {
				return nil, errorz.Errorf("checkPeekToken: %w", err)
			}
			p.nextToken() // current = KEY
			constraints = append(constraints, &PrimaryKeyConstraint{
				Name:    NewIdent(fmt.Sprintf("%s_pkey", tableName.Name)),
				Columns: []*Ident{column.Name},
			})
		case TOKEN_REFERENCES:
			if err := p.checkPeekToken(TOKEN_IDENT); err != nil {
				return nil, errorz.Errorf("checkPeekToken: %w", err)
			}
			p.nextToken() // current = table_name
			constraint := &ForeignKeyConstraint{
				Name:    NewIdent(fmt.Sprintf("%s_%s_fkey", tableName.Name, column.Name.Name)),
				Ref:     NewIdent(p.currentToken.Literal.Str),
				Columns: []*Ident{column.Name},
			}
			p.nextToken() // current = (
			idents, err := p.parseIdents()
			if err != nil {
				return nil, errorz.Errorf("parseIdents: %w", err)
			}
			constraint.RefColumns = idents
			constraints = append(constraints, constraint)
		case TOKEN_UNIQUE:
			constraints = append(constraints, &UniqueConstraint{
				Name:    NewIdent(fmt.Sprintf("%s_unique_%s", tableName.Name, column.Name.Name)),
				Columns: []*Ident{column.Name},
			})
		case TOKEN_CHECK:
			if err := p.checkPeekToken(TOKEN_OPEN_PAREN); err != nil {
				return nil, errorz.Errorf("checkPeekToken: %w", err)
			}
			p.nextToken() // current = (
			constraint := &CheckConstraint{
				Name: NewIdent(fmt.Sprintf("%s_%s_check", tableName.Name, column.Name.Name)),
			}
		LabelCheck:
			for {
				switch p.currentToken.Type { //nolint:exhaustive
				case TOKEN_OPEN_PAREN:
					// do nothing
				case TOKEN_IDENT:
					constraint.Expr = append(constraint.Expr, NewIdent(p.currentToken.Literal.Str))
				case TOKEN_EQUAL, TOKEN_GREATER, TOKEN_LESS:
					value := p.currentToken.Literal.Str
					switch p.peekToken.Type { //nolint:exhaustive
					case TOKEN_EQUAL, TOKEN_GREATER, TOKEN_LESS:
						value += p.peekToken.Literal.Str
						p.nextToken()
					}
					constraint.Expr = append(constraint.Expr, NewIdent(value))
				case TOKEN_CLOSE_PAREN:
					break LabelCheck
				default:
					return nil, errorz.Errorf("currentToken=%#v: %w", p.currentToken, ddl.ErrUnexpectedToken)
				}

				p.nextToken()
			}
			constraints = append(constraints, constraint)
		case TOKEN_IDENT, TOKEN_COMMA, TOKEN_CLOSE_PAREN:
			break LabelConstraints
		default:
			return nil, errorz.Errorf("currentToken=%#v: %w", p.currentToken, ddl.ErrUnexpectedToken)
		}

		p.nextToken()
	}

	return constraints, nil
}

//nolint:funlen,cyclop,gocognit
func (p *Parser) parseTableConstraint(tableName *Ident) (Constraint, error) { //nolint:ireturn
	var constraintName *Ident
	if p.currentToken.Type == TOKEN_CONSTRAINT {
		p.nextToken() // current = constraint_name
		if p.currentToken.Type != TOKEN_IDENT {
			return nil, errorz.Errorf("token=%s: %w", p.currentToken.Type, ddl.ErrUnexpectedToken)
		}
		constraintName = NewIdent(p.currentToken.Literal.Str)
		p.nextToken() // current = PRIMARY or UNIQUE or CHECK
	}

	switch p.currentToken.Type { //nolint:exhaustive
	case TOKEN_PRIMARY:
		if err := p.checkPeekToken(TOKEN_KEY); err != nil {
			return nil, errorz.Errorf("checkPeekToken: %w", err)
		}
		p.nextToken() // current = KEY
		if err := p.checkPeekToken(TOKEN_OPEN_PAREN); err != nil {
			return nil, errorz.Errorf("checkPeekToken: %w", err)
		}
		p.nextToken() // current = (
		idents, err := p.parseIdents()
		if err != nil {
			return nil, errorz.Errorf("parseIdents: %w", err)
		}
		if constraintName == nil {
			constraintName = NewIdent(fmt.Sprintf("%s_pkey", tableName.Name))
		}
		return &PrimaryKeyConstraint{
			Name:    constraintName,
			Columns: idents,
		}, nil
	case TOKEN_FOREIGN:
		if err := p.checkPeekToken(TOKEN_KEY); err != nil {
			return nil, errorz.Errorf("checkPeekToken: %w", err)
		}
		p.nextToken() // current = KEY
		if err := p.checkPeekToken(TOKEN_OPEN_PAREN); err != nil {
			return nil, errorz.Errorf("checkPeekToken: %w", err)
		}
		p.nextToken() // current = (
		idents, err := p.parseIdents()
		if err != nil {
			return nil, errorz.Errorf("parseIdents: %w", err)
		}
		if err := p.checkPeekToken(TOKEN_REFERENCES); err != nil {
			return nil, errorz.Errorf("checkPeekToken: %w", err)
		}
		p.nextToken() // current = REFERENCES
		if err := p.checkPeekToken(TOKEN_OPEN_PAREN); err != nil {
			return nil, errorz.Errorf("checkPeekToken: %w", err)
		}

		p.nextToken() // current = table_name
		if p.currentToken.Type != TOKEN_IDENT {
			return nil, errorz.Errorf("token=%s: %w", p.currentToken.Type, ddl.ErrUnexpectedToken)
		}
		refName := NewIdent(p.currentToken.Literal.Str)

		p.nextToken() // current = (
		identsRef, err := p.parseIdents()
		if err != nil {
			return nil, errorz.Errorf("parseIdents: %w", err)
		}
		if constraintName == nil {
			name := tableName.Name
			for _, ident := range idents {
				name += fmt.Sprintf("_%s", ident.Name)
			}
			name += "_fkey"
			constraintName = NewIdent(name)
		}
		return &ForeignKeyConstraint{
			Name:       constraintName,
			Columns:    idents,
			Ref:        refName,
			RefColumns: identsRef,
		}, nil

	case TOKEN_UNIQUE:
		idents, err := p.parseIdents()
		if err != nil {
			return nil, errorz.Errorf("parseIdents: %w", err)
		}
		if constraintName == nil {
			name := fmt.Sprintf("%s_unique", tableName.Name)
			for _, ident := range idents {
				name += fmt.Sprintf("_%s", ident.Name)
			}
			constraintName = NewIdent(name)
		}
		return &UniqueConstraint{
			Name:    constraintName,
			Columns: idents,
		}, nil
	default:
		return nil, errorz.Errorf("currentToken=%s: %w", p.currentToken.Type, ddl.ErrUnexpectedToken)
	}
}

func (p *Parser) parseIdents() ([]*Ident, error) {
	idents := make([]*Ident, 0)

LabelIdents:
	for {
		switch p.currentToken.Type { //nolint:exhaustive
		case TOKEN_OPEN_PAREN:
			// do nothing
		case TOKEN_IDENT:
			idents = append(idents, NewIdent(p.currentToken.Literal.Str))
		case TOKEN_COMMA:
			// do nothing
		case TOKEN_CLOSE_PAREN:
			p.nextToken()
			break LabelIdents
		default:
			return nil, errorz.Errorf("token=%s: %w", p.currentToken.Type, ddl.ErrUnexpectedToken)
		}
		p.nextToken()
	}

	return idents, nil
}

func isDataType(tokenType TokenType) bool {
	switch tokenType { //nolint:exhaustive
	case TOKEN_INTEGER,
		TOKEN_UUID,
		TOKEN_VARYING, TOKEN_TEXT,
		TOKEN_TIMESTAMP, TOKEN_TIMESTAMPZ:
		return true
	default:
		return false
	}
}

func isConstraint(tokenType TokenType) bool {
	switch tokenType { //nolint:exhaustive
	case TOKEN_CONSTRAINT,
		TOKEN_PRIMARY, TOKEN_KEY,
		TOKEN_FOREIGN, TOKEN_REFERENCES,
		TOKEN_UNIQUE,
		TOKEN_CHECK:
		return true
	default:
		return false
	}
}

//nolint:unused
func (p *Parser) checkCurrentToken(expected TokenType) error {
	if expected != p.currentToken.Type {
		return errorz.Errorf("current token: expected=%s, but got=%s: %w", expected, p.currentToken.Type, ddl.ErrUnexpectedToken)
	}
	return nil
}

func (p *Parser) checkPeekToken(expected TokenType) error {
	if expected != p.peekToken.Type {
		return errorz.Errorf("peek token: expected=%s, but got=%s: %w", expected, p.peekToken.Type, ddl.ErrUnexpectedToken)
	}
	return nil
}
