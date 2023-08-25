package httpz

import (
	"context"
	"net"
	"net/http"
	"strings"

	netz "github.com/kunitsucom/util.go/net"
)

const (
	HeaderXForwardedFor = "X-Forwarded-For"
	HeaderXRealIP       = "X-Real-IP"
)

func DefaultSetRealIPFrom() []*net.IPNet {
	return []*net.IPNet{
		netz.LoopbackAddress,
		netz.LinkLocalAddress,
		netz.PrivateIPAddressClassA,
		netz.PrivateIPAddressClassB,
		netz.PrivateIPAddressClassC,
	}
}

// XRealIP returns X-Real-IP value from real_ip_header.
// If real_ip_header is X-Forwarded-For and it has below values:
//
//	X-Forwarded-For: <SpoofingIP>, <ClientIP>, <ProxyIP>, <Proxy2IP>
//
// XRealIP returns <ClientIP>.
//
// NOTE: Argument naming conforms to NGINX configuration naming.
//
// Example:
//
//	realip := httpz.XRealIP(
//		r,
//		httpz.DefaultSetRealIPFrom(),
//		httpz.HeaderXForwardedFor,
//		true,
//	)
//
//nolint:revive,stylecheck
func XRealIP(r *http.Request, set_real_ip_from []*net.IPNet, real_ip_header string, real_ip_recursive bool) string {
	xff := strings.Split(r.Header.Get(real_ip_header), ",")

	// NOTE: If real_ip_recursive=off, return X-Forwarded-For tail value.
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
		return RemoteIP(r)
	}

	return xRealIP.String()
}

func ContextXRealIP(ctx context.Context) (xRealIP string) {
	v, ok := ctx.Value(keyXRealIP).(string)
	if ok {
		return v
	}

	return ""
}

func ContextWithXRealIP(parent context.Context, xRealIP string) context.Context {
	return context.WithValue(parent, keyXRealIP, xRealIP)
}

type xRealIPHandlerConfig struct {
	//nolint:revive,stylecheck
	set_real_ip_from []*net.IPNet
	//nolint:revive,stylecheck
	real_ip_header string
	//nolint:revive,stylecheck
	real_ip_recursive bool

	clientIPAddressHeader string
}

type XRealIPHandlerOption interface {
	apply(*xRealIPHandlerConfig)
}
type xRealIPHandlerOption func(h *xRealIPHandlerConfig)

func (f xRealIPHandlerOption) apply(h *xRealIPHandlerConfig) { f(h) }

// NewXRealIPHandler returns realip handler that appends X-Real-IP header.
// If set_real_ip_from is X-Forwarded-For and it has below values:
//
//	X-Forwarded-For: <SpoofingIP>, <ClientIP>, <ProxyIP>, <Proxy2IP>
//
// realip handler set <ClientIP> to X-Real-IP header.
//
// NOTE: Argument naming conforms to NGINX configuration naming.
//
// Example:
//
//	httpz.NewXRealIPHandler(
//		httpz.DefaultSetRealIPFrom(),
//		httpz.HeaderXForwardedFor,
//		true,
//	)
//
//nolint:revive,stylecheck
func NewXRealIPHandler(set_real_ip_from []*net.IPNet, real_ip_header string, real_ip_recursive bool, opts ...XRealIPHandlerOption) func(next http.Handler) http.Handler {
	c := &xRealIPHandlerConfig{
		set_real_ip_from:      set_real_ip_from,
		real_ip_header:        real_ip_header,
		real_ip_recursive:     real_ip_recursive,
		clientIPAddressHeader: HeaderXRealIP,
	}

	for _, opt := range opts {
		opt.apply(c)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			xRealIP := XRealIP(r, c.set_real_ip_from, c.real_ip_header, c.real_ip_recursive)

			r.Header.Set(c.clientIPAddressHeader, xRealIP)

			next.ServeHTTP(rw, r.WithContext(ContextWithXRealIP(r.Context(), xRealIP)))
		})
	}
}

func WithClientIPAddressHeader(header string) XRealIPHandlerOption {
	return xRealIPHandlerOption(func(h *xRealIPHandlerConfig) {
		h.clientIPAddressHeader = header
	})
}

func RemoteIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
