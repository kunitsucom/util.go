package ddlz

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type DefaultLogger struct {
	enable bool
	l      *log.Logger
}

func (l *DefaultLogger) Printf(format string, v ...interface{}) {
	if !l.enable {
		return
	}
	l.l.Printf(format, v...)
}

var DebugLogger Logger = &DefaultLogger{enable: true, l: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lmicroseconds)}

// NOTE: MySQL :: MySQL 8.0 リファレンスマニュアル :: 13.1.20 CREATE TABLE ステートメント https://dev.mysql.com/doc/refman/8.0/ja/create-table.html

type VerbType string

const (
	StatementTypeUnknown VerbType = "UNKNOWN"
	StatementTypeCreate  VerbType = "CREATE"
	StatementTypeAlter   VerbType = "ALTER"
	StatementTypeDrop    VerbType = "DROP"
)

type StmtType string

const (
	StmtTypeUnknown StmtType = "UNKNOWN"
	StmtTypeTable   StmtType = "TABLE"
)

type ColumnOption string

type Column struct {
	Name          string
	Type          string
	AutoIncrement bool
	PrimaryKey    bool
	NotNull       bool
	Default       []string
	OnUpdate      []string
	Options       []ColumnOption
	Comments      []string
}

type TableOption string

type Stmt struct {
	VerbType VerbType
	StmtType StmtType
	Name     string
	Columns  []Column
	Options  []TableOption
	Comments []string
}

type DDL struct {
	Stmts    []Stmt
	Comments []string
}

type Token string

//nolint:gochecknoglobals
var (
	tokenTypes = map[Token]struct{}{
		Token("INTEGER"):   {},
		Token("INT"):       {},
		Token("TINYINT"):   {},
		Token("TEXT"):      {},
		Token("CHAR"):      {},
		Token("VARCHAR"):   {},
		Token("TIMESTAMP"): {},
	}
	parseTokenFuncs = []func(token Token) bool{
		func(token Token) bool {
			if _, ok := tokenTypes[token]; ok {
				return true
			}
			for k := range tokenTypes {
				if strings.HasPrefix(string(token), string(k)) {
					return true
				}
			}
			return false
		},
	}
)

const (
	TokenUnknown          Token = "UNKNOWN"
	TokenEOF              Token = "EOF"
	TokenWS               Token = " "
	TokenComma            Token = ","
	TokenSingleQuote      Token = "'"
	TokenDoubleQuote      Token = "\""
	TokenBackQuote        Token = "`"
	TokenParenthesisOpen  Token = "("
	TokenParenthesisClose Token = ")"
	TokenSemicolon        Token = ";"
	TokenComment          Token = "--"
	TokenIdent            Token = "IDENT" // Literals. fields, table_name.
	TokenCreate           Token = "CREATE"
	TokenAlter            Token = "ALTER"
	TokenDrop             Token = "DROP"
	TokenTable            Token = "TABLE"
	TokenIf               Token = "IF"
	TokenExists           Token = "EXISTS"
	TokenPrimary          Token = "PRIMARY"
	TokenKey              Token = "KEY"
	TokenNot              Token = "NOT"
	TokenNull             Token = "NULL"
	TokenDefault          Token = "DEFAULT"
	TokenOn               Token = "ON"
	TokenUpdate           Token = "UPDATE"
	TokenAutoIncrement    Token = "AUTOINCREMENT"
)

const (
	eof              = rune(0)
	comma            = ','
	singleQuote      = '\''
	doubleQuote      = '"'
	backQuote        = '`'
	parenthesisOpen  = '('
	parenthesisClose = ')'
	hyphen           = '-'
	semicolon        = ';'
)

func isWhitespace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func isNewLine(c rune) bool {
	return c == '\n' || c == '\r'
}

