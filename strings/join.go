package stringz

import (
	"fmt"
	"strings"
)

func Join(sep string, a ...string) (s string) {
	return strings.Join(a, sep)
}

func JoinStringers[stringer fmt.Stringer](sep string, a ...stringer) (s string) {
	strs := make([]string, len(a))
	for i, v := range a {
		strs[i] = v.String()
	}

	return strings.Join(strs, sep)
}
