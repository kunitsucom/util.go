package genericz

func Pointer[T any](v T) *T {
	return &v
}

func IsZero[T comparable](v T) bool {
	var zero T
	return v == zero
}

func Zero[T comparable](_ T) T { //nolint:ireturn
	var zero T
	return zero
}
