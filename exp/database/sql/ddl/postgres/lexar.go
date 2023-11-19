package postgres

import (
	"strings"
)

// Token はSQL文のトークンを表す型です。
type Token struct {
	Type    TokenType
	Literal Literal
}

type Literal struct {
	Str           string
	Quoted        bool
	QuotationMark string
}

func (l Literal) String() string {
	if l.Quoted {
		return l.QuotationMark + l.Str + l.QuotationMark
	}

	return l.Str
}

type TokenType string

//nolint:revive,stylecheck
const (
	TOKEN_ILLEGAL     TokenType = "ILLEGAL"
	TOKEN_EOF         TokenType = "EOF"
	TOKEN_CREATE      TokenType = "CREATE"
	TOKEN_IF          TokenType = "IF"
	TOKEN_NOT         TokenType = "NOT"
	TOKEN_EXISTS      TokenType = "EXISTS"
	TOKEN_TABLE       TokenType = "TABLE"
	TOKEN_INT         TokenType = "INT"
	TOKEN_VARCHAR     TokenType = "VARCHAR"
	TOKEN_TEXT        TokenType = "TEXT"
	TOKEN_UUID        TokenType = "UUID"
	TOKEN_NULL        TokenType = "NULL"
	TOKEN_OPEN_PAREN  TokenType = "OPEN_PAREN"
	TOKEN_CLOSE_PAREN TokenType = "CLOSE_PAREN"
	TOKEN_COMMA       TokenType = "COMMA"
	TOKEN_PRIMARY     TokenType = "PRIMARY"
	TOKEN_KEY         TokenType = "KEY"
	TOKEN_UNIQUE      TokenType = "UNIQUE"
	TOKEN_IDENT       TokenType = "IDENT" // 識別子
)

func lookupIdent(ident string) TokenType {
	switch strings.ToUpper(ident) {
	case "CREATE":
		return TOKEN_CREATE
	case "IF":
		return TOKEN_IF
	case "NOT":
		return TOKEN_NOT
	case "EXISTS":
		return TOKEN_EXISTS
	case "TABLE":
		return TOKEN_TABLE
	case "INT", "INTEGER":
		return TOKEN_INT
	case "VARCHAR":
		return TOKEN_VARCHAR
	case "TEXT":
		return TOKEN_TEXT
	case "UUID":
		return TOKEN_UUID
	case "NULL":
		return TOKEN_NULL
	case "PRIMARY":
		return TOKEN_PRIMARY
	case "KEY":
		return TOKEN_KEY
	default:
		return TOKEN_IDENT
	}
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
	case 0:
		tok.Literal = Literal{}
		tok.Type = TOKEN_EOF
	default:
		if isLetter(l.ch) {
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
	for isLetter(l.ch) {
		l.readChar()
	}
	str := l.input[position:l.position]

	if q := `"`; strings.HasPrefix(str, q) && strings.HasSuffix(str, q) {
		// TODO: エスケープ文字の処理が必要？
		return Literal{Str: strings.Trim(str, q), Quoted: true, QuotationMark: q}
	}

	return Literal{Str: str}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' ||
		'A' <= ch && ch <= 'Z' ||
		ch == '_' ||
		ch == '"' // TODO: 一考の余地
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
