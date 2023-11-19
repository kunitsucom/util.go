package postgres

import (
	"fmt"
	"log"
	"runtime"

	errorz "github.com/kunitsucom/util.go/errors"
	"github.com/kunitsucom/util.go/exp/database/sql/ddl"
	filepathz "github.com/kunitsucom/util.go/path/filepath"
)

// Parser はSQL文を解析するパーサーです。
type Parser struct {
	l            *Lexer
	errors       []string
	currentToken Token
	peekToken    Token
}

// NewParser は新しいParserを生成します。
func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	return p
}

// nextToken は次のトークンを読み込みます。
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()

	_, file, line, _ := runtime.Caller(1)
	log.Printf("⏭️: caller=%s:%d currentToken: %+v, peekToken: %+v", filepathz.Short(file), line, p.currentToken, p.peekToken)
}

// ParseStatement はSQL文を解析します。
func (p *Parser) ParseStatement() (Stmt, error) { //nolint:ireturn
	p.nextToken()
	p.nextToken()

	switch p.currentToken.Type { //nolint:exhaustive
	case TOKEN_CREATE:
		return p.parseCreateStatement()
	default:
		return nil, errorz.Errorf("token=%s: %w", p.currentToken.Type, ddl.ErrUnexpectedToken)
	}
}

func (p *Parser) parseCreateStatement() (Stmt, error) { //nolint:ireturn
	p.nextToken()

	switch p.currentToken.Type { //nolint:exhaustive
	case TOKEN_TABLE:
		return p.parseCreateTableStatement()
	default:
		return nil, errorz.Errorf("token=%s: %w", p.currentToken.Type, ddl.ErrUnexpectedToken)
	}
}

// parseCreateTableStatement はCREATE TABLE文を解析します。
func (p *Parser) parseCreateTableStatement() (*CreateTableStmt, error) {
	p.nextToken()

	stmt := &CreateTableStmt{}

	if p.currentTokenIs(TOKEN_IF) {
		if err := p.expectPeekToken(TOKEN_NOT); err != nil {
			return nil, errorz.Errorf("expectPeekToken: %w", err)
		}
		p.nextToken()
		if err := p.expectPeekToken(TOKEN_EXISTS); err != nil {
			return nil, errorz.Errorf("expectPeekToken: %w", err)
		}
		p.nextToken()
		stmt.IfNotExists = true
		p.nextToken()
	}

	// テーブル名の解析
	if err := p.expectCurrentToken(TOKEN_IDENT); err != nil {
		return nil, errorz.Errorf("expectCurrentToken: %w", err)
	}
	stmt.TableName = p.currentToken.Literal

	// カラム定義の開始
	if err := p.expectPeekToken(TOKEN_OPEN_PAREN); err != nil {
		return nil, errorz.Errorf("expected=%s, actual=%s: %w", TOKEN_OPEN_PAREN, p.currentToken.Type, ddl.ErrUnexpectedToken)
	}
	p.nextToken()

	// カラムの解析
	var err error
	stmt.Columns, stmt.Constraints, err = p.parseColumns()
	if err != nil {
		return nil, errorz.Errorf("parseColumns: %w", err)
	}

	return stmt, nil
}

// parseColumns はカラム定義とテーブル制約を解析します。
func (p *Parser) parseColumns() ([]*TableColumn, []*TableConstraint, error) {
	var columns []*TableColumn
	var constraints []*TableConstraint

	for !p.peekTokenIs(TOKEN_CLOSE_PAREN) {
		p.nextToken()

		// 制約かカラム定義かを判断
		if p.currentTokenIs(TOKEN_IDENT) {
			column, err := p.parseColumn()
			if err != nil {
				return nil, nil, errorz.Errorf("parseColumn: %w", err)
			}
			if column != nil {
				columns = append(columns, column)
			}
		} else {
			constraint, err := p.parseConstraint()
			if err != nil {
				return nil, nil, errorz.Errorf("parseConstraint: %w", err)
			}
			if constraint != nil {
				constraints = append(constraints, constraint)
			}
		}

		if err := p.expectPeekToken(TOKEN_CLOSE_PAREN); err == nil {
			break
		}

		if err := p.expectPeekToken(TOKEN_EOF); err == nil {
			return nil, nil, errorz.Errorf("expectPeekToken: %w", ddl.ErrUnexpectedToken)
		}
	}

	return columns, constraints, nil
}