func isLetter(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func parseToken(s string) (Token, bool) {
	switch kw := Token(strings.ToUpper(s)); kw {
	case TokenCreate, TokenAlter, TokenDrop,
		TokenTable, TokenIf, TokenExists,
		TokenPrimary, TokenKey,
		TokenNot, TokenNull,
		TokenDefault, TokenAutoIncrement: // キーワード
		return kw, true
	case TokenUnknown, TokenEOF, TokenWS, TokenComma, TokenParenthesisOpen, TokenParenthesisClose, TokenComment, TokenSemicolon: // メタ文字
		return kw, false
	case TokenIdent: // 以前に isLetter を通っているのでここには来ないはず
		return kw, false
	default: // 型など
		for i := range parseTokenFuncs {
			if parseTokenFuncs[i](kw) {
				return kw, true
			}
		}
		return kw, false
	}
}

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) readRune() rune {
	c, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return c
}

func (s *Scanner) unreadRune() {
	_ = s.r.UnreadRune()
}

func (s *Scanner) Scan() (token Token, literal string) {
	c := s.readRune()

	if isWhitespace(c) {
		s.unreadRune()
		return s.scanWhitespace()
	} else if isLetter(c) || isDigit(c) { // 文字or数字なら IDENT 読み込み開始
		s.unreadRune()
		return s.scanIdent()
	}

	switch c {
	case eof:
		return TokenEOF, ""
	case singleQuote:
		return TokenSingleQuote, string(c)
	case doubleQuote:
		return TokenDoubleQuote, string(c)
	case backQuote:
		return TokenBackQuote, string(c)
	case comma:
		return TokenComma, string(c)
	case parenthesisOpen:
		return TokenParenthesisOpen, string(c)
	case parenthesisClose:
		return TokenParenthesisClose, string(c)
	case semicolon:
		return TokenSemicolon, string(c)
	case hyphen:
		c = s.readRune()
		if c == hyphen {
			return s.scanComment()
		}
		s.unreadRune()
	}

	return TokenUnknown, string(c)
}

func (s *Scanner) scanWhitespace() (token Token, literal string) {
	var buf bytes.Buffer

	for {
		if c := s.readRune(); c == eof {
			break
		} else if isWhitespace(c) {
			buf.WriteRune(c)
		} else {
			s.unreadRune()
			break
		}
	}

	return TokenWS, buf.String()
}

func (s *Scanner) scanIdent() (token Token, literal string) {
	var buf bytes.Buffer

	for {
		if c := s.readRune(); c == eof {
			break
		} else if isLetter(c) || isDigit(c) || c == parenthesisOpen || c == parenthesisClose || c == '_' {
			_, _ = buf.WriteRune(c)
		} else {
			s.unreadRune()
			break
		}
	}

	if kw, ok := parseToken(buf.String()); ok {
		return kw, buf.String()
	}

	return TokenIdent, buf.String()
}

func (s *Scanner) scanComment() (token Token, literal string) {
	var buf bytes.Buffer

	for {
		// --      this is comment
		//   ^^^^^^ <- skip whitespace
		if c := s.readRune(); !(isWhitespace(c) && !isNewLine(c)) {
			s.unreadRune()
			break
		}
	}

	for {
		if c := s.readRune(); c == eof {
			break
		} else if isNewLine(c) {
			break
		} else {
			_, _ = buf.WriteRune(c)
		}
	}

	return TokenComment, buf.String()
}

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		token   Token  // last read token
		literal string // last read literal
		n       int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func (p *Parser) readToken() (token Token, literal string) {
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.token, p.buf.literal
	}

	token, literal = p.s.Scan()

	p.buf.token, p.buf.literal = token, literal

	return
}

func (p *Parser) unreadToken() {
	DebugLogger.Printf("🪲: unreadToken")
	p.buf.n = 1
}

func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.readToken()
	if tok == TokenWS {
		tok, lit = p.readToken()
	}
	return
}

func (p *Parser) Parse() (*DDL, error) {
	var (
		ddl  = &DDL{}
		verb Token
	)
	for {
		token, literal := p.scanIgnoreWhitespace()
		DebugLogger.Printf("🪲: Parse: token: %q, literal: %q", token, literal)
		switch token {
		case TokenCreate, TokenAlter, TokenDrop:
			verb = token
		case TokenTable:
			switch verb {
			case TokenCreate:
				if err := p.parseCreateTable(ddl); err != nil {
					return nil, fmt.Errorf("p.parseCreateTable: %w", err)
				}
			case TokenAlter:
				// TODO: implement
			case TokenDrop:
				// TODO: implement
			}
		case TokenComment:
			ddl.Comments = append(ddl.Comments, literal)
		case TokenEOF:
			return ddl, nil
		}
	}
}

