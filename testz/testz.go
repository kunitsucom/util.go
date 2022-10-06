package testz

import (
	"errors"
)

var ErrTestError = errors.New("testz: test error")

type ReadWriter struct {
	N   int
	Err error
}

func NewReadWriter(n int, err error) *ReadWriter {
	return &ReadWriter{
		N:   n,
		Err: err,
	}
}

func (t *ReadWriter) Read(p []byte) (n int, err error) {
	return t.N, t.Err
}

func (t *ReadWriter) Write(p []byte) (n int, err error) {
	return t.N, t.Err
}
