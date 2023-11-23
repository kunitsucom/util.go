package postgres

import (
	"testing"

	"github.com/kunitsucom/util.go/exp/database/sql/ddl"
	"github.com/kunitsucom/util.go/testing/assert"
	"github.com/kunitsucom/util.go/testing/require"
)

//nolint:paralleltest,tparallel
func TestDiff(t *testing.T) {
	tests := []struct {
		name    string
		before  string
		after   string
		want    *DDL
		wantErr error
	}{
		{
			name:    "failure,ddl.ErrNoDifference",
			before:  `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`,
			after:   `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`,
			want:    nil,
			wantErr: ddl.ErrNoDifference,
		},
		{
			name:   "success,ADD_COLUMN",
			before: `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, description TEXT, PRIMARY KEY ("id"));`,
			after:  `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INTEGER DEFAULT 0 NOT NULL CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`,
			want: &DDL{
				Stmts: []Stmt{
					&AlterTableStmt{
						TableName: &Ident{
							Name:          "users",
							QuotationMark: `"`,
							Raw:           `"users"`,
						},
						Action: &AddColumn{
							Column: &Column{
								Name: &Ident{
									Name:          "age",
									QuotationMark: `"`,
									Raw:           `"age"`,
								},
								DataType: &DataType{
									Name: "INTEGER",
								},
								Default: &Default{
									Value: &Ident{
										Name: "0",
										Raw:  "0",
									},
								},
								NotNull: true,
							},
						},
					},
					&AlterTableStmt{
						TableName: &Ident{
							Name:          "users",
							QuotationMark: `"`,
							Raw:           `"users"`,
						},
						Action: &AddConstraint{
							Constraint: &CheckConstraint{
								Name: &Ident{
									Name:          "users_age_check",
									QuotationMark: ``,
									Raw:           "users_age_check",
								},
								Expr: []*Ident{
									{
										Name:          "age",
										QuotationMark: `"`,
										Raw:           `"age"`,
									},
									{
										Name: ">=",
										Raw:  ">=",
									},
									{
										Name: "0",
										Raw:  "0",
									},
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name:   "success,DROP_COLUMN",
			before: `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INTEGER DEFAULT 0 NOT NULL CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`,
			after:  `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, description TEXT, PRIMARY KEY ("id"));`,
			want: &DDL{
				Stmts: []Stmt{
					&AlterTableStmt{
						TableName: &Ident{
							Name:          "users",
							QuotationMark: `"`,
							Raw:           `"users"`,
						},
						Action: &DropConstraint{
							Name: &Ident{
								Name:          "users_age_check",
								QuotationMark: ``,
								Raw:           "users_age_check",
							},
						},
					},
					&AlterTableStmt{
						TableName: &Ident{
							Name:          "users",
							QuotationMark: `"`,
							Raw:           `"users"`,
						},
						Action: &DropColumn{
							Name: &Ident{
								Name:          "age",
								QuotationMark: `"`,
								Raw:           `"age"`,
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name:   "success,ALTER_COLUMN_SET_DATA_TYPE",
			before: `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`,
			after:  `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" TEXT NOT NULL UNIQUE, "age" BIGINT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`,
			want: &DDL{
				Stmts: []Stmt{
					&AlterTableStmt{
						TableName: &Ident{
							Name:          "users",
							QuotationMark: `"`,
							Raw:           `"users"`,
						},
						Action: &AlterColumn{
							Name: &Ident{
								Name:          "name",
								QuotationMark: `"`,
								Raw:           `"name"`,
							},
							Action: &AlterColumnSetDataType{
								DataType: &DataType{
									Name: "TEXT",
								},
							},
						},
					},
					&AlterTableStmt{
						TableName: &Ident{
							Name:          "users",
							QuotationMark: `"`,
							Raw:           `"users"`,
						},
						Action: &AlterColumn{
							Name: &Ident{
								Name:          "age",
								QuotationMark: `"`,
								Raw:           `"age"`,
							},
							Action: &AlterColumnSetDataType{
								DataType: &DataType{
									Name: "BIGINT",
								},
							},
						},
					},
				},
			},
		},
		{
			name:   "success,ALTER_COLUMN_DROP_DEFAULT",
			before: `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`,
			after:  `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`,
			want: &DDL{
				Stmts: []Stmt{
					&AlterTableStmt{
						TableName: &Ident{
							Name:          "users",
							QuotationMark: `"`,
							Raw:           `"users"`,
						},
						Action: &AlterColumn{
							Name: &Ident{
								Name:          "age",
								QuotationMark: `"`,
								Raw:           `"age"`,
							},
							Action: &AlterColumnDropDefault{},
						},
					},
				},
			},
		},
		{
			name:   "success,ALTER_COLUMN_SET_DEFAULT",
			before: `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`,
			after:  `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`,
			want: &DDL{
				Stmts: []Stmt{
					&AlterTableStmt{
						TableName: &Ident{
							Name:          "users",
							QuotationMark: `"`,
							Raw:           `"users"`,
						},
						Action: &AlterColumn{
							Name: &Ident{
								Name:          "age",
								QuotationMark: `"`,
								Raw:           `"age"`,
							},
							Action: &AlterColumnSetDefault{
								Default: &Default{
									Value: &Ident{
										Name: "0",
										Raw:  "0",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "failure,ddl.ErrNotSupported",
			before:  `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`,
			after:   `CREATE TABLE "app_users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`,
			want:    nil,
			wantErr: ddl.ErrNotSupported,
		},
		{
			name:   "success,SET_NOT_NULL",
			before: `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`,
			after:  `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INTEGER DEFAULT 0 NOT NULL CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`,
			want: &DDL{
				Stmts: []Stmt{
					&AlterTableStmt{
						TableName: &Ident{
							Name:          "users",
							QuotationMark: `"`,
							Raw:           `"users"`,
						},
						Action: &AlterColumn{
							Name: &Ident{
								Name:          "age",
								QuotationMark: `"`,
								Raw:           `"age"`,
							},
							Action: &AlterColumnSetNotNull{},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name:   "success,DROP_NOT_NULL",
			before: `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 NOT NULL CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`,
			after:  `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`,
			want: &DDL{
				Stmts: []Stmt{
					&AlterTableStmt{
						TableName: &Ident{
							Name:          "users",
							QuotationMark: `"`,
							Raw:           `"users"`,
						},
						Action: &AlterColumn{
							Name: &Ident{
								Name:          "age",
								QuotationMark: `"`,
								Raw:           `"age"`,
							},
							Action: &AlterColumnDropNotNull{},
						},
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			before, err := NewParser(NewLexer(tt.before)).Parse()
			require.NoError(t, err)

			after, err := NewParser(NewLexer(tt.after)).Parse()
			require.NoError(t, err)

			ddls, err := DiffCreateTable(
				before.Stmts[0].(*CreateTableStmt),
				after.Stmts[0].(*CreateTableStmt),
				DiffCreateTableUseAlterTableAddConstraintNotValid(false),
			)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, ddls)

			t.Logf("✅:\n%s", ddls)
		})
	}

	t.Run("success,DiffCreateTableUseAlterTableAddConstraintNotValid", func(t *testing.T) {
		t.Parallel()

		before := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0, description TEXT, PRIMARY KEY ("id"));`
		beforeDDL, err := NewParser(NewLexer(before)).Parse()
		require.NoError(t, err)

		after := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`

		afterDDL, err := NewParser(NewLexer(after)).Parse()
		require.NoError(t, err)

		ddls, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(true),
		)

		assert.NoError(t, err)
		assert.Equal(t, &DDL{
			Stmts: []Stmt{
				&AlterTableStmt{
					TableName: &Ident{
						Name:          "users",
						QuotationMark: `"`,
						Raw:           `"users"`,
					},
					Action: &AddConstraint{
						Constraint: &CheckConstraint{
							Name: &Ident{
								Name:          "users_age_check",
								QuotationMark: ``,
								Raw:           "users_age_check",
							},
							Expr: []*Ident{
								{
									Name:          "age",
									QuotationMark: `"`,
									Raw:           `"age"`,
								},
								{
									Name: ">=",
									Raw:  ">=",
								},
								{
									Name: "0",
									Raw:  "0",
								},
							},
						},
						NotValid: true,
					},
				},
			},
		}, ddls)

		t.Logf("✅: %s:\n%s", t.Name(), ddls)
	})

	t.Run("success,CREATE_TABLE", func(t *testing.T) {
		t.Parallel()

		after := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`

		afterDDL, err := NewParser(NewLexer(after)).Parse()
		require.NoError(t, err)

		ddls, err := DiffCreateTable(
			nil,
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(true),
		)

		assert.NoError(t, err)
		assert.Equal(t, &DDL{
			Stmts: []Stmt{
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
								Name: "VARYING",
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
								Value: &Ident{
									Name:          "0",
									QuotationMark: "",
									Raw:           "0",
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
						&ForeignKeyConstraint{
							Name: &Ident{
								Name:          "users_group_id_fkey",
								QuotationMark: ``,
								Raw:           "users_group_id_fkey",
							},
							Columns: []*Ident{
								{Name: "group_id", QuotationMark: "", Raw: "group_id"},
							},
							Ref: &Ident{
								Name:          "groups",
								QuotationMark: `"`,
								Raw:           `"groups"`,
							},
							RefColumns: []*Ident{
								{Name: "id", QuotationMark: `"`, Raw: `"id"`},
							},
						},
						&UniqueConstraint{
							Name: &Ident{
								Name:          "users_unique_name",
								QuotationMark: ``,
								Raw:           "users_unique_name",
							},
							Columns: []*Ident{
								{Name: "name", QuotationMark: `"`, Raw: `"name"`},
							},
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
						&PrimaryKeyConstraint{
							Name: &Ident{
								Name:          "users_pkey",
								QuotationMark: ``,
								Raw:           "users_pkey",
							},
							Columns: []*Ident{
								{Name: "id", QuotationMark: `"`, Raw: `"id"`},
							},
						},
					},
				},
			},
		}, ddls)

		t.Logf("✅: %s:\n%s", t.Name(), ddls)
	})

	t.Run("success,DROP_TABLE", func(t *testing.T) {
		t.Parallel()

		before := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`

		beforeDDL, err := NewParser(NewLexer(before)).Parse()
		require.NoError(t, err)

		ddls, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			nil,
			DiffCreateTableUseAlterTableAddConstraintNotValid(true),
		)

		assert.NoError(t, err)
		assert.Equal(t, &DDL{
			Stmts: []Stmt{
				&DropTableStmt{
					Name: &Ident{
						Name:          "users",
						QuotationMark: `"`,
						Raw:           `"users"`,
					},
				},
			},
		}, ddls)

		t.Logf("✅: %s:\n%s", t.Name(), ddls)
	})
}
