package sqlz //nolint:testpackage

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"strconv"
	"testing"
)

func Test_DB_QueryContext(t *testing.T) {
	t.Parallel()
	t.Run("failure,sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID   string `db:"user_id"`
			Username string `db:"username"`
		}
		var u []*user
		if err := NewDB(&sqlDBMock{Rows: nil, Error: sql.ErrNoRows}).QueryContext(context.Background(), &u, "SELECT * FROM users"); !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("❌: QueryContext: %v", err)
		}
	})
}

func Test_DB_queryContext(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID   string `db:"user_id"`
			Username string `db:"username"`
		}
		var u user
		db := newDB(&sqlDBMock{}, WithNewDBOptionStructTag("db"))
		i := 0
		rows := &sqlRowsMock{
			NextFunc: func() bool {
				i++
				return i < 2
			},

			ColumnsFunc: func() ([]string, error) { return []string{"user_id", "username"}, nil },
			ScanFunc: func(dest ...interface{}) error {
				for i := range dest {
					reflect.ValueOf(dest[i]).Elem().SetString("column" + strconv.Itoa(i))
				}
				return nil
			},
		}
		if err := db.queryContext(rows, nil, &u); err != nil {
			t.Fatalf("❌: queryContext: %v", err)
		}
		t.Logf("✅: queryContext: %+v", u)
	})
}

func Test_DB_QueryRowContext(t *testing.T) {
	t.Parallel()
	t.Run("failure,sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID   string `db:"user_id"`
			Username string `db:"username"`
		}
		var u user
		if err := NewDB(&sqlDBMock{Rows: nil, Error: sql.ErrNoRows}).QueryRowContext(context.Background(), &u, "SELECT * FROM users"); !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("❌: QueryRowContext: %v", err)
		}
	})
}

func Test_DB_queryRowContext(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID   string `db:"user_id"`
			Username string `db:"username"`
		}
		var u user
		db := newDB(&sqlDBMock{}, WithNewDBOptionStructTag("db"))
		i := 0
		rows := &sqlRowsMock{
			NextFunc: func() bool {
				i++
				return i < 2
			},
			ColumnsFunc: func() ([]string, error) { return []string{"user_id", "username"}, nil },
			ScanFunc: func(dest ...interface{}) error {
				for i := range dest {
					reflect.ValueOf(dest[i]).Elem().SetString("column" + strconv.Itoa(i))
				}
				return nil
			},
		}
		if err := db.queryRowContext(rows, nil, &u); err != nil {
			t.Fatalf("❌: queryRowContext: err != nil: %v", err)
		}
		t.Logf("✅: queryContext: %+v", u)
	})
	t.Run("failure,sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID   string `db:"user_id"`
			Username string `db:"username"`
		}
		var u user
		db := newDB(&sqlDBMock{}, WithNewDBOptionStructTag("db"))
		rows := &sqlRowsMock{
			NextFunc: func() bool { return false },
			ErrFunc:  func() error { return nil },
		}
		if err := db.queryRowContext(rows, nil, &u); !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("❌: queryRowContext: expect(%v) != actual(%v)", sql.ErrNoRows, err)
		}
	})
	t.Run("failure,context.Canceled", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID   string `db:"user_id"`
			Username string `db:"username"`
		}
		var u user
		db := newDB(&sqlDBMock{}, WithNewDBOptionStructTag("db"))
		rows := &sqlRowsMock{
			NextFunc: func() bool { return false },
			ErrFunc:  func() error { return context.Canceled },
		}
		if err := db.queryRowContext(rows, nil, &u); !errors.Is(err, context.Canceled) {
			t.Fatalf("❌: queryRowContext: expect(%v) != actual(%v)", context.Canceled, err)
		}
	})
}
