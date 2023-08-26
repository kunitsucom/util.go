package sqlz

import (
	"fmt"
	"reflect"
)

// ScanRows scans rows to dst.
//
// dst must be a pointer.
func ScanRows(rows sqlRows, structTag string, dst interface{}) error {
	pointer := reflect.ValueOf(dst) // expect *Type or *[]Type or *[]*Type
	if pointer.Kind() != reflect.Ptr {
		return fmt.Errorf("structSlicePointer.Kind=%s: %w", pointer.Kind(), ErrMustBePointer)
	}
	if pointer.IsNil() {
		return fmt.Errorf("structSlicePointer.IsNil: %w", ErrMustNotNil)
	}

	deref := pointer.Elem()
	switch deref.Kind() { //nolint:exhaustive
	case reflect.Slice:
		if err := scanRowsToSlice(rows, structTag, deref); err != nil { // []Type (or []*Type)
			return fmt.Errorf("scanRowsToSlice: type=%T: %w", dst, err)
		}
	case reflect.Struct:
		if err := scanRowsToStruct(rows, structTag, deref); err != nil { // Type (or *Type)
			return fmt.Errorf("scanRowsToStruct: type=%T: %w", dst, err)
		}
	default:
		return fmt.Errorf("type=%T: %w", dst, ErrDataTypeNotSupported)
	}
	return nil
}

func scanRowsToSlice(rows sqlRows, structTag string, destStructSlice reflect.Value) error { // destStructSlice: []Type (or []*Type)
	sliceContentType := destStructSlice.Type().Elem() // sliceContentType: Type (or *Type)
	var sliceContentIsPointer bool
	if sliceContentType.Kind() == reflect.Ptr {
		sliceContentIsPointer = true
		sliceContentType = sliceContentType.Elem() // sliceContentType: Type
	}

	if sliceContentType.Kind() != reflect.Struct {
		// TODO: support other types
		return fmt.Errorf("destStructSlice.Kind=%s: %w", destStructSlice.Kind(), ErrDataTypeNotSupported)
	}

	destStructSlice.SetLen(0)
	for rows.Next() {
		v := reflect.New(sliceContentType).Elem()
		if err := scanRowsToStruct(rows, structTag, v); err != nil {
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

func scanRowsToStruct(rows sqlRows, structTag string, destStruct reflect.Value) error {
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
