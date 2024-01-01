package postgres

import (
	"testing"

	"github.com/kunitsucom/util.go/testing/require"
)

func Test_lookuplookupIdent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  TokenType
	}{
		{name: "success,CREATE", input: "CREATE", want: TOKEN_CREATE},
		{name: "success,ALTER", input: "ALTER", want: TOKEN_ALTER},
		{name: "success,DROP", input: "DROP", want: TOKEN_DROP},
		{name: "success,RENAME", input: "RENAME", want: TOKEN_RENAME},
		{name: "success,CREATE", input: "CREATE", want: TOKEN_CREATE},
		{name: "success,ALTER", input: "ALTER", want: TOKEN_ALTER},
		{name: "success,DROP", input: "DROP", want: TOKEN_DROP},
		{name: "success,RENAME", input: "RENAME", want: TOKEN_RENAME},
		{name: "success,TRUNCATE", input: "TRUNCATE", want: TOKEN_TRUNCATE},
		{name: "success,TABLE", input: "TABLE", want: TOKEN_TABLE},
		{name: "success,INDEX", input: "INDEX", want: TOKEN_INDEX},
		{name: "success,VIEW", input: "VIEW", want: TOKEN_VIEW},
		{name: "success,IF", input: "IF", want: TOKEN_IF},
		{name: "success,EXISTS", input: "EXISTS", want: TOKEN_EXISTS},
		{name: "success,INT", input: "INT", want: TOKEN_INTEGER},
		{name: "success,INTEGER", input: "INTEGER", want: TOKEN_INTEGER},
		{name: "success,UUID", input: "UUID", want: TOKEN_UUID},
		{name: "success,VARCHAR", input: "VARCHAR", want: TOKEN_VARYING},
		{name: "success,TEXT", input: "TEXT", want: TOKEN_TEXT},
		{name: "success,TIMESTAMP", input: "TIMESTAMP", want: TOKEN_TIMESTAMP},
		{name: "success,TIMESTAMPZ", input: "TIMESTAMPZ", want: TOKEN_TIMESTAMPZ},
		{name: "success,NOT", input: "NOT", want: TOKEN_NOT},
		{name: "success,NULL", input: "NULL", want: TOKEN_NULL},
		{name: "success,PRIMARY", input: "PRIMARY", want: TOKEN_PRIMARY},
		{name: "success,KEY", input: "KEY", want: TOKEN_KEY},
		{name: "success,UNIQUE", input: "UNIQUE", want: TOKEN_UNIQUE},
		{name: "success,IDENT", input: "users", want: TOKEN_IDENT},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := lookupIdent(tt.input)

			if !require.Equal(t, tt.want, got) {
				t.FailNow()
			}
		})
	}
}

