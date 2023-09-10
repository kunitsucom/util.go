package sqlz

import "testing"

//nolint:paralleltest
func TestResetGlobalColumnsCache(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		qb := NewQueryBuilder(WithNewQueryBuilderOptionStructTag("testdb"))
		qb.private()
		columns := qb.Columns(testUser{})
		if len(columns) != 3 {
			t.Errorf("❌: len(columns): expect(%v) != actual(%v)", 3, len(columns))
		}
		t.Logf("✅: columns: %#v", columns)
		t.Logf("✅: len(columns): %d", len(columns))

		ResetGlobalColumnsCache()

		notCachedColumns := qb.Columns(testUser{})
		if len(notCachedColumns) != 3 {
			t.Errorf("❌: len(cachedColumns): expect(%v) != actual(%v)", 3, len(notCachedColumns))
		}
		t.Logf("✅: cachedColumns: %#v", notCachedColumns)
		t.Logf("✅: len(cachedColumns): %d", len(notCachedColumns))
	})
}

func TestQueryBuilder_Columns(t *testing.T) {
	t.Parallel()

	t.Run("success,interface", func(t *testing.T) {
		t.Parallel()

		qb := NewQueryBuilder(WithNewQueryBuilderOptionStructTag("testdb"))
		columns := qb.Columns(&testUser{})
		if len(columns) != 3 {
			t.Errorf("❌: len(columns): expect(%v) != actual(%v)", 3, len(columns))
		}
		t.Logf("✅: columns: %#v", columns)
		t.Logf("✅: len(columns): %d", len(columns))

		cachedColumns := qb.Columns(&testUser{})
		if len(cachedColumns) != 3 {
			t.Errorf("❌: len(cachedColumns): expect(%v) != actual(%v)", 3, len(cachedColumns))
		}
		t.Logf("✅: cachedColumns: %#v", cachedColumns)
		t.Logf("✅: len(cachedColumns): %d", len(cachedColumns))
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		qb := NewQueryBuilder(WithNewQueryBuilderOptionStructTag("testdb"))
		columns := qb.Columns(testUser{})
		if len(columns) != 3 {
			t.Errorf("❌: len(columns): expect(%v) != actual(%v)", 3, len(columns))
		}
		t.Logf("✅: columns: %#v", columns)
		t.Logf("✅: len(columns): %d", len(columns))

		cachedColumns := qb.Columns(testUser{})
		if len(cachedColumns) != 3 {
			t.Errorf("❌: len(cachedColumns): expect(%v) != actual(%v)", 3, len(cachedColumns))
		}
		t.Logf("✅: cachedColumns: %#v", cachedColumns)
		t.Logf("✅: len(cachedColumns): %d", len(cachedColumns))
	})
}

func TestQueryBuilder_TableName(t *testing.T) {
	t.Parallel()

	t.Run("success,interface", func(t *testing.T) {
		t.Parallel()

		qb := NewQueryBuilder(WithNewQueryBuilderOptionStructTag("testdb"))
		tableName := qb.TableName(&testUser{})
		if expect, actual := "test_user", tableName; expect != actual {
			t.Errorf("❌: tableName: expect(%v) != actual(%v)", expect, actual)
		}
		t.Logf("✅: tableName: %#v", tableName)

		cachedTableName := qb.TableName(&testUser{})
		if expect, actual := "test_user", tableName; expect != actual {
			t.Errorf("❌: tableName: expect(%v) != actual(%v)", expect, actual)
		}
		t.Logf("✅: cachedTableName: %#v", cachedTableName)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		qb := NewQueryBuilder(WithNewQueryBuilderOptionStructTag("testdb"))
		tableName := qb.TableName(testUser{})
		if expect, actual := "testUser", tableName; expect != actual {
			t.Errorf("❌: tableName: expect(%v) != actual(%v)", expect, actual)
		}
		t.Logf("✅: tableName: %#v", tableName)

		cachedTableName := qb.TableName(testUser{})
		if expect, actual := "testUser", tableName; expect != actual {
			t.Errorf("❌: tableName: expect(%v) != actual(%v)", expect, actual)
		}
		t.Logf("✅: cachedTableName: %#v", cachedTableName)
	})
}
