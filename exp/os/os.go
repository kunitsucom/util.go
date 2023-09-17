package osz

import (
	"errors"
	"fmt"
	"os"
)

func ReadlinkAndReadFile(path string) (resolved string, bytes []byte, err error) {
	if p, err := os.Readlink(path); errors.Is(err, nil) {
		path = p
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return "", nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	return path, b, nil
}
