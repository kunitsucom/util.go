package testingz

import (
	"io"
)

var _ io.ReadWriter = (*ReadWriter)(nil)

type ReadWriter struct {
	ReadFunc  func(p []byte) (n int, err error)
	WriteFunc func(p []byte) (n int, err error)
}

func (rw *ReadWriter) Read(p []byte) (n int, err error) {
	return rw.ReadFunc(p)
}

func (rw *ReadWriter) Write(p []byte) (n int, err error) {
	return rw.WriteFunc(p)
}
