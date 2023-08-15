package urlz

import "net/url"

type Value func(url.Values)

func Add(key, value string) func(url.Values) {
	return func(v url.Values) {
		v.Add(key, value)
	}
}

func Set(key, value string) func(url.Values) {
	return func(v url.Values) {
		v.Set(key, value)
	}
}

func NewValues(values ...Value) url.Values {
	v := make(url.Values)

	for _, value := range values {
		value(v)
	}

	return v
}
