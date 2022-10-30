package httpz

type contextKey int

const (
	_ contextKey = iota
	keyXRealIP
	keyRequestBodyBuffer
)
