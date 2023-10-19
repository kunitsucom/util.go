package filepathz

import "path/filepath"

func Short(path string) string {
	dirname := filepath.Base(filepath.Dir(path))
	basename := filepath.Base(path)

	return filepath.Join(dirname, basename)
}
