package sqlz //nolint:testpackage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/kunitsucom/util.go/pointer"
	slicez "github.com/kunitsucom/util.go/slices"
)

func Test_DB_QueryContext(t *testing.T) {
	t.Parallel()
	t.Run("failure,sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID     int     `testdb:"user_id"`
			Username   string  `testdb:"username"`
			NullString *string `testdb:"null_string"`
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
			UserID     int     `testdb:"user_id"`
			Username   string  `testdb:"username"`
			NullString *string `testdb:"null_string"`
		}
		var u []user
		db := newDB(&sqlDBMock{}, WithNewDBOptionStructTag("testdb"))
		i := 0
		columns := []string{"user_id", "username", "null_string"}
		rows := &sqlRowsMock{
			NextFunc:    func() bool { i++; return i < 51 },
			ColumnsFunc: func() ([]string, error) { return slicez.Copy(columns), nil },
			ScanFunc: func(dest ...interface{}) error {
				for dstIdx := range dest {
					for colIdx := range columns {
						if dstIdx == colIdx {
							switch columns[colIdx] {
							case "user_id":
								reflect.ValueOf(dest[dstIdx]).Elem().SetInt(int64(i))
							case "username":
								reflect.ValueOf(dest[dstIdx]).Elem().SetString(columns[colIdx] + "_" + fmt.Sprintf("%03d", i))
							case "null_string":
								reflect.ValueOf(dest[dstIdx]).Elem().Set(reflect.ValueOf(pointer.Pointer(columns[colIdx] + "_" + fmt.Sprintf("%03d", i))))
							}
						}
					}
				}
				return nil
			},
		}
		if err := db.queryContext(rows, nil, &u); err != nil {
			t.Fatalf("❌: queryContext: %v", err)
		}
		if expect, actual := 50, len(u); expect != actual {
			t.Errorf("❌: len(u): expect(%v) != actual(%v)", expect, actual)
		}
		if expect, actual := 1, u[0].UserID; expect != actual {
			t.Errorf("❌: u[0].UserID: expect(%v) != actual(%v)", expect, actual)
		}
		if expect, actual := 50, u[len(u)-1].UserID; expect != actual {
			t.Errorf("❌: u[0].UserID: expect(%v) != actual(%v)", expect, actual)
		}
		if len(u) > 0 {
			t.Logf("✅: ScanRows: u[0]: %#v", u[0])
			t.Logf("✅: ScanRows: u[len(u)-1]: %#v", u[len(u)-1])
		} else {
			t.Logf("✅: ScanRows: u: %#v", u)
		}
	})
}

func Test_DB_QueryRowContext(t *testing.T) {
	t.Parallel()
	t.Run("failure,sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID     int     `testdb:"user_id"`
			Username   string  `testdb:"username"`
			NullString *string `testdb:"null_string"`
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
			UserID     int     `testdb:"user_id"`
			Username   string  `testdb:"username"`
			NullString *string `testdb:"null_string"`
		}
		var u user
		db := newDB(&sqlDBMock{}, WithNewDBOptionStructTag("testdb"))
		i := 0
		columns := []string{"user_id", "username", "null_string"}
		rows := &sqlRowsMock{
			NextFunc:    func() bool { i++; return i < 2 },
			ColumnsFunc: func() ([]string, error) { return slicez.Copy(columns), nil },
			ScanFunc: func(dest ...interface{}) error {
				for dstIdx := range dest {
					for colIdx := range columns {
						if dstIdx == colIdx {
							switch columns[colIdx] {
							case "user_id":
								reflect.ValueOf(dest[dstIdx]).Elem().SetInt(int64(i))
							case "username":
								reflect.ValueOf(dest[dstIdx]).Elem().SetString(columns[colIdx] + "_" + fmt.Sprintf("%03d", i))
							case "null_string":
								reflect.ValueOf(dest[dstIdx]).Elem().Set(reflect.ValueOf(pointer.Pointer(columns[colIdx] + "_" + fmt.Sprintf("%03d", i))))
							}
						}
					}
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
			UserID     int     `testdb:"user_id"`
			Username   string  `testdb:"username"`
			NullString *string `testdb:"null_string"`
		}
		var u user
		db := newDB(&sqlDBMock{}, WithNewDBOptionStructTag("testdb"))
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
			UserID     int     `testdb:"user_id"`
			Username   string  `testdb:"username"`
			NullString *string `testdb:"null_string"`
		}
		var u user
		db := newDB(&sqlDBMock{}, WithNewDBOptionStructTag("testdb"))
		rows := &sqlRowsMock{
			NextFunc: func() bool { return false },
			ErrFunc:  func() error { return context.Canceled },
		}
		if err := db.queryRowContext(rows, nil, &u); !errors.Is(err, context.Canceled) {
			t.Fatalf("❌: queryRowContext: expect(%v) != actual(%v)", context.Canceled, err)
		}
	})
}
