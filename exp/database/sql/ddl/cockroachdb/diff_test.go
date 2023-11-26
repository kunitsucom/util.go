package cockroachdb

import (
	"fmt"
	"testing"

	"github.com/kunitsucom/util.go/exp/database/sql/ddl"
	"github.com/kunitsucom/util.go/testing/assert"
	"github.com/kunitsucom/util.go/testing/require"
)

func TestDiff(t *testing.T) {
	t.Parallel()

	t.Run("failure,ddl.ErrNoDifference", func(t *testing.T) {
		t.Parallel()

		before := &DDL{}
		after := &DDL{}
		_, err := Diff(before, after)
		require.ErrorIs(t, err, ddl.ErrNoDifference)
	})

	t.Run("failure,ddl.ErrNotSupported,DropTableStmt", func(t *testing.T) {
		t.Parallel()

		{
			before := &DDL{
				Stmts: []Stmt{
					&DropTableStmt{Name: &ObjectName{Name: &Ident{Name: "table_name", Raw: "table_name"}}},
				},
			}
			after := (*DDL)(nil)
			_, err := Diff(before, after)
			require.ErrorIs(t, err, ddl.ErrNotSupported)
		}
		{
			before := &DDL{
				Stmts: []Stmt{
					&DropTableStmt{Name: &ObjectName{Name: &Ident{Name: "table_name", Raw: "table_name"}}},
				},
			}
			after := &DDL{}
			_, err := Diff(before, after)
			require.ErrorIs(t, err, ddl.ErrNotSupported)
		}
		{
			before := &DDL{}
			after := &DDL{
				Stmts: []Stmt{
					&DropTableStmt{Name: &ObjectName{Name: &Ident{Name: "table_name", Raw: "table_name"}}},
				},
			}
			_, err := Diff(before, after)
			require.ErrorIs(t, err, ddl.ErrNotSupported)
		}
	})

	t.Run("success,after", func(t *testing.T) {
		t.Parallel()

		before := (*DDL)(nil)
		after := &DDL{
			Stmts: []Stmt{
				&CreateTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "table_name", Raw: "table_name"}},
					Columns: []*Column{
						{
							Name: &Ident{Name: "column_name", Raw: "column_name"},
							DataType: &DataType{
								Name: "STRING",
							},
							NotNull: true,
						},
					},
					Constraints: []Constraint{
						&PrimaryKeyConstraint{
							Columns: []*ColumnIdent{
								{
									Ident: &Ident{Name: "column_name", Raw: "column_name"},
								},
							},
						},
					},
				},
			},
		}
		actual, err := Diff(before, after)
		require.NoError(t, err)
		if !assert.Equal(t, after, actual) {
			assert.Equal(t, fmt.Sprintf("%#v", after), fmt.Sprintf("%#v", actual))
		}
		assert.Equal(t, `CREATE TABLE table_name (
    column_name STRING NOT NULL,
    PRIMARY KEY (column_name)
);
`, actual.String())
	})

	t.Run("success,before,nil,Table", func(t *testing.T) {
		t.Parallel()

		before := &DDL{
			Stmts: []Stmt{
				&CreateTableStmt{
					Name: &ObjectName{Schema: &Ident{Name: "public", Raw: "public"}, Name: &Ident{Name: "table_name", Raw: "table_name"}},
					Columns: []*Column{
						{
							Name: &Ident{Name: "column_name", Raw: "column_name"},
						},
					},
				},
			},
		}
		after := (*DDL)(nil)
		actual, err := Diff(before, after)
		require.NoError(t, err)
		expected := &DDL{
			Stmts: []Stmt{
				&DropTableStmt{
					Name: &ObjectName{Schema: &Ident{Name: "public", Raw: "public"}, Name: &Ident{Name: "table_name", Raw: "table_name"}},
				},
			},
		}
		if !assert.Equal(t, expected, actual) {
			assert.Equal(t, fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", actual))
		}
		assert.Equal(t, `DROP TABLE public.table_name;
`, actual.String())
	})

	t.Run("success,before,Table", func(t *testing.T) {
		t.Parallel()

		before := &DDL{
			Stmts: []Stmt{
				&CreateTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "table_name", Raw: "table_name"}},
					Columns: []*Column{
						{
							Name: &Ident{Name: "column_name", Raw: "column_name"},
						},
					},
				},
			},
		}
		after := &DDL{}
		actual, err := Diff(before, after)
		require.NoError(t, err)
		expected := &DDL{
			Stmts: []Stmt{
				&DropTableStmt{
					Name: &ObjectName{Name: &Ident{Name: "table_name", Raw: "table_name"}},
				},
			},
		}
		if !assert.Equal(t, expected, actual) {
			assert.Equal(t, fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", actual))
		}
		assert.Equal(t, `DROP TABLE table_name;
`, actual.String())
	})

	t.Run("success,before,nil,Index", func(t *testing.T) {
		t.Parallel()

		before := &DDL{
			Stmts: []Stmt{
				&CreateIndexStmt{
					Name: &ObjectName{Name: &Ident{Name: "table_name_idx_column_name", Raw: "table_name_idx_column_name"}},
					Columns: []*ColumnIdent{
						{
							Ident: &Ident{Name: "column_name", Raw: "column_name"},
						},
					},
				},
			},
		}
		after := (*DDL)(nil)
		actual, err := Diff(before, after)
		require.NoError(t, err)
		expected := &DDL{
			Stmts: []Stmt{
				&DropIndexStmt{
					Name: &ObjectName{Name: &Ident{Name: "table_name_idx_column_name", Raw: "table_name_idx_column_name"}},
				},
			},
		}
		if !assert.Equal(t, expected, actual) {
			assert.Equal(t, fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", actual))
		}
		assert.Equal(t, `DROP INDEX table_name_idx_column_name;
`, actual.String())
	})

	t.Run("success,before,Index", func(t *testing.T) {
		t.Parallel()

		before := &DDL{
			Stmts: []Stmt{
				&CreateIndexStmt{
					Name: &ObjectName{Name: &Ident{Name: "table_name_idx_column_name", Raw: "table_name_idx_column_name"}},
					Columns: []*ColumnIdent{
						{
							Ident: &Ident{Name: "column_name", Raw: "column_name"},
						},
					},
				},
			},
		}
		after := &DDL{}
		actual, err := Diff(before, after)
		require.NoError(t, err)
		expected := &DDL{
			Stmts: []Stmt{
				&DropIndexStmt{
					Name: &ObjectName{Name: &Ident{Name: "table_name_idx_column_name", Raw: "table_name_idx_column_name"}},
				},
			},
		}
		if !assert.Equal(t, expected, actual) {
			assert.Equal(t, fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", actual))
		}
		assert.Equal(t, `DROP INDEX table_name_idx_column_name;
`, actual.String())
	})

	t.Run("success,before,Table", func(t *testing.T) {
		t.Parallel()

		before := &DDL{}
		after := &DDL{
			Stmts: []Stmt{
				&CreateTableStmt{
					Name: &ObjectName{Schema: &Ident{Name: "public", Raw: "public"}, Name: &Ident{Name: "table_name", Raw: "table_name"}},
					Columns: []*Column{
						{
							Name: &Ident{Name: "column_name", Raw: "column_name"},
							DataType: &DataType{
								Name: "STRING",
							},
							NotNull: true,
						},
					},
					Constraints: []Constraint{
						&PrimaryKeyConstraint{
							Columns: []*ColumnIdent{
								{
									Ident: &Ident{Name: "column_name", Raw: "column_name"},
								},
							},
						},
					},
				},
			},
		}
		actual, err := Diff(before, after)
		require.NoError(t, err)
		if !assert.Equal(t, after, actual) {
			assert.Equal(t, fmt.Sprintf("%#v", after), fmt.Sprintf("%#v", actual))
		}
		assert.Equal(t, `CREATE TABLE public.table_name (
    column_name STRING NOT NULL,
    PRIMARY KEY (column_name)
);
`, actual.String())
	})

	t.Run("success,before,Index", func(t *testing.T) {
		t.Parallel()

		before := &DDL{}
		after := &DDL{
			Stmts: []Stmt{
				&CreateIndexStmt{
					Name:      &ObjectName{Name: &Ident{Name: "table_name_idx_column_name", Raw: "table_name_idx_column_name"}},
					TableName: &ObjectName{Name: &Ident{Name: "table_name", Raw: "table_name"}},
					Columns: []*ColumnIdent{
						{
							Ident: &Ident{Name: "column_name", Raw: "column_name"},
						},
					},
				},
			},
		}
		actual, err := Diff(before, after)
		require.NoError(t, err)
		if !assert.Equal(t, after, actual) {
			assert.Equal(t, fmt.Sprintf("%#v", after), fmt.Sprintf("%#v", actual))
		}
		assert.Equal(t, `CREATE INDEX table_name_idx_column_name ON table_name (column_name);
`, actual.String())
	})

	t.Run("success,before,after,Table", func(t *testing.T) {
		t.Parallel()

		before, err := NewParser(NewLexer(`CREATE TABLE public.users (
    user_id UUID NOT NULL,
    username VARCHAR(256) NOT NULL,
    is_verified BOOL NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('UTC':::STRING, current_timestamp():::TIMESTAMPTZ),
    CONSTRAINT users_pkey PRIMARY KEY (user_id ASC),
    INDEX users_idx_by_username (username DESC)
);
`)).Parse()
		require.NoError(t, err)

		after, err := NewParser(NewLexer(`CREATE TABLE public.users (
    user_id UUID NOT NULL,
    username VARCHAR(256) NOT NULL,
    is_verified BOOL NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('UTC':::STRING, current_timestamp():::TIMESTAMPTZ),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('UTC':::STRING, current_timestamp():::TIMESTAMPTZ),
    CONSTRAINT users_pkey PRIMARY KEY (user_id ASC),
    INDEX users_idx_by_username (username DESC)
);
`)).Parse()
		require.NoError(t, err)

		expected := &DDL{
			Stmts: []Stmt{
				&AlterTableStmt{
					Name: &ObjectName{Schema: &Ident{Name: "public", Raw: "public"}, Name: &Ident{Name: "users", Raw: "users"}},
					Action: &AddColumn{
						Column: &Column{
							Name:     &Ident{Name: "updated_at", Raw: "updated_at"},
							DataType: &DataType{Name: "TIMESTAMPTZ", Size: ""},
							NotNull:  true,
							Default: &Default{
								Value: &DefaultValue{
									Idents: []*Ident{
										{Name: "timezone", Raw: "timezone"},
										{Name: "(", Raw: "("},
										{Name: "'UTC'", Raw: "'UTC'"},
										{Name: ":::", Raw: ":::"},
										{Name: "STRING", Raw: "STRING"},
										{Name: ",", Raw: ","},
										{Name: "current_timestamp", Raw: "current_timestamp"},
										{Name: "(", Raw: "("},
										{Name: ")", Raw: ")"},
										{Name: ":::", Raw: ":::"},
										{Name: "TIMESTAMPTZ", Raw: "TIMESTAMPTZ"},
										{Name: ")", Raw: ")"},
									},
								},
							},
						},
					},
				},
			},
		}
		actual, err := Diff(before, after)
		require.NoError(t, err)
		if !assert.Equal(t, expected, actual) {
			assert.Equal(t, fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", actual))
		}
		assert.Equal(t, `ALTER TABLE public.users ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('UTC':::STRING, current_timestamp():::TIMESTAMPTZ);
`, actual.String())
	})

	t.Run("success,before,after,Table,Asc", func(t *testing.T) {
		t.Parallel()

		before, err := NewParser(NewLexer(`CREATE TABLE public.users (
    user_id UUID NOT NULL,
    username VARCHAR(256) NOT NULL,
    is_verified BOOL NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('UTC':::STRING, current_timestamp():::TIMESTAMPTZ),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('UTC':::STRING, current_timestamp():::TIMESTAMPTZ),
    CONSTRAINT users_pkey PRIMARY KEY (user_id ASC),
    INDEX users_idx_by_username (username DESC)
);
`)).Parse()
		require.NoError(t, err)

		after, err := NewParser(NewLexer(`CREATE TABLE public.users (
    user_id UUID NOT NULL,
    username VARCHAR(256) NOT NULL,
    is_verified BOOL NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('UTC':::STRING, current_timestamp():::TIMESTAMPTZ),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('UTC':::STRING, current_timestamp():::TIMESTAMPTZ),
    CONSTRAINT users_pkey PRIMARY KEY (user_id ASC),
    INDEX users_idx_by_username (username ASC)
);
`)).Parse()
		require.NoError(t, err)

		expected := &DDL{
			Stmts: []Stmt{
				&DropIndexStmt{
					Name: &ObjectName{Schema: &Ident{Name: "public", Raw: "public"}, Name: &Ident{Name: "users_idx_by_username", Raw: "users_idx_by_username"}},
				},
				&CreateIndexStmt{
					Unique:    false,
					Name:      &ObjectName{Schema: &Ident{Name: "public", Raw: "public"}, Name: &Ident{Name: "users_idx_by_username", Raw: "users_idx_by_username"}},
					TableName: &ObjectName{Schema: &Ident{Name: "public", Raw: "public"}, Name: &Ident{Name: "users", Raw: "users"}},
					Columns: []*ColumnIdent{
						{
							Ident: &Ident{Name: "username", Raw: "username"},
							Order: &Order{Desc: false},
						},
					},
				},
			},
		}
		actual, err := Diff(before, after)
		require.NoError(t, err)
		if !assert.Equal(t, expected, actual) {
			assert.Equal(t, fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", actual))
		}
		assert.Equal(t, `DROP INDEX public.users_idx_by_username;
CREATE INDEX public.users_idx_by_username ON public.users (username ASC);
`, actual.String())
	})

	t.Run("success,before,after,Index", func(t *testing.T) {
		t.Parallel()

		before, err := NewParser(NewLexer(`CREATE UNIQUE INDEX IF NOT EXISTS public.users_idx_by_username ON public.users (username DESC);`)).Parse()
		require.NoError(t, err)

		after, err := NewParser(NewLexer(`CREATE UNIQUE INDEX IF NOT EXISTS public.users_idx_by_username ON public.users (username ASC, age ASC);`)).Parse()
		require.NoError(t, err)

		expected := &DDL{
			Stmts: []Stmt{
				&DropIndexStmt{
					Name: &ObjectName{Schema: &Ident{Name: "public", Raw: "public"}, Name: &Ident{Name: "users_idx_by_username", Raw: "users_idx_by_username"}},
				},
				&CreateIndexStmt{
					Unique:      true,
					IfNotExists: true,
					Name:        &ObjectName{Schema: &Ident{Name: "public", Raw: "public"}, Name: &Ident{Name: "users_idx_by_username", Raw: "users_idx_by_username"}},
					TableName:   &ObjectName{Schema: &Ident{Name: "public", Raw: "public"}, Name: &Ident{Name: "users", Raw: "users"}},
					Columns: []*ColumnIdent{
						{
							Ident: &Ident{Name: "username", Raw: "username"},
							Order: &Order{Desc: false},
						},
						{
							Ident: &Ident{Name: "age", Raw: "age"},
							Order: &Order{Desc: false},
						},
					},
				},
			},
		}
		actual, err := Diff(before, after)
		require.NoError(t, err)
		if !assert.Equal(t, expected, actual) {
			assert.Equal(t, fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", actual))
		}
		assert.Equal(t, `DROP INDEX public.users_idx_by_username;
CREATE UNIQUE INDEX IF NOT EXISTS public.users_idx_by_username ON public.users (username ASC, age ASC);
`, actual.String())
	})
}
