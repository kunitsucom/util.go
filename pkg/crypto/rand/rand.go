package randz

import (
	crypto_rand "crypto/rand"
	"io"
)

const DefaultRandomSource = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

//nolint:gochecknoglobals
var (
	DefaultReader io.Reader = NewReader()
)

type Reader struct {
	RandomSource string
	RandomReader io.Reader
}

type ReaderOption func(r *Reader)

func WithRandomSource(str string) ReaderOption {
	return func(r *Reader) {
		r.RandomSource = str
	}
}

func WithRandomReader(reader io.Reader) ReaderOption {
	return func(r *Reader) {
		r.RandomReader = reader
	}
}

func NewReader(opts ...ReaderOption) *Reader {
	r := &Reader{
		RandomSource: DefaultRandomSource,
		RandomReader: crypto_rand.Reader,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = io.ReadFull(r.RandomReader, p)

	randomSourceLength := len(r.RandomSource)
	for i := range p {
		p[i] = r.RandomSource[int(p[i])%randomSourceLength]
	}

	return n, err //nolint:wrapcheck
}

func (r *Reader) ReadString(length int) (random string, err error) {
	b := make([]byte, length)

	if _, err := io.ReadFull(r, b); err != nil {
		return "", err //nolint:wrapcheck
	}

	return string(b), nil
}

func ReadString(length int) (random string, err error) {
	b := make([]byte, length)

	if _, err := io.ReadFull(DefaultReader, b); err != nil {
		return "", err //nolint:wrapcheck
	}

	return string(b), nil
}
