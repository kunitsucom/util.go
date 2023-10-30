package ioz

type WriteFunc func(p []byte) (n int, err error)

func (f WriteFunc) Write(p []byte) (n int, err error) {
	return f(p)
}
