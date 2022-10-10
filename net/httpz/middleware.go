package httpz

import "net/http"

type Middleware func(next http.Handler) http.Handler

func Middlewares(middlewares ...func(next http.Handler) http.Handler) Middleware {
	return func(next http.Handler) http.Handler {
		for i := range middlewares {
			next = middlewares[i](next)
		}

		return next
	}
}

func (m Middleware) Middlewares(middlewares ...func(next http.Handler) http.Handler) Middleware {
	return Middlewares(m, func(next http.Handler) http.Handler {
		for i := range middlewares {
			next = middlewares[i](next)
		}

		return next
	})
}
