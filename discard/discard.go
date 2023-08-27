package discard

func Discard(discard any) {
	_ = discard
}

func One[T any](v1 T, _ any) T { //nolint:ireturn
	return v1
}

func Two[T1, T2 any](v1 T1, v2 T2, _ any) (T1, T2) { //nolint:ireturn
	return v1, v2
}

func Three[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, _ any) (T1, T2, T3) { //nolint:ireturn
	return v1, v2, v3
}

func Four[T1, T2, T3, T4 any](v1 T1, v2 T2, v3 T3, v4 T4, _ any) (T1, T2, T3, T4) { //nolint:ireturn
	return v1, v2, v3, v4
}

func Five[T1, T2, T3, T4, T5 any](v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, _ any) (T1, T2, T3, T4, T5) { //nolint:ireturn
	return v1, v2, v3, v4, v5
}
