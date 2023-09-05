package stringz

func Join(sep string, a ...string) (s string) {
	for i, v := range a {
		if i > 0 {
			s += sep
		}
		s += v
	}
	return s
}
