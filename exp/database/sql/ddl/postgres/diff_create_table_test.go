package postgres

import (
	"testing"

	"github.com/kunitsucom/util.go/exp/database/sql/ddl"
	"github.com/kunitsucom/util.go/testing/assert"
	"github.com/kunitsucom/util.go/testing/require"
)

//nolint:paralleltest,tparallel
func TestDiffCreateTable(t *testing.T) {
	t.Run("failure,ddl.ErrNoDifference", func(t *testing.T) {
		t.Parallel()

		before := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, description TEXT, PRIMARY KEY ("id"));`
		beforeDDL, err := NewParser(NewLexer(before)).Parse()
		require.NoError(t, err)

		after := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, description TEXT, PRIMARY KEY ("id"));`

		afterDDL, err := NewParser(NewLexer(after)).Parse()
		require.NoError(t, err)

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)

		assert.ErrorIs(t, err, ddl.ErrNoDifference)
		assert.Nil(t, actual)

		t.Logf("✅: %s:\n%s", t.Name(), actual)
	})

	t.Run("success,ADD_COLUMN", func(t *testing.T) {
		t.Parallel()

		before := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, description TEXT, PRIMARY KEY ("id"));`
		beforeDDL, err := NewParser(NewLexer(before)).Parse()
		require.NoError(t, err)

		after := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INTEGER DEFAULT 0 NOT NULL CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`

		afterDDL, err := NewParser(NewLexer(after)).Parse()
		require.NoError(t, err)

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)

		expectedStr := `ALTER TABLE "users" ADD COLUMN "age" INTEGER DEFAULT 0 NOT NULL;
ALTER TABLE "users" ADD CONSTRAINT users_age_check CHECK ("age" >= 0);
`

		assert.NoError(t, err)
		assert.Equal(t, expectedStr, actual.String())

		t.Logf("✅: %s:\n%s", t.Name(), actual)
	})

	t.Run("success,DROP_COLUMN", func(t *testing.T) {
		t.Parallel()

		before := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INTEGER DEFAULT 0 NOT NULL CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`
		beforeDDL, err := NewParser(NewLexer(before)).Parse()
		require.NoError(t, err)

		after := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL, description TEXT, PRIMARY KEY ("id"));`

		afterDDL, err := NewParser(NewLexer(after)).Parse()
		require.NoError(t, err)

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)

		expectedStr := `ALTER TABLE "users" DROP CONSTRAINT users_unique_name;
ALTER TABLE "users" DROP CONSTRAINT users_age_check;
ALTER TABLE "users" DROP COLUMN "age";
`

		assert.NoError(t, err)
		assert.Equal(t, expectedStr, actual.String())

		t.Logf("✅: %s:\n%s", t.Name(), actual)
	})

	t.Run("success,ALTER_COLUMN_SET_DATA_TYPE", func(t *testing.T) {
		t.Parallel()

		before := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`
		beforeDDL, err := NewParser(NewLexer(before)).Parse()
		require.NoError(t, err)

		after := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" TEXT NOT NULL UNIQUE, "age" BIGINT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`

		afterDDL, err := NewParser(NewLexer(after)).Parse()
		require.NoError(t, err)

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)

		expectedStr := `ALTER TABLE "users" ALTER COLUMN "name" SET DATA TYPE TEXT;
ALTER TABLE "users" ALTER COLUMN "age" SET DATA TYPE BIGINT;
ALTER TABLE "users" ADD CONSTRAINT users_unique_name UNIQUE ("name");
`

		assert.NoError(t, err)
		assert.Equal(t, expectedStr, actual.String())

		t.Logf("✅: %s:\n%s", t.Name(), actual)
	})

	t.Run("success,ALTER_COLUMN_DROP_DEFAULT", func(t *testing.T) {
		before := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`
		beforeDDL, err := NewParser(NewLexer(before)).Parse()
		require.NoError(t, err)

		after := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`
		afterDDL, err := NewParser(NewLexer(after)).Parse()
		require.NoError(t, err)

		expectedStr := `ALTER TABLE "users" ALTER COLUMN "age" DROP DEFAULT;
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)
		assert.NoError(t, err)
		assert.Equal(t, expectedStr, actual.String())

		t.Logf("✅: %s:\n%s", t.Name(), actual)
	})

	t.Run("success,ALTER_COLUMN_SET_DEFAULT", func(t *testing.T) {
		before := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`
		beforeDDL, err := NewParser(NewLexer(before)).Parse()
		require.NoError(t, err)

		after := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" <> 0), description TEXT, PRIMARY KEY (id));`
		afterDDL, err := NewParser(NewLexer(after)).Parse()
		require.NoError(t, err)

		expectedStr := `ALTER TABLE "users" ALTER COLUMN "age" SET DEFAULT 0;
ALTER TABLE "users" DROP CONSTRAINT users_age_check;
ALTER TABLE "users" ADD CONSTRAINT users_age_check CHECK ("age" <> 0);
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)
		assert.NoError(t, err)
		assert.Equal(t, expectedStr, actual.String())

		t.Logf("✅: %s:\n%s", t.Name(), actual)
	})

	t.Run("success,ALTER_TABLE_RENAME_TO", func(t *testing.T) {
		t.Parallel()

		before := `CREATE TABLE "public.users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`
		beforeDDL, err := NewParser(NewLexer(before)).Parse()
		require.NoError(t, err)

		after := `CREATE TABLE "app_users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`
		afterDDL, err := NewParser(NewLexer(after)).Parse()
		require.NoError(t, err)

		expectedStr := `ALTER TABLE "public.users" RENAME TO "public.app_users";
ALTER TABLE "public.app_users" DROP CONSTRAINT users_group_id_fkey;
ALTER TABLE "public.app_users" DROP CONSTRAINT users_unique_name;
ALTER TABLE "public.app_users" DROP CONSTRAINT users_age_check;
ALTER TABLE "public.app_users" DROP CONSTRAINT users_pkey;
ALTER TABLE "public.app_users" ADD CONSTRAINT app_users_group_id_fkey FOREIGN KEY (group_id) REFERENCES "groups" ("id");
ALTER TABLE "public.app_users" ADD CONSTRAINT app_users_unique_name UNIQUE ("name");
ALTER TABLE "public.app_users" ADD CONSTRAINT app_users_age_check CHECK ("age" >= 0);
ALTER TABLE "public.app_users" ADD CONSTRAINT app_users_pkey PRIMARY KEY ("id");
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)
		assert.NoError(t, err)
		assert.Equal(t, expectedStr, actual.String())

		t.Logf("✅: %s:\n%s", t.Name(), actual)
	})

	t.Run("success,SET_NOT_NULL", func(t *testing.T) {
		t.Parallel()

		before := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`
		beforeDDL, err := NewParser(NewLexer(before)).Parse()
		require.NoError(t, err)

		after := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INTEGER DEFAULT 0 NOT NULL CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`
		afterDDL, err := NewParser(NewLexer(after)).Parse()
		require.NoError(t, err)

		expectedStr := `ALTER TABLE "users" ALTER COLUMN "age" SET NOT NULL;
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)
		assert.NoError(t, err)
		assert.Equal(t, expectedStr, actual.String())

		t.Logf("✅: %s:\n%s", t.Name(), actual)
	})

	t.Run("success,DROP_NOT_NULL", func(t *testing.T) {
		t.Parallel()

		before := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 NOT NULL CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`
		beforeDDL, err := NewParser(NewLexer(before)).Parse()
		require.NoError(t, err)

		after := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`
		afterDDL, err := NewParser(NewLexer(after)).Parse()
		require.NoError(t, err)

		expectedStr := `ALTER TABLE "users" ALTER COLUMN "age" DROP NOT NULL;
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)
		assert.NoError(t, err)
		assert.Equal(t, expectedStr, actual.String())

		t.Logf("✅: %s:\n%s", t.Name(), actual)
	})

	t.Run("success,DROP_ADD_PRIMARY_KEY", func(t *testing.T) {
		t.Parallel()

		before := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 NOT NULL CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`
		beforeDDL, err := NewParser(NewLexer(before)).Parse()
		require.NoError(t, err)

		after := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 NOT NULL CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id", name));`
		afterDDL, err := NewParser(NewLexer(after)).Parse()
		require.NoError(t, err)

		expectedStr := `ALTER TABLE "users" DROP CONSTRAINT users_pkey;
ALTER TABLE "users" ADD CONSTRAINT users_pkey PRIMARY KEY ("id", name);
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)
		assert.NoError(t, err)
		assert.Equal(t, expectedStr, actual.String())

		t.Logf("✅: %s:\n%s", t.Name(), actual)
	})

	t.Run("success,DROP_ADD_FOREIGN_KEY", func(t *testing.T) {
		t.Parallel()

		before := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL, "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 NOT NULL CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"), CONSTRAINT users_group_id_fkey FOREIGN KEY (group_id) REFERENCES "groups" ("id"));`
		beforeDDL, err := NewParser(NewLexer(before)).Parse()
		require.NoError(t, err)

		after := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL, "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 NOT NULL CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"), CONSTRAINT users_group_id_fkey FOREIGN KEY (group_id, name) REFERENCES "groups" ("id", name));`
		afterDDL, err := NewParser(NewLexer(after)).Parse()
		require.NoError(t, err)

		expectedStr := `ALTER TABLE "users" DROP CONSTRAINT users_group_id_fkey;
