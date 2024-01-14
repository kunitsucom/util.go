package stringz

import (
	"strings"
)

func ReadLine(content string, lineSeparator string, f func(line string, lineSeparator string, lastLine bool) (treated string)) string {
	var result string
	lines := strings.Split(content, lineSeparator)
	lastLine := len(lines) - 1
	for i, line := range lines {
		treated := f(line, lineSeparator, i == lastLine)
		result += treated
	}
	return result
}

func ReadLineFuncRemoveCommentLine(commentPrefix string) func(line string, lineSeparator string, lastLine bool) (treated string) {
	return func(line string, lineSeparator string, lastLine bool) (treated string) {
		trimmed := strings.TrimSpace(line)

		if lastLine && trimmed == "" {
			return ""
		}

		if strings.HasPrefix(trimmed, commentPrefix) {
			return ""
		}

		return line + lineSeparator
	}
}
