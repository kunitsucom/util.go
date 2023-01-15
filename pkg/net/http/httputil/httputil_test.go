// nolint: testpackage
package httputilz

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"reflect"
	"testing"
)

func TestDumpResponse(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expectBody := "test"
		expectDump := "HTTP/1.1 200 OK\r\n" +
			"Connection: close\r\n" +
			"Content-Type: text/plain\r\n" +
			"\r\n" +
			expectBody
		expectResponse := &http.Response{
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Header:     http.Header{"Content-Type": []string{"text/plain"}},
			Body:       io.NopCloser(bytes.NewBufferString(expectBody)),
		}
		actualDump, actualBody, err := DumpResponse(expectResponse)
		if err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
		if !reflect.DeepEqual(expectDump, string(actualDump)) {
			t.Errorf("❌: expectDump != actualDump: %s", actualDump)
		}
		if expectBody != actualBody.String() {
			t.Errorf("❌: string(expectDump) != string(actualDump): %s", actualDump)
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		testResponse := &http.Response{
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Header:     http.Header{"Content-Type": []string{"text/plain"}},
			Body:       io.NopCloser(bytes.NewBufferString("test")),
		}
		_, _, err := dumpResponse(testResponse, func(dst io.Writer, src io.Reader) (written int64, err error) { return 0, io.EOF }, httputil.DumpResponse)
		if err == nil {
			t.Errorf("❌: err == nil")
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		testResponse := &http.Response{
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Header:     http.Header{"Content-Type": []string{"text/plain"}},
			Body:       io.NopCloser(bytes.NewBufferString("test")),
		}
		if _, _, err := dumpResponse(testResponse, io.Copy, func(resp *http.Response, body bool) ([]byte, error) { return nil, io.EOF }); err == nil {
			t.Errorf("❌: err == nil")
		}
	})
}
