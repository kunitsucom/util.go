package sqlz

import (
	"context"
	"database/sql"
)

type SQLQueryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type SQLExecuter interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type SQLRows interface {
	Close() error
	Columns() ([]string, error)
	Next() bool
	Scan(...interface{}) error
	Err() error
}

type SQLTxBeginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}
