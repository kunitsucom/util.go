// nolint: testpackage
package bytez

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"
)

func TestReadSeekBuffer_Read(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		buf := NewReadSeekBuffer(bytes.NewBufferString("TestString"))
		{
			expect := []byte("TestS")
			actual := make([]byte, 5)
			_, err := buf.Read(actual)
			if err != nil {
				t.Errorf("err != nil: %v", err)
			}
			if !reflect.DeepEqual(expect, actual) {
				t.Errorf("expect != actual: %s != %s", expect, actual)
			}
		}
		{
			expect := []byte("tring")
			actual := make([]byte, 5)
			_, err := buf.Read(actual)
			if err != nil {
				t.Errorf("err != nil: %v", err)
			}
			if !reflect.DeepEqual(expect, actual) {
				t.Errorf("expect != actual: %s != %s", expect, actual)
			}
		}
		{
			_, err := buf.Seek(0, io.SeekStart)
			if !errors.Is(err, ErrOriginalIOReaderHasNotCompletedReading) {
				t.Errorf("err != nil: %v", err)
			}
		}
		{
			expect := make([]byte, 5)
			actual := make([]byte, 5)
			_, err := buf.Read(actual)
			if !errors.Is(err, io.EOF) {
				t.Errorf("err != io.EOF: %v", err)
			}
			if !reflect.DeepEqual(expect, actual) {
				t.Errorf("expect != actual: %s != %s", expect, actual)
			}
		}
		{
			if _, err := buf.Seek(0, io.SeekStart); err != nil {
				t.Errorf("err != nil: %v", err)
			}
			expect := []byte("TestS")
			actual := make([]byte, 5)
			_, err := buf.Read(actual)
			if err != nil {
				t.Errorf("err != nil: %v", err)
			}
			if !reflect.DeepEqual(expect, actual) {
				t.Errorf("expect != actual: %s != %s", expect, actual)
			}
		}
		{
			if _, err := buf.Seek(-6, io.SeekEnd); err != nil {
				t.Errorf("err != nil: %v", err)
			}
			expect := []byte("String")
			actual := make([]byte, 6)
			_, err := buf.Read(actual)
			if err != nil {
				t.Errorf("err != nil: %v", err)
			}
			if !reflect.DeepEqual(expect, actual) {
				t.Errorf("expect != actual: %s != %s", expect, actual)
			}
		}
		{
			expect := []byte("")
			actual, err := io.ReadAll(buf)
			if err != nil {
				t.Errorf("err != nil: %v", err)
			}
			if !reflect.DeepEqual(expect, actual) {
				t.Errorf("expect != actual: %s != %s", expect, actual)
			}
		}
		{
			if _, err := buf.Seek(0, io.SeekStart); err != nil {
				t.Errorf("err != nil: %v", err)
			}
			expect := []byte("TestString")
			actual, err := io.ReadAll(buf)
			if err != nil {
				t.Errorf("err != nil: %v", err)
			}
			if !reflect.DeepEqual(expect, actual) {
				t.Errorf("expect != actual: %s != %s", expect, actual)
			}
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		buf := NewReadSeekBuffer(bytes.NewBuffer(nil))
		expect := []byte{}
		actual := []byte{}
		buf.completed = true
		buf.off = len(buf.buf)
		n, err := buf.Read(actual)
		if !errors.Is(err, nil) {
			t.Errorf("err != io.EOF: %v", err)
		}
		if n != 0 {
			t.Errorf("n != 0: %v", n)
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %s != %s", expect, actual)
		}
	})
}

