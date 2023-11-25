package cockroachdb

import (
	"testing"

	"github.com/kunitsucom/util.go/testing/require"
)

func Test_isConstraint(t *testing.T) {
	t.Parallel()

	(&PrimaryKeyConstraint{}).isConstraint()
	(&ForeignKeyConstraint{}).isConstraint()
	(&UniqueConstraint{}).isConstraint()
	(&CheckConstraint{}).isConstraint()
}

func TestPrimaryKeyConstraint(t *testing.T) {
	t.Parallel()

	t.Run("success,PrimaryKeyConstraint", func(t *testing.T) {
		t.Parallel()

		primaryKeyConstraint := &PrimaryKeyConstraint{Name: &Ident{Name: "pk_users", QuotationMark: `"`, Raw: `"pk_users"`}, Columns: []*ColumnIdent{{Ident: &Ident{Name: "id", QuotationMark: `"`, Raw: `"id"`}}}}
		expected := "CONSTRAINT \"pk_users\" PRIMARY KEY (\"id\")"
		actual := primaryKeyConstraint.String()
		require.Equal(t, expected, actual)

		t.Logf("✅: primaryKeyConstraint: %#v", primaryKeyConstraint)
	})
	t.Run("success,PrimaryKeyConstraint,empty", func(t *testing.T) {
		t.Parallel()

		primaryKeyConstraint := &PrimaryKeyConstraint{}
		expected := "PRIMARY KEY ()"
		actual := primaryKeyConstraint.String()
		require.Equal(t, expected, actual)

		t.Logf("✅: primaryKeyConstraint: %#v", primaryKeyConstraint)
	})
}

func TestForeignKeyConstraint(t *testing.T) {
	t.Parallel()
	t.Run("success,ForeignKeyConstraint", func(t *testing.T) {
		t.Parallel()

		foreignKeyConstraint := &ForeignKeyConstraint{
			Name:       &Ident{Name: "fk_users_groups", QuotationMark: `"`, Raw: `"fk_users_groups"`},
			Columns:    []*ColumnIdent{{Ident: &Ident{Name: "group_id", QuotationMark: `"`, Raw: `"group_id"`}}},
			Ref:        &Ident{Name: "groups", QuotationMark: `"`, Raw: `"groups"`},
			RefColumns: []*ColumnIdent{{Ident: &Ident{Name: "id", QuotationMark: `"`, Raw: `"id"`}}},
		}

		expected := `CONSTRAINT "fk_users_groups" FOREIGN KEY ("group_id") REFERENCES "groups" ("id")`
		actual := foreignKeyConstraint.String()
		require.Equal(t, expected, actual)

		t.Logf("✅: foreignKeyConstraint: %#v", foreignKeyConstraint)
	})
}

func TestUniqueConstraint(t *testing.T) {
	t.Parallel()
	t.Run("success,UniqueConstraint", func(t *testing.T) {
		t.Parallel()

		uniqueConstraint := &UniqueConstraint{
			Name:    &Ident{Name: "uq_users_email", QuotationMark: `"`, Raw: `"uq_users_email"`},
			Columns: []*ColumnIdent{{Ident: &Ident{Name: "email", QuotationMark: `"`, Raw: `"email"`}}},
		}

		expected := `CONSTRAINT "uq_users_email" UNIQUE ("email")`
		actual := uniqueConstraint.String()
		require.Equal(t, expected, actual)

		t.Logf("✅: uniqueConstraint: %#v", uniqueConstraint)
	})
}

func TestCheckConstraint(t *testing.T) {
	t.Parallel()
	t.Run("success,CheckConstraint", func(t *testing.T) {
		t.Parallel()

		checkConstraint := &CheckConstraint{
			Name: &Ident{Name: "users_check_age", QuotationMark: `"`, Raw: `"users_check_age"`},
			Expr: []*Ident{{Name: "age", QuotationMark: `"`, Raw: `"age"`}},
		}

		expected := `CONSTRAINT "users_check_age" CHECK ("age")`
		actual := checkConstraint.String()
		require.Equal(t, expected, actual)

		t.Logf("✅: checkConstraint: %#v", checkConstraint)
	})
}

func TestDefault_String(t *testing.T) {
	t.Parallel()

	t.Run("success,String,nil", func(t *testing.T) {
		t.Parallel()

		d := (*Default)(nil)
		expected := ""
		actual := d.String()
		require.Equal(t, expected, actual)
	})
	t.Run("success,String,nilnil", func(t *testing.T) {
		t.Parallel()

		d := &Default{}
		expected := ""
		actual := d.String()
		require.Equal(t, expected, actual)
	})
	t.Run("success,PlainString,nilnil", func(t *testing.T) {
		t.Parallel()

		d := &Default{}
		expected := ""
		actual := d.PlainString()
		require.Equal(t, expected, actual)
	})
	t.Run("success,DEFAULT_VALUE", func(t *testing.T) {
		t.Parallel()

		d := &Default{Value: &DefaultValue{[]*Ident{{Name: "now()", Raw: "now()"}}}}
		expected := "DEFAULT now()"
		actual := d.String()
		require.Equal(t, expected, actual)

		t.Logf("✅: d: %#v", d)
	})
	t.Run("success,DEFAULT_VALUE,empty", func(t *testing.T) {
		t.Parallel()

		d := (*DefaultValue)(nil)
		expected := ""
		actual := d.String()
		require.Equal(t, expected, actual)
	})
	t.Run("success,DEFAULT_EXPR", func(t *testing.T) {
		t.Parallel()

		d := &Default{Value: &DefaultValue{[]*Ident{{Name: "(", Raw: "("}, {Name: "age", Raw: "age"}, {Name: ">=", Raw: ">="}, {Name: "0", Raw: "0"}, {Name: ")", Raw: ")"}}}}
		expected := "DEFAULT (age >= 0)"
		actual := d.String()
		require.Equal(t, expected, actual)

		t.Logf("✅: d: %#v", d)
	})
}

func TestColumn(t *testing.T) {
	t.Parallel()

	t.Run("success,Column", func(t *testing.T) {
		t.Parallel()

		column := &Column{
			Name:     &Ident{Name: "id", QuotationMark: `"`, Raw: `"id"`},
			DataType: &DataType{Name: "INTEGER"},
		}

		expected := `"id" INTEGER`
		actual := column.String()
		require.Equal(t, expected, actual)

		t.Logf("✅: column: %#v", column)
	})
}

func TestOption(t *testing.T) {
	t.Parallel()

	t.Run("success,Option", func(t *testing.T) {
		t.Parallel()

		option := &Option{Name: "TABLESPACE", Value: &Ident{Name: "pg_default", QuotationMark: `"`, Raw: `"pg_default"`}}

		expected := `TABLESPACE "pg_default"`
		actual := option.String()
		require.Equal(t, expected, actual)

		t.Logf("✅: option: %#v", option)
	})

	t.Run("success,Option,empty", func(t *testing.T) {
		t.Parallel()

		option := &Option{}
		expected := ""
		actual := option.String()
		require.Equal(t, expected, actual)

		t.Logf("✅: option: %#v", option)
	})
}
