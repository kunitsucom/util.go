package httpz

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	HeaderOrigin                             = "Origin"
	HeaderVary                               = "Vary"
	HeaderAccessControlAllowOrigin           = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods          = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders          = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials      = "Access-Control-Allow-Credentials"
	HeaderAccessControlAllowPrivateNetwork   = "Access-Control-Allow-Private-Network"
	HeaderAccessControlExposeHeaders         = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge                = "Access-Control-Max-Age"
	HeaderAccessControlRequestMethod         = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders        = "Access-Control-Request-Headers"
	HeaderAccessControlRequestPrivateNetwork = "Access-Control-Request-Private-Network"
)

type CORSConfig struct {
	AllowOrigin          string
	AllowMethods         []string
	AllowHeaders         []string
	ExposeHeaders        []string
	AllowCredentials     bool
	AllowPrivateNetwork  bool
	MaxAge               int
	OptionsPassthrough   bool
	OptionsSuccessStatus int
}

type CORS struct {
	configs []*CORSConfig
	next    http.Handler
}

func NewCORSHandler(configs ...*CORSConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		cors := &CORS{
			configs: configs,
			next:    next,
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(HeaderOrigin) == "" {
				// NOTE: Not a request from a web browser
				next.ServeHTTP(w, r)
				return
			}
			if r.Method == http.MethodOptions && r.Header.Get(HeaderAccessControlRequestMethod) != "" {
				cors.handlePreflightRequest(w, r)
			} else {
				cors.handleRequest(w, r)
			}
		})
	}
}

//nolint:gocognit,cyclop
func (cors *CORS) handlePreflightRequest(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get(HeaderOrigin)
	responseHeader := w.Header()

	for _, c := range cors.configs {
		if c.AllowOrigin == origin || c.AllowOrigin == "*" { //nolint:nestif
			// handle Vary header
			responseHeader.Add(HeaderVary, HeaderOrigin)
			responseHeader.Add(HeaderVary, HeaderAccessControlRequestMethod)
			responseHeader.Add(HeaderVary, HeaderAccessControlRequestHeaders)
			if c.AllowPrivateNetwork {
				responseHeader.Add(HeaderVary, HeaderAccessControlRequestPrivateNetwork)
			}

			// handle Access-Control-Request-*
			requestMethod := r.Header.Get(HeaderAccessControlRequestMethod)
			if !c.methodAllowed(requestMethod) {
				cors.next.ServeHTTP(w, r)
				return
			}
			requestHeaders := canonicalHeaders(
				// ref. https://github.com/rs/cors/blob/e90f167479505c4dbe1161306c3c977f162c1442/cors.go#L311-L314
				r.Header.Values(HeaderAccessControlRequestHeaders),
			)
			if !c.headersAllowed(requestHeaders) {
				cors.next.ServeHTTP(w, r)
				return
			}

			// handle Access-Control-Allow-*
			responseHeader.Set(HeaderAccessControlAllowOrigin, c.AllowOrigin)
			responseHeader.Set(HeaderAccessControlAllowMethods, strings.ToUpper(strings.Join(c.AllowMethods, ", ")))
			if len(requestHeaders) > 0 {
				responseHeader.Add(HeaderAccessControlAllowHeaders, strings.Join(requestHeaders, ", "))
			}
			if c.AllowCredentials {
				responseHeader.Set(HeaderAccessControlAllowCredentials, "true")
			}
			if c.AllowPrivateNetwork && r.Header.Get(HeaderAccessControlRequestPrivateNetwork) == "true" {
				responseHeader.Set(HeaderAccessControlAllowPrivateNetwork, "true")
			}
			if c.MaxAge > 0 {
				responseHeader.Set("Access-Control-Max-Age", strconv.Itoa(c.MaxAge))
			}

			// handle response
			if c.OptionsPassthrough {
				cors.next.ServeHTTP(w, r)
				return
			}
			if c.OptionsSuccessStatus > 0 {
				w.WriteHeader(c.OptionsSuccessStatus)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	// if no CORSConfig:
	cors.next.ServeHTTP(w, r)
}

func (cors *CORS) handleRequest(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get(HeaderOrigin)
	responseHeader := w.Header()

	for _, c := range cors.configs {
		if c.AllowOrigin == origin || c.AllowOrigin == "*" { //nolint:nestif
			// handle Vary header
			responseHeader.Add(HeaderVary, HeaderOrigin)

			// handle Access-Control-Allow-*
			if !c.methodAllowed(r.Method) {
				cors.next.ServeHTTP(w, r)
				return
			}
			responseHeader.Set(HeaderAccessControlAllowOrigin, c.AllowOrigin)
			if len(c.ExposeHeaders) > 0 {
				responseHeader.Add(HeaderAccessControlExposeHeaders, strings.Join(c.ExposeHeaders, ", "))
			}
			if c.AllowCredentials {
				responseHeader.Set(HeaderAccessControlAllowCredentials, "true")
			}

			// handle response
			cors.next.ServeHTTP(w, r)
			return
		}
	}

	// if no CORSConfig:
	cors.next.ServeHTTP(w, r)
}

func canonicalHeaders(headers []string) []string {
	h := make([]string, len(headers))
	for i := range headers {
		h[i] = http.CanonicalHeaderKey(headers[i])
	}
	return h
}

func (c *CORSConfig) methodAllowed(method string) bool {
	if len(c.AllowMethods) == 0 {
		// If no method allowed, always return false, even for preflight request
		return false
	}
	method = strings.ToUpper(method)
	if method == http.MethodOptions {
		return true
	}
	for _, m := range c.AllowMethods {
		if m == method {
			return true
		}
	}
	return false
}

// a cross-domain request.
func (c *CORSConfig) headersAllowed(requestedHeaders []string) bool {
	if len(requestedHeaders) == 0 {
		return true
	}
	for _, allowed := range c.AllowHeaders {
		if allowed == "*" {
			return true
		}
		for _, header := range requestedHeaders {
			header = http.CanonicalHeaderKey(header)
			if allowed == header {
				return true
			}
		}
	}
	return false
}