func TestLex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  []Token
	}{
		{
			name: "success,CREATE_TABLE",
			input: `CREATE TABLE IF NOT EXISTS "users" (
    "user_id"    UUID         NOT NULL,
    "name"       VARCHAR(255) NOT NULL,
    "email"      VARCHAR(255) NOT NULL,
    "password"   VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMPZ   NOT NULL,
    "updated_at" TIMESTAMPZ   NOT NULL,
    PRIMARY KEY ("user_id"),
    UNIQUE ("email")
);`,
			want: []Token{
				{Type: TOKEN_CREATE, Literal: Literal{Str: "CREATE"}},
				{Type: TOKEN_TABLE, Literal: Literal{Str: "TABLE"}},
				{Type: TOKEN_IF, Literal: Literal{Str: "IF"}},
				{Type: TOKEN_NOT, Literal: Literal{Str: "NOT"}},
				{Type: TOKEN_EXISTS, Literal: Literal{Str: "EXISTS"}},
				{Type: TOKEN_IDENT, Literal: Literal{Str: `"users"`}},
				{Type: TOKEN_OPEN_PAREN, Literal: Literal{Str: "("}},
				{Type: TOKEN_IDENT, Literal: Literal{Str: `"user_id"`}},
				{Type: TOKEN_UUID, Literal: Literal{Str: "UUID"}},
				{Type: TOKEN_NOT, Literal: Literal{Str: "NOT"}},
				{Type: TOKEN_NULL, Literal: Literal{Str: "NULL"}},
				{Type: TOKEN_COMMA, Literal: Literal{Str: ","}},
				{Type: TOKEN_IDENT, Literal: Literal{Str: `"name"`}},
				{Type: TOKEN_VARYING, Literal: Literal{Str: "VARCHAR"}},
				{Type: TOKEN_OPEN_PAREN, Literal: Literal{Str: "("}},
				{Type: TOKEN_IDENT, Literal: Literal{Str: "255"}},
				{Type: TOKEN_CLOSE_PAREN, Literal: Literal{Str: ")"}},
				{Type: TOKEN_NOT, Literal: Literal{Str: "NOT"}},
				{Type: TOKEN_NULL, Literal: Literal{Str: "NULL"}},
				{Type: TOKEN_COMMA, Literal: Literal{Str: ","}},
				{Type: TOKEN_IDENT, Literal: Literal{Str: `"email"`}},
				{Type: TOKEN_VARYING, Literal: Literal{Str: "VARCHAR"}},
				{Type: TOKEN_OPEN_PAREN, Literal: Literal{Str: "("}},
				{Type: TOKEN_IDENT, Literal: Literal{Str: "255"}},
				{Type: TOKEN_CLOSE_PAREN, Literal: Literal{Str: ")"}},
				{Type: TOKEN_NOT, Literal: Literal{Str: "NOT"}},
				{Type: TOKEN_NULL, Literal: Literal{Str: "NULL"}},
				{Type: TOKEN_COMMA, Literal: Literal{Str: ","}},
				{Type: TOKEN_IDENT, Literal: Literal{Str: `"password"`}},
				{Type: TOKEN_VARYING, Literal: Literal{Str: "VARCHAR"}},
				{Type: TOKEN_OPEN_PAREN, Literal: Literal{Str: "("}},
				{Type: TOKEN_IDENT, Literal: Literal{Str: "255"}},
				{Type: TOKEN_CLOSE_PAREN, Literal: Literal{Str: ")"}},
				{Type: TOKEN_NOT, Literal: Literal{Str: "NOT"}},
				{Type: TOKEN_NULL, Literal: Literal{Str: "NULL"}},
				{Type: TOKEN_COMMA, Literal: Literal{Str: ","}},
				{Type: TOKEN_IDENT, Literal: Literal{Str: `"created_at"`}},
				{Type: TOKEN_TIMESTAMPZ, Literal: Literal{Str: "TIMESTAMPZ"}},
				{Type: TOKEN_NOT, Literal: Literal{Str: "NOT"}},
				{Type: TOKEN_NULL, Literal: Literal{Str: "NULL"}},
				{Type: TOKEN_COMMA, Literal: Literal{Str: ","}},
				{Type: TOKEN_IDENT, Literal: Literal{Str: `"updated_at"`}},
				{Type: TOKEN_TIMESTAMPZ, Literal: Literal{Str: "TIMESTAMPZ"}},
				{Type: TOKEN_NOT, Literal: Literal{Str: "NOT"}},
				{Type: TOKEN_NULL, Literal: Literal{Str: "NULL"}},
				{Type: TOKEN_COMMA, Literal: Literal{Str: ","}},
				{Type: TOKEN_PRIMARY, Literal: Literal{Str: "PRIMARY"}},
				{Type: TOKEN_KEY, Literal: Literal{Str: "KEY"}},
				{Type: TOKEN_OPEN_PAREN, Literal: Literal{Str: "("}},
				{Type: TOKEN_IDENT, Literal: Literal{Str: `"user_id"`}},
				{Type: TOKEN_CLOSE_PAREN, Literal: Literal{Str: ")"}},
				{Type: TOKEN_COMMA, Literal: Literal{Str: ","}},
				{Type: TOKEN_UNIQUE, Literal: Literal{Str: "UNIQUE"}},
				{Type: TOKEN_OPEN_PAREN, Literal: Literal{Str: "("}},
				{Type: TOKEN_IDENT, Literal: Literal{Str: `"email"`}},
				{Type: TOKEN_CLOSE_PAREN, Literal: Literal{Str: ")"}},
				{Type: TOKEN_CLOSE_PAREN, Literal: Literal{Str: ")"}},
				{Type: TOKEN_SEMICOLON, Literal: Literal{Str: ";"}},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			l := NewLexer(tt.input)
			got := make([]Token, 0)
			for {
				tok := l.NextToken()
				if tok.Type == TOKEN_EOF {
					break
				}
				got = append(got, tok)
			}

			if !require.Equal(t, tt.want, got) {
				t.FailNow()
			}

			for i := range got {
				if !require.Equal(t, got[i].Type, tt.want[i].Type) {
					t.Fail()
				}

				if !require.Equal(t, got[i].Literal, tt.want[i].Literal) {
					t.Fail()
				}
			}
		})
	}
}

func TestLexer_NextToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  Token
	}{
		{
			name:  "failure,|",
			input: `|`,
			want: Token{
				Type:    TOKEN_ILLEGAL,
				Literal: Literal{Str: "|"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			l := NewLexer(tt.input)
			got := l.NextToken()

			if !require.Equal(t, tt.want, got) {
				t.FailNow()
			}
		})
	}
}

func TestLiteral_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		literal Literal
		want    string
	}{
		{
			name:    "success,CREATE",
			literal: Literal{Str: "CREATE"},
			want:    "CREATE",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.literal.String()

			if !require.Equal(t, tt.want, got) {
				t.FailNow()
			}
		})
	}
}
