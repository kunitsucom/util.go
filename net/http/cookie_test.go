package httpz_test

import (
	"net/http"
	"testing"

	httpz "github.com/kunitsucom/util.go/net/http"
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
