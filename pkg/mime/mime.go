package mime

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func DetectContentType(reader io.Reader) (contentType string, err error) {
	// https://github.com/golang/go/blob/70fd4e47d73b92fe90e44ac785e2f98f9df0ab67/src/net/http/sniff.go#L12-L13
	// The algorithm uses at most sniffLen bytes to make its decision.
	const sniffLen = 512
	buffer := make([]byte, sniffLen)

	// Read sniff
	bytesRead, err := reader.Read(buffer)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("❌: reader.Read: %w", err)
	}

	// Worked around a problem where any content would be marked as application/octet-stream if the content length was less than a sniff.
	buffer = buffer[:bytesRead]

	return http.DetectContentType(buffer), nil
}
