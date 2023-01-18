package osz

import (
	"errors"
	"fmt"
	"os"
)

var ErrPathIsNotDirectory = errors.New("osz: path is not directory")

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func CheckDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("os.Stat: %w", err)
	}
	if info.IsDir() {
		return nil
	}

	return ErrPathIsNotDirectory
}

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
