package v1

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

	go_grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	v1 "github.com/nedson202/auth-manager/api/proto/v1"
	"github.com/nedson202/auth-manager/pkg/api/v1/middleware"
	grpc_middleware "github.com/nedson202/auth-manager/pkg/api/v1/middleware/grpc"
	"github.com/nedson202/auth-manager/pkg/logger"
	"google.golang.org/grpc"
)

func unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return handler(ctx, req)
		}
		h = middleware.ChainMiddleware(grpc_middleware.Validator())(h)
		reply, err := h(ctx, req)
		return reply, err
	}
}

// RunServer runs gRPC service to publish Auth service
func RunGrpcServer(ctx context.Context, v1API v1.AuthServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// gRPC server statup options
	opts := []grpc.ServerOption{}

	// add middleware
	unaryServerInterceptors, streamServerInterceptors := grpc_middleware.AddLogging(logger.Log)
	unaryServerInterceptors = append(unaryServerInterceptors, unaryServerInterceptor())

	opts = append(opts, go_grpc_middleware.WithUnaryServerChain(unaryServerInterceptors...))
	opts = append(opts, go_grpc_middleware.WithStreamServerChain(streamServerInterceptors...))

	// register service
	server := grpc.NewServer(opts...)
	v1.RegisterAuthServiceServer(server, v1API)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			logger.Log.Warn("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	message := fmt.Sprintf("starting gRPC server on port: %s...", port)
	logger.Log.Info(message)
	return server.Serve(listen)
}
