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
