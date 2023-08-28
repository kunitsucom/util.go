package testingz

import "io"

var _ io.ReadCloser = (*ReadCloser)(nil)

type ReadCloser struct {
	ReadFunc  func(p []byte) (n int, err error)
	CloseFunc func() error
}

func (rc *ReadCloser) Read(p []byte) (n int, err error) {
	return rc.ReadFunc(p)
}

func (rc *ReadCloser) Close() error {
	return rc.CloseFunc()
}
