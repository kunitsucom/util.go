package postgres

import (
	"strings"
)

// MEMO: https://www.postgresql.jp/docs/11/datatype.html

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

func (l *Literal) String() string {
	return l.Str
}

func (l *Literal) PlainString() string {
	return l.Str
}

type TokenType string

//nolint:revive,stylecheck
const (
	// SPECIAL TOKENS.
	TOKEN_ILLEGAL TokenType = "ILLEGAL"
	TOKEN_EOF     TokenType = "EOF"

	// SPECIAL CHARACTERS.
	TOKEN_OPEN_PAREN    TokenType = "OPEN_PAREN"    // (
	TOKEN_CLOSE_PAREN   TokenType = "CLOSE_PAREN"   // )
	TOKEN_COMMA         TokenType = "COMMA"         // ,
	TOKEN_SEMICOLON     TokenType = "SEMICOLON"     // ;
	TOKEN_EQUAL         TokenType = "EQUAL"         // =
	TOKEN_GREATER       TokenType = "GREATER"       // >
	TOKEN_LESS          TokenType = "LESS"          // <
	TOKEN_PLUS          TokenType = "PLUS"          // +
	TOKEN_MINUS         TokenType = "MINUS"         // -
	TOKEN_ASTERISK      TokenType = "ASTERISK"      // *
	TOKEN_SLASH         TokenType = "SLASH"         // /
	TOKEN_STRING_CONCAT TokenType = "STRING_CONCAT" //nolint:gosec // ||
	TOKEN_TYPECAST      TokenType = "TYPECAST"      // ::

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
	TOKEN_ON     TokenType = "ON"
	TOKEN_TO     TokenType = "TO"

	// DATA TYPE.
	TOKEN_BOOLEAN     TokenType = "BOOLEAN"
	TOKEN_SMALLINT    TokenType = "SMALLINT"
	TOKEN_INTEGER     TokenType = "INTEGER"
	TOKEN_BIGINT      TokenType = "BIGINT"
	TOKEN_DECIMAL     TokenType = "DECIMAL"
	TOKEN_NUMERIC     TokenType = "NUMERIC"
	TOKEN_REAL        TokenType = "REAL"
	TOKEN_DOUBLE      TokenType = "DOUBLE"
	TOKEN_PRECISION   TokenType = "PRECISION"
	TOKEN_SMALLSERIAL TokenType = "SMALLSERIAL"
	TOKEN_SERIAL      TokenType = "SERIAL"
	TOKEN_BIGSERIAL   TokenType = "BIGSERIAL"
	TOKEN_UUID        TokenType = "UUID"
	TOKEN_JSONB       TokenType = "JSONB"
	TOKEN_CHARACTER   TokenType = "CHARACTER"
	TOKEN_VARYING     TokenType = "VARYING"
	TOKEN_TEXT        TokenType = "TEXT"
	TOKEN_TIMESTAMP   TokenType = "TIMESTAMP"
	TOKEN_TIMESTAMPTZ TokenType = "TIMESTAMPTZ"
	TOKEN_WITH        TokenType = "WITH"
	TOKEN_TIME        TokenType = "TIME"
	TOKEN_ZONE        TokenType = "ZONE"

	// COLUMN.
	TOKEN_DEFAULT TokenType = "DEFAULT"
	TOKEN_NOT     TokenType = "NOT"
	TOKEN_NULL    TokenType = "NULL"
	TOKEN_ASC     TokenType = "ASC"
	TOKEN_DESC    TokenType = "DESC"

	// CONSTRAINT.
	TOKEN_CONSTRAINT TokenType = "CONSTRAINT"
	TOKEN_PRIMARY    TokenType = "PRIMARY"
	TOKEN_KEY        TokenType = "KEY"
	TOKEN_FOREIGN    TokenType = "FOREIGN"
	TOKEN_REFERENCES TokenType = "REFERENCES"
	TOKEN_UNIQUE     TokenType = "UNIQUE"
	TOKEN_CHECK      TokenType = "CHECK"

	// FUNCTION.
	TOKEN_NULLIF TokenType = "NULLIF"

	// LITERAL.
	TOKEN_LITERAL TokenType = "LITERAL"

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
	case "ON":
		return TOKEN_ON
	case "TO":
		return TOKEN_TO
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
	case "SMALLSERIAL":
		return TOKEN_SMALLSERIAL
	case "SERIAL":
		return TOKEN_SERIAL
	case "BIGSERIAL":
		return TOKEN_BIGSERIAL
	case "UUID":
		return TOKEN_UUID
	case "JSONB":
		return TOKEN_JSONB
	case "CHARACTER":
		return TOKEN_CHARACTER
	case "VARYING", "VARCHAR":
		return TOKEN_VARYING
	case "TEXT":
		return TOKEN_TEXT
	case "TIMESTAMP":
		return TOKEN_TIMESTAMP
	case "TIMESTAMPTZ":
		return TOKEN_TIMESTAMPTZ
	case "WITH":
		return TOKEN_WITH
	case "TIME":
		return TOKEN_TIME
	case "ZONE":
		return TOKEN_ZONE
	case "DEFAULT":
		return TOKEN_DEFAULT
	case "NOT":
		return TOKEN_NOT
	case "NULL":
		return TOKEN_NULL
	case "ASC":
		return TOKEN_ASC
	case "DESC":
		return TOKEN_DESC
	case "CONSTRAINT":
		return TOKEN_CONSTRAINT
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
	case "CHECK":
		return TOKEN_CHECK
	case "NULLIF":
		return TOKEN_NULLIF
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
//nolint:funlen,cyclop
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	if l.ch == '-' && l.peekChar() == '-' {
		l.skipComment()
		return l.NextToken()
	}

	switch l.ch {
	case '"', '\'':
		tok.Type = TOKEN_IDENT
		tok.Literal = Literal{Str: l.readQuotedLiteral(l.ch)}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: TOKEN_STRING_CONCAT, Literal: Literal{Str: literal}}
		} else {
			tok = newToken(TOKEN_ILLEGAL, l.ch)
		}
	case ':':
		if l.peekChar() == ':' {
			l.readChar()
			tok = Token{Type: TOKEN_TYPECAST, Literal: Literal{Str: "::"}}
		} else {
			tok = newToken(TOKEN_ILLEGAL, l.ch)
		}
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
	case '+':
		tok = newToken(TOKEN_PLUS, l.ch)
	case '-':
		tok = newToken(TOKEN_MINUS, l.ch)
	case '*':
		tok = newToken(TOKEN_ASTERISK, l.ch)
	case '/':
		tok = newToken(TOKEN_SLASH, l.ch)
	case 0:
		tok.Literal = Literal{}
		tok.Type = TOKEN_EOF
	default:
		if isLiteral(l.ch) {
			lit := l.readIdentifier()
			tok.Type = lookupIdent(lit)
			tok.Literal = Literal{Str: lit}
			return tok
		}
		tok = newToken(TOKEN_ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

// readQuotedLiteral はクォーテーションで囲まれた文字列を読み込みます。
func (l *Lexer) readQuotedLiteral(quote byte) string {
	// position := l.position + 1 // クォーテーションの次の文字から開始
	position := l.position // クォーテーションの文字から開始
	for {
		l.readChar()
		if l.ch == quote || l.ch == 0 {
			break
		}
	}
	return l.input[position : l.position+1]
}

// peekChar は次の文字を覗き見ますが、現在の位置は進めません。
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: Literal{Str: string(ch)}}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLiteral(l.ch) {
		l.readChar()
	}
	str := l.input[position:l.position]

	return str
}

func isLiteral(ch byte) bool {
	return 'A' <= ch && ch <= 'Z' ||
		'a' <= ch && ch <= 'z' ||
		'0' <= ch && ch <= '9' ||
		ch == '_' ||
		ch == '.'
}

func (l *Lexer) skipWhitespace() (skipped bool) {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		skipped = true || skipped
		l.readChar()
	}
	return skipped
}

func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}
