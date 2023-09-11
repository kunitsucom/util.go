package sqlz

import (
	"errors"
	"testing"

	sqlz "github.com/kunitsucom/util.go/database/sql"
)

const defaultTestStructTag = "testdb"

func TestColumns(t *testing.T) {
	t.Parallel()

	t.Run("success,interface", func(t *testing.T) {
		t.Parallel()

		columns := Columns(&testUser{}, defaultTestStructTag)
		if len(columns) != 3 {
			t.Errorf("❌: len(columns): expect(%v) != actual(%v)", 3, len(columns))
		}
		t.Logf("✅: columns: %#v", columns)
		t.Logf("✅: len(columns): %d", len(columns))

		cachedColumns := Columns(&testUser{}, defaultTestStructTag)
		if len(cachedColumns) != 3 {
			t.Errorf("❌: len(cachedColumns): expect(%v) != actual(%v)", 3, len(cachedColumns))
		}
		t.Logf("✅: cachedColumns: %#v", cachedColumns)
		t.Logf("✅: len(cachedColumns): %d", len(cachedColumns))
	})

	t.Run("success,notInterface", func(t *testing.T) {
		t.Parallel()

		columns := Columns(testUser{}, defaultTestStructTag)
		if len(columns) != 3 {
			t.Errorf("❌: len(columns): expect(%v) != actual(%v)", 3, len(columns))
		}
		t.Logf("✅: columns: %#v", columns)
		t.Logf("✅: len(columns): %d", len(columns))

		cachedColumns := Columns(testUser{}, defaultTestStructTag)
		if len(cachedColumns) != 3 {
			t.Errorf("❌: len(cachedColumns): expect(%v) != actual(%v)", 3, len(cachedColumns))
		}
		t.Logf("✅: cachedColumns: %#v", cachedColumns)
		t.Logf("✅: len(cachedColumns): %d", len(cachedColumns))
	})
}

func TestTableName(t *testing.T) {
	t.Parallel()

	t.Run("success,interface", func(t *testing.T) {
		t.Parallel()

		tableName := TableName(&testUser{})
		if expect, actual := "test_user", tableName; expect != actual {
			t.Errorf("❌: tableName: expect(%v) != actual(%v)", expect, actual)
		}
		t.Logf("✅: tableName: %#v", tableName)

		cachedTableName := TableName(&testUser{})
		if expect, actual := "test_user", tableName; expect != actual {
			t.Errorf("❌: tableName: expect(%v) != actual(%v)", expect, actual)
		}
		t.Logf("✅: cachedTableName: %#v", cachedTableName)
	})

	t.Run("success,notInterface", func(t *testing.T) {
		t.Parallel()

		tableName := TableName(testUser{})
		if expect, actual := "testUser", tableName; expect != actual {
			t.Errorf("❌: tableName: expect(%v) != actual(%v)", expect, actual)
		}
		t.Logf("✅: tableName: %#v", tableName)

		cachedTableName := TableName(testUser{})
		if expect, actual := "testUser", tableName; expect != actual {
			t.Errorf("❌: tableName: expect(%v) != actual(%v)", expect, actual)
		}
		t.Logf("✅: cachedTableName: %#v", cachedTableName)
	})

	t.Run("failure,panic,notStruct", func(t *testing.T) {
		t.Parallel()

		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("❌: TableName: expect panic, but not")
				return
			}
			err, ok := r.(error)
			if !ok {
				t.Errorf("❌: TableName: expect error, but not")
				return
			}
			if expect, actual := sqlz.ErrDataTypeNotSupported, err; !errors.Is(actual, expect) {
				t.Errorf("❌: TableName: expect(%v) != actual(%v)", expect, actual)
				return
			}
		}()

		willBePanic := 1
		TableName(&willBePanic)
	})
}
