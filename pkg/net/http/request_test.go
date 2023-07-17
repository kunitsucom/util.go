package httpz_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httpz "github.com/kunitsuinc/util.go/pkg/net/http"
	urlz "github.com/kunitsuinc/util.go/pkg/net/url"
)

func TestDoRequest(t *testing.T) {
	t.Parallel()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := io.Copy(w, r.Body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	t.Cleanup(s.Close)

	t.Run("normal", func(t *testing.T) {
		t.Parallel()

		response, err := httpz.DoRequest(
			context.Background(),
			http.DefaultClient,
			http.MethodPost,
			s.URL,
			httpz.NewHeader(httpz.Set("Content-Type", "application/x-www-form-urlencoded")),
			strings.NewReader(urlz.NewValues(urlz.Add("key1", "value1"), urlz.Add("key2", "value2")).Encode()),
		)
		if err != nil {
			t.Errorf("❌: httpz.DoRequest: err != nil: %v", err)
		}
		if response.StatusCode != http.StatusOK {
			t.Errorf("❌: response.StatusCode != http.StatusOK: %d != %d", response.StatusCode, http.StatusOK)
		}
		buf := bytes.NewBuffer(nil)
		if _, err := buf.ReadFrom(response.Body); err != nil {
			t.Errorf("❌: buf.ReadFrom: err != nil: %v", err)
		}
		const expect = "key1=value1&key2=value2"
		actual := buf.String()
		if actual != expect {
			t.Errorf("❌: expect != actual: %s != %s", expect, actual)
		}
	})

	t.Run("abnormal", func(t *testing.T) {
		t.Parallel()

		if _, err := httpz.DoRequest(
			context.Background(),
			http.DefaultClient,
			string(byte(0x7f)),
			s.URL,
			httpz.NewHeader(httpz.Set("Content-Type", "application/x-www-form-urlencoded")),
			strings.NewReader(urlz.NewValues(urlz.Add("key1", "value1"), urlz.Add("key2", "value2")).Encode()),
		); err == nil {
			t.Errorf("❌: httpz.DoRequest: err == nil")
		}
	})
}
