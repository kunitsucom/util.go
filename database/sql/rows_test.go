package sqlz //nolint:testpackage

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/kunitsucom/util.go/pointer"
	slicez "github.com/kunitsucom/util.go/slices"
)

func Test_ScanRows(t *testing.T) {
	t.Parallel()
	t.Run("success,reflect.Slice", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID     int     `testdb:"user_id"`
			Username   string  `testdb:"username"`
			NullString *string `testdb:"null_string"`
		}
		var u []user
		i := 0
		columns := []string{"user_id", "username", "null_string"}
		rows := &sqlRowsMock{
			NextFunc:    func() bool { i++; return i < 101 },
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
		if err := ScanRows(rows, "testdb", &u); err != nil {
			t.Fatalf("❌: ScanRows: err != nil: %v", err)
		}
		if expect, actual := 100, len(u); expect != actual {
			t.Errorf("❌: len(u): expect(%v) != actual(%v)", expect, actual)
		}
		if expect, actual := 1, u[0].UserID; expect != actual {
			t.Errorf("❌: u[0].UserID: expect(%v) != actual(%v)", expect, actual)
		}
		if expect, actual := 100, u[len(u)-1].UserID; expect != actual {
			t.Errorf("❌: u[0].UserID: expect(%v) != actual(%v)", expect, actual)
		}
		if len(u) > 0 {
			t.Logf("✅: ScanRows: u[0]: %#v", u[0])
			t.Logf("✅: ScanRows: u[len(u)-1]: %#v", u[len(u)-1])
		} else {
			t.Logf("✅: ScanRows: u: %#v", u)
		}
	})
	t.Run("success,reflect.Slice_pointer_slice", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID     int     `testdb:"user_id"`
			Username   string  `testdb:"username"`
			NullString *string `testdb:"null_string"`
		}
		var u []*user
		i := 0
		columns := []string{"user_id", "username", "null_string"}
		rows := &sqlRowsMock{
			NextFunc:    func() bool { i++; return i < 3 },
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
		if err := ScanRows(rows, "testdb", &u); err != nil {
			t.Fatalf("❌: ScanRows: err != nil: %v", err)
		}
		t.Logf("✅: ScanRows: %+v", u)
	})
	t.Run("success,reflect.Struct", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID     int     `testdb:"user_id"`
			Username   string  `testdb:"username"`
			NullString *string `testdb:"null_string"`
		}
		var u user
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
		if err := ScanRows(rows, "testdb", &u); err != nil {
			t.Fatalf("❌: ScanRows: err != nil: %v", err)
		}
		t.Logf("✅: ScanRows: %+v", u)
	})
	t.Run("failure,ErrMustBePointer", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID     int     `testdb:"user_id"`
			Username   string  `testdb:"username"`
			NullString *string `testdb:"null_string"`
		}
		var notPointer user
		i := 0
		rows := &sqlRowsMock{
			NextFunc: func() bool { i++; return i < 2 },
		}
		if err := ScanRows(rows, "testdb", notPointer); !errors.Is(err, ErrMustBePointer) {
			t.Fatalf("❌: ScanRows: expect(%v) != actual(%v)", ErrMustBePointer, err)
		}
	})
	t.Run("failure,ErrMustNotNil", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID     int     `testdb:"user_id"`
			Username   string  `testdb:"username"`
			NullString *string `testdb:"null_string"`
		}
		var nilPointer *user
		i := 0
		rows := &sqlRowsMock{
			NextFunc: func() bool { i++; return i < 2 },
		}
		if err := ScanRows(rows, "testdb", nilPointer); !errors.Is(err, ErrMustNotNil) {
			t.Fatalf("❌: ScanRows: expect(%v) != actual(%v)", ErrMustNotNil, err)
		}
	})
	t.Run("failure,reflect.Slice,Scan", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID     int     `testdb:"user_id"`
			Username   string  `testdb:"username"`
			NullString *string `testdb:"null_string"`
		}
		var u []*user
		i := 0
		columns := []string{"user_id", "username", "null_string"}
		rows := &sqlRowsMock{
			NextFunc:    func() bool { i++; return i < 3 },
			ColumnsFunc: func() ([]string, error) { return slicez.Copy(columns), nil },
			ScanFunc: func(dest ...interface{}) error {
				return sql.ErrConnDone
			},
		}
		if err := ScanRows(rows, "testdb", &u); !errors.Is(err, sql.ErrConnDone) {
			t.Fatalf("❌: queryStructContext: expect(%v) != actual(%v)", sql.ErrConnDone, err)
		}
		i = 0 // for getStructTags Load
		if err := ScanRows(rows, "testdb", &u); !errors.Is(err, sql.ErrConnDone) {
			t.Fatalf("❌: queryStructContext: expect(%v) != actual(%v)", sql.ErrConnDone, err)
		}
	})
	t.Run("failure,reflect.Slice,Columns", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID     int     `testdb:"user_id"`
			Username   string  `testdb:"username"`
			NullString *string `testdb:"null_string"`
		}
		var u []*user
		i := 0
		rows := &sqlRowsMock{
			NextFunc:    func() bool { i++; return i < 2 },
			ColumnsFunc: func() ([]string, error) { return nil, sql.ErrNoRows },
		}
		if err := ScanRows(rows, "testdb", &u); !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("❌: queryStructContext: expect(%v) != actual(%v)", sql.ErrNoRows, err)
		}
	})
	t.Run("failure,reflect.Slice_ErrDataTypeNotSupported", func(t *testing.T) {
		t.Parallel()
		var u []string
		if err := ScanRows(&sqlRowsMock{
			NextFunc: func() bool { return true },
		}, "testdb", &u); !errors.Is(err, ErrDataTypeNotSupported) {
			t.Fatalf("❌: queryStructContext: expect(%v) != actual(%v)", ErrDataTypeNotSupported, err)
		}
	})
	t.Run("failure,reflect.Struct_Scan", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID     int     `testdb:"user_id"`
			Username   string  `testdb:"username"`
			NullString *string `testdb:"null_string"`
		}
		var u user
		i := 0
		columns := []string{"user_id", "username", "null_string"}
		rows := &sqlRowsMock{
			NextFunc:    func() bool { i++; return i < 2 },
			ColumnsFunc: func() ([]string, error) { return slicez.Copy(columns), nil },
			ScanFunc: func(dest ...interface{}) error {
				return sql.ErrConnDone
			},
		}
		if err := ScanRows(rows, "testdb", &u); !errors.Is(err, sql.ErrConnDone) {
			t.Fatalf("❌: queryStructContext: expect(%v) != actual(%v)", sql.ErrConnDone, err)
		}
	})
	t.Run("failure,reflect.Struct_Scan", func(t *testing.T) {
		t.Parallel()
		type user struct {
			UserID     int     `testdb:"user_id"`
			Username   string  `testdb:"username"`
			NullString *string `testdb:"null_string"`
		}
		var u user
		i := 0
		rows := &sqlRowsMock{
			NextFunc:    func() bool { i++; return i < 2 },
			ColumnsFunc: func() ([]string, error) { return nil, sql.ErrNoRows },
		}
		if err := ScanRows(rows, "testdb", &u); !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("❌: queryStructContext: expect(%v) != actual(%v)", sql.ErrNoRows, err)
		}
	})
	t.Run("failure,ErrDataTypeNotSupported", func(t *testing.T) {
		t.Parallel()
		var user string
		i := 0
		rows := &sqlRowsMock{
			NextFunc: func() bool { i++; return i < 2 },
		}
		if err := ScanRows(rows, "testdb", &user); !errors.Is(err, ErrDataTypeNotSupported) {
			t.Fatalf("❌: ScanRows: expect(%v) != actual(%v)", ErrDataTypeNotSupported, err)
		}
	})
}
