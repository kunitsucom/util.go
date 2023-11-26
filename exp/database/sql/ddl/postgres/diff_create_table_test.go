package postgres

import (
	"fmt"
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

		expected := &DDL{
			Stmts: []Stmt{
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
					Action: &AddColumn{
						Column: &Column{
							Name:     &Ident{Name: "age", QuotationMark: `"`, Raw: `"age"`},
							DataType: &DataType{Name: "INTEGER"},
							Default: &Default{
								Value: &DefaultValue{
									[]*Ident{
										{Name: "0", QuotationMark: "", Raw: "0"},
									},
								},
							},
							NotNull: true,
						},
					},
				},
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
					Action: &AddConstraint{
						Constraint: &CheckConstraint{
							Name: &Ident{Name: "users_age_check", QuotationMark: ``, Raw: "users_age_check"},
							Expr: []*Ident{
								{Name: "age", QuotationMark: `"`, Raw: `"age"`},
								{Name: ">=", Raw: ">="},
								{Name: "0", Raw: "0"},
							},
						},
					},
				},
			},
		}
		expectedStr := `ALTER TABLE "users" ADD COLUMN "age" INTEGER DEFAULT 0 NOT NULL;
ALTER TABLE "users" ADD CONSTRAINT users_age_check CHECK ("age" >= 0);
`

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
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

		expected := &DDL{
			Stmts: []Stmt{
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
					Action: &DropConstraint{
						Name: &Ident{Name: "users_unique_name", QuotationMark: ``, Raw: "users_unique_name"},
					},
				},
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
					Action: &DropConstraint{
						Name: &Ident{Name: "users_age_check", QuotationMark: ``, Raw: "users_age_check"},
					},
				},
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
					Action: &DropColumn{
						Name: &Ident{Name: "age", QuotationMark: `"`, Raw: `"age"`},
					},
				},
			},
		}
		expectedStr := `ALTER TABLE "users" DROP CONSTRAINT users_unique_name;
ALTER TABLE "users" DROP CONSTRAINT users_age_check;
ALTER TABLE "users" DROP COLUMN "age";
`

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
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

		expected := &DDL{
			Stmts: []Stmt{
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
					Action: &AlterColumn{
						Name: &Ident{Name: "name", QuotationMark: `"`, Raw: `"name"`},
						Action: &AlterColumnSetDataType{
							DataType: &DataType{Name: "TEXT"},
						},
					},
				},
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
					Action: &AlterColumn{
						Name: &Ident{Name: "age", QuotationMark: `"`, Raw: `"age"`},
						Action: &AlterColumnSetDataType{
							DataType: &DataType{Name: "BIGINT"},
						},
					},
				},
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
					Action: &AddConstraint{
						Constraint: &UniqueConstraint{
							Name: &Ident{
								Name:          "users_unique_name",
								QuotationMark: ``,
								Raw:           "users_unique_name",
							},
							Columns: []*ColumnIdent{
								{Ident: &Ident{Name: "name", QuotationMark: `"`, Raw: `"name"`}},
							},
						},
					},
				},
			},
		}
		expectedStr := `ALTER TABLE "users" ALTER COLUMN "name" SET DATA TYPE TEXT;
ALTER TABLE "users" ALTER COLUMN "age" SET DATA TYPE BIGINT;
ALTER TABLE "users" ADD CONSTRAINT users_unique_name UNIQUE ("name");
`

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
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

		expected := &DDL{
			Stmts: []Stmt{
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
					Action: &AlterColumn{
						Name:   &Ident{Name: "age", QuotationMark: `"`, Raw: `"age"`},
						Action: &AlterColumnDropDefault{},
					},
				},
			},
		}
		expectedStr := `ALTER TABLE "users" ALTER COLUMN "age" DROP DEFAULT;
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
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

		expected := &DDL{
			Stmts: []Stmt{
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
					Action: &AlterColumn{
						Name: &Ident{Name: "age", QuotationMark: `"`, Raw: `"age"`},
						Action: &AlterColumnSetDefault{
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
					},
				},
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
					Action: &DropConstraint{
						Name: &Ident{
							Name:          "users_age_check",
							QuotationMark: ``,
							Raw:           "users_age_check",
						},
					},
				},
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
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
									Name: "<>",
									Raw:  "<>",
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
		}
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
		assert.Equal(t, expected, actual)
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

		expected := &DDL{
			Stmts: []Stmt{
				&AlterTableStmt{
					Name: &ObjectName{Schema: &Ident{Name: "public", QuotationMark: `"`, Raw: `"public"`}, Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
					Action: &RenameTable{
						NewName: &ObjectName{Schema: &Ident{Name: "public", QuotationMark: `"`, Raw: `"public"`}, Name: &Ident{Name: "app_users", QuotationMark: `"`, Raw: `"app_users"`}},
					},
				},
				&AlterTableStmt{
					Name: &ObjectName{Schema: &Ident{Name: "public", QuotationMark: `"`, Raw: `"public"`}, Name: &Ident{Name: "app_users", QuotationMark: `"`, Raw: `"app_users"`}},
					Action: &DropConstraint{
						Name: &Ident{Name: "users_group_id_fkey", QuotationMark: ``, Raw: "users_group_id_fkey"},
					},
				},
				&AlterTableStmt{
					Name: &ObjectName{Schema: &Ident{Name: "public", QuotationMark: `"`, Raw: `"public"`}, Name: &Ident{Name: "app_users", QuotationMark: `"`, Raw: `"app_users"`}},
					Action: &DropConstraint{
						Name: &Ident{Name: "users_unique_name", QuotationMark: ``, Raw: "users_unique_name"},
					},
				},
				&AlterTableStmt{
					Name: &ObjectName{Schema: &Ident{Name: "public", QuotationMark: `"`, Raw: `"public"`}, Name: &Ident{Name: "app_users", QuotationMark: `"`, Raw: `"app_users"`}},
					Action: &DropConstraint{
						Name: &Ident{Name: "users_age_check", QuotationMark: ``, Raw: "users_age_check"},
					},
				},
				&AlterTableStmt{
					Name: &ObjectName{Schema: &Ident{Name: "public", QuotationMark: `"`, Raw: `"public"`}, Name: &Ident{Name: "app_users", QuotationMark: `"`, Raw: `"app_users"`}},
					Action: &DropConstraint{
						Name: &Ident{Name: "users_pkey", QuotationMark: ``, Raw: "users_pkey"},
					},
				},
				&AlterTableStmt{
					Name: &ObjectName{Schema: &Ident{Name: "public", QuotationMark: `"`, Raw: `"public"`}, Name: &Ident{Name: "app_users", QuotationMark: `"`, Raw: `"app_users"`}},
					Action: &AddConstraint{
						Constraint: &ForeignKeyConstraint{
							Name: &Ident{Name: "app_users_group_id_fkey", QuotationMark: ``, Raw: "app_users_group_id_fkey"},
							Columns: []*ColumnIdent{
								{Ident: &Ident{Name: "group_id", QuotationMark: "", Raw: "group_id"}},
							},
							Ref: &Ident{Name: "groups", QuotationMark: `"`, Raw: `"groups"`},
							RefColumns: []*ColumnIdent{
								{Ident: &Ident{Name: "id", QuotationMark: `"`, Raw: `"id"`}},
							},
						},
					},
				},
				&AlterTableStmt{
					Name: &ObjectName{Schema: &Ident{Name: "public", QuotationMark: `"`, Raw: `"public"`}, Name: &Ident{Name: "app_users", QuotationMark: `"`, Raw: `"app_users"`}},
					Action: &AddConstraint{
						Constraint: &UniqueConstraint{
							Name: &Ident{Name: "app_users_unique_name", QuotationMark: ``, Raw: "app_users_unique_name"},
							Columns: []*ColumnIdent{
								{Ident: &Ident{Name: "name", QuotationMark: `"`, Raw: `"name"`}},
							},
						},
					},
				},
				&AlterTableStmt{
					Name: &ObjectName{Schema: &Ident{Name: "public", QuotationMark: `"`, Raw: `"public"`}, Name: &Ident{Name: "app_users", QuotationMark: `"`, Raw: `"app_users"`}},
					Action: &AddConstraint{
						Constraint: &CheckConstraint{
							Name: &Ident{Name: "app_users_age_check", QuotationMark: ``, Raw: "app_users_age_check"},
							Expr: []*Ident{
								{Name: "age", QuotationMark: `"`, Raw: `"age"`},
								{Name: ">=", Raw: ">="},
								{Name: "0", Raw: "0"},
							},
						},
					},
				},
				&AlterTableStmt{
					Name: &ObjectName{Schema: &Ident{Name: "public", QuotationMark: `"`, Raw: `"public"`}, Name: &Ident{Name: "app_users", QuotationMark: `"`, Raw: `"app_users"`}},
					Action: &AddConstraint{
						Constraint: &PrimaryKeyConstraint{
							Name: &Ident{Name: "app_users_pkey", QuotationMark: ``, Raw: "app_users_pkey"},
							Columns: []*ColumnIdent{
								{Ident: &Ident{Name: "id", QuotationMark: `"`, Raw: `"id"`}},
							},
						},
					},
				},
			},
		}
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
		assert.Equal(t, expected, actual)
		assert.Equal(t, expectedStr, actual.String())
		assert.Equal(t, fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", actual))

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

		expected := &DDL{
			Stmts: []Stmt{
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
					Action: &AlterColumn{
						Name:   &Ident{Name: "age", QuotationMark: `"`, Raw: `"age"`},
						Action: &AlterColumnSetNotNull{},
					},
				},
			},
		}
		expectedStr := `ALTER TABLE "users" ALTER COLUMN "age" SET NOT NULL;
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
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

		expected := &DDL{
			Stmts: []Stmt{
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
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
		}
		expectedStr := `ALTER TABLE "users" ALTER COLUMN "age" DROP NOT NULL;
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
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

		expected := &DDL{
			Stmts: []Stmt{
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
					Action: &DropConstraint{
						Name: &Ident{
							Name:          "users_pkey",
							QuotationMark: ``,
							Raw:           "users_pkey",
						},
					},
				},
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
					Action: &AddConstraint{
						Constraint: &PrimaryKeyConstraint{
							Name: &Ident{
								Name:          "users_pkey",
								QuotationMark: ``,
								Raw:           "users_pkey",
							},
							Columns: []*ColumnIdent{
								{Ident: &Ident{Name: "id", QuotationMark: `"`, Raw: `"id"`}},
								{Ident: &Ident{Name: "name", Raw: `name`}},
							},
						},
					},
				},
			},
		}
		expectedStr := `ALTER TABLE "users" DROP CONSTRAINT users_pkey;
ALTER TABLE "users" ADD CONSTRAINT users_pkey PRIMARY KEY ("id", name);
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
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

		expected := &DDL{
			Stmts: []Stmt{
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
					Action: &DropConstraint{
						Name: &Ident{
							Name:          "users_group_id_fkey",
							QuotationMark: ``,
							Raw:           "users_group_id_fkey",
						},
					},
				},
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
					Action: &AddConstraint{
						Constraint: &ForeignKeyConstraint{
							Name: &Ident{
								Name:          "users_group_id_fkey",
								QuotationMark: ``,
								Raw:           "users_group_id_fkey",
							},
							Columns: []*ColumnIdent{
								{Ident: &Ident{Name: "group_id", Raw: `group_id`}},
								{Ident: &Ident{Name: "name", Raw: `name`}},
							},
							Ref: &Ident{
								Name:          "groups",
								QuotationMark: `"`,
								Raw:           `"groups"`,
							},
							RefColumns: []*ColumnIdent{
								{Ident: &Ident{Name: "id", QuotationMark: `"`, Raw: `"id"`}},
								{Ident: &Ident{Name: "name", Raw: `name`}},
							},
						},
					},
				},
			},
		}
		expectedStr := `ALTER TABLE "users" DROP CONSTRAINT users_group_id_fkey;
ALTER TABLE "users" ADD CONSTRAINT users_group_id_fkey FOREIGN KEY (group_id, name) REFERENCES "groups" ("id", name);
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
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

		expected := &DDL{
			Stmts: []Stmt{
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
					Action: &DropConstraint{
						Name: &Ident{
							Name:          "users_unique_name",
							QuotationMark: ``,
							Raw:           "users_unique_name",
						},
					},
				},
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
					Action: &AddConstraint{
						Constraint: &UniqueConstraint{
							Name: &Ident{
								Name:          "users_unique_name",
								QuotationMark: ``,
								Raw:           "users_unique_name",
							},
							Columns: []*ColumnIdent{
								{Ident: &Ident{Name: "id", QuotationMark: `"`, Raw: `"id"`}},
								{Ident: &Ident{Name: "name", Raw: `name`}},
							},
						},
					},
				},
			},
		}
		expectedStr := `ALTER TABLE "users" DROP CONSTRAINT users_unique_name;
ALTER TABLE "users" ADD CONSTRAINT users_unique_name UNIQUE ("id", name);
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
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

		expected := &DDL{
			Stmts: []Stmt{
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
					Action: &AlterColumn{
						Name: &Ident{
							Name:          "age",
							QuotationMark: `"`,
							Raw:           `"age"`,
						},
						Action: &AlterColumnSetDefault{
							Default: &Default{
								Value: &DefaultValue{
									[]*Ident{
										{Name: "(", Raw: "("},
										{Name: "(", Raw: "("},
										{Name: "0", Raw: "0"},
										{Name: "+", Raw: "+"},
										{Name: "3", Raw: "3"},
										{Name: ")", Raw: ")"},
										{Name: "-", Raw: "-"},
										{Name: "1", Raw: "1"},
										{Name: "*", Raw: "*"},
										{Name: "4", Raw: "4"},
										{Name: "/", Raw: "/"},
										{Name: "2", Raw: "2"},
										{Name: ")", Raw: ")"},
									},
								},
							},
						},
					},
				},
			},
		}
		expectedStr := `ALTER TABLE "users" ALTER COLUMN "age" SET DEFAULT ((0 + 3) - 1 * 4 / 2);
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
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

		expected := &DDL{
			Stmts: []Stmt{
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "complex_defaults", Raw: "complex_defaults"}},
					Action: &AlterColumn{
						Name: &Ident{
							Name: "unique_code",
							Raw:  "unique_code",
						},
						Action: &AlterColumnSetDefault{
							Default: &Default{
								Value: &DefaultValue{
									[]*Ident{{Name: "'CODE-'", Raw: "'CODE-'"}, {Name: "||", Raw: "||"}, {Name: "TO_CHAR", Raw: "TO_CHAR"}, {Name: "(", Raw: "("}, {Name: "NOW", Raw: "NOW"}, {Name: "(", Raw: "("}, {Name: ")", Raw: ")"}, {Name: ",", Raw: ","}, {Name: "'YYYYMMDDHH24MISS'", Raw: "'YYYYMMDDHH24MISS'"}, {Name: ")", Raw: ")"}, {Name: "||", Raw: "||"}, {Name: "'-'", Raw: "'-'"}, {Name: "||", Raw: "||"}, {Name: "LPAD", Raw: "LPAD"}, {Name: "(", Raw: "("}, {Name: "TO_CHAR", Raw: "TO_CHAR"}, {Name: "(", Raw: "("}, {Name: "NEXTVAL", Raw: "NEXTVAL"}, {Name: "(", Raw: "("}, {Name: "'seq_complex_default'", Raw: "'seq_complex_default'"}, {Name: ")", Raw: ")"}, {Name: ")", Raw: ")"}, {Name: ",", Raw: ","}, {Name: "5", Raw: "5"}, {Name: ",", Raw: ","}, {Name: "'0'", Raw: "'0'"}, {Name: ")", Raw: ")"}},
								},
							},
						},
					},
				},
			},
		}
		expectedStr := `ALTER TABLE complex_defaults ALTER COLUMN unique_code SET DEFAULT 'CODE-' || TO_CHAR(NOW(), 'YYYYMMDDHH24MISS') || '-' || LPAD(TO_CHAR(NEXTVAL('seq_complex_default')), 5, '0');
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(false),
		)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
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

		expected := &DDL{
			Stmts: []Stmt{
				&AlterTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
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
		}
		expectedStr := `ALTER TABLE "users" ADD CONSTRAINT users_age_check CHECK ("age" >= 0) NOT VALID;
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			afterDDL.Stmts[0].(*CreateTableStmt),
			DiffCreateTableUseAlterTableAddConstraintNotValid(true),
		)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
		assert.Equal(t, expectedStr, actual.String())

		t.Logf("✅: %s:\n%s", t.Name(), actual)
	})

	t.Run("success,CREATE_TABLE", func(t *testing.T) {
		t.Parallel()

		after := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`

		afterDDL, err := NewParser(NewLexer(after)).Parse()
		require.NoError(t, err)

		expected := &DDL{
			Stmts: []Stmt{
				&CreateTableStmt{
					Indent: Indent,
					Name:   &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
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
						&ForeignKeyConstraint{
							Name: &Ident{
								Name:          "users_group_id_fkey",
								QuotationMark: ``,
								Raw:           "users_group_id_fkey",
							},
							Columns: []*ColumnIdent{
								{Ident: &Ident{Name: "group_id", QuotationMark: "", Raw: "group_id"}},
							},
							Ref: &Ident{
								Name:          "groups",
								QuotationMark: `"`,
								Raw:           `"groups"`,
							},
							RefColumns: []*ColumnIdent{
								{Ident: &Ident{Name: "id", QuotationMark: `"`, Raw: `"id"`}},
							},
						},
						&UniqueConstraint{
							Name: &Ident{
								Name:          "users_unique_name",
								QuotationMark: ``,
								Raw:           "users_unique_name",
							},
							Columns: []*ColumnIdent{
								{Ident: &Ident{Name: "name", QuotationMark: `"`, Raw: `"name"`}},
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
							Columns: []*ColumnIdent{
								{Ident: &Ident{Name: "id", QuotationMark: `"`, Raw: `"id"`}},
							},
						},
					},
				},
			},
		}
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
		assert.Equal(t, expected, actual)
		assert.Equal(t, expectedStr, actual.String())

		t.Logf("✅: %s:\n%s", t.Name(), actual)
	})

	t.Run("success,DROP_TABLE", func(t *testing.T) {
		t.Parallel()

		before := `CREATE TABLE "users" (id UUID NOT NULL, group_id UUID NOT NULL REFERENCES "groups" ("id"), "name" VARCHAR(255) NOT NULL UNIQUE, "age" INT DEFAULT 0 CHECK ("age" >= 0), description TEXT, PRIMARY KEY ("id"));`

		beforeDDL, err := NewParser(NewLexer(before)).Parse()
		require.NoError(t, err)

		expected := &DDL{
			Stmts: []Stmt{
				&DropTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "users", QuotationMark: `"`, Raw: `"users"`}},
				},
			},
		}
		expectedStr := `DROP TABLE "users";
`

		actual, err := DiffCreateTable(
			beforeDDL.Stmts[0].(*CreateTableStmt),
			nil,
			DiffCreateTableUseAlterTableAddConstraintNotValid(true),
		)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
		assert.Equal(t, expectedStr, actual.String())

		t.Logf("✅: %s:\n%s", t.Name(), actual)
	})
}
