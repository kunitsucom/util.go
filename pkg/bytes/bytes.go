package bytez

import (
	"errors"
	"io"
)

var _ io.ReadSeeker = (*ReadSeekBuffer)(nil)

type ReadSeekBuffer struct {
	buf       []byte
	r         io.Reader
	completed bool
	off       int
}

func NewReadSeekBuffer(r io.Reader) *ReadSeekBuffer {
	return &ReadSeekBuffer{
		buf:       make([]byte, 0, 64),
		r:         r,
		completed: false,
		off:       0,
	}
}

func NewReadSeekBufferBytes(buf []byte) *ReadSeekBuffer {
	return &ReadSeekBuffer{
		buf:       buf,
		r:         nil,
		completed: true,
		off:       0,
	}
}

func NewReadSeekBufferString(s string) *ReadSeekBuffer {
	return &ReadSeekBuffer{
		buf:       []byte(s),
		r:         nil,
		completed: true,
		off:       0,
	}
}

func (b *ReadSeekBuffer) Read(p []byte) (n int, err error) {
	if !b.completed {
		n, err = b.r.Read(p)

		if shortage := b.off + n - len(b.buf); shortage > 0 {
			b.buf = append(b.buf, make([]byte, shortage)...)
		}

		_ = copy(b.buf[b.off:b.off+n], p)
		b.off += n

		if errors.Is(err, io.EOF) {
			b.completed = true
		}

		//nolint:wrapcheck
		return n, err
	}

	if len(b.buf) <= b.off {
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}

	n = copy(p, b.buf[b.off:])
	b.off += n

	return n, nil
}

const (
	intSize = 32 << (^uint(0) >> 63) // 32 or 64
	maxInt  = 1<<(intSize-1) - 1
)

var (
	ErrOriginalIOReaderHasNotCompletedReading = errors.New("bytez: original io.Reader has not completed reading")
	ErrTooLargeOffset                         = errors.New("bytez: too large offset")
	ErrUnexpectedWhence                       = errors.New("bytez: unexpected whence")
	ErrOffsetExceedsBufferSize                = errors.New("bytez: offset exceeds buffer size")
)

func (b *ReadSeekBuffer) Seek(offset int64, whence int) (newoffset int64, err error) {
	if !b.completed {
		return int64(b.off), ErrOriginalIOReaderHasNotCompletedReading
	}
	if int64(len(b.buf)) < offset {
		return int64(b.off), ErrOffsetExceedsBufferSize
	}

	off := int(offset)
	switch whence {
	case io.SeekStart:
		b.off = off
	case io.SeekCurrent:
		b.off += off
	case io.SeekEnd:
		b.off = len(b.buf) + off
	default:
		return int64(b.off), ErrUnexpectedWhence
	}

	if len(b.buf) < b.off {
		return int64(b.off), ErrOffsetExceedsBufferSize
	}

	return int64(b.off), nil
}

func (b *ReadSeekBuffer) Bytes() []byte {
	return b.buf[b.off:]
}

func (b *ReadSeekBuffer) String() string {
	if b == nil {
		// Special case, useful in debugging.
		return "<nil>"
	}
	return string(b.buf[b.off:])
}
