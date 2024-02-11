package errorz

import (
	"errors"
	"net"
	"strings"
)

func HasPrefix(err error, prefix string) bool {
	return err != nil && strings.HasPrefix(err.Error(), prefix)
}

func HasSuffix(err error, suffix string) bool {
	return err != nil && strings.HasSuffix(err.Error(), suffix)
}

func Contains(err error, substr string) bool {
	return err != nil && strings.Contains(err.Error(), substr)
}

func IsNetTimeout(err error) bool {
	if netErr := (net.Error)(nil); errors.As(err, &netErr) {
		return netErr.Timeout()
	}

	return false
}
