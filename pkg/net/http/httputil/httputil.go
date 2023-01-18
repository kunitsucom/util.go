package httputilz

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

func DumpResponse(resp *http.Response) (dump []byte, body *bytes.Buffer, err error) {
	return dumpResponse(resp, io.Copy, httputil.DumpResponse)
}

// nolint: revive,stylecheck
func dumpResponse(
	resp *http.Response,
	io_Copy func(dst io.Writer, src io.Reader) (written int64, err error),
	httputil_DumpResponse func(resp *http.Response, body bool) ([]byte, error),
) (dump []byte, body *bytes.Buffer, err error) {
	bodyBuf := bytes.NewBuffer(nil)

	if _, err := io_Copy(bodyBuf, resp.Body); err != nil {
		return nil, nil, fmt.Errorf("(*bytes.Buffer).ReadFrom: %w", err)
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBuf.Bytes()))

	d, err := httputil_DumpResponse(resp, true)
	if err != nil {
		return nil, nil, fmt.Errorf("httputil.DumpResponse: %w", err)
	}

	return d, bodyBuf, nil
}
