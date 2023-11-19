package postgres

import (
	"log"
	"testing"

	"github.com/kunitsucom/util.go/testing/require"
)

func TestAll(t *testing.T) {
	log.SetFlags(log.Flags() | log.Lshortfile)

	input := `
CREATE TABLE IF NOT EXISTS "users" (
    "user_id"  INT NOT NULL,
    "username" TEXT,
	PRIMARY KEY ("user_id")
)`
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	stmt, err := parser.ParseStatement()
	if err != nil {
		require.NoError(t, err)
	}

	t.Logf("✅: stmt:\n%+v\n", stmt)
}
