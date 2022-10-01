package slice

import (
	"reflect"
)

func Contains[T comparable](s []T, value T) bool {
	for _, elem := range s {
		if value == elem {
			return true
		}
	}

	return false
}

func DeepContains[T any](s []T, value T) bool {
	for _, elem := range s {
		if reflect.DeepEqual(value, elem) {
			return true
		}
	}

	return false
}

func Equal[T comparable](sliceA, sliceB []T) bool {
	if len(sliceA) != len(sliceB) {
		return false
	}

	for index, elem := range sliceA {
		if elem != sliceB[index] {
			return false
		}
	}

	return true
}

func Exclude[T comparable](s []T, excludes []T) (excluded []T) {
	for _, exclude := range excludes {
		for _, v := range s {
			if v == exclude {
				continue
			}

			excluded = append(excluded, v)
		}
	}

	return excluded
}

func DeepExclude[T any](s []T, exclude T) (excluded []T) {
	for _, v := range s {
		if reflect.DeepEqual(v, exclude) {
			continue
		}

		excluded = append(excluded, v)
	}

	return excluded
}

func Filter[T any](s []T, filter func(index int, elem T) bool) []T {
	filtered := make([]T, 0, len(s))
	for idx, e := range s {
		if filter(idx, e) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

func ToMap[Key comparable, Value any](s []Value, getKey func(index int, value Value) (key Key)) map[Key]Value {
	m := make(map[Key]Value, len(s))
	for idx, elem := range s {
		m[getKey(idx, elem)] = elem
	}
	return m
}

func Generate[Source, Generated any](s []Source, generator func(index int, source Source) (generated Generated)) []Generated {
	gen := make([]Generated, 0, len(s))
	for idx, e := range s {
		gen = append(gen, generator(idx, e))
	}
	return gen
}

func First[T any](s []T) T {
	if len(s) == 0 {
		var zero T
		return zero
	}
	return s[0]
}

func Last[T any](s []T) T {
	if len(s) == 0 {
		var zero T
		return zero
	}
	return s[len(s)-1]
}

func Split[T any](source []T, size int) [][]T {
	sourceLength := len(source)
	splits := make([][]T, 0, sourceLength/size+1)
	for i := 0; i < sourceLength; i += size {
		next := i + size
		if sourceLength < next {
			next = sourceLength
		}
		splits = append(splits, source[i:next])
	}
	return splits
}

func Uniq[T comparable](source []T) []T {
	m := make(map[T]bool)
	uniq := []T{}

	for _, elem := range source {
		if !m[elem] {
			m[elem] = true
			uniq = append(uniq, elem)
		}
	}

	return uniq
}