ALTER TABLE "users" ADD CONSTRAINT users_group_id_fkey FOREIGN KEY (group_id, name) REFERENCES "groups" ("id", name);
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)
		assert.NoError(t, err)
		assert.Equal(t, expectedStr, actual.String())

		t.Logf("✅: %s:\n%s", t.Name(), actual)
	})

	t.Run("success,DROP_ADD_UNIQUE", func(t *testing.T) {
		t.Parallel()

		before := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL, "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 NOT NULL CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`
		beforeDDL, err := NewParser(NewLexer(before)).Parse()
		require.NoError(t, err)

		after := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL, "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 NOT NULL CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"), CONSTRAINT users_unique_name UNIQUE ("id", name));`
		afterDDL, err := NewParser(NewLexer(after)).Parse()
		require.NoError(t, err)

		expectedStr := `ALTER TABLE "users" DROP CONSTRAINT users_unique_name;
ALTER TABLE "users" ADD CONSTRAINT users_unique_name UNIQUE ("id", name);
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)
		assert.NoError(t, err)
		assert.Equal(t, expectedStr, actual.String())

		t.Logf("✅: %s:\n%s", t.Name(), actual)
	})

	t.Run("success,ALTER_COLUMN_SET_DEFAULT_OVERWRITE", func(t *testing.T) {
		t.Parallel()

		before := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL, "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 NOT NULL CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`
		beforeDDL, err := NewParser(NewLexer(before)).Parse()
		require.NoError(t, err)

		after := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL, "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT ( (0 + 3) - 1 * 4 / 2 ) NOT NULL CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`
		afterDDL, err := NewParser(NewLexer(after)).Parse()
		require.NoError(t, err)

		expectedStr := `ALTER TABLE "users" ALTER COLUMN "age" SET DEFAULT ((0 + 3) - 1 * 4 / 2);
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)
		assert.NoError(t, err)
		assert.Equal(t, expectedStr, actual.String())

		t.Logf("✅: %s:\n%s", t.Name(), actual)
	})

	t.Run("success,ALTER_COLUMN_SET_DEFAULT_complex", func(t *testing.T) {
		t.Parallel()

		before := `CREATE TABLE complex_defaults (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    unique_code TEXT,
    status TEXT DEFAULT 'pending',
    random_number INTEGER DEFAULT FLOOR(RANDOM() * 100)::INTEGER,
    json_data JSONB DEFAULT '{}',
    calculated_value INTEGER DEFAULT (SELECT COUNT(*) FROM another_table)
);
`
		beforeDDL, err := NewParser(NewLexer(before)).Parse()
		require.NoError(t, err)

		after := `CREATE TABLE complex_defaults (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    unique_code TEXT DEFAULT 'CODE-' || TO_CHAR(NOW(), 'YYYYMMDDHH24MISS') || '-' || LPAD(TO_CHAR(NEXTVAL('seq_complex_default')), 5, '0'),
    status TEXT DEFAULT 'pending',
    random_number INTEGER DEFAULT FLOOR(RANDOM() * 100)::INTEGER,
    json_data JSONB DEFAULT '{}',
    calculated_value INTEGER DEFAULT (SELECT COUNT(*) FROM another_table)
);
`
		afterDDL, err := NewParser(NewLexer(after)).Parse()
		require.NoError(t, err)

		expectedStr := `ALTER TABLE complex_defaults ALTER COLUMN unique_code SET DEFAULT 'CODE-' || TO_CHAR(NOW(), 'YYYYMMDDHH24MISS') || '-' || LPAD(TO_CHAR(NEXTVAL('seq_complex_default')), 5, '0');
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)
		assert.NoError(t, err)
		assert.Equal(t, expectedStr, actual.String())

		t.Logf("✅: %s:\n%s", t.Name(), actual)
	})

	t.Run("success,DiffCreateTableUseAlterTableAddConstraintNotValid", func(t *testing.T) {
		t.Parallel()

		before := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0, description TEXT, PRIMARY KEY ("id"));`
		beforeDDL, err := NewParser(NewLexer(before)).Parse()
		require.NoError(t, err)

		after := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`
		afterDDL, err := NewParser(NewLexer(after)).Parse()
		require.NoError(t, err)

		expectedStr := `ALTER TABLE "users" ADD CONSTRAINT users_age_check CHECK ("age" >= 0) NOT VALID;
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(true),
		)

		assert.NoError(t, err)
		assert.Equal(t, expectedStr, actual.String())

		t.Logf("✅: %s:\n%s", t.Name(), actual)
	})

	t.Run("success,CREATE_TABLE", func(t *testing.T) {
		t.Parallel()

		after := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`

		afterDDL, err := NewParser(NewLexer(after)).Parse()
		require.NoError(t, err)

		expectedStr := `CREATE TABLE "users" (
    id UUID NOT NULL,
    group_id UUID NOT NULL,
    "name" VARYING(255) NOT NULL,
    "age" INTEGER DEFAULT 0,
    description TEXT,
    CONSTRAINT users_group_id_fkey FOREIGN KEY (group_id) REFERENCES "groups" ("id"),
    CONSTRAINT users_unique_name UNIQUE ("name"),
    CONSTRAINT users_age_check CHECK ("age" >= 0),
    CONSTRAINT users_pkey PRIMARY KEY ("id")
);
`

		actual, err := DiffCreateTable(
			nil,
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(true),
		)
		assert.NoError(t, err)
		assert.Equal(t, expectedStr, actual.String())

		t.Logf("✅: %s:\n%s", t.Name(), actual)
	})

	t.Run("success,DROP_TABLE", func(t *testing.T) {
		t.Parallel()

		before := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`

		beforeDDL, err := NewParser(NewLexer(before)).Parse()
		require.NoError(t, err)

		expectedStr := `DROP TABLE "users";
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			nil,
			DiffCreateTableUseAlterTableAddConstraintNotValid(true),
		)

		assert.NoError(t, err)
		assert.Equal(t, expectedStr, actual.String())

		t.Logf("✅: %s:\n%s", t.Name(), actual)
	})
}