package testingz

import "fmt"

var _ fmt.Stringer = (*Stringer)(nil)

type Stringer struct {
	StringFunc func() string
}

func (r *Stringer) String() string {
	return r.StringFunc()
}
