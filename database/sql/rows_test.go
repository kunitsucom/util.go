package sqlz //nolint:testpackage

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"testing"

	genericz "github.com/kunitsucom/util.go/generics"
	slicez "github.com/kunitsucom/util.go/slices"
)

func Test_ScanRows(t *testing.T) {
	t.Parallel()
	t.Run("success,intSlice", func(t *testing.T) {
		t.Parallel()
		var userIDs []int
		i := 0
		rows := &sqlRowsMock{
			NextFunc: func() bool { i++; return i < 101 },
			ScanFunc: func(dest ...interface{}) error {
				reflect.ValueOf(dest[0]).Elem().SetInt(int64(i))
				return nil
			},
		}
		if err := ScanRows(rows, "testdb", &userIDs); err != nil {
			t.Fatalf("❌: ScanRows: err != nil: %v", err)
		}
		if expect, actual := 100, len(userIDs); expect != actual {
			t.Errorf("❌: len(userIDs): expect(%v) != actual(%v)", expect, actual)
		}
		if expect, actual := 1, userIDs[0]; expect != actual {
			t.Errorf("❌: userIDs[0]: expect(%v) != actual(%v)", expect, actual)
		}
		if expect, actual := 100, userIDs[len(userIDs)-1]; expect != actual {
			t.Errorf("❌: userIDs[0]: expect(%v) != actual(%v)", expect, actual)
		}
		if len(userIDs) > 0 {
			t.Logf("✅: ScanRows: userIDs[0]: %v", userIDs[0])
			t.Logf("✅: ScanRows: userIDs[len(userIDs)-1]: %v", userIDs[len(userIDs)-1])
		} else {
			t.Logf("✅: ScanRows: userIDs: %v", userIDs)
		}
	})
	t.Run("success,intSlice", func(t *testing.T) {
		t.Parallel()
		var userIDs []*int
		i := 0
		rows := &sqlRowsMock{
			NextFunc: func() bool { i++; return i < 101 },
			ScanFunc: func(dest ...interface{}) error {
				reflect.ValueOf(dest[0]).Elem().SetInt(int64(i))
				return nil
			},
		}
		if err := ScanRows(rows, "testdb", &userIDs); err != nil {
			t.Fatalf("❌: ScanRows: err != nil: %v", err)
		}
		if expect, actual := 100, len(userIDs); expect != actual {
			t.Errorf("❌: len(userIDs): expect(%v) != actual(%v)", expect, actual)
		}
		if expect, actual := 1, userIDs[0]; expect != *actual {
			t.Errorf("❌: userIDs[0]: expect(%v) != actual(%v)", expect, actual)
		}
		if expect, actual := 100, userIDs[len(userIDs)-1]; expect != *actual {
			t.Errorf("❌: userIDs[0]: expect(%v) != actual(%v)", expect, actual)
		}
		if len(userIDs) > 0 {
			t.Logf("✅: ScanRows: userIDs[0]: %v", userIDs[0])
			t.Logf("✅: ScanRows: userIDs[len(userIDs)-1]: %v", userIDs[len(userIDs)-1])
		} else {
			t.Logf("✅: ScanRows: userIDs: %v", userIDs)
		}
	})
	t.Run("success,structSlice", func(t *testing.T) {
		t.Parallel()
		var u []testUser
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
								reflect.ValueOf(dest[dstIdx]).Elem().Set(reflect.ValueOf(genericz.Pointer(columns[colIdx] + "_" + fmt.Sprintf("%03d", i))))
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
	t.Run("success,pointerSlice", func(t *testing.T) {
		t.Parallel()
		var u []*testUser
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
								reflect.ValueOf(dest[dstIdx]).Elem().Set(reflect.ValueOf(genericz.Pointer(columns[colIdx] + "_" + fmt.Sprintf("%03d", i))))
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
	t.Run("success,int", func(t *testing.T) {
		t.Parallel()
		var userID int
		i := 0
		rows := &sqlRowsMock{
			NextFunc: func() bool { i++; return i < 2 },
			ScanFunc: func(dest ...interface{}) error {
				reflect.ValueOf(dest[0]).Elem().SetInt(int64(i))
				return nil
			},
		}
		if err := ScanRows(rows, "testdb", &userID); err != nil {
			t.Fatalf("❌: ScanRows: err != nil: %v", err)
		}
		if expect, actual := 1, userID; expect != actual {
			t.Errorf("❌: userID: expect(%v) != actual(%v)", expect, actual)
		}
	})
	t.Run("success,struct", func(t *testing.T) {
		t.Parallel()
		var u testUser
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
								reflect.ValueOf(dest[dstIdx]).Elem().Set(reflect.ValueOf(genericz.Pointer(columns[colIdx] + "_" + fmt.Sprintf("%03d", i))))
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
		if expect, actual := 1, u.UserID; expect != actual {
			t.Errorf("❌: u.UserID: expect(%v) != actual(%v)", expect, actual)
		}
		t.Logf("✅: ScanRows: %+v", u)
	})
	t.Run("failure,ErrMustBePointer", func(t *testing.T) {
		t.Parallel()
		var notPointer testUser
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
		var nilPointer *testUser
		i := 0
		rows := &sqlRowsMock{
			NextFunc: func() bool { i++; return i < 2 },
		}
		if err := ScanRows(rows, "testdb", nilPointer); !errors.Is(err, ErrMustNotNil) {
			t.Fatalf("❌: ScanRows: expect(%v) != actual(%v)", ErrMustNotNil, err)
		}
	})
	t.Run("failure,pointerStructSlice,Scan", func(t *testing.T) {
		t.Parallel()
		var u []*testUser
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
	t.Run("failure,pointerStructSlice,Columns", func(t *testing.T) {
		t.Parallel()
		var u []*testUser
		i := 0
		rows := &sqlRowsMock{
			NextFunc:    func() bool { i++; return i < 2 },
			ColumnsFunc: func() ([]string, error) { return nil, sql.ErrNoRows },
		}
		if err := ScanRows(rows, "testdb", &u); !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("❌: queryStructContext: expect(%v) != actual(%v)", sql.ErrNoRows, err)
		}
	})
	t.Run("failure,intSlice", func(t *testing.T) {
		t.Parallel()
		var userIDs []int
		i := 0
		rows := &sqlRowsMock{
			NextFunc: func() bool { i++; return i < 101 },
			ScanFunc: func(dest ...interface{}) error {
				return driver.ErrBadConn
			},
		}
		if err := ScanRows(rows, "testdb", &userIDs); !errors.Is(err, driver.ErrBadConn) {
			t.Fatalf("❌: ScanRows: expect(%v) != actual(%v)", driver.ErrBadConn, err)
		}
	})
	t.Run("failure,struct,Scan", func(t *testing.T) {
		t.Parallel()
		var u testUser
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
	t.Run("failure,struct,Columns", func(t *testing.T) {
		t.Parallel()
		var u testUser
		i := 0
		rows := &sqlRowsMock{
			NextFunc:    func() bool { i++; return i < 2 },
			ColumnsFunc: func() ([]string, error) { return nil, sql.ErrNoRows },
		}
		if err := ScanRows(rows, "testdb", &u); !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("❌: queryStructContext: expect(%v) != actual(%v)", sql.ErrNoRows, err)
		}
	})
	t.Run("failure,int,rows.Err", func(t *testing.T) {
		t.Parallel()
		var userID int
		if err := ScanRows(&sqlRowsMock{
			NextFunc: func() bool { return false },
			ErrFunc:  func() error { return driver.ErrBadConn },
		}, "testdb", &userID); !errors.Is(err, driver.ErrBadConn) {
			t.Fatalf("❌: queryStructContext: expect(%v) != actual(%v)", driver.ErrBadConn, err)
		}
	})
	t.Run("failure,int,rows.Err", func(t *testing.T) {
		t.Parallel()
		var userID int
		if err := ScanRows(&sqlRowsMock{
			NextFunc: func() bool { return false },
			ErrFunc:  func() error { return driver.ErrBadConn },
		}, "testdb", &userID); !errors.Is(err, driver.ErrBadConn) {
			t.Fatalf("❌: queryStructContext: expect(%v) != actual(%v)", driver.ErrBadConn, err)
		}
	})
	t.Run("failure,int,sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()
		var userID int
		if err := ScanRows(&sqlRowsMock{
			NextFunc: func() bool { return false },
		}, "testdb", &userID); !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("❌: queryStructContext: expect(%v) != actual(%v)", sql.ErrNoRows, err)
		}
	})
	t.Run("failure,int,Scan", func(t *testing.T) {
		t.Parallel()
		var userID int
		if err := ScanRows(&sqlRowsMock{
			NextFunc: func() bool { return true },
			ScanFunc: func(dest ...interface{}) error {
				return sql.ErrNoRows
			},
		}, "testdb", &userID); !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("❌: queryStructContext: expect(%v) != actual(%v)", sql.ErrNoRows, err)
		}
	})
	t.Run("failure,array,ErrDataTypeNotSupported", func(t *testing.T) {
		t.Parallel()
		var u [16]string
		if err := ScanRows(&sqlRowsMock{
			NextFunc: func() bool { return true },
		}, "testdb", &u); !errors.Is(err, ErrDataTypeNotSupported) {
			t.Fatalf("❌: queryStructContext: expect(%v) != actual(%v)", ErrDataTypeNotSupported, err)
		}
	})
	t.Run("failure,slice,ErrDataTypeNotSupported", func(t *testing.T) {
		t.Parallel()
		var user [][]string
		i := 0
		rows := &sqlRowsMock{
			NextFunc: func() bool { i++; return i < 2 },
		}
		if err := ScanRows(rows, "testdb", &user); !errors.Is(err, ErrDataTypeNotSupported) {
			t.Fatalf("❌: ScanRows: expect(%v) != actual(%v)", ErrDataTypeNotSupported, err)
		}
	})
}
