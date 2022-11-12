package discard

//nolint:ireturn
func Discard(discard any) {
	_ = discard
}

//nolint:ireturn
func One[T any](v1 T, _ any) T {
	return v1
}

//nolint:ireturn
func Two[T1, T2 any](v1 T1, v2 T2, _ any) (T1, T2) {
	return v1, v2
}

//nolint:ireturn
func Three[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, _ any) (T1, T2, T3) {
	return v1, v2, v3
}

//nolint:ireturn
func Four[T1, T2, T3, T4 any](v1 T1, v2 T2, v3 T3, v4 T4, _ any) (T1, T2, T3, T4) {
	return v1, v2, v3, v4
}

//nolint:ireturn
func Five[T1, T2, T3, T4, T5 any](v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, _ any) (T1, T2, T3, T4, T5) {
	return v1, v2, v3, v4, v5
}
