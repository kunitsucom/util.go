package httpz

import "net/http"

const (
	HeaderCookie = "Cookie"
)

func ParseCookies(rawCookies string) []*http.Cookie {
	header := http.Header{}
	header.Add(HeaderCookie, rawCookies)
	request := http.Request{Header: header}
	return request.Cookies()
}

type Cookies []*http.Cookie

func (cookies Cookies) Get(name string) (cookie *http.Cookie, ok bool) {
	for _, c := range cookies {
		if c == nil {
			continue
		}
		if c.Name == name {
			return c, true
		}
	}
	return nil, false
}

func NewCookieHandler(cookieName string, handler func(next http.Handler, w http.ResponseWriter, r *http.Request, cookie *http.Cookie)) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if rawCookies := r.Header.Get(HeaderCookie); rawCookies != "" {
				cookies := Cookies(ParseCookies(rawCookies))
				if cookie, ok := cookies.Get(cookieName); ok {
					handler(next, rw, r, cookie)
					return
				}
			}

			// noop
			next.ServeHTTP(rw, r)
		})
	}
}
