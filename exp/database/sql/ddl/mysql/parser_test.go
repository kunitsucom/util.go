//nolint:testpackage
package mysql

import (
	"log"
	"os"
	"testing"

	"github.com/kunitsucom/util.go/exp/database/sql/ddl"
	"github.com/kunitsucom/util.go/exp/database/sql/ddl/internal"
	"github.com/kunitsucom/util.go/testing/assert"
	"github.com/kunitsucom/util.go/testing/require"
)

//nolint:paralleltest,tparallel
func TestParser_Parse(t *testing.T) {
	backup := internal.TraceLog
	t.Cleanup(func() {
		internal.TraceLog = backup
	})
	internal.TraceLog = log.New(os.Stderr, "TRACE: ", log.LstdFlags|log.Lshortfile)

	successTests := []struct {
		name    string
		input   string
		wantErr error
		wantStr string
	}{
		{
			name:    "success,CREATE_TABLE",
			input:   "CREATE TABLE `groups` (`group_id` VARCHAR(36) NOT NULL PRIMARY KEY, description TEXT) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci; CREATE TABLE `users` (user_id VARCHAR(36) NOT NULL, group_id VARCHAR(36) NOT NULL REFERENCES `groups` (`group_id`), `name` VARCHAR(255) NOT NULL UNIQUE, `age` INT DEFAULT 0 CHECK (`age` >= 0), description TEXT, PRIMARY KEY (`user_id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;",
			wantErr: nil,
			wantStr: `CREATE TABLE ` + "`" + `groups` + "`" + ` (
    ` + "`" + `group_id` + "`" + ` VARCHAR(36) NOT NULL,
    description TEXT,
    CONSTRAINT groups_pkey PRIMARY KEY (` + "`" + `group_id` + "`" + `)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
CREATE TABLE ` + "`" + `users` + "`" + ` (
    user_id VARCHAR(36) NOT NULL,
    group_id VARCHAR(36) NOT NULL,
    ` + "`" + `name` + "`" + ` VARCHAR(255) NOT NULL,
    ` + "`" + `age` + "`" + ` INT DEFAULT 0,
    description TEXT,
    CONSTRAINT users_pkey PRIMARY KEY (` + "`" + `user_id` + "`" + `),
    CONSTRAINT users_group_id_fkey FOREIGN KEY (group_id) REFERENCES ` + "`" + `groups` + "`" + ` (` + "`" + `group_id` + "`" + `),
    UNIQUE INDEX users_unique_name (` + "`" + `name` + "`" + `),
    CONSTRAINT users_age_check CHECK (` + "`" + `age` + "`" + ` >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
`,
		},
		// 		{
		// 			name: "success,complex_defaults",
		// 			input: `-- table: complex_defaults
		// CREATE TABLE IF NOT EXISTS complex_defaults (
		//     -- id is the primary key.
		//     id SERIAL PRIMARY KEY,
		//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		//     updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		//     unique_code TEXT DEFAULT 'CODE-' || TO_CHAR(NOW(), 'YYYYMMDDHH24MISS') || '-' || LPAD(TO_CHAR(NEXTVAL('seq_complex_default')), 5, '0'),
		//     status CHARACTER VARYING DEFAULT 'pending',
		//     random_number INTEGER DEFAULT FLOOR(RANDOM() * 100::INTEGER)::INTEGER,
		//     json_data JSONB DEFAULT '{}',
		//     calculated_value INTEGER DEFAULT (SELECT COUNT(*) FROM another_table)
		// );
		// `,
		// 			wantErr: nil,
		// 			wantStr: `CREATE TABLE IF NOT EXISTS complex_defaults (
		//     id SERIAL,
		//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		//     updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		//     unique_code TEXT DEFAULT 'CODE-' || TO_CHAR(NOW(), 'YYYYMMDDHH24MISS') || '-' || LPAD(TO_CHAR(NEXTVAL('seq_complex_default')), 5, '0'),
		//     status CHARACTER VARYING DEFAULT 'pending',
		//     random_number INTEGER DEFAULT FLOOR(RANDOM() * 100::INTEGER)::INTEGER,
		//     json_data JSONB DEFAULT '{}',
		//     calculated_value INTEGER DEFAULT (SELECT COUNT(*) FROM another_table),
		//     CONSTRAINT complex_defaults_pkey PRIMARY KEY (id)
		// );
		// `,
		// 		},
		{
			name: "success,CREATE_TABLE_TYPE_ANNOTATION",
			input: `CREATE TABLE IF NOT EXISTS public.users (
    user_id VARCHAR(36) NOT NULL,
    username VARCHAR(256) NOT NULL,
    is_verified BOOL NOT NULL DEFAULT false,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_pkey PRIMARY KEY (user_id),
    INDEX users_idx_by_username (username DESC)
);
`,
			wantErr: nil,
			wantStr: `CREATE TABLE IF NOT EXISTS public.users (
    user_id VARCHAR(36) NOT NULL,
    username VARCHAR(256) NOT NULL,
    is_verified BOOL NOT NULL DEFAULT false,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_pkey PRIMARY KEY (user_id),
    INDEX users_idx_by_username (username DESC)
);
`,
		},
	}

	for _, tt := range successTests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()

			l := NewLexer(tt.input)
			p := NewParser(l)
			actual, err := p.Parse()
			require.ErrorIs(t, err, tt.wantErr)

			if !assert.Equal(t, tt.wantStr, actual.String()) {
				t.Fail()
			}

			t.Logf("✅: %s: actual: %%#v: \n%#v", t.Name(), actual)
			t.Logf("✅: %s: actual: %%s: \n%s", t.Name(), actual)
		})
	}

	failureTests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "failure,invalid",
			input:   `)invalid`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_INVALID",
			input:   `CREATE INVALID;`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_INVALID",
			input:   `CREATE TABLE;`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_IF_INVALID",
			input:   `CREATE TABLE IF;`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_IF_NOT_INVALID",
			input:   `CREATE TABLE IF NOT;`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_INVALID",
			input:   `CREATE TABLE "users";`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_INVALID",
			input:   `CREATE TABLE "users" ("id";`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_data_type_INVALID",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36);`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_INVALID",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36), CONSTRAINT "invalid" NOT;`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_INVALID",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36))(;`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_COMMA_INVALID",
			input:   `CREATE TABLE "users" ("id" TIMESTAMP CREATE`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_DATA_TYPE_INVALID",
			input:   `CREATE TABLE "users" ("id" VARYING();`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_INVALID_NOT",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36) NULL NOT;`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_INVALID_DEFAULT",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36) DEFAULT ("id")`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_INVALID_DEFAULT_OPEN_PAREN",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36) DEFAULT ("id",`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_INVALID_PRIMARY_KEY",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36) PRIMARY NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_INVALID_REFERENCES",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36) REFERENCES NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_INVALID_REFERENCES_IDENTS",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36) REFERENCES "groups" (NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_INVALID_CHECK",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36) CHECK NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CHECK_INVALID_IDENTS",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36) CHECK (NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_INVALID_IDENT",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36), CONSTRAINT NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_INVALID_PRIMARY",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36), PRIMARY NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_INVALID_PRIMARY_KEY",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36), PRIMARY KEY NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_INVALID_PRIMARY_KEY_OPEN_PAREN",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36), PRIMARY KEY (NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_INVALID_FOREIGN",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36), FOREIGN NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_INVALID_FOREIGN_KEY",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36), FOREIGN KEY NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_INVALID_FOREIGN_KEY_OPEN_PAREN",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36), FOREIGN KEY (NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_FOREIGN_KEY_IDENTS_INVALID",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36), FOREIGN KEY ("group_id") NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_FOREIGN_KEY_IDENTS_REFERENCES_INVALID",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36), FOREIGN KEY ("group_id") REFERENCES `,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_FOREIGN_KEY_IDENTS_REFERENCES_INVALID_IDENTS",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36), FOREIGN KEY ("group_id") REFERENCES "groups" NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_FOREIGN_KEY_IDENTS_REFERENCES_INVALID_CLOSE_PAREN",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36), FOREIGN KEY ("group_id") REFERENCES "groups" ("id")`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_UNIQUE_INVALID",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36), UNIQUE NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_UNIQUE_INVALID",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36), UNIQUE NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_UNIQUE_INDEX_INVALID",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36), UNIQUE INDEX users_idx_name NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_UNIQUE_INDEX_COLUMN_INVALID",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36), UNIQUE INDEX users_idx_name (NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_UNIQUE_INVALID",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36), UNIQUE INDEX NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_UNIQUE_IDENTS_INVALID",
			input:   `CREATE TABLE "users" ("id" VARCHAR(36), name TEXT, UNIQUE ("id", name)`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_INDEX_INVALID",
			input:   `CREATE INDEX NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_INDEX_IF_INVALID",
			input:   `CREATE INDEX IF;`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_INDEX_IF_NOT_INVALID",
			input:   `CREATE INDEX IF NOT;`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_INDEX_IF_NOT_EXISTS_INVALID",
			input:   `CREATE INDEX IF NOT EXISTS;`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_INDEX_index_name_INVALID",
			input:   `CREATE INDEX users_idx_username NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_INDEX_index_name_ON_INVALID",
			input:   `CREATE INDEX users_idx_username ON NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_INDEX_index_name_ON_table_name_INVALID",
			input:   `CREATE INDEX users_idx_username ON users NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_INDEX_index_name_ON_table_name_USING_INVALID",
			input:   `CREATE INDEX users_idx_username ON users USING NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_INDEX_index_name_ON_table_name_USING_method_INVALID",
			input:   `CREATE INDEX users_idx_username ON users USING btree NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_INDEX_index_name_ON_table_name_USING_method_OPEN_INVALID",
			input:   `CREATE INDEX users_idx_username ON users USING btree (NOT)`,
			wantErr: ddl.ErrUnexpectedToken,
		},
	}

	for _, tt := range failureTests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewParser(NewLexer(tt.input)).Parse()
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestParser_parseColumn(t *testing.T) {
	t.Parallel()

	t.Run("failure,invalid", func(t *testing.T) {
		t.Parallel()

		_, _, err := NewParser(NewLexer(`NOT`)).parseColumn(&Ident{Name: "table_name", QuotationMark: `"`, Raw: `"table_name"`})
		require.ErrorIs(t, err, ddl.ErrUnexpectedToken)
	})
}