// parseColumn は個々のカラム定義を解析します。
func (p *Parser) parseColumn() (*TableColumn, error) {
	column := &TableColumn{}
	defer log.Printf("defer: column: %+v", column)

	for {
		if p.peekTokenIs(TOKEN_COMMA) || p.peekTokenIs(TOKEN_CLOSE_PAREN) {
			break
		}

		if p.currentTokenIs(TOKEN_IDENT) {
			column.Name = p.currentToken.Literal
			p.nextToken()
		}

		column.ColumnConstraint += " " + p.currentToken.Literal.Str
	}

	if err := p.expectCurrentToken(TOKEN_IDENT); err != nil {
		return nil, errorz.Errorf("expectCurrentToken: %w", err)
	}
	column.Name = p.currentToken.Literal

	p.nextToken()
	column.DataType = p.currentToken.Literal.Str

	p.nextToken()
	if err := p.expectCurrentToken(TOKEN_IDENT); err == nil {
		p.nextToken()
		column.DataType += " " + p.currentToken.Literal.Str
	}

	for !(p.peekTokenIs(TOKEN_COMMA) || p.peekTokenIs(TOKEN_CLOSE_PAREN)) {
		p.nextToken()
		column.ColumnConstraint += " " + p.currentToken.Literal.Str
	}

	if p.peekTokenIs(TOKEN_COMMA) || p.peekTokenIs(TOKEN_CLOSE_PAREN) {
		p.nextToken()
	}

	return column, nil
}

// parseConstraint はテーブルの制約を解析します。
func (p *Parser) parseConstraint() (*TableConstraint, error) {
	constraint := &TableConstraint{}
	defer log.Printf("defer: constraint: %+v", constraint)

	switch {
	case p.currentTokenIs(TOKEN_PRIMARY):
		p.nextToken()
		switch {
		case p.currentTokenIs(TOKEN_KEY):
			constraint.ConstraintType = "PRIMARY KEY"
			p.nextToken()
			if err := p.expectCurrentToken(TOKEN_OPEN_PAREN); err != nil {
				return nil, errorz.Errorf("expectCurrentToken: %w", err)
			}
			p.nextToken()
			constraint.Columns = append(constraint.Columns, p.currentToken.Literal)

			// 複合キーの場合（例：PRIMARY KEY (id, name)）
			for p.peekTokenIs(TOKEN_COMMA) {
				p.nextToken()
				p.nextToken()
				constraint.Columns = append(constraint.Columns, p.currentToken.Literal)
			}

			if err := p.expectPeekToken(TOKEN_CLOSE_PAREN); err != nil {
				return nil, errorz.Errorf("expectPeekToken: %w", err)
			}
		}
	// 他の制約タイプ（例：FOREIGN KEY, UNIQUE）もここに追加...
	default:
		return nil, errorz.Errorf("currentToken=%s, peekToken=%s: %w", p.currentToken.Type, p.peekToken.Type, ddl.ErrUnexpectedToken)
	}

	return constraint, nil
}

// Helper functions.
func (p *Parser) currentTokenIs(t TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) expectCurrentToken(t TokenType) error {
	if p.currentToken.Type == t {
		return nil
	}

	return errorz.Errorf("expected=%s, actual=%s: %w", t, p.currentToken.Type, ddl.ErrUnexpectedToken)
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeekToken(t TokenType) error {
	if p.peekToken.Type == t {
		return nil
	}

	return errorz.Errorf("expected=%s, actual=%s: %w", t, p.currentToken.Type, ddl.ErrUnexpectedToken)
}

// peekError は次のトークンが期待するものでない場合にエラーを追加します。
func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
