package sql_test

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"testing"

	sqlz "github.com/kunitsucom/util.go/database/sql"
	errorz "github.com/kunitsucom/util.go/errors"
	"github.com/kunitsucom/util.go/tests/database/sql/mysql"
)

const (
	MySQLCreateTableTestUser = `
CREATE TABLE IF NOT EXISTS test_user (
    id          INT          NOT NULL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
	profile     VARCHAR(255) NOT NULL,
	null_string VARCHAR(255) NULL
)
`
	MySQLInsertTestUser = `
INSERT IGNORE INTO test_user
	(id, name, profile, null_string)
VALUES
	(1,   'test_user_001', 'hello', NULL),
	(2,   'test_user_002', 'hello', NULL),
	(3,   'test_user_003', 'hello', NULL),
	(4,   'test_user_004', 'hello', NULL),
	(5,   'test_user_005', 'hello', NULL),
	(6,   'test_user_006', 'hello', NULL),
	(7,   'test_user_007', 'hello', NULL),
	(8,   'test_user_008', 'hello', NULL),
	(9,   'test_user_009', 'hello', NULL),
	(10,  'test_user_010', 'hello', NULL),
	(11,  'test_user_011', 'hello', NULL),
	(12,  'test_user_012', 'hello', NULL),
	(13,  'test_user_013', 'hello', NULL),
	(14,  'test_user_014', 'hello', NULL),
	(15,  'test_user_015', 'hello', NULL),
	(16,  'test_user_016', 'hello', NULL),
	(17,  'test_user_017', 'hello', NULL),
	(18,  'test_user_018', 'hello', NULL),
	(19,  'test_user_019', 'hello', NULL),
	(20,  'test_user_020', 'hello', NULL),
	(21,  'test_user_021', 'hello', NULL),
	(22,  'test_user_022', 'hello', NULL),
	(23,  'test_user_023', 'hello', NULL),
	(24,  'test_user_024', 'hello', NULL),
	(25,  'test_user_025', 'hello', NULL),
	(26,  'test_user_026', 'hello', NULL),
	(27,  'test_user_027', 'hello', NULL),
	(28,  'test_user_028', 'hello', NULL),
	(29,  'test_user_029', 'hello', NULL),
	(30,  'test_user_030', 'hello', NULL),
	(31,  'test_user_031', 'hello', NULL),
	(32,  'test_user_032', 'hello', NULL),
	(33,  'test_user_033', 'hello', NULL),
	(34,  'test_user_034', 'hello', NULL),
	(35,  'test_user_035', 'hello', NULL),
	(36,  'test_user_036', 'hello', NULL),
	(37,  'test_user_037', 'hello', NULL),
	(38,  'test_user_038', 'hello', NULL),
	(39,  'test_user_039', 'hello', NULL),
	(40,  'test_user_040', 'hello', NULL),
	(41,  'test_user_041', 'hello', NULL),
	(42,  'test_user_042', 'hello', NULL),
	(43,  'test_user_043', 'hello', NULL),
	(44,  'test_user_044', 'hello', NULL),
	(45,  'test_user_045', 'hello', NULL),
	(46,  'test_user_046', 'hello', NULL),
	(47,  'test_user_047', 'hello', NULL),
	(48,  'test_user_048', 'hello', NULL),
	(49,  'test_user_049', 'hello', NULL),
	(50,  'test_user_050', 'hello', NULL),
	(51,  'test_user_051', 'hello', NULL),
	(52,  'test_user_052', 'hello', NULL),
	(53,  'test_user_053', 'hello', NULL),
	(54,  'test_user_054', 'hello', NULL),
	(55,  'test_user_055', 'hello', NULL),
	(56,  'test_user_056', 'hello', NULL),
	(57,  'test_user_057', 'hello', NULL),
	(58,  'test_user_058', 'hello', NULL),
	(59,  'test_user_059', 'hello', NULL),
	(60,  'test_user_060', 'hello', NULL),
	(61,  'test_user_061', 'hello', NULL),
	(62,  'test_user_062', 'hello', NULL),
	(63,  'test_user_063', 'hello', NULL),
	(64,  'test_user_064', 'hello', NULL),
	(65,  'test_user_065', 'hello', NULL),
	(66,  'test_user_066', 'hello', NULL),
	(67,  'test_user_067', 'hello', NULL),
	(68,  'test_user_068', 'hello', NULL),
	(69,  'test_user_069', 'hello', NULL),
	(70,  'test_user_070', 'hello', NULL),
	(71,  'test_user_071', 'hello', NULL),
	(72,  'test_user_072', 'hello', NULL),
	(73,  'test_user_073', 'hello', NULL),
	(74,  'test_user_074', 'hello', NULL),
	(75,  'test_user_075', 'hello', NULL),
	(76,  'test_user_076', 'hello', NULL),
	(77,  'test_user_077', 'hello', NULL),
	(78,  'test_user_078', 'hello', NULL),
	(79,  'test_user_079', 'hello', NULL),
	(80,  'test_user_080', 'hello', NULL),
	(81,  'test_user_081', 'hello', NULL),
	(82,  'test_user_082', 'hello', NULL),
	(83,  'test_user_083', 'hello', NULL),
	(84,  'test_user_084', 'hello', NULL),
	(85,  'test_user_085', 'hello', NULL),
	(86,  'test_user_086', 'hello', NULL),
	(87,  'test_user_087', 'hello', NULL),
	(88,  'test_user_088', 'hello', NULL),
	(89,  'test_user_089', 'hello', NULL),
	(90,  'test_user_090', 'hello', NULL),
	(91,  'test_user_091', 'hello', NULL),
	(92,  'test_user_092', 'hello', NULL),
	(93,  'test_user_093', 'hello', NULL),
	(94,  'test_user_094', 'hello', NULL),
	(95,  'test_user_095', 'hello', NULL),
	(96,  'test_user_096', 'hello', NULL),
	(97,  'test_user_097', 'hello', NULL),
	(98,  'test_user_098', 'hello', NULL),
	(99,  'test_user_099', 'hello', NULL),
	(100, 'test_user_100', 'hello', NULL)
	;
`
	MySQLSelectAllFromTestUser = `
SELECT * FROM test_user
`
	MySQLSelectAllFromTestUserWhereIDEq1 = `
SELECT * FROM test_user WHERE id = 1
`
	MySQLSelectIDFromTestUser = `
SELECT id FROM test_user
`
	MySQLSelectIDFromTestUserWhereIDEq1 = `
SELECT id FROM test_user WHERE id = 1
`
	MySQLSelectAllFromTestUserWhereIDEq999999 = `
SELECT * FROM test_user WHERE id = 999999
`
)

