package sqlz

import (
	"context"
	"database/sql"
)

type sqlQueryerContext interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type sqlRows interface {
	Close() error
	Columns() ([]string, error)
	Next() bool
	Scan(...interface{}) error
	Err() error
}

type sqlTxBeginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}
