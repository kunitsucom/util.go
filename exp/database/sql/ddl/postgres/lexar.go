package postgres

import (
	"strings"
)

const (
	QuotationChar = '"'
	QuotationStr  = string(QuotationChar)
)

// Token はSQL文のトークンを表す型です。
type Token struct {
	Type    TokenType
	Literal Literal
}

type Literal struct {
	Str string
}

func (l Literal) String() string {
	return l.Str
}

type TokenType string

//nolint:revive,stylecheck
const (
	// SPECIAL TOKENS.
	TOKEN_ILLEGAL TokenType = "ILLEGAL"
	TOKEN_EOF     TokenType = "EOF"

	// SPECIAL CHARACTERS.
	TOKEN_OPEN_PAREN  TokenType = "OPEN_PAREN"  // (
	TOKEN_CLOSE_PAREN TokenType = "CLOSE_PAREN" // )
	TOKEN_COMMA       TokenType = "COMMA"       // ,
	TOKEN_SEMICOLON   TokenType = "SEMICOLON"   // ;
	TOKEN_EQUAL       TokenType = "EQUAL"       // =
	TOKEN_GREATER     TokenType = "GREATER"     // >
	TOKEN_LESS        TokenType = "LESS"        // <

	// VERB.
	TOKEN_CREATE   TokenType = "CREATE"
	TOKEN_ALTER    TokenType = "ALTER"
	TOKEN_DROP     TokenType = "DROP"
	TOKEN_RENAME   TokenType = "RENAME"
	TOKEN_TRUNCATE TokenType = "TRUNCATE"

	// OBJECT.
	TOKEN_TABLE  TokenType = "TABLE"
	TOKEN_INDEX  TokenType = "INDEX"
	TOKEN_VIEW   TokenType = "VIEW"
	TOKEN_IF     TokenType = "IF"
	TOKEN_EXISTS TokenType = "EXISTS"

	// DATA TYPE.
	TOKEN_BOOLEAN    TokenType = "BOOLEAN"
	TOKEN_SMALLINT   TokenType = "SMALLINT"
	TOKEN_INTEGER    TokenType = "INTEGER"
	TOKEN_BIGINT     TokenType = "BIGINT"
	TOKEN_DECIMAL    TokenType = "DECIMAL"
	TOKEN_NUMERIC    TokenType = "NUMERIC"
	TOKEN_REAL       TokenType = "REAL"
	TOKEN_DOUBLE     TokenType = "DOUBLE"
	TOKEN_PRECISION  TokenType = "PRECISION"
	TOKEN_UUID       TokenType = "UUID"
	TOKEN_VARYING    TokenType = "VARYING"
	TOKEN_TEXT       TokenType = "TEXT"
	TOKEN_TIMESTAMP  TokenType = "TIMESTAMP"
	TOKEN_TIMESTAMPZ TokenType = "TIMESTAMPZ"

	// CONSTRAINT.
	TOKEN_CONSTRAINT TokenType = "CONSTRAINT"
	TOKEN_NOT        TokenType = "NOT"
	TOKEN_NULL       TokenType = "NULL"
	TOKEN_PRIMARY    TokenType = "PRIMARY"
	TOKEN_KEY        TokenType = "KEY"
	TOKEN_FOREIGN    TokenType = "FOREIGN"
	TOKEN_REFERENCES TokenType = "REFERENCES"
	TOKEN_UNIQUE     TokenType = "UNIQUE"
	TOKEN_DEFAULT    TokenType = "DEFAULT"
	TOKEN_CHECK      TokenType = "CHECK"

	// IDENTIFIER.
	TOKEN_IDENT TokenType = "IDENT"
)

