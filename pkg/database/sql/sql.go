package sqlz

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
)

type SQLQueryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type Rows interface {
	Close() error
	Columns() ([]string, error)
	Next() bool
	Scan(...interface{}) error
	Err() error
}

type Queryer interface {
	QueryStructSliceContext(ctx context.Context, structTag string, destStructSlicePointer interface{}, query string, args ...any) error
	QueryStructContext(ctx context.Context, structTag string, destStructPointer interface{}, query string, args ...any) error
}

type _DB struct {
	SQLQueryer
}

func NewDB(db SQLQueryer) Queryer { //nolint:ireturn
	return &_DB{
		SQLQueryer: db,
	}
}

func (s *_DB) QueryStructSliceContext(ctx context.Context, structTag string, destStructSlicePointer interface{}, query string, args ...any) error {
	rows, err := s.QueryContext(ctx, query, args...) //nolint:rowserrcheck
	return s.queryStructSliceContext(rows, err, structTag, destStructSlicePointer)
}

func (s *_DB) queryStructSliceContext(rows Rows, queryContextErr error, structTag string, destStructSlicePointer interface{}) error {
	if queryContextErr != nil {
		return fmt.Errorf("QueryContext: %w", queryContextErr)
	}
	defer rows.Close()

	return ScanRows(rows, structTag, destStructSlicePointer)
}

func (s *_DB) QueryStructContext(ctx context.Context, structTag string, destStructPointer interface{}, query string, args ...any) error {
	rows, err := s.QueryContext(ctx, query, args...) //nolint:rowserrcheck
	return s.queryStructContext(rows, err, structTag, destStructPointer)
}

func (s *_DB) queryStructContext(rows Rows, queryContextErr error, structTag string, destStructPointer interface{}) error {
	if queryContextErr != nil {
		return fmt.Errorf("QueryContext: %w", queryContextErr)
	}
	defer rows.Close()

	if !rows.Next() {
		return sql.ErrNoRows
	}

	return ScanRows(rows, structTag, destStructPointer)
}

func ScanRows(rows Rows, structTag string, destPointer interface{}) error {
	pointer := reflect.ValueOf(destPointer) // *Type or *[]Type or *[]*Type
	if pointer.Kind() != reflect.Ptr {
		return fmt.Errorf("structSlicePointer.Kind=%s: %w", pointer.Kind(), ErrMustBePointer)
	}
	if pointer.IsNil() {
		return fmt.Errorf("structSlicePointer.IsNil: %w", ErrMustNotNil)
	}

	deref := pointer.Elem()
	switch deref.Kind() { //nolint:exhaustive
	case reflect.Slice:
		if err := scanRowsToStructSlice(rows, deref, structTag); err != nil { // []Type (or []*Type)
			return fmt.Errorf("type=%T: %w", destPointer, err)
		}
	case reflect.Struct:
		if err := scanRowsToStruct(rows, deref, structTag); err != nil { // Type (or *Type)
			return fmt.Errorf("type=%T: %w", destPointer, err)
		}
	default:
		return fmt.Errorf("type=%T: %w", destPointer, ErrDataTypeNotSupported)
	}
	return nil
}

func scanRowsToStructSlice(rows Rows, destStructSlice reflect.Value, structTag string) error { // destStructSlice: []Type (or []*Type)
	sliceContentType := destStructSlice.Type().Elem() // sliceContentType: Type (or *Type)
	var sliceContentIsPointer bool
	if sliceContentType.Kind() == reflect.Ptr {
		sliceContentIsPointer = true
		sliceContentType = sliceContentType.Elem() // sliceContentType: Type
	}

	if sliceContentType.Kind() != reflect.Struct {
		return fmt.Errorf("destStructSlice.Kind=%s: %w", destStructSlice.Kind(), ErrDataTypeNotSupported)
	}

	destStructSlice.SetLen(0)
	for rows.Next() {
		v := reflect.New(sliceContentType).Elem()
		if err := scanRowsToStruct(rows, v, structTag); err != nil {
			return fmt.Errorf("scanRowsToStruct: %w", err)
		}

		if sliceContentIsPointer {
			destStructSlice.Set(reflect.Append(destStructSlice, v.Addr()))
		} else {
			destStructSlice.Set(reflect.Append(destStructSlice, v))
		}
	}

	return nil
}

func scanRowsToStruct(rows Rows, destStruct reflect.Value, structTag string) error {
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("rows.Columns: %w", err)
	}

	structType := destStruct.Type()
	tags := make([]string, structType.NumField())
	values := make([]reflect.Value, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		tags[i] = structType.Field(i).Tag.Get(structTag)
		values[i] = reflect.New(structType.Field(i).Type)
	}

	sqlRows := make([]interface{}, len(columns))
	for i, column := range columns {
		for j, tag := range tags {
			if column == tag {
				sqlRows[i] = values[j].Interface()
			}
		}
	}

	if err := rows.Scan(sqlRows...); err != nil {
		return fmt.Errorf("rows.Scan: %w", err)
	}

	for i := 0; i < structType.NumField(); i++ {
		destStruct.Field(i).Set(values[i].Elem())
	}

	return nil
}

type SQLBeginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

func MustBeginTx(ctx context.Context, db SQLBeginner, opts *sql.TxOptions) *sql.Tx {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		panic(err)
	}
	return tx
}
