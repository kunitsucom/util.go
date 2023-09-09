package sqlz

import (
	"fmt"
	"reflect"
	"sync"
)

// ScanRows scans rows to dst.
//
// dst must be a pointer.
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
		columns, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("rows.Columns: %w", err)
		}
		dests := make([]interface{}, len(columns))
		tags := getStructTags(deref.Type(), structTag)
		if err := scanRowsToStruct(rows, columns, dests, tags, deref); err != nil { // expect Type (or *Type)
			return fmt.Errorf("scanRowsToStruct: type=%T: %w", dst, err)
		}
	default:
		return fmt.Errorf("type=%T: %w", dst, ErrDataTypeNotSupported)
	}
	return nil
}

func scanRowsToSlice(rows sqlRows, structTag string, destStructSlice reflect.Value) error { // destStructSlice: []Type (or []*Type)
	structType := destStructSlice.Type().Elem() // Type (or *Type) <- []Type (or []*Type)
	var sliceContentIsPointer bool
	if structType.Kind() == reflect.Pointer {
		sliceContentIsPointer = true
		structType = structType.Elem() // Type <- *Type
	}

	if structType.Kind() != reflect.Struct {
		// TODO: support other types
		return fmt.Errorf("elem=%s, expected=%s: %w", structType.Kind(), reflect.Struct, ErrDataTypeNotSupported)
	}

	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("rows.Columns: %w", err)
	}
	dests := make([]interface{}, len(columns))
	tags := getStructTags(structType, structTag)
	slice, put := getReflectValueSlice()
	defer put()
	for rows.Next() {
		v := reflect.Indirect(reflect.New(structType))
		if err := scanRowsToStruct(rows, columns, dests, tags, v); err != nil {
			return fmt.Errorf("scanRowsToStruct: %w", err)
		}
		if sliceContentIsPointer {
			slice.Slice = append(slice.Slice, v.Addr())
		} else {
			slice.Slice = append(slice.Slice, v)
		}
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
var tagsMap sync.Map

func getStructTags(t reflect.Type, structTag string) []string {
	if tags, ok := tagsMap.Load(t); ok {
		return tags.([]string) //nolint:forcetypeassert
	}

	tags := make([]string, t.NumField())
	for i := 0; t.NumField() > i; i++ {
		tags[i] = t.Field(i).Tag.Get(structTag)
	}
	tagsMap.Store(t, tags)
	return tags
}

type (
	_ReflectValueSliceType = []reflect.Value
	_ReflectValueSlice     struct{ Slice _ReflectValueSliceType }
)

//nolint:gochecknoglobals
var _ReflectValueSlicePool = &sync.Pool{New: func() interface{} { return &_ReflectValueSlice{make(_ReflectValueSliceType, 0, 128)} }} // NOTE: both len and cap are needed.

func getReflectValueSlice() (v *_ReflectValueSlice, put func()) {
	b := _ReflectValueSlicePool.Get().(*_ReflectValueSlice) //nolint:forcetypeassert
	b.Slice = b.Slice[:0]
	return b, func() { _ReflectValueSlicePool.Put(b) }
}
