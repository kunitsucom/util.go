package sqlz //nolint:testpackage

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"strconv"
	"testing"
)

type mockDB struct {
	SQLQueryer
	SQLTxBeginner

	Rows  *sql.Rows
	Error error

	BeginTxFunc func(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

func (m *mockDB) QueryContext(_ context.Context, _ string, _ ...interface{}) (*sql.Rows, error) {
	return m.Rows, m.Error
}

func (m *mockDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return m.BeginTxFunc(ctx, opts)
}

type mockRows struct {
	SQLRows
	CloseError    error
	ColumnsReturn []string
	ColumnsError  error
	NextFunc      func() bool
	ScanFunc      func(dest ...interface{}) error
}

func (m *mockRows) Close() error {
	return m.CloseError
}

func (m *mockRows) Columns() ([]string, error) {
	return m.ColumnsReturn, m.ColumnsError
}

func (m *mockRows) Next() bool {
	return m.NextFunc()
}

func (m *mockRows) Scan(dest ...interface{}) error {
	return m.ScanFunc(dest...)
}

func Test_DB_QueryStructSliceContext(t *testing.T) {
	t.Parallel()
	t.Run("failure,sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID   string `db:"user_id"`
			Username string `db:"username"`
		}
		var u []*user
		if err := NewDB(&mockDB{Rows: nil, Error: sql.ErrNoRows}).QueryStructSliceContext(context.Background(), "db", &u, "SELECT * FROM users"); !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("❌: QueryStructSliceContext: %v", err)
		}
	})
}

func Test_DB_queryStructSliceContext(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID   string `db:"user_id"`
			Username string `db:"username"`
		}
		var u user
		db := &_Queryer{}
		i := 0
		rows := &mockRows{
			NextFunc: func() bool {
				i++
				return i < 2
			},
			ColumnsReturn: []string{"user_id", "username"},
			ScanFunc: func(dest ...interface{}) error {
				for i := range dest {
					reflect.ValueOf(dest[i]).Elem().SetString("column" + strconv.Itoa(i))
				}
				return nil
			},
		}
		if err := db.queryStructSliceContext(rows, nil, "db", &u); err != nil {
			t.Fatalf("❌: queryStructSliceContext: %v", err)
		}
		t.Logf("✅: queryStructSliceContext: %+v", u)
	})
}

func Test_DB_QueryStructContext(t *testing.T) {
	t.Parallel()
	t.Run("failure,sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID   string `db:"user_id"`
			Username string `db:"username"`
		}
		var u user
		if err := NewDB(&mockDB{Rows: nil, Error: sql.ErrNoRows}).QueryStructContext(context.Background(), "db", &u, "SELECT * FROM users"); !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("❌: QueryStructContext: %v", err)
		}
	})
}

func Test_DB_queryStructContext(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID   string `db:"user_id"`
			Username string `db:"username"`
		}
		var u user
		db := &_Queryer{}
		i := 0
		rows := &mockRows{
			NextFunc: func() bool {
				i++
				return i < 2
			},
			ColumnsReturn: []string{"user_id", "username"},
			ScanFunc: func(dest ...interface{}) error {
				for i := range dest {
					reflect.ValueOf(dest[i]).Elem().SetString("column" + strconv.Itoa(i))
				}
				return nil
			},
		}
		if err := db.queryStructContext(rows, nil, "db", &u); err != nil {
			t.Fatalf("❌: queryStructContext: err != nil: %v", err)
		}
		t.Logf("✅: queryStructSliceContext: %+v", u)
	})
	t.Run("failure,sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID   string `db:"user_id"`
			Username string `db:"username"`
		}
		var u user
		db := &_Queryer{}
		rows := &mockRows{
			NextFunc: func() bool { return false },
		}
		if err := db.queryStructContext(rows, nil, "db", &u); !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("❌: queryStructContext: expect(%v) != actual(%v)", sql.ErrNoRows, err)
		}
	})
}
