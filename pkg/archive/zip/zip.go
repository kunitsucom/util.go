package zipz

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var ErrDoNotUnzipAFileAtRiskOfZipSlip = errors.New("zipz: do not unzip a file at risk of zip slip")

func Unzip(srcZipFilePath, dstDir string) (paths []string, err error) {
	r, err := zip.OpenReader(srcZipFilePath)
	if err != nil {
		return nil, fmt.Errorf("zip.OpenReader: %w", err)
	}
	defer r.Close()

	paths = make([]string, len(r.File))
	for i, f := range r.File {
		path, err := unzip(f, dstDir)
		if err != nil {
			return nil, fmt.Errorf("unzip: %w", err)
		}

		paths[i] = path
	}

	return paths, nil
}

func unzip(zipfile *zip.File, dstDir string) (path string, err error) {
	if strings.Contains(zipfile.Name, "..") {
		return "", fmt.Errorf("zipfile.Name=%s: %w", zipfile.Name, ErrDoNotUnzipAFileAtRiskOfZipSlip)
	}

	r, err := zipfile.Open()
	if err != nil {
		return "", fmt.Errorf("zipfile.Open: %w", err)
	}
	defer r.Close()

	path = filepath.Join(dstDir, filepath.Clean(zipfile.Name))

	if zipfile.FileInfo().IsDir() {
		if err := os.MkdirAll(path, zipfile.Mode()); err != nil {
			return "", fmt.Errorf("name=%s: os.MkdirAll: %w", path, err)
		}
		return path, nil
	}

	w, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipfile.Mode())
	if err != nil {
		return "", fmt.Errorf("name=%s: os.OpenFile: %w", path, err)
	}
	defer w.Close()

	for {
		if _, err := io.CopyN(w, r, 2048); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return "", fmt.Errorf("name=%s: io.CopyN: %w", path, err)
		}
	}

	return path, nil
}