func TestReadSeekBuffer_Seek(t *testing.T) {
	t.Parallel()

	t.Run("success(whence=0", func(t *testing.T) {
		t.Parallel()
		buf := NewReadSeekBuffer(bytes.NewBufferString("TestString"))
		{
			expect := []byte("TestString")
			actual, err := io.ReadAll(buf)
			if err != nil {
				t.Errorf("err != nil: %v", err)
			}
			if !reflect.DeepEqual(expect, actual) {
				t.Errorf("expect != actual: %s != %s", expect, actual)
			}
		}
		_, err := buf.Seek(0, io.SeekStart)
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		expect := []byte("TestString")
		actual, err := io.ReadAll(buf)
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %s != %s", expect, actual)
		}
		if _, err := buf.Seek(0, io.SeekStart); err != nil {
			t.Errorf("err != nil: %v", err)
		}
	})

	t.Run("success(whence=1)", func(t *testing.T) {
		t.Parallel()
		buf := NewReadSeekBuffer(bytes.NewBufferString("TestString"))
		{
			expect := []byte("TestString")
			actual, err := io.ReadAll(buf)
			if err != nil {
				t.Errorf("err != nil: %v", err)
			}
			if !reflect.DeepEqual(expect, actual) {
				t.Errorf("expect != actual: %s != %s", expect, actual)
			}
			if _, err := buf.Seek(0, io.SeekStart); err != nil {
				t.Errorf("err != nil: %v", err)
			}
		}
		if _, err := buf.Seek(1, io.SeekCurrent); err != nil {
			t.Errorf("err != nil: %v", err)
		}
		expect := []byte("estString")
		actual, err := io.ReadAll(buf)
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %s != %s", expect, actual)
		}
	})

	t.Run("success(whence=2)", func(t *testing.T) {
		t.Parallel()
		buf := NewReadSeekBuffer(bytes.NewBufferString("TestString"))
		{
			expect := []byte("TestString")
			actual, err := io.ReadAll(buf)
			if err != nil {
				t.Errorf("err != nil: %v", err)
			}
			if !reflect.DeepEqual(expect, actual) {
				t.Errorf("expect != actual: %s != %s", expect, actual)
			}
			if _, err := buf.Seek(0, io.SeekStart); err != nil {
				t.Errorf("err != nil: %v", err)
			}
		}
		if _, err := buf.Seek(-9, io.SeekEnd); err != nil {
			t.Errorf("err != nil: %v", err)
		}
		expect := []byte("estString")
		actual, err := io.ReadAll(buf)
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %s != %s", expect, actual)
		}
	})

	t.Run("failure(ErrOffsetExceedsBufferSize)", func(t *testing.T) {
		t.Parallel()
		buf := NewReadSeekBuffer(bytes.NewBufferString("TestString"))
		expect := []byte("TestString")
		actual, err := io.ReadAll(buf)
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %s != %s", expect, actual)
		}
		if _, err := buf.Seek(maxInt, io.SeekStart); !errors.Is(err, ErrOffsetExceedsBufferSize) {
			t.Errorf("err != ErrOffsetExceedsBufferSize: %v", err)
		}
	})

	t.Run("failure(ErrUnexpectedWhence)", func(t *testing.T) {
		t.Parallel()
		buf := NewReadSeekBuffer(bytes.NewBufferString("TestString"))
		expect := []byte("TestString")
		actual, err := io.ReadAll(buf)
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %s != %s", expect, actual)
		}
		if _, err := buf.Seek(0, 3); !errors.Is(err, ErrUnexpectedWhence) {
			t.Errorf("err != ErrUnexpectedWhence: %v", err)
		}
	})

	t.Run("failure(ErrOffsetExceedsBufferSize)", func(t *testing.T) {
		t.Parallel()
		buf := NewReadSeekBuffer(bytes.NewBufferString("TestString"))
		expect := []byte("TestString")
		actual, err := io.ReadAll(buf)
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %s != %s", expect, actual)
		}
		if _, err := buf.Seek(10, io.SeekCurrent); !errors.Is(err, ErrOffsetExceedsBufferSize) {
			t.Errorf("err != ErrOffsetExceedsBufferSize: %v", err)
		}
	})
}
