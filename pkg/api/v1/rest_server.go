package v1

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	v1 "github.com/nedson202/auth-manager/api/proto/v1"
	middleware "github.com/nedson202/auth-manager/pkg/api/v1/middleware/rest"
	"github.com/nedson202/auth-manager/pkg/logger"
)

// RunServer runs HTTP/REST gateway
func RunRestServer(ctx context.Context, grpcPort, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	rmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := v1.RegisterAuthServiceHandlerFromEndpoint(ctx, rmux, "localhost:"+grpcPort, opts); err != nil {
		logger.Log.Fatal("failed to start HTTP gateway", zap.String("reason", err.Error()))
	}

	srv := &http.Server{
		Addr: ":" + httpPort,
		Handler: middleware.AddRequestID(
			middleware.AddLogger(logger.Log, rmux)),
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			_ = srv.Shutdown(ctx)
		}

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()

	message := fmt.Sprintf("starting HTTP/REST gateway on port: %s...", httpPort)
	logger.Log.Info(message)
	return srv.ListenAndServe()
}
