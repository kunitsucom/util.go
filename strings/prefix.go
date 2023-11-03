package stringz

import "strings"

func AddPrefix(prefix string, s string, sep string) (added string) {
	if prefix == "" {
		return s
	}
	if s == "" {
		return prefix
	}
	if sep == "" {
		return prefix + s
	}
	ss := strings.Split(s, sep)
	for i := range ss {
		if lastIndex := len(ss) - 1; i == lastIndex && ss[i] == "" {
			continue
		}
		added += prefix + ss[i] + sep
	}
	return added
}