func (p *Parser) parseCreateTable(ddl *DDL) error {
	stmt := Stmt{
		VerbType: StatementTypeCreate,
		StmtType: StmtTypeTable,
	}
	for {
		token, literal := p.scanIgnoreWhitespace()
		DebugLogger.Printf("🪲: parseCreateTable: token: %q, literal: %q", token, literal)
		switch token {
		case TokenIdent:
			if stmt.Name == "" {
				stmt.Name = literal
			} else {
				stmt.Options = append(stmt.Options, TableOption(literal))
			}
		case TokenParenthesisOpen:
			if err := p.parseCreateTableColumns(&stmt); err != nil {
				return fmt.Errorf("p.parseCreateTableColumns: %w", err)
			}
		case TokenSemicolon:
			ddl.Stmts = append(ddl.Stmts, stmt)
			return nil
		case TokenComment:
			stmt.Comments = append(stmt.Comments, literal)
		case TokenEOF:
			p.unreadToken()
			return nil
		}
	}
}

func (p *Parser) parseCreateTableColumns(out *Stmt) error {
	column := Column{}
	for {
		token, literal := p.scanIgnoreWhitespace()
		DebugLogger.Printf("🪲: parseCreateTableColumns: token: %q, literal: %q", token, literal)
		switch token {
		case TokenIdent, TokenSingleQuote, TokenDoubleQuote, TokenBackQuote:
			p.unreadToken()
			if err := p.parseCreateTableColumn(&column); err != nil {
				return fmt.Errorf("p.parseCreateTableColumn: %w", err)
			}
		case TokenExists:
			// TODO: implement
		case TokenComma:
			out.Columns = append(out.Columns, column)
			column = Column{}
		case TokenComment:
			column.Comments = append(column.Comments, literal)
		default:
			for i := range parseTokenFuncs {
				if parseTokenFuncs[i](token) {
					column.Type = literal
					continue
				}
			}
			return fmt.Errorf("unexpected token: %q", token)
		case TokenParenthesisClose, TokenSemicolon, TokenEOF:
			out.Columns = append(out.Columns, column)
			p.unreadToken()
			return nil
		}
	}
}

func (p *Parser) parseCreateTableColumn(out *Column) error {
	var (
		not         bool
		on          bool
		onUpdate    bool
		defaultVal  bool
		doubleQuote bool
	)
Parser:
	for {
		token, literal := p.scanIgnoreWhitespace()
		DebugLogger.Printf("🪲: parseCreateTableColumn: token: %q, literal: %q", token, literal)
		switch token {
		case TokenIdent:
			switch {
			case out.Name == "":
				out.Name = literal
			case defaultVal:
				out.Default = append(out.Default, literal)
			case onUpdate:
				out.OnUpdate = append(out.OnUpdate, literal)
			default:
				out.Options = append(out.Options, ColumnOption(literal))
			}
		case TokenPrimary, TokenKey:
			out.PrimaryKey = true
		case TokenNot:
			not = true
		case TokenNull:
			if not {
				out.NotNull = not
				not = false
			}
		case TokenAutoIncrement:
			out.AutoIncrement = true
		case TokenDefault:
			defaultVal = true
		case TokenOn:
			on = true
		case TokenUpdate:
			if on {
				onUpdate = true
				on = false
			}
		case TokenParenthesisOpen:
			// TODO: implement
		case TokenParenthesisClose:
			// TODO: implement
			p.unreadToken()
			return nil
		case TokenDoubleQuote:
			doubleQuote = !doubleQuote
			// TODO: implement
		case TokenComma:
			p.unreadToken()
			return nil
		case TokenComment:
			out.Comments = append(out.Comments, literal)
		case TokenSemicolon, TokenEOF:
			p.unreadToken()
			return fmt.Errorf("unexpected token: %q", token)
		default:
			for i := range parseTokenFuncs {
				if parseTokenFuncs[i](token) {
					out.Type = literal
					continue Parser
				}
			}
			return fmt.Errorf("unexpected token: %q", token)
		}
	}
}
