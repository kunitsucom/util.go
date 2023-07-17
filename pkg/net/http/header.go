package httpz

import "net/http"

type Header func(http.Header)

func Add(key, value string) func(http.Header) {
	return func(h http.Header) {
		h.Add(key, value)
	}
}

func Set(key, value string) func(http.Header) {
	return func(h http.Header) {
		h.Set(key, value)
	}
}

func NewHeader(headers ...Header) http.Header {
	h := make(http.Header)

	for _, header := range headers {
		header(h)
	}

	return h
}
