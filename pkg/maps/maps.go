package mapz

import (
	"sort"
)

func SortStringKey[V any](m map[string]V) (sorted []struct {
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
