//nolint:testpackage
package cockroachdb

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/kunitsucom/util.go/exp/database/sql/ddl"
	"github.com/kunitsucom/util.go/exp/database/sql/ddl/internal"
	"github.com/kunitsucom/util.go/exp/diff/simplediff"
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
		want    *DDL
		wantErr error
		wantStr string
	}{
		{
			name:  "success,CREATE_TABLE",
			input: `CREATE TABLE "groups" ("id" UUID NOT NULL PRIMARY KEY, description TEXT); CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`,
			want: &DDL{
				Stmts: []Stmt{
					&CreateTableStmt{
						Indent: Indent,
						Name: &Ident{
							Name:          "groups",
							QuotationMark: `"`,
							Raw:           `"groups"`,
						},
						Columns: []*Column{
							{
								Name: &Ident{
									Name:          "id",
									QuotationMark: `"`,
									Raw:           `"id"`,
								},
								DataType: &DataType{
									Name: "UUID",
									Size: "",
								},
								NotNull: true,
							},
							{
								Name: &Ident{
									Name:          "description",
									QuotationMark: "",
									Raw:           "description",
								},
								DataType: &DataType{
									Name: "TEXT",
									Size: "",
								},
							},
						},
						Constraints: []Constraint{
							&PrimaryKeyConstraint{
								Name: &Ident{
									Name:          "groups_pkey",
									QuotationMark: ``,
									Raw:           "groups_pkey",
								},
								Columns: []*ColumnIdent{{Ident: &Ident{Name: "id", QuotationMark: `"`, Raw: `"id"`}}},
							},
						},
					},
					&CreateTableStmt{
						Indent: Indent,
						Name: &Ident{
							Name:          "users",
							QuotationMark: `"`,
							Raw:           `"users"`,
						},
						Columns: []*Column{
							{
								Name: &Ident{
									Name:          "id",
									QuotationMark: "",
									Raw:           "id",
								},
								DataType: &DataType{
									Name: "UUID",
									Size: "",
								},
								NotNull: true,
							},
							{
								Name: &Ident{
									Name:          "group_id",
									QuotationMark: "",
									Raw:           "group_id",
								},
								DataType: &DataType{
									Name: "UUID",
									Size: "",
								},
								NotNull: true,
							},
							{
								Name: &Ident{
									Name:          "name",
									QuotationMark: `"`,
									Raw:           `"name"`,
								},
								DataType: &DataType{
									Name: "VARCHAR",
									Size: "255",
								},
								NotNull: true,
							},
							{
								Name: &Ident{
									Name:          "age",
									QuotationMark: `"`,
									Raw:           `"age"`,
								},
								DataType: &DataType{
									Name: "INTEGER",
									Size: "",
								},
								Default: &Default{
									Value: &DefaultValue{
										[]*Ident{
											{
												Name:          "0",
												QuotationMark: "",
												Raw:           "0",
											},
										},
									},
								},
							},
							{
								Name: &Ident{
									Name:          "description",
									QuotationMark: "",
									Raw:           "description",
								},
								DataType: &DataType{
									Name: "TEXT",
									Size: "",
								},
							},
						},
						Constraints: []Constraint{
							&PrimaryKeyConstraint{
								Name: &Ident{
									Name:          "users_pkey",
									QuotationMark: ``,
									Raw:           "users_pkey",
								},
								Columns: []*ColumnIdent{{Ident: &Ident{Name: "id", QuotationMark: `"`, Raw: `"id"`}}},
							},
							&ForeignKeyConstraint{
								Name: &Ident{
									Name:          "users_group_id_fkey",
									QuotationMark: ``,
									Raw:           "users_group_id_fkey",
								},
								Columns: []*ColumnIdent{{Ident: &Ident{Name: "group_id", QuotationMark: "", Raw: "group_id"}}},
								Ref: &Ident{
									Name:          "groups",
									QuotationMark: `"`,
									Raw:           `"groups"`,
								},
								RefColumns: []*ColumnIdent{{Ident: &Ident{Name: "id", QuotationMark: `"`, Raw: `"id"`}}},
							},
							&IndexConstraint{
								Unique: true,
								Name: &Ident{
									Name:          "users_unique_name",
									QuotationMark: ``,
									Raw:           "users_unique_name",
								},
								Columns: []*ColumnIdent{{Ident: &Ident{Name: "name", QuotationMark: `"`, Raw: `"name"`}}},
							},
							&CheckConstraint{
								Name: &Ident{
									Name:          "users_age_check",
									QuotationMark: ``,
									Raw:           "users_age_check",
								},
								Expr: []*Ident{
									{Name: "age", QuotationMark: `"`, Raw: `"age"`},
									{Name: ">=", Raw: ">="},
									{Name: "0", Raw: "0"},
								},
							},
						},
					},
				},
			},
			wantErr: nil,
			wantStr: `CREATE TABLE "groups" (
    "id" UUID NOT NULL,
    description TEXT,
    CONSTRAINT groups_pkey PRIMARY KEY ("id")
);
CREATE TABLE "users" (
    id UUID NOT NULL,
    group_id UUID NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "age" INTEGER DEFAULT 0,
    description TEXT,
    CONSTRAINT users_pkey PRIMARY KEY ("id"),
    CONSTRAINT users_group_id_fkey FOREIGN KEY (group_id) REFERENCES "groups" ("id"),
    UNIQUE INDEX users_unique_name ("name"),
    CONSTRAINT users_age_check CHECK ("age" >= 0)
);
`,
		},
		{
			name: "success,complex_defaults",
			input: `CREATE TABLE complex_defaults (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    unique_code TEXT DEFAULT 'CODE-' || TO_CHAR(NOW(), 'YYYYMMDDHH24MISS') || '-' || LPAD(TO_CHAR(NEXTVAL('seq_complex_default')), 5, '0'),
    status TEXT DEFAULT 'pending',
    random_number INTEGER DEFAULT FLOOR(RANDOM() * 100::INTEGER)::INTEGER,
    json_data JSONB DEFAULT '{}',
    calculated_value INTEGER DEFAULT (SELECT COUNT(*) FROM another_table)
);
`,
			want: &DDL{
				Stmts: []Stmt{
					&CreateTableStmt{
						Indent: Indent,
						Name: &Ident{
							Name:          "complex_defaults",
							QuotationMark: "",
							Raw:           "complex_defaults",
						},
						Columns: []*Column{
							{
								Name:     &Ident{Name: "id", Raw: "id"},
								DataType: &DataType{Name: "SERIAL", Size: ""},
							},
							{
								Name:     &Ident{Name: "created_at", Raw: "created_at"},
								DataType: &DataType{Name: "TIMESTAMP WITH TIME ZONE", Size: ""},
								Default:  &Default{Value: &DefaultValue{[]*Ident{{Name: "CURRENT_TIMESTAMP", Raw: "CURRENT_TIMESTAMP"}}}},
							},
							{
								Name:     &Ident{Name: "updated_at", Raw: "updated_at"},
								DataType: &DataType{Name: "TIMESTAMP WITH TIME ZONE", Size: ""},
								Default:  &Default{Value: &DefaultValue{[]*Ident{{Name: "CURRENT_TIMESTAMP", Raw: "CURRENT_TIMESTAMP"}}}},
							},
							{
								Name:     &Ident{Name: "unique_code", Raw: "unique_code"},
								DataType: &DataType{Name: "TEXT", Size: ""},
								Default:  &Default{Value: &DefaultValue{[]*Ident{{Name: "'CODE-'", Raw: "'CODE-'"}, {Name: "||", Raw: "||"}, {Name: "TO_CHAR", Raw: "TO_CHAR"}, {Name: "(", Raw: "("}, {Name: "NOW", Raw: "NOW"}, {Name: "(", Raw: "("}, {Name: ")", Raw: ")"}, {Name: ",", Raw: ","}, {Name: "'YYYYMMDDHH24MISS'", Raw: "'YYYYMMDDHH24MISS'"}, {Name: ")", Raw: ")"}, {Name: "||", Raw: "||"}, {Name: "'-'", Raw: "'-'"}, {Name: "||", Raw: "||"}, {Name: "LPAD", Raw: "LPAD"}, {Name: "(", Raw: "("}, {Name: "TO_CHAR", Raw: "TO_CHAR"}, {Name: "(", Raw: "("}, {Name: "NEXTVAL", Raw: "NEXTVAL"}, {Name: "(", Raw: "("}, {Name: "'seq_complex_default'", Raw: "'seq_complex_default'"}, {Name: ")", Raw: ")"}, {Name: ")", Raw: ")"}, {Name: ",", Raw: ","}, {Name: "5", Raw: "5"}, {Name: ",", Raw: ","}, {Name: "'0'", Raw: "'0'"}, {Name: ")", Raw: ")"}}}},
							},
							{
								Name:     &Ident{Name: "status", Raw: "status"},
								DataType: &DataType{Name: "TEXT", Size: ""},
								Default:  &Default{Value: &DefaultValue{[]*Ident{{Name: "'pending'", Raw: "'pending'"}}}},
							},
							{
								Name:     &Ident{Name: "random_number", Raw: "random_number"},
								DataType: &DataType{Name: "INTEGER", Size: ""},
								Default:  &Default{Value: &DefaultValue{[]*Ident{{Name: "FLOOR", Raw: "FLOOR"}, {Name: "(", Raw: "("}, {Name: "RANDOM", Raw: "RANDOM"}, {Name: "(", Raw: "("}, {Name: ")", Raw: ")"}, {Name: "*", Raw: "*"}, {Name: "100", Raw: "100"}, {Name: "::", Raw: "::"}, {Name: "INTEGER", Raw: "INTEGER"}, {Name: ")", Raw: ")"}, {Name: "::", Raw: "::"}, {Name: "INTEGER", Raw: "INTEGER"}}}},
							},
							{
								Name:     &Ident{Name: "json_data", Raw: "json_data"},
								DataType: &DataType{Name: "JSONB", Size: ""},
								Default:  &Default{Value: &DefaultValue{[]*Ident{{Name: "'{}'", Raw: "'{}'"}}}},
							},
							{
								Name:     &Ident{Name: "calculated_value", Raw: "calculated_value"},
								DataType: &DataType{Name: "INTEGER", Size: ""},
								Default:  &Default{Value: &DefaultValue{[]*Ident{{Name: "(", Raw: "("}, {Name: "SELECT", Raw: "SELECT"}, {Name: "COUNT", Raw: "COUNT"}, {Name: "(", Raw: "("}, {Name: "*", Raw: "*"}, {Name: ")", Raw: ")"}, {Name: "FROM", Raw: "FROM"}, {Name: "another_table", Raw: "another_table"}, {Name: ")", Raw: ")"}}}},
							},
						},
						Constraints: []Constraint{
							&PrimaryKeyConstraint{
								Name:    &Ident{Name: "complex_defaults_pkey", Raw: "complex_defaults_pkey"},
								Columns: []*ColumnIdent{{Ident: &Ident{Name: "id", Raw: "id"}}},
							},
						},
					},
				},
			},
			wantErr: nil,
			wantStr: `CREATE TABLE complex_defaults (
    id SERIAL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    unique_code TEXT DEFAULT 'CODE-' || TO_CHAR(NOW(), 'YYYYMMDDHH24MISS') || '-' || LPAD(TO_CHAR(NEXTVAL('seq_complex_default')), 5, '0'),
    status TEXT DEFAULT 'pending',
    random_number INTEGER DEFAULT FLOOR(RANDOM() * 100::INTEGER)::INTEGER,
    json_data JSONB DEFAULT '{}',
    calculated_value INTEGER DEFAULT (SELECT COUNT(*) FROM another_table),
    CONSTRAINT complex_defaults_pkey PRIMARY KEY (id)
);
`,
		},
		{
			name: "success,CREATE_TABLE_TYPE_ANNOTATION",
			input: `CREATE TABLE public.users (
    user_id UUID NOT NULL,
    username VARCHAR(256) NOT NULL,
    is_verified BOOL NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('UTC':::STRING, current_timestamp():::TIMESTAMPTZ),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('UTC':::STRING, current_timestamp():::TIMESTAMPTZ),
    CONSTRAINT users_pkey PRIMARY KEY (user_id ASC),
    INDEX users_idx_by_username (username ASC)
);
`,
			want: &DDL{
				Stmts: []Stmt{
					&CreateTableStmt{
						Indent: Indent,
						Schema: &Ident{Name: "public", Raw: "public"},
						Name:   &Ident{Name: "users", Raw: "users"},
						Columns: []*Column{
							{Name: &Ident{Name: "user_id", Raw: "user_id"}, DataType: &DataType{Name: "UUID", Size: ""}, NotNull: true},
							{Name: &Ident{Name: "username", Raw: "username"}, DataType: &DataType{Name: "VARCHAR", Size: "256"}, NotNull: true},
							{Name: &Ident{Name: "is_verified", Raw: "is_verified"}, DataType: &DataType{Name: "BOOL", Size: ""}, NotNull: true, Default: &Default{Value: &DefaultValue{[]*Ident{{Name: "false", Raw: "false"}}}}},
							{Name: &Ident{Name: "created_at", Raw: "created_at"}, DataType: &DataType{Name: "TIMESTAMPTZ", Size: ""}, NotNull: true, Default: &Default{Value: &DefaultValue{[]*Ident{{Name: "timezone", Raw: "timezone"}, {Name: "(", Raw: "("}, {Name: "'UTC'", Raw: "'UTC'"}, {Name: ":::", Raw: ":::"}, {Name: "STRING", Raw: "STRING"}, {Name: ",", Raw: ","}, {Name: "current_timestamp", Raw: "current_timestamp"}, {Name: "(", Raw: "("}, {Name: ")", Raw: ")"}, {Name: ":::", Raw: ":::"}, {Name: "TIMESTAMPTZ", Raw: "TIMESTAMPTZ"}, {Name: ")", Raw: ")"}}}}},
							{Name: &Ident{Name: "updated_at", Raw: "updated_at"}, DataType: &DataType{Name: "TIMESTAMPTZ", Size: ""}, NotNull: true, Default: &Default{Value: &DefaultValue{[]*Ident{{Name: "timezone", Raw: "timezone"}, {Name: "(", Raw: "("}, {Name: "'UTC'", Raw: "'UTC'"}, {Name: ":::", Raw: ":::"}, {Name: "STRING", Raw: "STRING"}, {Name: ",", Raw: ","}, {Name: "current_timestamp", Raw: "current_timestamp"}, {Name: "(", Raw: "("}, {Name: ")", Raw: ")"}, {Name: ":::", Raw: ":::"}, {Name: "TIMESTAMPTZ", Raw: "TIMESTAMPTZ"}, {Name: ")", Raw: ")"}}}}},
						},
						Constraints: []Constraint{
							&PrimaryKeyConstraint{
								Name:    &Ident{Name: "users_pkey", Raw: "users_pkey"},
								Columns: []*ColumnIdent{{Ident: &Ident{Name: "user_id", Raw: "user_id"}, Order: &Order{Desc: false}}},
							},
							&IndexConstraint{
								Name:    &Ident{Name: "users_idx_by_username", Raw: "users_idx_by_username"},
								Columns: []*ColumnIdent{{Ident: &Ident{Name: "username", Raw: "username"}, Order: &Order{Desc: false}}},
							},
						},
					},
				},
			},
			wantErr: nil,
			wantStr: `CREATE TABLE public.users (
    user_id UUID NOT NULL,
    username VARCHAR(256) NOT NULL,
    is_verified BOOL NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('UTC':::STRING, current_timestamp():::TIMESTAMPTZ),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('UTC':::STRING, current_timestamp():::TIMESTAMPTZ),
    CONSTRAINT users_pkey PRIMARY KEY (user_id ASC),
    INDEX users_idx_by_username (username ASC)
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
			stmt, err := p.Parse()
			require.ErrorIs(t, err, tt.wantErr)

			if !assert.Equal(t, tt.want, stmt) {
				t.Logf("❌: %s: expected != actual:\n--- EXPECTED\n+++ ACTUAL\n%s", t.Name(), simplediff.Diff(fmt.Sprintf("%#v", tt.want), fmt.Sprintf("%#v", stmt)))
			}

			if !assert.Equal(t, tt.wantStr, stmt.String()) {
				t.Fail()
			}

			t.Logf("✅: %s: stmt: %%#v: \n%#v", t.Name(), stmt)
			t.Logf("✅: %s: stmt: %%s: \n%s", t.Name(), stmt)
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
			name:    "failure,CREATE",
			input:   `CREATE INVALID;`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE",
			input:   `CREATE TABLE;`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name",
			input:   `CREATE TABLE "users";`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name",
			input:   `CREATE TABLE "users" ("id";`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_data_type",
			input:   `CREATE TABLE "users" ("id" UUID;`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT",
			input:   `CREATE TABLE "users" ("id" UUID, CONSTRAINT "invalid" NOT;`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_INVALID",
			input:   `CREATE TABLE "users" ("id" UUID)(;`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_COMMA_INVALID",
			input:   `CREATE TABLE "users" ("id" UUID,(;`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_DATA_TYPE_INVALID",
			input:   `CREATE TABLE "users" ("id" VARYING();`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_INVALID_NOT",
			input:   `CREATE TABLE "users" ("id" UUID NULL NOT;`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_INVALID_DEFAULT",
			input:   `CREATE TABLE "users" ("id" UUID DEFAULT ("id")`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_INVALID_DEFAULT_OPEN_PAREN",
			input:   `CREATE TABLE "users" ("id" UUID DEFAULT ("id",`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_INVALID_PRIMARY_KEY",
			input:   `CREATE TABLE "users" ("id" UUID PRIMARY NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_INVALID_REFERENCES",
			input:   `CREATE TABLE "users" ("id" UUID REFERENCES NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_INVALID_REFERENCES_IDENTS",
			input:   `CREATE TABLE "users" ("id" UUID REFERENCES "groups" (NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_INVALID_CHECK",
			input:   `CREATE TABLE "users" ("id" UUID CHECK NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CHECK_INVALID_IDENTS",
			input:   `CREATE TABLE "users" ("id" UUID CHECK (NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_INVALID_IDENT",
			input:   `CREATE TABLE "users" ("id" UUID, CONSTRAINT NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_INVALID_PRIMARY",
			input:   `CREATE TABLE "users" ("id" UUID, PRIMARY NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_INVALID_PRIMARY_KEY",
			input:   `CREATE TABLE "users" ("id" UUID, PRIMARY KEY NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_INVALID_PRIMARY_KEY_OPEN_PAREN",
			input:   `CREATE TABLE "users" ("id" UUID, PRIMARY KEY (NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_INVALID_FOREIGN",
			input:   `CREATE TABLE "users" ("id" UUID, FOREIGN NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_INVALID_FOREIGN_KEY",
			input:   `CREATE TABLE "users" ("id" UUID, FOREIGN KEY NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_INVALID_FOREIGN_KEY_OPEN_PAREN",
			input:   `CREATE TABLE "users" ("id" UUID, FOREIGN KEY (NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_FOREIGN_KEY_IDENTS_INVALID",
			input:   `CREATE TABLE "users" ("id" UUID, FOREIGN KEY ("group_id") NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_FOREIGN_KEY_IDENTS_REFERENCES_INVALID",
			input:   `CREATE TABLE "users" ("id" UUID, FOREIGN KEY ("group_id") REFERENCES `,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_FOREIGN_KEY_IDENTS_REFERENCES_INVALID_IDENTS",
			input:   `CREATE TABLE "users" ("id" UUID, FOREIGN KEY ("group_id") REFERENCES "groups" NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_FOREIGN_KEY_IDENTS_REFERENCES_INVALID_CLOSE_PAREN",
			input:   `CREATE TABLE "users" ("id" UUID, FOREIGN KEY ("group_id") REFERENCES "groups" ("id")`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_UNIQUE_INVALID",
			input:   `CREATE TABLE "users" ("id" UUID, UNIQUE NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_UNIQUE_IDENTS_INVALID",
			input:   `CREATE TABLE "users" ("id" UUID, UNIQUE (NOT`,
			wantErr: ddl.ErrUnexpectedToken,
		},
		{
			name:    "failure,CREATE_TABLE_table_name_column_name_CONSTRAINT_UNIQUE_IDENTS_INVALID",
			input:   `CREATE TABLE "users" ("id" UUID, name TEXT, UNIQUE ("id", name)`,
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

	t.Run("failure,TIMESTAMP_WITH_NOT", func(t *testing.T) {
		t.Parallel()

		p := NewParser(NewLexer(`TIMESTAMP WITH NOT`))
		p.nextToken()
		p.nextToken()
		_, err := p.parseDataType()
		require.ErrorIs(t, err, ddl.ErrUnexpectedToken)
	})

	t.Run("failure,TIMESTAMP_WITH_TIME_NOT", func(t *testing.T) {
		t.Parallel()

		p := NewParser(NewLexer(`TIMESTAMP WITH TIME NOT`))
		p.nextToken()
		p.nextToken()
		_, err := p.parseDataType()
		require.ErrorIs(t, err, ddl.ErrUnexpectedToken)
	})

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
