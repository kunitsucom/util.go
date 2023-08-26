package sqlz

import (
	"context"
	"database/sql"
)

func MustOpen(ctx context.Context, driverName string, dataSourceName string) *sql.DB {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}

	if err := db.PingContext(ctx); err != nil {
		panic(err)
	}

	return db
}
