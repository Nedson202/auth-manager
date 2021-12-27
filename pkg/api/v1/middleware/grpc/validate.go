package middleware

import (
	"context"

	"github.com/nedson202/auth-manager/pkg/api/v1/middleware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type validator interface {
	Validate() error
}

// Validator is a validator middleware.
func Validator() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if v, ok := req.(validator); ok {
				if err := v.Validate(); err != nil {
					return nil, status.Error(codes.Code(400), "validator-> "+err.Error())
				}
			}
			return handler(ctx, req)
		}
	}
}
