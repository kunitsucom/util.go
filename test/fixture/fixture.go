package fixture

import (
	"errors"
)

var ErrFixtureError = errors.New("fixture error")

type Reader struct {
	N   int
	Err error
}

func (t *Reader) Read(p []byte) (n int, err error) {
	return t.N, t.Err
}
