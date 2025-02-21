package handlers

import (
	"context"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func NewMiddleWareChain(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}

func ContextMW(handler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			r = r.WithContext(ctx)
			handler.ServeHTTP(w, r)
		})
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
