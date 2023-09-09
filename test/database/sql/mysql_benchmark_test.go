package sql_test

import (
	"context"
	"log"
	"testing"

	"github.com/jmoiron/sqlx"

	sqlz "github.com/kunitsucom/util.go/database/sql"
	errorz "github.com/kunitsucom/util.go/errors"
	"github.com/kunitsucom/util.go/test/database/sql/mysql"
)

func BenchmarkQueryerContext_QueryContext(b *testing.B) {
	ctx := context.Background()
	dsn, cleanup, err := mysql.NewTestDB(ctx)
	if err != nil {
		b.Fatalf("❌: mysql.NewTestDB: %v", err)
	}
	b.Cleanup(func() {
		if err := cleanup(ctx); err != nil {
			err = errorz.Errorf("cleanup: %w", err)
			log.Print(err)
		}
	})
	db, err := sqlz.OpenContext(ctx, "mysql", dsn)
	if err != nil {
		b.Fatalf("❌: sqlz.OpenContext: %v", err)
	}
	if err := InitTestDB(ctx, db); err != nil {
		b.Fatalf("❌: InitTestDB: %v", err)
	}

	dbx := sqlx.NewDb(db, "mysql")
	_ = dbx
	b.Run("sqlx", func(b *testing.B) {
		b.Logf("🚀: %s: %q", b.Name(), MySQLSelectTestUser)
		b.ResetTimer()
		var u []*TestUser
		for i := 0; b.N > i; i++ {
			if err := dbx.SelectContext(ctx, &u, MySQLSelectTestUser); err != nil {
				b.Fatalf("❌: %s: dbx.SelectContext: %v", b.Name(), err)
			}
		}
	})

	dbz := sqlz.NewDB(db)
	_ = dbz
	b.Run("sqlz", func(b *testing.B) {
		b.Logf("🚀: %s: %q", b.Name(), MySQLSelectTestUser)
		b.ResetTimer()
		var u []*TestUser
		for i := 0; b.N > i; i++ {
			if err := dbz.QueryContext(ctx, &u, MySQLSelectTestUser); err != nil {
				b.Fatalf("❌: %s: dbz.QueryContext: %v", b.Name(), err)
			}
		}
	})
}
