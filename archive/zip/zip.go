package zipz

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	osz "github.com/kunitsucom/util.go/os"
)

var ErrDoNotUnzipAFileAtRiskOfZipSlip = errors.New("zipz: do not unzip a file at risk of zip slip")

type (
	zipDirConfig struct {
		walkFunc             func(path string, info os.FileInfo, err error) error
		pathInZipHandlerFunc func(path string) string
	}

	ZipDirOption interface{ apply(*zipDirConfig) }

	zipDirOptionWalkFunc struct {
		f func(path string, info os.FileInfo, err error) error
	}
	zipDirOptionPathInZipHandlerFunc struct {
		f func(path string) string
	}
)

func (f zipDirOptionWalkFunc) apply(cfg *zipDirConfig) { cfg.walkFunc = f.f }
func WithZipDirOptionWalkFunc(f func(path string, info os.FileInfo, err error) error) ZipDirOption { //nolint:ireturn
	return zipDirOptionWalkFunc{f}
}

func (f zipDirOptionPathInZipHandlerFunc) apply(cfg *zipDirConfig) { cfg.pathInZipHandlerFunc = f.f }

// WithZipDirOptionPathInZipHandlerFunc is a function to specify the path in the zip file.
func WithZipDirOptionPathInZipHandlerFunc(f func(path string) string) ZipDirOption { //nolint:ireturn
	return zipDirOptionPathInZipHandlerFunc{f}
}

//nolint:cyclop
func ZipDir(dst io.Writer, srcDir string, opts ...ZipDirOption) error {
	cfg := new(zipDirConfig)
	for _, opt := range opts {
		opt.apply(cfg)
	}

	zipWriter := zip.NewWriter(dst)
	defer zipWriter.Close()

	if err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if cfg.walkFunc != nil {
			if err := cfg.walkFunc(path, info, err); err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		cleaned := filepath.Clean(path)
		file, err := os.Open(cleaned)
		if err != nil {
			return fmt.Errorf("os.Open: name=%s: %w", cleaned, err)
		}
		defer file.Close()

		pathInZip := cleaned
		if cfg.pathInZipHandlerFunc != nil {
			pathInZip = cfg.pathInZipHandlerFunc(pathInZip)
		}
		f, err := zipWriter.Create(pathInZip)
		if err != nil {
			return fmt.Errorf("(*zip.Writer).Create: name=%s: %w", pathInZip, err)
		}

		if _, err = io.Copy(f, file); err != nil {
			return fmt.Errorf("io.Copy: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("filepath.Walk: root=%s: %w", srcDir, err)
	}

	return nil
}

type (
	unzipFileConfig struct {
		unzipFileFileInZipHandler func(zipfile *zip.File, dstDir string) error
	}

	UnzipFileOption interface{ apply(*unzipFileConfig) }

	unzipFileOptionFileInZipHandler struct {
		f func(zipfile *zip.File, dstDir string) error
	}
)

func (f unzipFileOptionFileInZipHandler) apply(cfg *unzipFileConfig) {
	cfg.unzipFileFileInZipHandler = f.f
}

func WithUnzipFileOptionFileInZipHandler(f func(zipfile *zip.File, dstDir string) error) UnzipFileOption { //nolint:ireturn
	return unzipFileOptionFileInZipHandler{f}
}

func UnzipFile(srcZipFilePath, dstDir string, opts ...UnzipFileOption) (paths []string, err error) {
	cfg := new(unzipFileConfig)
	for _, opt := range opts {
		opt.apply(cfg)
	}

	if !osz.IsDir(dstDir) {
		if err := os.MkdirAll(dstDir, 0o755); err != nil {
			return nil, fmt.Errorf("os.MkdirAll: path=%s: %w", dstDir, err)
		}
	}

	r, err := zip.OpenReader(srcZipFilePath)
	if err != nil {
		return nil, fmt.Errorf("zip.OpenReader: %w", err)
	}
	defer r.Close()

	for _, fileInZip := range r.File {
		if cfg.unzipFileFileInZipHandler != nil {
			if err := cfg.unzipFileFileInZipHandler(fileInZip, dstDir); err != nil {
				return nil, fmt.Errorf("unzipFileZipFileHandler: %w", err)
			}
		}
		path, err := unzip(fileInZip, dstDir)
		if err != nil {
			return nil, fmt.Errorf("unzip: %w", err)
		}

		paths = append(paths, path)
	}

	return paths, nil
}

//nolint:cyclop
func unzip(fileInZip *zip.File, dstDir string) (path string, err error) {
	if strings.Contains(fileInZip.Name, "..") {
		return "", fmt.Errorf("(*zip.File).Name=%s: %w", fileInZip.Name, ErrDoNotUnzipAFileAtRiskOfZipSlip)
	}

	r, err := fileInZip.Open()
	if err != nil {
		return "", fmt.Errorf("(*zip.File).Open: %w", err)
	}
	defer r.Close()

	path = filepath.Join(dstDir, filepath.Clean(fileInZip.Name))

	if fileInZip.FileInfo().IsDir() {
		if err := os.MkdirAll(path, fileInZip.Mode()); err != nil {
			return "", fmt.Errorf("os.MkdirAll: path=%s: %w", path, err)
		}
		return path, nil
	}

	if dir := filepath.Dir(path); !osz.IsDir(dir) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return "", fmt.Errorf("os.MkdirAll: path=%s: %w", dir, err)
		}
	}

	w, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileInZip.Mode())
	if err != nil {
		return "", fmt.Errorf("os.OpenFile: name=%s: %w", path, err)
	}
	defer w.Close()

	for {
		if _, err := io.CopyN(w, r, 2048); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return "", fmt.Errorf("io.CopyN: name=%s: %w", path, err)
		}
	}

	return path, nil
}
