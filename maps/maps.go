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

type ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

func SortMapByKey[Key ordered, Value any](m map[Key]Value) (sorted []struct {
	Key   Key
	Value Value
},
) {
	sorted = make([]struct {
		Key   Key
		Value Value
	}, len(m))

	sortedKeys := make([]Key, len(m))
	i := 0
	for key := range m {
		sortedKeys[i] = key
		i++
	}
	sort.SliceStable(sortedKeys, func(i, j int) bool {
		return func(k1, k2 Key) bool {
			return k1 < k2
		}(sortedKeys[i], sortedKeys[j])
	})

	for i, key := range sortedKeys {
		sorted[i] = struct {
			Key   Key
			Value Value
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
