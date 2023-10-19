package sqlz

import (
	"context"
	"database/sql"
	"fmt"
)

func OpenContext(ctx context.Context, driverName string, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db.PingContext: %w", err)
	}

	return db, nil
}

func MustOpenContext(ctx context.Context, driverName string, dataSourceName string) *sql.DB {
	db, err := OpenContext(ctx, driverName, dataSourceName)
	if err != nil {
		err = fmt.Errorf("OpenContext: %w", err)
		panic(err)
	}

	return db
}
