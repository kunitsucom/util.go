package testingz

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
)

func NewFileWriter(tb testing.TB) (writer *os.File, closeFunc func() (result *bytes.Buffer), err error) {
	tb.Helper()

	buf := bytes.NewBuffer(nil)
	r, w, err := os.Pipe()
	if err != nil {
		return nil, nil, fmt.Errorf("os.Pipe: %w", err)
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := io.Copy(buf, r); err != nil {
			tb.Logf("io.Copy: %v", err)
		}
	}()

	closeFunc = func() (result *bytes.Buffer) {
		if err := w.Close(); err != nil {
			tb.Logf("w.Close: %v", err)
		}
		wg.Wait()
		return buf
	}

	return w, closeFunc, nil
}
