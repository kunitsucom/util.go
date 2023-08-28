package testingz

import "io"

var _ io.Writer = (*Writer)(nil)

type Writer struct {
	WriteFunc func(p []byte) (n int, err error)
}

func (w *Writer) Write(p []byte) (n int, err error) {
	return w.WriteFunc(p)
}
