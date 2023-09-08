package sql_test

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/kunitsucom/ilog.go"

	sqlz "github.com/kunitsucom/util.go/database/sql"
	errorz "github.com/kunitsucom/util.go/errors"
	"github.com/kunitsucom/util.go/test/database/sql/mysql"
)

type TestUser struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func (u *TestUser) GoString() string {
	typ := reflect.TypeOf(*u)
	val := reflect.ValueOf(*u)
	elems := make([]string, typ.NumField())
	for i := 0; typ.NumField() > i; i++ {
		elems[i] = fmt.Sprintf("%s:%#v", typ.Field(i).Name, val.Field(i))
	}
	return fmt.Sprintf("&%s{%s}", typ, strings.Join(elems, ", "))
}

const (
	CreateTableTestUser = `
CREATE TABLE IF NOT EXISTS test_user (
    id INT,
    name VARCHAR(255)
)
`
	InsertTestUser = `
INSERT INTO test_user (id, name) VALUES (1, 'test_user_001'), (2, 'test_user_002');
`
	SelectTestUser = `
SELECT * FROM test_user
`
	SelectTestUserWhereIDEq1 = `
SELECT * FROM test_user WHERE id = 1
`
)

func TestQuery(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	l := ilog.L().Copy()
	ctx = ilog.WithContext(ctx, l)

	dsn, shutdown, err := mysql.NewTestDB(ctx)
	if err != nil {
		t.Fatalf("❌: mysql.NewTestDB: %v", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			err = errorz.Errorf("shutdown: %w", err)
			l.Err(err).Errorf(err.Error())
		}
	}()
	db, err := sqlz.OpenContext(ctx, "mysql", dsn)
	if err != nil {
		t.Fatalf("❌: sqlz.OpenContext: %v", err)
	}

	if _, err := db.Exec(CreateTableTestUser); err != nil {
		t.Fatalf("❌: db.Exec: q=%q: %v", CreateTableTestUser, err)
	}
	if _, err := db.Exec(InsertTestUser); err != nil {
		t.Fatalf("❌: db.Exec: q=%q: %v", InsertTestUser, err)
	}

	var testUsers []TestUser
	if err := sqlz.NewDB(db).QueryContext(ctx, &testUsers, SelectTestUser); err != nil {
		t.Fatalf("❌: sqlz.NewDB(db).QueryContext: %v", err)
	}
	t.Logf("✅: testUsers: %#v", testUsers)

	var testPointerUsers []*TestUser
	if err := sqlz.NewDB(db).QueryContext(ctx, &testPointerUsers, SelectTestUser); err != nil {
		t.Fatalf("❌: sqlz.NewDB(db).QueryContext: %v", err)
	}
	t.Logf("✅: testPointerUsers: %#v", testPointerUsers)

	var testUser TestUser
	if err := sqlz.NewDB(db).QueryRowContext(ctx, &testUser, SelectTestUserWhereIDEq1); err != nil {
		t.Fatalf("❌: sqlz.NewDB(db).QueryContext: %v", err)
	}
	t.Logf("✅: testUser: %#v", testUser)
	if expect, actual := 1, testUser.ID; expect != actual {
		t.Fatalf("❌: testUser.ID: expect(%v) != actual(%v)", expect, actual)
	}
}
