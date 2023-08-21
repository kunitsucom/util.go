package mapz

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
)

var (
	ErrKeyIsNotFound            = errors.New(`mapz: key is not found`)
	ErrVIsNotPointerOrInterface = errors.New(`mapz: v is not pointer or interface`)
	ErrValueTypeIsNotMatch      = errors.New(`mapz: value type is not match`)
)

func SortMapStringKey[V any](m map[string]V) (sorted []struct {
	Key   string
	Value V
},
) {
	sorted = make([]struct {
		Key   string
		Value V
	}, len(m))

	sortedKeys := make([]string, len(m))
	i := 0
	for key := range m {
		sortedKeys[i] = key
		i++
	}
	sort.Strings(sortedKeys)

	for i, key := range sortedKeys {
		sorted[i] = struct {
			Key   string
			Value V
		}{
			Key:   key,
			Value: m[key],
		}
	}

	return sorted
}

func SortMapIntKey[V any](m map[int]V) (sorted []struct {
	Key   int
	Value V
},
) {
	sorted = make([]struct {
		Key   int
		Value V
	}, len(m))

	sortedKeys := make([]int, len(m))
	i := 0
	for key := range m {
		sortedKeys[i] = key
		i++
	}
	sort.Ints(sortedKeys)

	for i, key := range sortedKeys {
		sorted[i] = struct {
			Key   int
			Value V
		}{
			Key:   key,
			Value: m[key],
		}
	}

	return sorted
}

func Get[Key comparable](m map[Key]any, key Key, v any) error {
	reflectValue := reflect.ValueOf(v)
	// NOTE: memo
	// if !reflectValue.IsValid() {
	// 	return errors.New("")
	// }
	if reflectValue.Kind() != reflect.Pointer && reflectValue.Kind() != reflect.Interface {
		return fmt.Errorf("v.(type)==%T: %w", v, ErrVIsNotPointerOrInterface)
	}
	reflectValueElem := reflectValue.Elem()
	// NOTE: memo
	// if !reflectValueElem.CanSet() {
	// 	return errors.New("")
	// }
	param, ok := m[key]
	if !ok {
		return fmt.Errorf("map[%v]: %w", key, ErrKeyIsNotFound)
	}
	paramReflectValue := reflect.ValueOf(param)
	if reflectValueElem.Type() != paramReflectValue.Type() {
		return fmt.Errorf("map[%v].(type)==%T, v.(type)==%T: %w", key, param, v, ErrValueTypeIsNotMatch)
	}
	reflectValueElem.Set(paramReflectValue)
	return nil
}
