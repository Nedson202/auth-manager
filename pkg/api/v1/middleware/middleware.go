package middleware

import (
	"context"
)

// Handler _
type Handler func(ctx context.Context, req interface{}) (interface{}, error)

// Middleware _
type Middleware func(Handler) Handler

func ChainMiddleware(m ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}
