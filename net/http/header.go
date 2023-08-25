package httpz

import "net/http"

type HeaderBuilder interface {
	Add(key, value string) HeaderBuilder
	Set(key, value string) HeaderBuilder
	Merge(header http.Header) HeaderBuilder
	Build() http.Header
}

type headerBuilder struct {
	header http.Header
}

func NewHeaderBuilder() HeaderBuilder {
	return &headerBuilder{
		header: make(http.Header),
	}
}

func (h *headerBuilder) Add(key, value string) HeaderBuilder {
	h.header.Add(key, value)
	return h
}

func (h *headerBuilder) Set(key, value string) HeaderBuilder {
	h.header.Set(key, value)
	return h
}

func (h *headerBuilder) Merge(header http.Header) HeaderBuilder {
	for key, values := range header {
		for _, value := range values {
			h.Add(key, value)
		}
	}
	return h
}

func (h *headerBuilder) Build() http.Header {
	return h.header
}
