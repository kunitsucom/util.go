package randz

import (
	crypto_rand "crypto/rand"
	"fmt"
	"io"
)

const DefaultStringReaderRandomSource = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

//nolint:gochecknoglobals
var (
	StringReader = NewReader()
)

type reader struct {
	randomSource string
	randomReader io.Reader
}

type NewReaderOption func(r *reader)

func WithNewReaderRandomSource(s string) NewReaderOption {
	return func(r *reader) {
		r.randomSource = s
	}
}

func WithNewReaderOptionRandomReader(random io.Reader) NewReaderOption {
	return func(r *reader) {
		r.randomReader = random
	}
}

func NewReader(opts ...NewReaderOption) io.Reader {
	r := &reader{
		randomSource: DefaultStringReaderRandomSource,
		randomReader: crypto_rand.Reader,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *reader) Read(p []byte) (n int, err error) {
	n, err = io.ReadFull(r.randomReader, p)
	if err != nil {
		return n, fmt.Errorf("io.ReadFull: %w", err)
	}

	randomSourceLength := len(r.randomSource)
	for i := range p {
		p[i] = r.randomSource[int(p[i])%randomSourceLength]
	}

	return n, nil
}

func ReadString(random io.Reader, length int) (string, error) {
	b := make([]byte, length)

	if _, err := io.ReadFull(random, b); err != nil {
		return "", fmt.Errorf("io.ReadFull: %w", err)
	}

	return string(b), nil
}
