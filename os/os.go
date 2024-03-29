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

func IsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
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
