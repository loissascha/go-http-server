package server

import "net/http"

func chainMiddleware(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	h = recoverMiddleware(h)
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	h = corsMiddleware(h)
	return h
}
