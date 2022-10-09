package realip

import (
	"net"
	"net/http"
	"strings"

	"github.com/kunitsuinc/util.go/netz"
)

const (
	HeaderXForwardedFor = "X-Forwarded-For"
	HeaderXRealIP       = "X-Real-IP"
)

// XRealIP returns X-Real-IP value from real_ip_header.
// If real_ip_header is X-Forwarded-For and it has below values:
//
//	X-Forwarded-For: <SpoofingIP>, <ClientIP>, <ProxyIP>, <Proxy2IP>
//
// XRealIP returns <ClientIP>.
// nolint: revive,stylecheck
func XRealIP(r *http.Request, set_real_ip_from []*net.IPNet, real_ip_header string, real_ip_recursive bool) string {
	xff := strings.Split(r.Header.Get(real_ip_header), ",")

	// NOTE: If real_ip_recursive=on, return X-Forwarded-For tail value.
	if !real_ip_recursive {
		return strings.TrimSpace(xff[len(xff)-1])
	}

	var xRealIP net.IP
	for idx := len(xff) - 1; idx >= 0; idx-- {
		ip := net.ParseIP(strings.TrimSpace(xff[idx]))
		// NOTE: If invalid ip, treat previous loop ip as X-Real-IP.
		if len(ip) == 0 {
			break
		}

		xRealIP = ip

		// NOTE: If set_real_ip_from does not contain ip, treat this loop ip as X-Real-IP.
		if !netz.IPNetSet(set_real_ip_from).Contains(ip) {
			break
		}
	}

	// NOTE: If X-Forwarded-For is invalid csv that has invalid IP string, return RemoteAddr as X-Real-IP.
	if len(xRealIP) == 0 {
		return remoteIP(r.RemoteAddr)
	}

	return xRealIP.String()
}

type Handler struct {
	next http.Handler
	// nolint: revive,stylecheck
	set_real_ip_from []*net.IPNet
	// nolint: revive,stylecheck
	real_ip_header string
	// nolint: revive,stylecheck
	real_ip_recursive bool

	clientIPAddressHeader string
}

type Option func(h *Handler) *Handler

// New returns *realip.Handler that appends X-Real-IP header.
// If set_real_ip_from is X-Forwarded-For and it has below values:
//
//	X-Forwarded-For: <SpoofingIP>, <ClientIP>, <ProxyIP>, <Proxy2IP>
//
// *realip.Handler set <ClientIP> to X-Real-IP header.
// nolint: revive,stylecheck
func New(next http.Handler, set_real_ip_from []*net.IPNet, real_ip_header string, real_ip_recursive bool, opts ...Option) *Handler {
	h := &Handler{
		next:                  next,
		set_real_ip_from:      set_real_ip_from,
		real_ip_header:        real_ip_header,
		real_ip_recursive:     real_ip_recursive,
		clientIPAddressHeader: HeaderXRealIP,
	}

	for _, opt := range opts {
		h = opt(h)
	}

	return h
}

func WithClientIPAddressHeader(header string) Option {
	return func(h *Handler) *Handler {
		h.clientIPAddressHeader = header
		return h
	}
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	r.Header.Set(h.clientIPAddressHeader, XRealIP(r, h.set_real_ip_from, h.real_ip_header, h.real_ip_recursive))

	h.next.ServeHTTP(rw, r)
}

func remoteIP(remoteAddr string) string {
	ip, _, _ := net.SplitHostPort(remoteAddr)
	return ip
}