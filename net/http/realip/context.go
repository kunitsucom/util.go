package realip

import "context"

type contextKey int

const (
	_ contextKey = iota
	keyXRealIP
)

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
