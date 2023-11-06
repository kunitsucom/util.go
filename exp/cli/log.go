package cliz

import (
	"io"
	"log"
	"os"
)

//nolint:gochecknoglobals
var (
	// TraceLog is the logger to be used for trace log.
	TraceLog = logger(Stderr, "UTIL_GO_CLI_TRACE", "TRACE: ")
	// DebugLog is the logger to be used for debug log.
	DebugLog = logger(Stderr, "UTIL_GO_CLI_DEBUG", "DEBUG: ")
)

func logger(w io.Writer, environ string, prefix string) *log.Logger {
	if v := os.Getenv(environ); v == "true" {
		return log.New(w, prefix, log.LstdFlags|log.Lshortfile)
	}

	return log.New(io.Discard, prefix, log.LstdFlags)
}
