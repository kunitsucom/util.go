package sqlz

import (
	"context"
	"fmt"
)

type QueryerContext interface {
	// QueryContext executes a query that returns rows, typically a SELECT.
	//
	// The dst must be a pointer.
	// The args are for any placeholder parameters in the query.
	QueryContext(ctx context.Context, dst interface{}, query string, args ...interface{}) error
	// QueryRowContext executes a query that is expected to return at most one row.
	// It always returns a non-nil value or an error.
	//
	// The dst must be a pointer.
	// The args are for any placeholder parameters in the query.
	QueryRowContext(ctx context.Context, dst interface{}, query string, args ...interface{}) error
}

type (
	queryerContext struct {
		sqlQueryer sqlQueryerContext
		// Options
		structTag string
	}

	NewDBOption interface{ apply(qc *queryerContext) }

	newDBOptionStructTag string
)

func (f newDBOptionStructTag) apply(qc *queryerContext) { qc.structTag = string(f) }
func WithNewDBOptionStructTag(structTag string) NewDBOption { //nolint:ireturn
	return newDBOptionStructTag(structTag)
}

func NewDB(db sqlQueryerContext, opts ...NewDBOption) QueryerContext { //nolint:ireturn
	return newDB(db, opts...)
}

const defaultStructTag = "db"

func newDB(db sqlQueryerContext, opts ...NewDBOption) *queryerContext {
	qc := &queryerContext{
		sqlQueryer: db,
		structTag:  defaultStructTag,
	}

	for _, opt := range opts {
		opt.apply(qc)
	}

	return qc
}

func (qc *queryerContext) QueryContext(ctx context.Context, dst interface{}, query string, args ...interface{}) error {
	rows, err := qc.sqlQueryer.QueryContext(ctx, query, args...) //nolint:rowserrcheck
	return qc.queryContext(rows, err, dst)
}

func (qc *queryerContext) queryContext(rows sqlRows, queryContextErr error, dst interface{}) (err error) {
	if queryContextErr != nil {
		return fmt.Errorf("QueryContext: %w", queryContextErr)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("rows.Close: %w", closeErr)
			return
		}
	}()

	return ScanRows(rows, qc.structTag, dst)
}

func (qc *queryerContext) QueryRowContext(ctx context.Context, dst interface{}, query string, args ...interface{}) error {
	rows, err := qc.sqlQueryer.QueryContext(ctx, query, args...) //nolint:rowserrcheck
	return qc.queryRowContext(rows, err, dst)
}

func (qc *queryerContext) queryRowContext(rows sqlRows, queryContextErr error, dst interface{}) (err error) {
	if queryContextErr != nil {
		return fmt.Errorf("QueryContext: %w", queryContextErr)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("rows.Close: %w", closeErr)
			return
		}
	}()

	return ScanRows(rows, qc.structTag, dst)
}
