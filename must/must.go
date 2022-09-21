package must

func One[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}

func Two[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err)
	}

	return v1, v2
}

func Three[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, err error) (T1, T2, T3) {
	if err != nil {
		panic(err)
	}

	return v1, v2, v3
}

func Four[T1, T2, T3, T4 any](v1 T1, v2 T2, v3 T3, v4 T4, err error) (T1, T2, T3, T4) {
	if err != nil {
		panic(err)
	}

	return v1, v2, v3, v4
}

func Five[T1, T2, T3, T4, T5 any](v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, err error) (T1, T2, T3, T4, T5) {
	if err != nil {
		panic(err)
	}

	return v1, v2, v3, v4, v5
}
