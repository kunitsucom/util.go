package strconvz

import (
	"fmt"
	"strconv"
)

// Atoi64 is equivalent to ParseInt(s, 10, 64).
func Atoi64(s string) (int64, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseInt: %w", err)
	}

	return v, nil
}

// Itoa64 is equivalent to strconv.FormatInt(i, 10).
func Itoa64(i int64) string {
	return strconv.FormatInt(i, 10)
}
