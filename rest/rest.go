package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"clouditor.io/clouditor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

//go:generate protoc -I ../proto -I ../third_party auth.proto --grpc-gateway_out=../ --grpc-gateway_opt logtostderr=true

func RunServer(ctx context.Context, grpcPort int, httpPort int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := clouditor.RegisterAuthenticationHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), opts); err != nil {
		return fmt.Errorf("failed to connect to authentication gRPC service %w", err)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: mux,
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
		}

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()

	log.Printf("Starting REST gateway on :%d ...\n", httpPort)

	return srv.ListenAndServe()
}
