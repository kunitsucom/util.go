package syncz

import (
	"sync"
	"sync/atomic"
)

type Once struct {
	i uint32
	m sync.Mutex
}

// Do is almost identical to (*sync.Once).Do, but counts as "once" only if f does not return an error.
func (o *Once) Do(f func() error) error {
	if o.incomplete() {
		if err := o.do(f); err != nil {
			return err
		}
	}
	return nil
}

func (o *Once) incomplete() bool {
	return atomic.LoadUint32(&o.i) == 0
}

func (o *Once) complete() {
	atomic.StoreUint32(&o.i, 1)
}

func (o *Once) do(f func() error) error {
	o.m.Lock()
	defer o.m.Unlock()
	if o.incomplete() {
		if err := f(); err != nil {
			return err
		}
		o.complete()
	}

	return nil
}

func (o *Once) Reset() {
	o.m.Lock()
	defer o.m.Unlock()
	atomic.StoreUint32(&o.i, 0)
}
