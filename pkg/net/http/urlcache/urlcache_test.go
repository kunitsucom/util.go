//nolint:testpackage
package urlcache

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	testz "github.com/kunitsuinc/util.go/pkg/test"
)

func TestClient_Get(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(strconv.FormatInt(time.Now().UnixNano(), 10)))
	})
	s := httptest.NewServer(mux)
	pathRoot := "http://" + s.Listener.Addr().String()
	defer t.Cleanup(s.Close)

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		c := NewClient(http.DefaultClient, WithCacheTTL[string](1*time.Minute))
		expect, err := c.Get(pathRoot, func(resp *http.Response) bool { return 200 <= resp.StatusCode && resp.StatusCode < 300 }, func(resp *http.Response) (string, error) {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", err //nolint:wrapcheck
			}
			return string(b), nil
		})
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}

		actual, err := c.Get(pathRoot, func(resp *http.Response) bool { return 200 <= resp.StatusCode && resp.StatusCode < 300 }, func(resp *http.Response) (string, error) {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", err //nolint:wrapcheck
			}
			return string(b), nil
		})
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}

		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestClient_Do(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(strconv.FormatInt(time.Now().UnixNano(), 10)))
	})
	mux.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	s := httptest.NewServer(mux)
	pathRoot := "http://" + s.Listener.Addr().String()
	path500 := "http://" + s.Listener.Addr().String() + "/500"
	defer t.Cleanup(s.Close)

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		c := NewClient(http.DefaultClient, WithCacheTTL[string](1*time.Minute))
		req1, _ := http.NewRequest(http.MethodGet, pathRoot, nil)
		expect, err := c.Do(req1, func(resp *http.Response) bool { return 200 <= resp.StatusCode && resp.StatusCode < 300 }, func(resp *http.Response) (string, error) {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", err //nolint:wrapcheck
			}
			return string(b), nil
		})
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}

		req2, _ := http.NewRequest(http.MethodGet, pathRoot, nil)
		actual, err := c.Do(req2, func(resp *http.Response) bool { return 200 <= resp.StatusCode && resp.StatusCode < 300 }, func(resp *http.Response) (string, error) {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", err //nolint:wrapcheck
			}
			return string(b), nil
		})
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}

		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(POST)", func(t *testing.T) {
		t.Parallel()

		c := NewClient(http.DefaultClient, WithCacheTTL[string](1*time.Minute))
		req1, _ := http.NewRequest(http.MethodPost, pathRoot, nil)
		_, err := c.Do(req1, func(resp *http.Response) bool { return 200 <= resp.StatusCode && resp.StatusCode < 300 }, func(resp *http.Response) (string, error) {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", err //nolint:wrapcheck
			}
			return string(b), nil
		})
		if err == nil {
			t.Errorf("err == nil")
		}
		expect := ErrNotSupportedExceptGet.Error()
		if !strings.Contains(err.Error(), expect) {
			t.Errorf("expect != actual: %v != %v", err.Error(), expect)
		}
	})

	t.Run("failure(network)", func(t *testing.T) {
		t.Parallel()

		c := NewClient(http.DefaultClient, WithCacheTTL[string](1*time.Minute))
		req1, _ := http.NewRequest(http.MethodGet, "http://__no_such_host__", nil)
		_, err := c.Do(req1, func(resp *http.Response) bool { return 200 <= resp.StatusCode && resp.StatusCode < 300 }, func(resp *http.Response) (string, error) {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", err //nolint:wrapcheck
			}
			return string(b), nil
		})
		if err == nil {
			t.Errorf("err == nil")
		}
	})

	t.Run("failure(responseConvert)", func(t *testing.T) {
		t.Parallel()

		c := NewClient(http.DefaultClient, WithCacheTTL[string](1*time.Minute))
		req1, _ := http.NewRequest(http.MethodGet, path500, nil)
		_, err := c.Do(req1, func(resp *http.Response) bool { return 200 <= resp.StatusCode && resp.StatusCode < 300 }, func(resp *http.Response) (string, error) {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", err //nolint:wrapcheck
			}
			return string(b), nil
		})
		if err == nil {
			t.Errorf("err == nil")
		}
		expect := ErrNotCacheable.Error()
		if !strings.Contains(err.Error(), expect) {
			t.Errorf("expect != actual: %v != %v", err.Error(), expect)
		}
	})

	t.Run("failure(not_cacheable)", func(t *testing.T) {
		t.Parallel()

		c := NewClient(http.DefaultClient, WithCacheTTL[string](1*time.Minute))
		req1, _ := http.NewRequest(http.MethodGet, pathRoot, nil)
		_, err := c.Do(req1, func(resp *http.Response) bool { return 200 <= resp.StatusCode && resp.StatusCode < 300 }, func(resp *http.Response) (string, error) {
			return "", io.ErrUnexpectedEOF
		})
		if err == nil {
			t.Errorf("err == nil")
		}
		expect := io.ErrUnexpectedEOF.Error()
		if !strings.Contains(err.Error(), expect) {
			t.Errorf("expect != actual: %v != %v", err.Error(), expect)
		}
	})
}

func TestClient_do(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(strconv.FormatInt(time.Now().UnixNano(), 10)))
	})
	s := httptest.NewServer(mux)
	pathRoot := "http://" + s.Listener.Addr().String()
	defer t.Cleanup(s.Close)

	t.Run("failure(not_cacheable)", func(t *testing.T) {
		t.Parallel()

		c := NewClient(http.DefaultClient, WithCacheTTL[string](1*time.Minute))
		req1, _ := http.NewRequest(http.MethodGet, pathRoot, nil)
		_, err := c.do(req1.URL.String(), func() (*http.Response, error) {
			return &http.Response{Body: io.NopCloser(testz.NewReadWriter(nil, 0, testz.ErrTestError))}, nil
		}, func(resp *http.Response) bool { return 200 <= resp.StatusCode && resp.StatusCode < 300 },
			func(resp *http.Response) (string, error) {
				b, err := io.ReadAll(resp.Body)
				if err != nil {
					return "", err //nolint:wrapcheck
				}
				return string(b), nil
			}, time.Now())
		if err == nil {
			t.Errorf("err == nil")
		}
		expect := testz.ErrTestError.Error()
		if !strings.Contains(err.Error(), expect) {
			t.Errorf("expect != actual: %v != %v", err.Error(), expect)
		}
	})
}
