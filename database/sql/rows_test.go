package sqlz //nolint:testpackage

import (
	"database/sql"
	"errors"
	"reflect"
	"strconv"
	"testing"
)

func Test_ScanRows(t *testing.T) {
	t.Parallel()
	t.Run("success,reflect.Slice", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID   string `db:"user_id"`
			Username string `db:"username"`
		}
		var u []user
		i := 0
		rows := &mockRows{
			NextFunc: func() bool {
				i++
				return i < 3
			},
			ColumnsReturn: []string{"user_id", "username"},
			ScanFunc: func(dest ...interface{}) error {
				for i := range dest {
					reflect.ValueOf(dest[i]).Elem().SetString("column" + strconv.Itoa(i))
				}
				return nil
			},
		}
		if err := ScanRows(rows, "db", &u); err != nil {
			t.Fatalf("❌: ScanRows: err != nil: %v", err)
		}
		t.Logf("✅: ScanRows: %+v", u)
	})
	t.Run("success,reflect.Slice_pointer_slice", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID   string `db:"user_id"`
			Username string `db:"username"`
		}
		var u []*user
		i := 0
		rows := &mockRows{
			NextFunc: func() bool {
				i++
				return i < 3
			},
			ColumnsReturn: []string{"user_id", "username"},
			ScanFunc: func(dest ...interface{}) error {
				for i := range dest {
					reflect.ValueOf(dest[i]).Elem().SetString("column" + strconv.Itoa(i))
				}
				return nil
			},
		}
		if err := ScanRows(rows, "db", &u); err != nil {
			t.Fatalf("❌: ScanRows: err != nil: %v", err)
		}
		t.Logf("✅: ScanRows: %+v", u)
	})
	t.Run("success,reflect.Struct", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID   string `db:"user_id"`
			Username string `db:"username"`
		}
		var u user
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
		if err := ScanRows(rows, "db", &u); err != nil {
			t.Fatalf("❌: ScanRows: err != nil: %v", err)
		}
		t.Logf("✅: ScanRows: %+v", u)
	})
	t.Run("failure,ErrMustBePointer", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID   string `db:"user_id"`
			Username string `db:"username"`
		}
		var notPointer user
		i := 0
		rows := &mockRows{
			NextFunc: func() bool {
				i++
				return i < 2
			},
		}
		if err := ScanRows(rows, "db", notPointer); !errors.Is(err, ErrMustBePointer) {
			t.Fatalf("❌: ScanRows: expect(%v) != actual(%v)", ErrMustBePointer, err)
		}
	})
	t.Run("failure,ErrMustNotNil", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID   string `db:"user_id"`
			Username string `db:"username"`
		}
		var nilPointer *user
		i := 0
		rows := &mockRows{
			NextFunc: func() bool {
				i++
				return i < 2
			},
		}
		if err := ScanRows(rows, "db", nilPointer); !errors.Is(err, ErrMustNotNil) {
			t.Fatalf("❌: ScanRows: expect(%v) != actual(%v)", ErrMustNotNil, err)
		}
	})
	t.Run("failure,reflect.Slice_Scan", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID   string `db:"user_id"`
			Username string `db:"username"`
		}
		var u []*user
		i := 0
		rows := &mockRows{
			NextFunc: func() bool {
				i++
				return i < 3
			},
			ColumnsReturn: []string{"user_id", "username"},
			ScanFunc: func(dest ...interface{}) error {
				return sql.ErrConnDone
			},
		}
		if err := ScanRows(rows, "db", &u); !errors.Is(err, sql.ErrConnDone) {
			t.Fatalf("❌: queryStructContext: expect(%v) != actual(%v)", sql.ErrConnDone, err)
		}
	})
	t.Run("failure,reflect.Slice_ErrDataTypeNotSupported", func(t *testing.T) {
		t.Parallel()
		var u []string
		if err := ScanRows(&mockRows{}, "db", &u); !errors.Is(err, ErrDataTypeNotSupported) {
			t.Fatalf("❌: queryStructContext: expect(%v) != actual(%v)", ErrDataTypeNotSupported, err)
		}
	})
	t.Run("failure,reflect.Struct_Scan", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID   string `db:"user_id"`
			Username string `db:"username"`
		}
		var u user
		i := 0
		rows := &mockRows{
			NextFunc: func() bool {
				i++
				return i < 2
			},
			ColumnsReturn: []string{"user_id", "username"},
			ScanFunc: func(dest ...interface{}) error {
				return sql.ErrConnDone
			},
		}
		if err := ScanRows(rows, "db", &u); !errors.Is(err, sql.ErrConnDone) {
			t.Fatalf("❌: queryStructContext: expect(%v) != actual(%v)", sql.ErrConnDone, err)
		}
	})
	t.Run("failure,reflect.Struct_Scan", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID   string `db:"user_id"`
			Username string `db:"username"`
		}
		var u user
		i := 0
		rows := &mockRows{
			NextFunc: func() bool {
				i++
				return i < 2
			},
			ColumnsError: sql.ErrNoRows,
		}
		if err := ScanRows(rows, "db", &u); !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("❌: queryStructContext: expect(%v) != actual(%v)", sql.ErrNoRows, err)
		}
	})
	t.Run("failure,ErrDataTypeNotSupported", func(t *testing.T) {
		t.Parallel()
		var user string
		i := 0
		rows := &mockRows{
			NextFunc: func() bool {
				i++
				return i < 2
			},
		}
		if err := ScanRows(rows, "db", &user); !errors.Is(err, ErrDataTypeNotSupported) {
			t.Fatalf("❌: ScanRows: expect(%v) != actual(%v)", ErrDataTypeNotSupported, err)
		}
	})
}