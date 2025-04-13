package sqlz

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// ScanRows scans rows to dst.
//
// structTag is used to get column names from struct tags.
//
// dst must be a pointer.
// If dst is a pointer to a struct or a slice of struct, column names are got from structTag.
// If dst is a pointer to a slice of primitive, ignore structTag.
//
//nolint:cyclop
func ScanRows(rows sqlRows, structTag string, dst interface{}) error {
	pointer := reflect.ValueOf(dst) // expect *Type or *[]Type or *[]*Type
	if pointer.Kind() != reflect.Pointer {
		return fmt.Errorf("structSlicePointer.Kind=%s: %w", pointer.Kind(), ErrMustBePointer)
	}
	if pointer.IsNil() {
		return fmt.Errorf("structSlicePointer.IsNil: %w", ErrMustNotNil)
	}

	deref := reflect.Indirect(pointer) // Type or []Type or []*Type <- *Type or *[]Type or *[]*Type
	switch deref.Kind() {              //nolint:exhaustive
	case reflect.Slice:
		if err := scanRowsToSlice(rows, structTag, deref); err != nil { // expect []Type (or []*Type)
			return fmt.Errorf("scanRowsToSlice: type=%T: %w", dst, err)
		}
	case reflect.Struct:
		// behaver like *sql.Row
		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return err //nolint:wrapcheck
			}
			return sql.ErrNoRows
		}
		columns, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("rows.Columns: %w", err)
		}
		dests := make([]interface{}, len(columns))
		tags := getStructTags(deref.Type(), structTag)
		if err := scanRowsToStruct(rows, columns, dests, tags, deref); err != nil { // expect Type (or *Type)
			return fmt.Errorf("scanRowsToStruct: type=%T: %w", dst, err)
		}
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String: // primitives
		// behaver like *sql.Row
		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return err //nolint:wrapcheck
			}
			return sql.ErrNoRows
		}
		if err := rows.Scan(dst); err != nil {
			return fmt.Errorf("rows.Scan: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("type=%T: %w", dst, ErrDataTypeNotSupported)
	}
	return nil
}

func scanRowsToSlice(rows sqlRows, structTag string, destStructSlice reflect.Value) error { // destStructSlice: []Type (or []*Type)
	elementType := destStructSlice.Type().Elem() // Type (or *Type) <- []Type (or []*Type)
	var elementIsPointer bool
	if elementType.Kind() == reflect.Pointer {
		elementIsPointer = true
		elementType = elementType.Elem() // Type <- *Type
	}

	switch elementType.Kind() { //nolint:exhaustive
	case reflect.Struct:
		if err := scanRowsToStructSlice(rows, structTag, elementType, elementIsPointer, destStructSlice); err != nil { // expect []Type (or []*Type)
			return fmt.Errorf("scanRowsToStructSlice: %w", err)
		}
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String: // primitives
		slice, put := getReflectValueSlice()
		defer put()
		for rows.Next() {
			v := reflect.Indirect(reflect.New(elementType))
			if err := rows.Scan(v.Addr().Interface()); err != nil {
				return fmt.Errorf("rows.Scan: %w", err)
			}
			if elementIsPointer {
				slice.Slice = append(slice.Slice, v.Addr())
			} else {
				slice.Slice = append(slice.Slice, v)
			}
		}
		if err := rows.Err(); err != nil {
			return fmt.Errorf("rows.Err: %w", err)
		}
		destStructSlice.Set(reflect.Append(destStructSlice, slice.Slice...))
		return nil
	default:
		// TODO: support other types
		return fmt.Errorf("elem=%s, expected=%s: %w", elementType.Kind(), reflect.Struct, ErrDataTypeNotSupported)
	}
	return nil
}

func scanRowsToStructSlice(rows sqlRows, structTag string, elementType reflect.Type, elementIsPointer bool, destStructSlice reflect.Value) error { // destStructSlice: []Type (or []*Type)
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("rows.Columns: %w", err)
	}
	dests := make([]interface{}, len(columns))
	tags := getStructTags(elementType, structTag)
	slice, put := getReflectValueSlice()
	defer put()
	for rows.Next() {
		v := reflect.Indirect(reflect.New(elementType))
		if err := scanRowsToStruct(rows, columns, dests, tags, v); err != nil {
			return fmt.Errorf("scanRowsToStruct: %w", err)
		}
		if elementIsPointer {
			slice.Slice = append(slice.Slice, v.Addr())
		} else {
			slice.Slice = append(slice.Slice, v)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows.Err: %w", err)
	}
	destStructSlice.Set(reflect.Append(destStructSlice, slice.Slice...))

	return nil
}

func scanRowsToStruct(rows sqlRows, columns []string, dests []interface{}, tags []string, destStruct reflect.Value) error {
	for clmIdx := range columns {
		for tagIdx := range tags {
			if columns[clmIdx] == tags[tagIdx] {
				dests[clmIdx] = destStruct.Field(tagIdx).Addr().Interface()
			}
		}
	}

	if err := rows.Scan(dests...); err != nil {
		return fmt.Errorf("rows.Scan: %w", err)
	}

	return nil
}

//nolint:gochecknoglobals
var (
	tagsCache sync.Map
)

func getStructTags(structType reflect.Type, structTag string) []string {
	if tags, ok := tagsCache.Load(structType); ok {
		return tags.([]string) //nolint:forcetypeassert
	}

	tags := make([]string, structType.NumField())
	for i := range structType.NumField() {
		rawTag := structType.Field(i).Tag.Get(structTag)
		tags[i] = strings.Split(rawTag, ",")[0]
	}
	tagsCache.Store(structType, tags)
	return tags
}

type (
	_ReflectValueSliceType = []reflect.Value
	_ReflectValueSlice     struct{ Slice _ReflectValueSliceType }
)

//nolint:gochecknoglobals
var _ReflectValueSlicePool = &sync.Pool{New: func() interface{} {
	const sliceCap = 128
	return &_ReflectValueSlice{make(_ReflectValueSliceType, 0, sliceCap)}
}} // NOTE: both len and cap are needed.

func getReflectValueSlice() (v *_ReflectValueSlice, put func()) {
	b := _ReflectValueSlicePool.Get().(*_ReflectValueSlice) //nolint:forcetypeassert
	b.Slice = b.Slice[:0]
	return b, func() { _ReflectValueSlicePool.Put(b) }
}