//nolint:funlen,cyclop,gocognit,gocyclo
func lookupIdent(ident string) TokenType {
	token := strings.ToUpper(ident)
	// MEMO: bash lexar-gen.sh lexar.go | pbcopy
	// START CASES DO NOT EDIT
	switch token {
	case "EQUAL":
		return TOKEN_EQUAL
	case "GREATER":
		return TOKEN_GREATER
	case "LESS":
		return TOKEN_LESS
	case "CREATE":
		return TOKEN_CREATE
	case "ALTER":
		return TOKEN_ALTER
	case "DROP":
		return TOKEN_DROP
	case "RENAME":
		return TOKEN_RENAME
	case "TRUNCATE":
		return TOKEN_TRUNCATE
	case "TABLE":
		return TOKEN_TABLE
	case "INDEX":
		return TOKEN_INDEX
	case "VIEW":
		return TOKEN_VIEW
	case "IF":
		return TOKEN_IF
	case "EXISTS":
		return TOKEN_EXISTS
	case "BOOLEAN":
		return TOKEN_BOOLEAN
	case "SMALLINT":
		return TOKEN_SMALLINT
	case "INTEGER", "INT":
		return TOKEN_INTEGER
	case "BIGINT":
		return TOKEN_BIGINT
	case "DECIMAL":
		return TOKEN_DECIMAL
	case "NUMERIC":
		return TOKEN_NUMERIC
	case "REAL":
		return TOKEN_REAL
	case "DOUBLE":
		return TOKEN_DOUBLE
	case "PRECISION":
		return TOKEN_PRECISION
	case "UUID":
		return TOKEN_UUID
	case "VARYING", "VARCHAR":
		return TOKEN_VARYING
	case "TEXT":
		return TOKEN_TEXT
	case "TIMESTAMP":
		return TOKEN_TIMESTAMP
	case "TIMESTAMPZ":
		return TOKEN_TIMESTAMPZ
	case "CONSTRAINT":
		return TOKEN_CONSTRAINT
	case "NOT":
		return TOKEN_NOT
	case "NULL":
		return TOKEN_NULL
	case "PRIMARY":
		return TOKEN_PRIMARY
	case "KEY":
		return TOKEN_KEY
	case "FOREIGN":
		return TOKEN_FOREIGN
	case "REFERENCES":
		return TOKEN_REFERENCES
	case "UNIQUE":
		return TOKEN_UNIQUE
	case "DEFAULT":
		return TOKEN_DEFAULT
	case "CHECK":
		return TOKEN_CHECK
	default:
		return TOKEN_IDENT
	}
	// END CASES DO NOT EDIT
}

// Lexer はSQL文をトークンに分割するレキサーです。
type Lexer struct {
	input        string
	position     int  // 現在の位置
	readPosition int  // 次の位置
	ch           byte // 現在の文字
}

// NewLexer は新しいLexerを生成します。
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}

	// 1文字読み込む
	l.readChar()

	return l
}

// readChar は入力から次の文字を読み込みます。
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		// 終端に達したら0を返す
		l.ch = 0
	} else {
		// 1文字読み込む
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// NextToken は次のトークンを返します。
//
//nolint:cyclop
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '(':
		tok = newToken(TOKEN_OPEN_PAREN, l.ch)
	case ')':
		tok = newToken(TOKEN_CLOSE_PAREN, l.ch)
	case ',':
		tok = newToken(TOKEN_COMMA, l.ch)
	case ';':
		tok = newToken(TOKEN_SEMICOLON, l.ch)
	case '=':
		tok = newToken(TOKEN_EQUAL, l.ch)
	case '>':
		tok = newToken(TOKEN_GREATER, l.ch)
	case '<':
		tok = newToken(TOKEN_LESS, l.ch)
	case 0:
		tok.Literal = Literal{}
		tok.Type = TOKEN_EOF
	default:
		if isLiteral(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal.Str)
			return tok
		}
		tok = newToken(TOKEN_ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

func newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: Literal{Str: string(ch)}}
}

func (l *Lexer) readIdentifier() Literal {
	position := l.position
	for isLiteral(l.ch) {
		l.readChar()
	}
	str := l.input[position:l.position]

	return Literal{Str: str}
}

func isLiteral(ch byte) bool {
	return 'A' <= ch && ch <= 'Z' ||
		'a' <= ch && ch <= 'z' ||
		'0' <= ch && ch <= '9' ||
		ch == '_' ||
		ch == QuotationChar // TODO: 一考の余地
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
