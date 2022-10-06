package errorz

import "strings"

func Contains(err error, substr string) bool {
	return err != nil && strings.Contains(err.Error(), substr)
}
