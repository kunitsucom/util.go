package ioz

import "io"

func ReadAllString(r io.Reader) (string, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err //nolint:wrapcheck
	}

	return string(b), nil
}
