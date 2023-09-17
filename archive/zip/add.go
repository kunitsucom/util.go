package zipz

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
)

type (
	AddFileToZipOption interface{ apply(*addFileToZipConfig) }

	addFileToZipDecompressionBombLimit int64
)

type addFileToZipConfig struct {
	decompressionBombLimit int64
}

func (f addFileToZipDecompressionBombLimit) apply(cfg *addFileToZipConfig) {
	cfg.decompressionBombLimit = int64(f)
}

func WithAddFileToZipOptionDecompressionBombLimit(decompressionBombLimit int64) AddFileToZipOption { //nolint:ireturn
	return addFileToZipDecompressionBombLimit(decompressionBombLimit)
}

// AddFileToZip is a function to add a file to an existing zip file.
//
// For Decompression Bomb protection, files larger than the size limit will cause an error using io.LimitReader.
// If a value less than or equal to 0 is specified for decompressionBombLimit, no Decompression Bomb protection will be performed.
//
//nolint:funlen,cyclop
func AddFileToZip(dstZipFile string, entryName string, src io.Reader, opts ...AddFileToZipOption) error {
	cfg := new(addFileToZipConfig)
	for _, opt := range opts {
		opt.apply(cfg)
	}

	fi, err := os.Stat(dstZipFile)
	if err != nil {
		return fmt.Errorf("os.Stat: zipFilename=%s: %w", dstZipFile, err)
	}

	// NOTE: 1. Read an existing zip file
	zr, err := zip.OpenReader(dstZipFile)
	if err != nil {
		return fmt.Errorf("zip.OpenReader: zipFilename=%s: %w", dstZipFile, err)
	}
	defer zr.Close()

	// NOTE: 2. Create a buffer
	buf := bytes.NewBuffer(nil)
	zw := zip.NewWriter(buf)

	// NOTE: 3. Copy all entries from the existing zip file to the buffer
	for _, zf := range zr.File {
		if err := func() error { // NOTE: anonymous function for defer zfr.Close()
			zfr, err := zf.Open()
			if err != nil {
				return fmt.Errorf("f.Open: f.Name=%s: %w", zf.Name, err)
			}
			defer zfr.Close()

			header := &zip.FileHeader{
				Name:     zf.Name,
				Method:   zf.Method,
				Modified: zf.Modified,
			}
			zfw, err := zw.CreateHeader(header)
			if err != nil {
				return fmt.Errorf("w.CreateHeader: f.Name=%s: %w", zf.Name, err)
			}

			// NOTE: Use io.LimitReader to make files larger than the size limit an error for Decompression Bomb protection.
			if cfg.decompressionBombLimit > 0 {
				zfr = io.NopCloser(io.LimitReader(zfr, cfg.decompressionBombLimit))
			}

			//nolint:gosec
			if _, err := io.Copy(zfw, zfr); err != nil {
				return fmt.Errorf("io.Copy: f.Name=%s: %w", zf.Name, err)
			}
			return nil
		}(); err != nil {
			return err
		}
	}

	// NOTE: 4. Add a new entry to the buffer
	newZipEntry, err := zw.Create(entryName)
	if err != nil {
		return fmt.Errorf("w.Create: entryName=%s: %w", entryName, err)
	}
	defer zw.Close()
	if _, err := io.Copy(newZipEntry, src); err != nil {
		return fmt.Errorf("io.Copy: entryName=%s: %w", entryName, err)
	}
	zw.Close()

	// NOTE: Overwrite the existing zip file with the buffer
	if err := os.WriteFile(dstZipFile, buf.Bytes(), fi.Mode()); err != nil {
		return fmt.Errorf("os.WriteFile: zipFilename=%s: %w", dstZipFile, err)
	}

	return nil
}
