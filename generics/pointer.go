package genericz

func Pointer[T any](v T) *T {
	return &v
}
