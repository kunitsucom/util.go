package httpz

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func NewRequest(ctx context.Context, method, url string, header http.Header, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	for k, v := range header {
		r.Header[k] = v
	}

	return r, nil
}

func DoRequest(ctx context.Context, client *http.Client, method, url string, header http.Header, body io.Reader) (*http.Response, error) {
	r, err := NewRequest(ctx, method, url, header, body)
	if err != nil {
		return nil, fmt.Errorf("NewRequest: %w", err)
	}

	return client.Do(r) //nolint:wrapcheck
}
