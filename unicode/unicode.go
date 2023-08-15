package unicodez

import (
	"strings"
	"unicode"
)

func TrimNonGraphic(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, s)
}
