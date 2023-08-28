package testingz

import "io"

var _ io.Reader = (*Reader)(nil)

type Reader struct {
	ReadFunc func(p []byte) (n int, err error)
}

func (r *Reader) Read(p []byte) (n int, err error) {
	return r.ReadFunc(p)
}
