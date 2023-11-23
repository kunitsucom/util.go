//nolint:testpackage
package postgres

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/kunitsucom/util.go/exp/database/sql/ddl/internal"
	"github.com/kunitsucom/util.go/exp/diff/simplediff"
	"github.com/kunitsucom/util.go/testing/assert"
	"github.com/kunitsucom/util.go/testing/require"
)

//nolint:paralleltest
func TestParser_Parse(t *testing.T) {
	backup := internal.TraceLog
	t.Cleanup(func() {
		internal.TraceLog = backup
	})
	internal.TraceLog = log.New(os.Stderr, "TRACE: ", log.LstdFlags|log.Lshortfile)

	tests := []struct {
		name    string
		input   string
		want    *DDL
		wantErr error
		wantStr string
	}{
		{
			name:  "success,CREATE TABLE",
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
								Columns: []*Ident{
									{Name: "id", QuotationMark: `"`, Raw: `"id"`},
								},
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
    "name" VARYING(255) NOT NULL,
    "age" INTEGER DEFAULT 0,
    description TEXT,
    CONSTRAINT users_group_id_fkey FOREIGN KEY (group_id) REFERENCES "groups" ("id"),
    CONSTRAINT users_unique_name UNIQUE ("name"),
    CONSTRAINT users_age_check CHECK ("age" >= 0),
    CONSTRAINT users_pkey PRIMARY KEY ("id")
);
`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()

			l := NewLexer(tt.input)
			p := NewParser(l)
			stmt, err := p.Parse()
			require.ErrorIs(t, err, tt.wantErr)

			if !assert.Equal(t, tt.want, stmt) {
				t.Logf("❌: expected != actual:\n--- EXPECTED\n+++ ACTUAL\n%s", simplediff.Diff(fmt.Sprintf("%#v", tt.want), fmt.Sprintf("%#v", stmt)))
			}

			if !assert.Equal(t, tt.wantStr, stmt.String()) {
				t.Fail()
			}

			t.Logf("✅:\n%s", stmt)
		})
	}
}