func InitTestDB(ctx context.Context, db *sql.DB) (err error) {
	if _, err := db.ExecContext(ctx, MySQLCreateTableTestUser); err != nil {
		return errorz.Errorf("db.Exec: q=%q: %w", MySQLCreateTableTestUser, err)
	}
	if _, err := db.ExecContext(ctx, MySQLInsertTestUser); err != nil {
		return errorz.Errorf("db.Exec: q=%q: %w", MySQLInsertTestUser, err)
	}

	return nil
}

func TestQuery(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dsn, cleanup, err := mysql.NewTestDBv8_1(ctx)
	if err != nil {
		t.Fatalf("❌: mysql.NewTestDBv8_1: %v", err)
	}
	t.Cleanup(func() {
		if err := cleanup(ctx); err != nil {
			err = errorz.Errorf("cleanup: %w", err)
			log.Print(err)
		}
	})
	db, err := sqlz.OpenContext(ctx, "mysql", dsn)
	if err != nil {
		t.Fatalf("❌: sqlz.OpenContext: %v", err)
	}
	if err := InitTestDB(ctx, db); err != nil {
		t.Fatalf("❌: InitTestDB: %v", err)
	}

	dbz := sqlz.NewDB(db)

	t.Run("success,QueryContext,slice", func(t *testing.T) {
		t.Parallel()
		var testUsers []TestUser
		if err := dbz.QueryContext(ctx, &testUsers, MySQLSelectAllFromTestUser); err != nil {
			t.Fatalf("❌: dbz.QueryContext: %v", err)
		}
		if expect, actual := 1, testUsers[0].ID; expect != actual {
			t.Errorf("❌: testUsers[0].ID: expect(%v) != actual(%v)", expect, actual)
		}
		if expect, actual := 2, testUsers[1].ID; expect != actual {
			t.Errorf("❌: testUsers[1].ID: expect(%v) != actual(%v)", expect, actual)
		}
		t.Logf("\n✅: %s: testUsers[0]: %#v", t.Name(), testUsers[0])
		t.Logf("\n✅: %s: testUsers[1]: %#v", t.Name(), testUsers[1])
		t.Logf("\n✅: %s: testUsers[len(testUsers)-1]: %#v", t.Name(), testUsers[len(testUsers)-1])
	})

	t.Run("success,QueryContext,pointerSlice", func(t *testing.T) {
		t.Parallel()
		var testPointerUsers []*TestUser
		if err := dbz.QueryContext(ctx, &testPointerUsers, MySQLSelectAllFromTestUser); err != nil {
			t.Fatalf("❌: dbz.QueryContext: %v", err)
		}
		if expect, actual := 1, testPointerUsers[0].ID; expect != actual {
			t.Errorf("❌: testPointerUsers[0].ID: expect(%v) != actual(%v)", expect, actual)
		}
		if expect, actual := 2, testPointerUsers[1].ID; expect != actual {
			t.Errorf("❌: testPointerUsers[1].ID: expect(%v) != actual(%v)", expect, actual)
		}
		t.Logf("\n✅: %s: testPointerUsers[0]: %#v", t.Name(), testPointerUsers[0])
		t.Logf("\n✅: %s: testPointerUsers[1]: %#v", t.Name(), testPointerUsers[1])
		t.Logf("\n✅: %s: testPointerUsers[len(testPointerUsers)-1]: %#v", t.Name(), testPointerUsers[len(testPointerUsers)-1])
	})

	t.Run("success,QueryContext,intSlice", func(t *testing.T) {
		t.Parallel()
		var testUserIDs []int
		if err := dbz.QueryContext(ctx, &testUserIDs, MySQLSelectIDFromTestUser); err != nil {
			t.Fatalf("❌: dbz.QueryContext: %v", err)
		}
		if expect, actual := 1, testUserIDs[0]; expect != actual {
			t.Errorf("❌: testUserIDs[0]: expect(%v) != actual(%v)", expect, actual)
		}
		if expect, actual := 2, testUserIDs[1]; expect != actual {
			t.Errorf("❌: testUserIDs[1]: expect(%v) != actual(%v)", expect, actual)
		}
		if expect, actual := 100, testUserIDs[len(testUserIDs)-1]; expect != actual {
			t.Errorf("❌: testUserIDs[1]: expect(%v) != actual(%v)", expect, actual)
		}
		t.Logf("\n✅: %s: testUserIDs[0]: %#v", t.Name(), testUserIDs[0])
		t.Logf("\n✅: %s: testUserIDs[1]: %#v", t.Name(), testUserIDs[1])
		t.Logf("\n✅: %s: testUserIDs[len(testUserIDs)-1]: %#v", t.Name(), testUserIDs[len(testUserIDs)-1])
	})

	t.Run("success,QueryContext,intPointerSlice", func(t *testing.T) {
		t.Parallel()
		var testUserIDs []*int
		if err := dbz.QueryContext(ctx, &testUserIDs, MySQLSelectIDFromTestUser); err != nil {
			t.Fatalf("❌: dbz.QueryContext: %v", err)
		}
		if expect, actual := 1, testUserIDs[0]; expect != *actual {
			t.Errorf("❌: testUserIDs[0]: expect(%v) != actual(%v)", expect, *actual)
		}
		if expect, actual := 2, testUserIDs[1]; expect != *actual {
			t.Errorf("❌: testUserIDs[1]: expect(%v) != actual(%v)", expect, *actual)
		}
		if expect, actual := 100, testUserIDs[len(testUserIDs)-1]; expect != *actual {
			t.Errorf("❌: testUserIDs[1]: expect(%v) != actual(%v)", expect, *actual)
		}
		t.Logf("\n✅: %s: testUserIDs[0]: %#v", t.Name(), *testUserIDs[0])
		t.Logf("\n✅: %s: testUserIDs[1]: %#v", t.Name(), *testUserIDs[1])
		t.Logf("\n✅: %s: testUserIDs[len(testUserIDs)-1]: %#v", t.Name(), *testUserIDs[len(testUserIDs)-1])
	})

	t.Run("failure,QueryContext", func(t *testing.T) {
		t.Parallel()

		var testUsers []TestUser
		if err := dbz.QueryContext(ctx, &testUsers, MySQLSelectAllFromTestUserWhereIDEq999999); err != nil {
			t.Fatalf("❌: dbz.QueryContext: %v", err)
		}
		if expect, actual := 0, len(testUsers); expect != actual {
			t.Errorf("❌: len(testUsers): expect(%v) != actual(%v)", expect, actual)
		}

		var testPointerUsers []*TestUser
		if err := dbz.QueryContext(ctx, &testPointerUsers, MySQLSelectAllFromTestUserWhereIDEq999999); err != nil {
			t.Fatalf("❌: dbz.QueryContext: %v", err)
		}
		if expect, actual := 0, len(testPointerUsers); expect != actual {
			t.Errorf("❌: len(testPointerUsers): expect(%v) != actual(%v)", expect, actual)
		}
	})

	t.Run("success,QueryRowContext,reflect.Struct", func(t *testing.T) {
		t.Parallel()
		var testUser TestUser
		if err := dbz.QueryRowContext(ctx, &testUser, MySQLSelectAllFromTestUserWhereIDEq1); err != nil {
			t.Fatalf("❌: dbz.QueryContext: %v", err)
		}
		if expect, actual := 1, testUser.ID; expect != actual {
			t.Errorf("❌: testUser.ID: expect(%v) != actual(%v)", expect, actual)
		}
		t.Logf("\n✅: %s: testUser: %#v", t.Name(), testUser)
	})

	t.Run("success,QueryRowContext,reflect.Int", func(t *testing.T) {
		t.Parallel()
		var testUserID int
		if err := dbz.QueryRowContext(ctx, &testUserID, MySQLSelectIDFromTestUserWhereIDEq1); err != nil {
			t.Fatalf("❌: dbz.QueryContext: %v", err)
		}
		if expect, actual := 1, testUserID; expect != actual {
			t.Errorf("❌: testUser.ID: expect(%v) != actual(%v)", expect, actual)
		}
		t.Logf("\n✅: %s: testUser: %#v", t.Name(), testUserID)
	})

	t.Run("failure,QueryRowContext", func(t *testing.T) {
		t.Parallel()
		var overflowTestUser TestUser
		if err := dbz.QueryRowContext(ctx, &overflowTestUser, MySQLSelectAllFromTestUser); err != nil {
			t.Fatalf("❌: dbz.QueryRowContext: %v", err)
		}
		if expect, actual := 1, overflowTestUser.ID; expect != actual {
			t.Errorf("❌: overflowTestUser.ID: expect(%v) != actual(%v)", expect, actual)
		}

		var notFoundTestUser TestUser
		if expect, actual := sql.ErrNoRows, dbz.QueryRowContext(ctx, &notFoundTestUser, MySQLSelectAllFromTestUserWhereIDEq999999); !errors.Is(actual, expect) {
			t.Errorf("❌: dbz.QueryRowContext: expect(%v) != actual(%v)", expect, actual)
		}
	})
}
