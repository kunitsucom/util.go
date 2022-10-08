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
		f := func(buf *ReadSeekBuffer) {
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
		}

		buf := NewReadSeekBuffer(bytes.NewBufferString("TestString"))
		f(buf)
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()

		f := func(buf *ReadSeekBuffer) {
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
		}

		buf := NewReadSeekBuffer(bytes.NewBuffer(nil))
		f(buf)
		buf2 := NewReadSeekBufferBytes([]byte("TestString"))
		f(buf2)
	})
}

func TestReadSeekBuffer_Seek(t *testing.T) {
	t.Parallel()

	t.Run("success(whence=0", func(t *testing.T) {
		t.Parallel()

		f := func(buf *ReadSeekBuffer) {
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
		}

		buf := NewReadSeekBuffer(bytes.NewBufferString("TestString"))
		f(buf)

		buf2 := NewReadSeekBufferBytes([]byte("TestString"))
		f(buf2)

		expect := []byte("TestString")
		actual := buf2.Bytes()
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %s != %s", expect, actual)
		}
	})

	t.Run("success(whence=1)", func(t *testing.T) {
		t.Parallel()
		f := func(buf *ReadSeekBuffer) {
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
		}

		buf := NewReadSeekBuffer(bytes.NewBufferString("TestString"))
		f(buf)
		buf2 := NewReadSeekBufferBytes([]byte("TestString"))
		f(buf2)
	})

	t.Run("success(whence=2)", func(t *testing.T) {
		t.Parallel()
		f := func(buf *ReadSeekBuffer) {
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
		}

		buf := NewReadSeekBuffer(bytes.NewBufferString("TestString"))
		f(buf)

		buf2 := NewReadSeekBufferBytes([]byte("TestString"))
		f(buf2)
	})

	t.Run("failure(ErrOffsetExceedsBufferSize)", func(t *testing.T) {
		t.Parallel()

		f := func(buf *ReadSeekBuffer) {
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
		}

		buf := NewReadSeekBuffer(bytes.NewBufferString("TestString"))
		f(buf)

		buf2 := NewReadSeekBufferBytes([]byte("TestString"))
		f(buf2)
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

func TestString(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := "TestString"
		buf := NewReadSeekBufferString(expect)
		actual := buf.String()

		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %s != %s", expect, actual)
		}
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var buf *ReadSeekBuffer
		actual := buf.String()

		expect := "<nil>"
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %s != %s", expect, actual)
		}
	})
}
