package sqlz

import (
	"context"
	"database/sql"
	"fmt"
)

type Queryer interface {
	QueryStructSliceContext(ctx context.Context, structTag string, destStructSlicePointer interface{}, query string, args ...any) error
	QueryStructContext(ctx context.Context, structTag string, destStructPointer interface{}, query string, args ...any) error
}

type _Queryer struct {
	SQLQueryer
}

func NewDB(db SQLQueryer) Queryer { //nolint:ireturn
	return &_Queryer{
		SQLQueryer: db,
	}
}

func (s *_Queryer) QueryStructSliceContext(ctx context.Context, structTag string, destStructSlicePointer interface{}, query string, args ...any) error {
	rows, err := s.QueryContext(ctx, query, args...) //nolint:rowserrcheck
	return s.queryStructSliceContext(rows, err, structTag, destStructSlicePointer)
}

func (s *_Queryer) queryStructSliceContext(rows SQLRows, queryContextErr error, structTag string, destStructSlicePointer interface{}) error {
	if queryContextErr != nil {
		return fmt.Errorf("QueryContext: %w", queryContextErr)
	}
	defer rows.Close()

	return ScanRows(rows, structTag, destStructSlicePointer)
}

func (s *_Queryer) QueryStructContext(ctx context.Context, structTag string, destStructPointer interface{}, query string, args ...any) error {
	rows, err := s.QueryContext(ctx, query, args...) //nolint:rowserrcheck
	return s.queryStructContext(rows, err, structTag, destStructPointer)
}

func (s *_Queryer) queryStructContext(rows SQLRows, queryContextErr error, structTag string, destStructPointer interface{}) error {
	if queryContextErr != nil {
		return fmt.Errorf("QueryContext: %w", queryContextErr)
	}
	defer rows.Close()

	if !rows.Next() {
		return sql.ErrNoRows
	}

	return ScanRows(rows, structTag, destStructPointer)
}
