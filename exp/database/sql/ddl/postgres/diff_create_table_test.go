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
						Indent: Indent,
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
						Indent: Indent,
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
						Indent: Indent,
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
						Indent: Indent,
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
						Indent: Indent,
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
						Indent: Indent,
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

			ddls, err := DiffCreateTable(before.Stmts[0].(*CreateTableStmt), after.Stmts[0].(*CreateTableStmt))

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, ddls)

			t.Logf("✅:\n%s", ddls)
		})
	}
}