func TestParser_parseExpr(t *testing.T) {
	t.Parallel()

	t.Run("failure,invalid", func(t *testing.T) {
		t.Parallel()

		p := NewParser(NewLexer(`NOT`))
		p.nextToken()
		p.nextToken()
		_, err := p.parseExpr()
		require.ErrorIs(t, err, ddl.ErrUnexpectedToken)
	})

	t.Run("failure,invalid2", func(t *testing.T) {
		t.Parallel()

		p := NewParser(NewLexer(`((NOT`))
		p.nextToken()
		p.nextToken()
		_, err := p.parseExpr()
		require.ErrorIs(t, err, ddl.ErrUnexpectedToken)
	})
}

func TestParser_parseDataType(t *testing.T) {
	t.Parallel()

	t.Run("failure,DOUBLE_NOT", func(t *testing.T) {
		t.Parallel()

		p := NewParser(NewLexer(`DOUBLE NOT`))
		p.nextToken()
		p.nextToken()
		_, err := p.parseDataType()
		require.ErrorIs(t, err, ddl.ErrUnexpectedToken)
	})

	t.Run("failure,DOUBLE_PRECISION", func(t *testing.T) {
		t.Parallel()

		p := NewParser(NewLexer(`DOUBLE PRECISION(NOT`))
		p.nextToken()
		p.nextToken()
		_, err := p.parseDataType()
		require.ErrorIs(t, err, ddl.ErrUnexpectedToken)
	})

	t.Run("failure,CHARACTER_NOT", func(t *testing.T) {
		t.Parallel()

		p := NewParser(NewLexer(`CHARACTER NOT`))
		p.nextToken()
		p.nextToken()
		_, err := p.parseDataType()
		require.ErrorIs(t, err, ddl.ErrUnexpectedToken)
	})

	t.Run("failure,CHARACTER_VARYING_NOT", func(t *testing.T) {
		t.Parallel()

		p := NewParser(NewLexer(`CHARACTER VARYING(NOT`))
		p.nextToken()
		p.nextToken()
		_, err := p.parseDataType()
		require.ErrorIs(t, err, ddl.ErrUnexpectedToken)
	})
}
