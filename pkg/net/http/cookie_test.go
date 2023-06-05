package httpz_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httpz "github.com/kunitsuinc/util.go/pkg/net/http"
)

func TestParseCookie(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		cookies := httpz.Cookies(httpz.ParseCookies("foo=bar; hoge=fuga"))
		if len(cookies) != 2 {
			t.Errorf("len(cookies) should be 2, but got %d", len(cookies))
		}
		if cookie, ok := cookies.Get("foo"); !ok || cookie.Value != "bar" {
			t.Errorf("cookie should be foo=bar, but got %v", cookie)
		}
		if cookie, ok := cookies.Get("hoge"); !ok || cookie.Value != "fuga" {
			t.Errorf("cookie should be hoge=fuga, but got %v", cookie)
		}
	})

	t.Run("success(nil)", func(t *testing.T) {
		t.Parallel()

		cookies := httpz.Cookies(append([]*http.Cookie{}, nil, nil))
		if len(cookies) != 2 {
			t.Errorf("len(cookies) should be 2, but got %d", len(cookies))
		}
		if cookie, ok := cookies.Get("NotFound"); ok {
			t.Errorf("cookie should be nil, but got %v", cookie)
		}
	})
}

func TestNewCookieHandler(t *testing.T) {
	t.Parallel()

	handler := httpz.NewCookieHandler("foo", func(next http.Handler, w http.ResponseWriter, r *http.Request, cookie *http.Cookie) {
		if cookie.Value != "bar" {
			t.Errorf("cookie.Value should be bar, but got %s", cookie.Value)
		}
	})

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		rw := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Add("Cookie", "foo=bar")
		handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// noop
		})).ServeHTTP(rw, r)
		if rw.Code != http.StatusOK {
			t.Errorf("rw.Code should be %d, but got %d", http.StatusOK, rw.Code)
		}
	})

	t.Run("success(noop)", func(t *testing.T) {
		t.Parallel()
		rw := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// noop
		})).ServeHTTP(rw, r)
		if rw.Code != http.StatusOK {
			t.Errorf("rw.Code should be %d, but got %d", http.StatusOK, rw.Code)
		}
	})
}
