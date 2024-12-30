package middleware

import "net/http"

type Middleware func(http.HandlerFunc) http.HandlerFunc

func MiddlewareChain(middlewares ...Middleware) Middleware {
	return func(finalHandler http.HandlerFunc) http.HandlerFunc {
		// Wrap the final handler with each middleware
		handler := finalHandler
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}
		return handler
	}
}
