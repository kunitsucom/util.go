package stringz_test

import (
	"testing"

	stringz "github.com/kunitsucom/util.go/strings"
	"github.com/kunitsucom/util.go/testing/require"
)

func TestReadLine(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		actual := stringz.ReadLine(`-- CREATE TABLE
CREATE TABLE IF NOT EXISTS "users" (
    -- id column
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "name" TEXT NOT NULL,
    "created_at" TEXT NOT NULL,
    "updated_at" TEXT NOT NULL
);

`, "\n", stringz.ReadLineFuncRemoveCommentLine("--"))

		expected := `CREATE TABLE IF NOT EXISTS "users" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "name" TEXT NOT NULL,
    "created_at" TEXT NOT NULL,
    "updated_at" TEXT NOT NULL
);

`

		require.Equal(t, expected, actual)
	})
}
